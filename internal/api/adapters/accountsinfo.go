package adapters

import (
	"database/sql"
	"time"

	acc "github.com/Ruscigno/ticker-signals/internal/transaction/accountsinfo"
	"github.com/Ruscigno/ticker-signals/internal/utils"

	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"
)

// ProtoToAccountInfo transform an proto to an account model
// TODO: try to use proto.Marshal and Unmarshal
func ProtoToAccountInfo(a *v1.Account) *acc.AccountInfo {
	rr := &acc.AccountInfo{
		AccountID:           a.AccountId,
		Balance:             a.Balance,
		Credit:              a.Credit,
		Profit:              a.Profit,
		Equity:              a.Equity,
		Margin:              a.Margin,
		FreeMargin:          a.FreeMargin,
		MarginLevel:         a.MarginLevel,
		MarginCall:          a.MarginCall,
		MarginStopout:       a.MarginStopout,
		TimeTradeServer:     sql.NullTime{Time: time.Unix(a.TimeTradeServer, 0), Valid: a.TimeTradeServer > utils.MT5_MIN_DATE},
		TimeCurrent:         sql.NullTime{Time: time.Unix(a.TimeCurrent, 0), Valid: a.TimeCurrent > utils.MT5_MIN_DATE},
		TimeLocal:           sql.NullTime{Time: time.Unix(a.TimeLocal, 0), Valid: a.TimeLocal > utils.MT5_MIN_DATE},
		TimeGMT:             sql.NullTime{Time: time.Unix(a.TimeGMT, 0), Valid: a.TimeGMT > utils.MT5_MIN_DATE},
		LocalTimeGMTOffset:  sql.NullInt64{Int64: a.LocalTimeGMTOffset, Valid: true},
		ServerTimeGMTOffset: sql.NullInt64{Int64: a.ServerTimeGMTOffset, Valid: true},
	}
	return rr
}
