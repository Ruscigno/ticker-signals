package signalrepo

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"strings"
	"time"

	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"
	"github.com/Ruscigno/ticker-signals/internal/tickerbeats/signal"
	tm "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction"
	tt "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction"
	"github.com/Ruscigno/ticker-signals/internal/utils"
	"github.com/blendle/zapdriver"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	ss "github.com/Ruscigno/ticker-signals/internal/tickerbeats/signal"
)

// NewSignalRepository creates a service to interact with PostgreSQL
func NewSignalRepository(ctx context.Context, dbCon *sqlx.DB) SignalRepository {
	return &signalRepository{
		ctx:   ctx,
		dbCon: dbCon,
	}
}

type signalRepository struct {
	ctx   context.Context
	dbCon *sqlx.DB
}

func (c *signalRepository) GetSignalByDestination(accountID int64) ([]*signal.Signal, error) {
	return c.getBySourceOrDestination(accountID, "destinationaccountid")
}

func (c *signalRepository) getBySourceOrDestination(accountID int64, filterField string) ([]*signal.Signal, error) {
	const SelectQuery string = "select %s from tickerbeats.signals where %s = $1 and active = true"

	result := []*signal.Signal{}
	_, fields := utils.StructToSlice(signal.Signal{}, nil)
	query := fmt.Sprintf(SelectQuery, strings.Join(fields[:], ","), filterField)
	//zap.L().Debug("getBySourceOrDestination query", zap.String("query", query))
	err := c.dbCon.Select(&result, query, accountID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return result, nil
		}
		return nil, err
	}
	return result, nil
}

func (c *signalRepository) ConfirmByExternalID(destinationAccountID, externalid, tickerBeatsID, positionID int64, status ss.SignalStatusEnum) error {
	const UpdateQuery string = `update tickerbeats.signalsresult set 
			%s
			signalstatus=signalstatus+$2,
			confirmationtime=$3,
			updated=$4
		where destinationaccountid=$1 
			  and externalid=$5;
	`
	const SelectQuery string = `
		SELECT signalstatus, externalid
		FROM tickerbeats.signalsresult s2
		where s2.destinationAccountID = $1
			  and s2.externalid = $2
		order by s2.externalid
	`
	resultSet := []*ss.SignalResult{}
	err := c.dbCon.Select(&resultSet, SelectQuery, destinationAccountID, externalid)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil
		}
		return err
	}
	now := time.Now().UTC()
	beatsID := sql.NullInt64{Int64: tickerBeatsID, Valid: tickerBeatsID > 0 && status != ss.DealConfirmed}
	var posID string = ""
	if positionID > 0 {
		posID = fmt.Sprintf("DestinationPositionID=%d,", positionID)
	}
	var sql string
	if beatsID.Valid {
		sql = fmt.Sprintf(UpdateQuery, posID+fmt.Sprintf("destinationbeatsid=%d,", tickerBeatsID))
	} else {
		sql = strings.ReplaceAll(UpdateQuery, "%s", posID)
	}
	for _, r := range resultSet {
		enabled := signal.CalcSignalStatusEnumEnabled(ss.SignalStatusEnum(r.SignalStatus))
		if enabled[status] {
			continue
		}
		rs, err := c.dbCon.Exec(sql, destinationAccountID, status, now, now, externalid)
		if err != nil {
			return err
		}
		ra, err := rs.RowsAffected()
		if err != nil {
			return err
		}
		if ra == 0 {
			continue
		}
		zap.L().Info("TickerBeats confirmed",
			zap.Int64("RowsAffected", ra),
			zap.Int64("destinationaccountid", r.DestinationAccountID),
			zap.Int64("externalid", r.ExternalID),
			zap.Int64("tickerBeatsID", tickerBeatsID),
			zap.Int64("positionID", positionID),
		)
	}
	return nil
}

func (c *signalRepository) UpdateStatus(sourceAccountID, destinationAccountID int64, status signal.SignalStatusEnum, groupID string) error {
	const UpdateQuery string = "update tickerbeats.signalsresult set signalstatus=signalstatus+$1,updated=$2,senttime=$3 where externalid=$4 and signalstatus+$1 < 1024"
	const SelectQuery string = `
		SELECT signalstatus, externalid
		FROM tickerbeats.signalsresult s2
		where s2.sourceaccountid = $1
			  and s2.destinationAccountID = $2
			  and s2.groupid = $3
		order by s2.externalid
	`
	resultSet := []*ss.SignalResult{}
	err := c.dbCon.Select(&resultSet, SelectQuery, sourceAccountID, destinationAccountID, groupID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil
		}
		return err
	}
	var raTotal int64 = 0
	now := time.Now().UTC()
	sentTime := sql.NullTime{Valid: status == signal.Sent, Time: now}
	for _, r := range resultSet {
		enabled := signal.CalcSignalStatusEnumEnabled(status)
		if !enabled[status] {
			continue
		}
		rs, err := c.dbCon.Exec(UpdateQuery, int(status), now, sentTime, r.ExternalID)
		if err != nil {
			return err
		}
		ra, err := rs.RowsAffected()
		if err != nil {
			return err
		}
		raTotal += ra
	}

	if raTotal < 1 {
		return nil
	}
	zap.L().Info("TickerBeats status updated",
		zap.Int64("RowsAffected", raTotal),
		zap.Int64("sourceAccountID", sourceAccountID),
		zap.Int64("destinationAccountID", destinationAccountID),
		zap.String("groupID", groupID),
	)
	return nil
}

func (c *signalRepository) CreateTickerBeats(sourceAccountID, destinationAccountID int64, groupID string, minToExpire int32) (int64, error) {
	const CreateQuery string = `
		insert into	tickerbeats.signalsresult(sourceaccountid,destinationaccountid,sourcebeatsid,signaltype,signalstatus,externalid,groupid,expireat,entry,updated,created,SourcePositionID)
			select distinct req.accountid,%d,req.orderid,'%s',%d,nextval('tickerbeats.tickerbeats_seq'),'%s', %s, d.entry, %s, %s, d.positionid
			from tickerbeats.traderequests req,
				 tickerbeats.signals s,
				 tickerbeats.tradetransactions tt,
				 tickerbeats.deals d,
				 tickerbeats.positions p
			where 
	              not exists (select 1 from tickerbeats.signalsresult sr 
							  where sr.sourceaccountid=%d and 
							        sr.destinationaccountid=%d and 
									sr.sourcebeatsid=req.orderid and 
									sr.signaltype='%s'
								)
				  and s.sourceaccountid = %d
				  and s.destinationaccountid = %d
				  and s.sourceaccountid = req.accountid
				  and tt.accountid = req.accountid
				  and tt.orderid = req.orderid
				  and tt.tradetype = '%s'
				  and p.accountid = req.accountid
				  and p.ticket = req.orderid
				  and p.deleted is null
				  and d.accountid = tt.accountid
				  and d.ticket = tt.dealid
				  and d.entry = '%s'
				  and req.created > %s
				  and d.deleted is null
				  and tt.deleted is null
				  and req.deleted is null;
	`
	processingTime := time.Now().UTC()
	sNow := fmt.Sprintf("to_timestamp('%s', 'YYYYMMDD HH24:MI:SS')", processingTime.Format("20060102 15:04:05"))
	// 1min aditional for the " ping-pong processing time"
	expireTime := processingTime.Add(time.Minute * time.Duration(int32(minToExpire/2))).Add(time.Minute * 1)
	expireAt := fmt.Sprintf("to_timestamp('%s', 'YYYYMMDD HH24:MI:SS')", expireTime.Format("20060102 15:04:05"))
	startAt := fmt.Sprintf("to_timestamp('%s', 'YYYYMMDD HH24:MI:SS')", processingTime.Add(-1*time.Minute*time.Duration(minToExpire)).Format("20060102 15:04:05"))
	sql := fmt.Sprintf(CreateQuery,
		destinationAccountID,
		v1.SignalType_name[int32(v1.SignalType_SIGINAL_TYPE_TRADE_REQUEST)],
		signal.Created,
		groupID,
		expireAt,
		sNow,
		sNow,
		sourceAccountID,
		destinationAccountID,
		v1.SignalType_name[int32(v1.SignalType_SIGINAL_TYPE_TRADE_REQUEST)],
		sourceAccountID,
		destinationAccountID,
		v1.TradeTransactionType_name[int32(v1.TradeTransactionType_TRADE_TRANSACTION_DEAL_ADD)],
		v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
		startAt,
	)
	rs, err := c.dbCon.Exec(sql)
	if err != nil {
		zap.L().Error("SignalService: CreateTickerBeats Error",
			zap.Error(err),
			zap.String("sql", sql),
			zapdriver.SourceLocation(runtime.Caller(0)))
		return 0, err
	}
	ra, err := rs.RowsAffected()
	if destinationAccountID == 64596499 {
		zap.L().Debug("CreateTickerBeats",
			zap.String("insert SQL", sql),
			zap.Int64("RowsAffected", ra),
		)
	}
	return ra, err
}

func (c *signalRepository) GetTradeRequesByGroupID(destinationAccountID int64, groupID string, entry v1.DealEntry) ([]*tt.TradeRequest, error) {
	const SelectQuery string = `
	    SELECT %d as accountid, req.action, s2.externalid as magic, req.symbol, 
		       req.volume, o.pricecurrent as price, req.stoplimit, req.stoploss, 
			   req.takeprofit, req.deviation, req.ordertype, req.typefilling, 
			   req.typetime, req.timeexpiration, s2.entry, 
			   s2.destinationpositionid as positionid
		FROM tickerbeats.traderequests req,
		     tickerbeats.signalsresult s2,
			 tickerbeats.orders o
		where s2.sourceaccountid = req.accountid 
			  and s2.sourcebeatsid = req.orderid
			  and req.accountid = o.accountid
			  and req.orderid = o.ticket
			  and s2.destinationaccountid = $1
			  and s2.groupid = $2
			  and req.deleted is null
			  and o.deleted is null
			  and s2.entry = $3
			  and s2.signalstatus = $4
		order by s2.externalid
	`
	sql := fmt.Sprintf(SelectQuery, destinationAccountID)
	result := []*tt.TradeRequest{}
	err := c.dbCon.Select(&result, sql, destinationAccountID, groupID, v1.DealEntry_name[int32(entry)], ss.Created)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return result, nil
		}
		return nil, err
	}
	return result, nil
}

func (c *signalRepository) RemoveDuplicatedSignals(destinationAccountID int64, groupID string) error {
	UpdateQuery := `
		update
			tickerbeats.signalsresult t1
		set
			signalstatus = signalstatus+$3
		from
			(select	b.sourceaccountid,b.destinationaccountid,b.sourcepositionid
			from tickerbeats.signalsresult b where b.destinationaccountid = $1 and b.groupid = $2
			group by b.sourceaccountid,b.destinationaccountid,b.sourcepositionid
			having	count(*) > 1) as t2
		where t1.sourceaccountid = t2.sourceaccountid
			and t1.destinationaccountid = t2.destinationaccountid
			and t1.sourcepositionid = t2.sourcepositionid
	`
	_, err := c.dbCon.Exec(UpdateQuery, destinationAccountID, groupID, ss.Expired)
	return err
}

func (c *signalRepository) UpdatePositionIdBeforeClose(sourceAccountID, destinationAccountID int64, groupID string) error {
	UpdateQuery := `
		update tickerbeats.signalsresult ss set destinationpositionid = tt.destinationpositionid
		from (select sr.destinationpositionid, sr.sourceaccountid, sr.sourcepositionid, sr.destinationaccountid from tickerbeats.signalsresult sr) as tt
		where ss.sourceaccountid = $1
			and ss.destinationaccountid = $2
			and ss.groupid = $3
			and ss.destinationpositionid is null
			and tt.sourceaccountid = ss.sourceaccountid
			and tt.destinationaccountid = ss.destinationaccountid
			and tt.sourcepositionid = ss.sourcepositionid
			and tt.destinationpositionid is not null
	`
	_, err := c.dbCon.Exec(UpdateQuery, sourceAccountID, destinationAccountID, groupID)
	if err != nil {
		return err
	}
	return c.deleteInvalidClosePositionBeats(sourceAccountID, destinationAccountID, groupID)
}

func (c *signalRepository) deleteInvalidClosePositionBeats(sourceAccountID, destinationAccountID int64, groupID string) error {
	DeleteQuery := `
		delete from tickerbeats.signalsresult ss
		where exists (select 1 from tickerbeats.signalsresult s2
					  where s2.sourceaccountid = ss.sourceaccountid
					    and s2.destinationaccountid = ss.destinationaccountid
					    and s2.sourcepositionid = ss.sourcepositionid
					  group by s2.sourceaccountid, s2.destinationaccountid, s2.sourcepositionid
					  having count(1) > 1
					 )
		  and ss.sourceaccountid = $1
		  and ss.destinationaccountid = $2
		  and ss.groupid = $3
		  and ss.signalstatus = $4
		  and ss.entry = $5
		  and ss.destinationpositionid is null
	`
	zap.L().Info("deleteInvalidClosePositionBeats",
		zap.String("DeleteQuery", DeleteQuery),
		zap.Int64("sourceAccountID", sourceAccountID),
		zap.Int64("destinationAccountID", destinationAccountID),
		zap.String("groupID", groupID),
		zap.Int64("ss.Created", int64(ss.Created)),
		zap.Int32("DealEntry", int32(v1.DealEntry_DEAL_ENTRY_OUT)),
	)
	_, err := c.dbCon.Exec(DeleteQuery, sourceAccountID, destinationAccountID, groupID, ss.Created, v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_OUT)])
	return err
}

func (c *signalRepository) CloseDeadPositions(sourceAccountID, destinationAccountID int64, groupID string, minToExpire int32) ([]*tt.TradeRequest, error) {
	const InsertQuery string = `
	INSERT INTO tickerbeats.signalsresult(sourceaccountid, destinationaccountid, sourcebeatsid, sourcepositionid, destinationpositionid, signaltype, signalstatus, externalid, groupid, entry, expireat, created, updated)
		select s1.sourceaccountid, s1.destinationaccountid, d1.ticket, s1.sourcepositionid, s1.destinationpositionid, '%s', %d, %s, '%s', '%s', %s, %s, %s
		from tickerbeats.positions p1,
			tickerbeats.signalsresult s1
			left join tickerbeats.deals d1 on d1.accountid = s1.sourceaccountid and d1.positionid = s1.sourcepositionid and d1.entry = '%s' and d1.deleted is null
		where s1.destinationpositionid is not null
		  and s1.sourceaccountid = %d
		  and s1.destinationaccountid = %d
		  and s1.destinationpositionid = p1.ticket
		  and s1.destinationaccountid = p1.accountid
		  and p1.deleted is null
		  and exists (select 1 from tickerbeats.positions p2 where p2.accountid = s1.sourceaccountid and p2.ticket = s1.sourcepositionid and p2.deleted is not null)
		  and not exists (select 1 from tickerbeats.signalsresult s2 where s2.sourceaccountid = s1.sourceaccountid 
																	   and s2.destinationaccountid = s1.destinationaccountid
																	   and s2.sourcepositionid = s1.sourcepositionid
																	   and s2.destinationpositionid = s1.destinationpositionid
																	   and s2.entry = '%s')
`
	processingTime := time.Now().UTC()
	sNow := fmt.Sprintf("to_timestamp('%s', 'YYYYMMDD HH24:MI:SS')", processingTime.Format("20060102 15:04:05"))
	// 1min aditional for the " ping-pong processing time"
	expireTime := processingTime.Add(time.Minute * time.Duration(minToExpire)).Add(time.Minute * 1)
	expireAt := fmt.Sprintf("to_timestamp('%s', 'YYYYMMDD HH24:MI:SS')", expireTime.Format("20060102 15:04:05"))
	sql := fmt.Sprintf(InsertQuery,
		v1.SignalType_name[int32(v1.SignalType_SIGINAL_TYPE_POSITION)],
		signal.Created,
		"nextval('tickerbeats.tickerbeats_seq')",
		groupID,
		v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_OUT)],
		expireAt,
		sNow,
		sNow,
		v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_OUT)],
		sourceAccountID,
		destinationAccountID,
		v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_OUT)],
	)
	rs, err := c.dbCon.Exec(sql)
	if err != nil {
		zap.L().Error("SignalService: CloseDeadPositions Error",
			zap.Error(err),
			zap.String("sql", sql),
			zapdriver.SourceLocation(runtime.Caller(0)))
		return []*tt.TradeRequest{}, err
	}
	ra, err := rs.RowsAffected()
	if err != nil {
		zap.L().Error("SignalService: CloseDeadPositions RowsAffected Error",
			zap.Error(err),
			zap.String("sql", sql),
			zapdriver.SourceLocation(runtime.Caller(0)))
		return []*tt.TradeRequest{}, err
	}
	if ra == 0 {
		return []*tt.TradeRequest{}, err
	}
	return c.GetByGroupID(destinationAccountID, groupID, v1.DealEntry_DEAL_ENTRY_OUT, v1.SignalType_SIGINAL_TYPE_POSITION)
}

func (c *signalRepository) GetByGroupID(destinationAccountID int64, groupID string, entry v1.DealEntry, stype v1.SignalType) ([]*tt.TradeRequest, error) {
	const SelectQuery string = `
	    SELECT sourceaccountid, destinationaccountid, sourcebeatsid, destinationbeatsid, sourcepositionid, destinationpositionid, 
		       signaltype, signalstatus, externalid, groupid, entry, senttime, confirmationtime, expireat, created, updated
		FROM tickerbeats.signalsresult s2
		where s2.destinationAccountID = $1
			  and s2.groupid = $2
			  and s2.entry = $3
			  and s2.signaltype = $4
		order by s2.externalid
	`
	resultSet := []*ss.SignalResult{}
	entryName := v1.DealEntry_name[int32(entry)]
	err := c.dbCon.Select(&resultSet, SelectQuery, destinationAccountID, groupID, entryName, v1.SignalType_name[int32(stype)])
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return []*tt.TradeRequest{}, nil
		}
		return []*tt.TradeRequest{}, err
	}
	result := make([]*tt.TradeRequest, len(resultSet))
	for i, r := range resultSet {
		result[i] = &tt.TradeRequest{
			AccountID:  r.DestinationAccountID,
			Magic:      sql.NullInt64{Int64: r.ExternalID, Valid: true},
			PositionID: r.DestinationPositionID,
			Entry:      entryName,
		}
	}
	return result, nil
}

func (c *signalRepository) NeedToCloseAllPositions(accountID int64, groupID string, stop int32) ([]*tm.TradeRequest, error) {
	// const SelectQuery string = "select %s from tickerbeats.positions where accountID = $1 and ticket = $2"
	// pos := pp.Position{}
	// _, fields := utils.StructToSlice(pos, false)
	// query := fmt.Sprintf(SelectQuery, strings.Join(fields[:], ","))
	// err := c.dbCon.Get(&pos, query, accountID, ticket)
	// if err != nil {
	// 	if strings.Contains(err.Error(), "no rows in result set") {
	// 		return &pos, nil
	// 	}
	// 	return nil, err
	// }
	return []*tm.TradeRequest{}, nil
}
