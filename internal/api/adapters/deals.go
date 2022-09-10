package adapters

import (
	"database/sql"
	"strings"
	"time"

	dea "github.com/Ruscigno/ticker-signals/internal/transaction/deals"
	"github.com/Ruscigno/ticker-signals/internal/utils"

	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"
)

// DealsToProto transform a list of deals model into a list of proto deals
// TODO: try to use proto.Marshal and Unmarshal
func DealsToProto(deals []*dea.Deal) []*v1.Deal {
	deaProto := make([]*v1.Deal, len(deals))
	for i, deaount := range deals {
		deaProto[i] = DealToProto(deaount)
	}
	return deaProto
}

// DealToProto transform an dea.Deal model into an proto v1.Deal
// TODO: try to use proto.Marshal and Unmarshal
func DealToProto(a *dea.Deal) *v1.Deal {
	var deleted int64
	if a.Deleted.Valid {
		deleted = a.Deleted.Time.Unix()
	}
	magic := a.Magic.Int64
	if !a.Magic.Valid {
		magic = 0
	}
	orderid := a.OrderID.Int64
	if !a.OrderID.Valid {
		orderid = 0
	}
	commission := a.Commission.Float64
	if !a.Commission.Valid {
		commission = 0
	}
	Profit := a.Profit.Float64
	if !a.Profit.Valid {
		Profit = 0
	}
	DealFee := a.DealFee.Float64
	if !a.DealFee.Valid {
		DealFee = 0
	}
	Symbol := a.Symbol.String
	if !a.Symbol.Valid {
		Symbol = ""
	}
	return &v1.Deal{
		DealId:     a.DealID,
		AccountId:  a.AccountID,
		Ticket:     a.Ticket,
		Magic:      magic,
		OrderId:    orderid,
		Symbol:     Symbol,
		DealTime:   a.DealTime.UnixNano(),
		DealType:   v1.DealType(v1.DealType_value[a.DealType]),
		Entry:      v1.DealEntry(v1.DealEntry_value[a.Entry]),
		PositionId: a.PositionId,
		Volume:     a.Volume,
		Price:      a.Price,
		Commission: commission,
		Swap:       a.Swap.Float64,
		Profit:     Profit,
		Comment:    a.Comment,
		ExternalId: a.ExternalId,
		Created:    a.Created.Unix(),
		Updated:    a.Updated.Unix(),
		Deleted:    deleted,
		Reason:     a.Reason,
		DealFee:    DealFee,
	}
}

// ProtoToDeal transform an proto to an deaount model
// TODO: try to use proto.Marshal and Unmarshal
func ProtoToDeal(a *v1.Deal) *dea.Deal {
	return &dea.Deal{
		DealID:     a.DealId,
		AccountID:  a.AccountId,
		Ticket:     a.Ticket,
		Magic:      sql.NullInt64{Int64: a.Magic, Valid: a.Magic != 0},
		OrderID:    sql.NullInt64{Int64: a.OrderId, Valid: a.OrderId != 0},
		Symbol:     sql.NullString{String: a.Symbol, Valid: strings.TrimSpace(a.Symbol) != ""},
		DealTime:   time.Unix(0, a.DealTime*int64(time.Millisecond)).Add(time.Millisecond * time.Duration(a.TimeGMTOffset)),
		DealType:   v1.DealType_name[int32(a.DealType)],
		Entry:      v1.DealEntry_name[int32(a.Entry)],
		PositionId: a.PositionId,
		Volume:     a.Volume,
		Price:      a.Price,
		Commission: sql.NullFloat64{Float64: a.Commission, Valid: a.Commission != 0},
		Swap:       sql.NullFloat64{Float64: a.Swap, Valid: a.Swap != 0},
		Profit:     sql.NullFloat64{Float64: a.Profit, Valid: a.Profit != 0},
		Comment:    a.Comment,
		ExternalId: a.ExternalId,
		Created:    time.Unix(a.Created*int64(time.Millisecond), 0).Add(time.Millisecond * time.Duration(a.TimeGMTOffset)),
		Updated:    time.Unix(a.Updated*int64(time.Millisecond), 0).Add(time.Millisecond * time.Duration(a.TimeGMTOffset)),
		Deleted:    sql.NullTime{Time: time.Unix(a.Deleted*int64(time.Millisecond), 0).Add(time.Millisecond * time.Duration(a.TimeGMTOffset)), Valid: a.Deleted > utils.MT5_MIN_DATE},
		Reason:     a.Reason,
		DealFee:    sql.NullFloat64{Float64: a.DealFee, Valid: a.DealFee != 0},
	}
}
