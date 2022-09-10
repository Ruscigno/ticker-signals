package adapters

import (
	"database/sql"
	"strings"
	"time"

	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"
	pos "github.com/Ruscigno/ticker-signals/internal/transaction/positions"
	"github.com/Ruscigno/ticker-signals/internal/utils"
)

// PositionsToProto transform a list of positions model into a list of proto positions
// TODO: try to use proto.Marshal and Unmarshal
func PositionsToProto(positions []*pos.Position) []*v1.Position {
	posProto := make([]*v1.Position, len(positions))
	for i, position := range positions {
		posProto[i] = PositionToProto(position)
	}
	return posProto
}

// PositionToProto transform an position model into an proto position
// TODO: try to use proto.Marshal and Unmarshal
func PositionToProto(a *pos.Position) *v1.Position {
	magic := a.Magic.Int64
	if !a.Magic.Valid {
		magic = 0
	}
	Symbol := a.Symbol.String
	if !a.Symbol.Valid {
		Symbol = ""
	}
	result := &v1.Position{
		PositionId:     a.PositionID,
		AccountId:      a.AccountID,
		Ticket:         a.Ticket,
		Symbol:         Symbol,
		PositionTime:   a.PositionTime.Time.UnixNano(),
		PositionType:   v1.PositionType(v1.PositionType_value[a.PositionType]),
		Volume:         a.Volume,
		PriceOpen:      a.PriceOpen,
		StopLoss:       a.StopLoss,
		TakeProfit:     a.TakeProfit,
		PriceCurrent:   a.PriceCurrent,
		Commission:     a.Commission,
		Swap:           a.Swap,
		Profit:         a.Profit,
		Comment:        a.Comment,
		Created:        a.Created.Unix(),
		Updated:        a.Updated.Unix(),
		PositionUpdate: a.PositionUpdate.Time.UnixNano(),
		Reason:         v1.PositionReason(v1.PositionReason_value[a.Reason]),
		ExternalId:     a.ExternalID,
		Magic:          magic,
	}
	if a.Deleted.Valid {
		result.Deleted = a.Deleted.Time.Unix()
	}
	return result
}

// ProtoToPosition transform an proto to an position model
// TODO: try to use proto.Marshal and Unmarshal
func ProtoToPosition(a *v1.Position) *pos.Position {
	positionTime := time.Unix(0, a.PositionTime*int64(time.Millisecond)).Add(time.Millisecond * time.Duration(a.TimeGMTOffset))
	created := positionTime
	if a.Created > utils.MT5_MIN_DATE {
		created = time.Unix(a.Created*int64(time.Millisecond), 0)
	}
	return &pos.Position{
		PositionID:     a.Identifier,
		AccountID:      a.AccountId,
		Ticket:         a.Ticket,
		Symbol:         sql.NullString{String: a.Symbol, Valid: strings.TrimSpace(a.Symbol) != ""},
		PositionTime:   sql.NullTime{Time: positionTime, Valid: a.PositionTime > utils.MT5_MIN_DATE},
		PositionType:   v1.PositionType_name[int32(a.PositionType)],
		Volume:         a.Volume,
		PriceOpen:      a.PriceOpen,
		StopLoss:       a.StopLoss,
		TakeProfit:     a.TakeProfit,
		PriceCurrent:   a.PriceCurrent,
		Commission:     a.Commission,
		Swap:           a.Swap,
		Profit:         a.Profit,
		Comment:        a.Comment,
		Created:        created,
		Updated:        time.Unix(a.Updated, 0),
		Deleted:        sql.NullTime{Time: time.Unix(0, a.Deleted*int64(time.Millisecond)).Add(time.Millisecond * time.Duration(a.TimeGMTOffset)), Valid: a.Deleted > utils.MT5_MIN_DATE},
		PositionUpdate: sql.NullTime{Time: time.Unix(0, a.PositionUpdate*int64(time.Millisecond)).Add(time.Millisecond * time.Duration(a.TimeGMTOffset)), Valid: a.PositionTime > utils.MT5_MIN_DATE},
		Reason:         v1.PositionReason_name[int32(a.Reason)],
		ExternalID:     a.ExternalId,
		Magic:          sql.NullInt64{Int64: a.Magic, Valid: a.Magic > 0},
	}
}
