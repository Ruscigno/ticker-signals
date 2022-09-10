package adapters

import (
	"database/sql"
	"strings"
	"time"

	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"
	ord "github.com/Ruscigno/ticker-signals/internal/transaction/orders"
	"github.com/Ruscigno/ticker-signals/internal/utils"
)

// OrdersToProto transform a list of orders model into a list of proto orders
// TODO: try to use proto.Marshal and Unmarshal
func OrdersToProto(orders []*ord.Order) []*v1.Order {
	ordProto := make([]*v1.Order, len(orders))
	for i, order := range orders {
		ordProto[i] = OrderToProto(order)
	}
	return ordProto
}

// OrderToProto transform an order model into an proto order
// TODO: try to use proto.Marshal and Unmarshal
func OrderToProto(a *ord.Order) *v1.Order {
	magic := a.Magic.Int64
	if !a.Magic.Valid {
		magic = 0
	}
	Symbol := a.Symbol.String
	if !a.Symbol.Valid {
		Symbol = ""
	}
	result := &v1.Order{
		OrderId:        a.OrderID,
		AccountId:      a.AccountID,
		Ticket:         a.Ticket,
		Symbol:         Symbol,
		TimeSetup:      a.TimeSetup.Time.UnixNano(),
		OrderType:      v1.OrderType(v1.OrderType_value[a.OrderType]),
		State:          v1.OrderState(v1.OrderState_value[a.State]),
		TimeExpiration: a.TimeExpiration.Time.Unix(),
		TimeDone:       a.TimeDone.Time.UnixNano(),
		TypeFilling:    v1.OrderFillingType(v1.OrderFillingType_value[a.TypeFilling]),
		TypeTime:       v1.OrderTypeTime(v1.OrderTypeTime_value[a.TypeTime]),
		Magic:          magic,
		PositionId:     a.PositionId,
		VolumeInitial:  a.VolumeInitial,
		VolumeCurrent:  a.VolumeCurrent.Float64,
		PriceOpen:      a.PriceOpen,
		StopLoss:       a.StopLoss.Float64,
		TakeProfit:     a.TakeProfit.Float64,
		PriceCurrent:   a.PriceCurrent,
		PriceStopLimit: a.PriceStopLimit.Float64,
		Comment:        a.Comment,
		ExternalId:     a.ExternalID,
		Reason:         v1.OrderReason(v1.OrderReason_value[a.Reason]),
		Created:        a.Created.Unix(),
		Updated:        a.Updated.Unix(),
		PositionById:   a.PositionByID.Int64,
		// Direction:      a.Direction,
	}
	if a.Deleted.Valid {
		result.Deleted = a.Deleted.Time.Unix()
	}
	return result
}

// ProtoToOrder transform an proto to an order model
// TODO: try to use proto.Marshal and Unmarshal
func ProtoToOrder(a *v1.Order) *ord.Order {
	var deleted sql.NullTime
	if a.Deleted > 0 {
		deleted.Time = time.Unix(a.Deleted, 0)
		deleted.Valid = true
	}

	ord := &ord.Order{
		OrderID:        a.OrderId,
		AccountID:      a.AccountId,
		Ticket:         a.Ticket,
		Symbol:         sql.NullString{String: a.Symbol, Valid: strings.TrimSpace(a.Symbol) != ""},
		TimeSetup:      sql.NullTime{Time: time.Unix(0, a.TimeSetup*int64(time.Millisecond)).Add(time.Millisecond * time.Duration(a.TimeGMTOffset)), Valid: a.TimeSetup > utils.MT5_MIN_DATE},
		OrderType:      v1.OrderType_name[int32(a.OrderType)],
		State:          v1.OrderState_name[int32(a.State)],
		TypeFilling:    v1.OrderFillingType_name[int32(a.TypeFilling)],
		TypeTime:       v1.OrderTypeTime_name[int32(a.TypeTime)],
		Reason:         v1.OrderReason_name[int32(a.Reason)],
		TimeExpiration: sql.NullTime{Time: time.Unix(a.TimeExpiration, 0).Add(time.Millisecond * time.Duration(a.TimeGMTOffset)), Valid: a.TimeExpiration > utils.MT5_MIN_DATE},
		TimeDone:       sql.NullTime{Time: time.Unix(0, a.TimeDone*int64(time.Millisecond)).Add(time.Millisecond * time.Duration(a.TimeGMTOffset)), Valid: a.TimeDone > utils.MT5_MIN_DATE},
		Magic:          sql.NullInt64{Int64: a.Magic, Valid: a.Magic != 0},
		PositionId:     a.PositionId,
		VolumeInitial:  a.VolumeInitial,
		VolumeCurrent:  sql.NullFloat64{Float64: a.VolumeCurrent, Valid: a.VolumeCurrent != 0},
		PriceOpen:      a.PriceOpen,
		StopLoss:       sql.NullFloat64{Float64: a.StopLoss, Valid: a.StopLoss != 0},
		TakeProfit:     sql.NullFloat64{Float64: a.TakeProfit, Valid: a.TakeProfit != 0},
		PriceCurrent:   a.PriceCurrent,
		PriceStopLimit: sql.NullFloat64{Float64: a.PriceStopLimit, Valid: a.PriceStopLimit != 0},
		Comment:        a.Comment,
		ExternalID:     a.ExternalId,
		Created:        time.Unix(a.Created, 0).Add(time.Millisecond * time.Duration(a.TimeGMTOffset)),
		Updated:        time.Unix(a.Updated, 0).Add(time.Millisecond * time.Duration(a.TimeGMTOffset)),
		Deleted:        deleted,
		PositionByID:   sql.NullInt64{Int64: a.PositionById, Valid: a.PositionById != 0},
	}
	return ord
}
