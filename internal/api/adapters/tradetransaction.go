package adapters

import (
	"database/sql"
	"strings"
	"time"

	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"
	model "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction"
	"github.com/Ruscigno/ticker-signals/internal/utils"
)

// ProtoToTradeTransaction transform an proto to an TradeTransaction model
func ProtoToTradeTransaction(accountID, creationOrder int64, a *v1.TradeTransaction) *model.TradeTransaction {
	if a == nil {
		return &model.TradeTransaction{AccountID: -1}
	}
	return &model.TradeTransaction{
		InternalID:      a.InternalId,
		AccountID:       accountID,
		OrderID:         a.OrderId,
		CreationOrder:   creationOrder,
		DealID:          sql.NullInt64{Int64: a.DealId, Valid: a.DealId > 0},
		Symbol:          sql.NullString{String: a.Symbol, Valid: strings.TrimSpace(a.Symbol) != ""},
		TradeType:       v1.TradeTransactionType_name[int32(a.TradeType)],
		OrderType:       v1.OrderType_name[int32(a.OrderType)],
		OrderState:      v1.OrderState_name[int32(a.OrderState)],
		DealType:        v1.DealType_name[int32(a.DealType)],
		TimeType:        v1.OrderTypeTime_name[int32(a.TimeType)],
		TimeExpiration:  sql.NullTime{Time: time.Unix(a.TimeExpiration, 0).Add(time.Millisecond * time.Duration(a.TimeGMTOffset)), Valid: a.TimeExpiration > utils.MT5_MIN_DATE},
		Price:           sql.NullFloat64{Float64: a.Price, Valid: a.Price != 0},
		PriceTrigger:    sql.NullFloat64{Float64: a.PriceTrigger, Valid: a.PriceTrigger != 0},
		PriceStopLoss:   sql.NullFloat64{Float64: a.PriceStopLoss, Valid: a.PriceStopLoss != 0},
		PriceTakeProfit: sql.NullFloat64{Float64: a.PriceStopLoss, Valid: a.PriceStopLoss != 0},
		Volume:          sql.NullFloat64{Float64: a.Volume, Valid: a.Volume != 0},
		PositionID:      sql.NullInt64{Int64: a.PositionId, Valid: a.PositionId > 0},
		PositionBy:      sql.NullInt64{Int64: a.PositionBy, Valid: a.PositionBy != 0},
		Created:         time.Unix(a.Created*int64(time.Millisecond), 0).Add(time.Millisecond * time.Duration(a.TimeGMTOffset)),
		Updated:         time.Unix(a.Updated*int64(time.Millisecond), 0).Add(time.Millisecond * time.Duration(a.TimeGMTOffset)),
		Deleted:         sql.NullTime{Time: time.Unix(a.Deleted*int64(time.Millisecond), 0).Add(time.Millisecond * time.Duration(a.TimeGMTOffset)), Valid: a.Deleted > utils.MT5_MIN_DATE},
	}
}

// ProtoToTradeTransaction transform an proto to an TradeTransaction model
func ProtoToTradeResult(accountID, creationOrder int64, a *v1.TradeResult) *model.TradeResult {
	if a == nil {
		return &model.TradeResult{AccountID: -1}
	}
	return &model.TradeResult{
		AccountID:       accountID,
		OrderID:         a.OrderId,
		CreationOrder:   creationOrder,
		RetCode:         a.Retcode,
		DealID:          sql.NullInt64{Int64: a.DealId, Valid: a.DealId > 0},
		Volume:          sql.NullFloat64{Float64: a.Volume, Valid: a.Volume > 0},
		Price:           sql.NullFloat64{Float64: a.Price, Valid: a.Price > 0},
		Bid:             sql.NullFloat64{Float64: a.Bid, Valid: a.Bid > 0},
		Ask:             sql.NullFloat64{Float64: a.Ask, Valid: a.Ask > 0},
		Comment:         a.Comment,
		RequestID:       a.RequestId,
		RetcodeExternal: sql.NullInt64{Int64: int64(a.RetcodeExternal), Valid: a.RetcodeExternal > 0},
	}
}

// ProtoToTradeTransaction transform an proto to an TradeTransaction model
func ProtoToTradeRequest(accountID, creationOrder int64, a *v1.TradeRequest) *model.TradeRequest {
	if a == nil {
		return &model.TradeRequest{AccountID: -1}
	}
	return &model.TradeRequest{
		AccountID:      accountID,
		OrderID:        a.OrderId,
		CreationOrder:  creationOrder,
		Action:         v1.TradeRequestActions_name[int32(a.Action)],
		Magic:          sql.NullInt64{Int64: a.Magic, Valid: a.Magic != 0},
		Symbol:         sql.NullString{String: a.Symbol, Valid: strings.TrimSpace(a.Symbol) != ""},
		Volume:         a.Volume,
		Price:          sql.NullFloat64{Float64: a.Price, Valid: a.Price != 0},
		StopLimit:      sql.NullFloat64{Float64: a.StopLimit, Valid: a.StopLimit != 0},
		StopLoss:       sql.NullFloat64{Float64: a.StopLoss, Valid: a.StopLoss != 0},
		TakeProfit:     sql.NullFloat64{Float64: a.TakeProfit, Valid: a.TakeProfit != 0},
		Deviation:      sql.NullInt64{Int64: a.Deviation, Valid: a.Deviation != 0},
		OrderType:      v1.OrderType_name[int32(a.OrderType)],
		TypeFilling:    v1.OrderFillingType_name[int32(a.TypeFilling)],
		TypeTime:       v1.OrderTypeTime_name[int32(a.TypeTime)],
		TimeExpiration: sql.NullTime{Time: time.Unix(a.TimeExpiration, 0).Add(time.Millisecond * time.Duration(a.TimeGMTOffset)), Valid: a.TimeExpiration > utils.MT5_MIN_DATE},
		Comment:        a.Comment,
		PositionID:     sql.NullInt64{Int64: a.PositionId, Valid: a.PositionId > 0},
		PositionBy:     sql.NullInt64{Int64: a.PositionBy, Valid: a.PositionBy > 0},
		Created:        time.Unix(a.Created*int64(time.Millisecond), 0).Add(time.Millisecond * time.Duration(a.TimeGMTOffset)),
		Updated:        time.Unix(a.Updated*int64(time.Millisecond), 0).Add(time.Millisecond * time.Duration(a.TimeGMTOffset)),
		Deleted:        sql.NullTime{Time: time.Unix(a.Deleted*int64(time.Millisecond), 0).Add(time.Millisecond * time.Duration(a.TimeGMTOffset)), Valid: a.Deleted > utils.MT5_MIN_DATE},
	}
}

// ProtoToTradeTransaction transform an proto to an TradeTransaction model
func TradeRequestToProto(a *model.TradeRequest) *v1.TradeRequest {
	deleted := a.Deleted.Time.UTC().Unix()
	if !a.Deleted.Valid {
		deleted = 0
	}
	price := a.Price.Float64
	if !a.Price.Valid {
		price = 0
	}
	stopLimit := a.StopLimit.Float64
	if !a.StopLimit.Valid {
		stopLimit = 0
	}
	stopLoss := a.StopLoss.Float64
	if !a.StopLoss.Valid {
		stopLoss = 0
	}
	takeProfit := a.TakeProfit.Float64
	if !a.TakeProfit.Valid {
		takeProfit = 0
	}
	deviation := a.Deviation.Int64
	if !a.Deviation.Valid {
		deviation = 0
	}
	positionId := a.PositionID.Int64
	if !a.PositionID.Valid {
		positionId = 0
	}
	positionBy := a.PositionBy.Int64
	if !a.PositionBy.Valid {
		positionBy = 0
	}
	magic := a.Magic.Int64
	if !a.Magic.Valid {
		magic = 0
	}
	Symbol := a.Symbol.String
	if !a.Symbol.Valid {
		Symbol = ""
	}
	return &v1.TradeRequest{
		AccountId:      a.AccountID,
		OrderId:        a.OrderID,
		CreationOrder:  a.CreationOrder,
		Action:         v1.TradeRequestActions(v1.TradeRequestActions_value[a.Action]),
		Magic:          magic,
		Symbol:         Symbol,
		Volume:         a.Volume,
		Price:          price,
		StopLimit:      stopLimit,
		StopLoss:       stopLoss,
		TakeProfit:     takeProfit,
		Deviation:      deviation,
		OrderType:      v1.OrderType(v1.OrderType_value[a.OrderType]),
		TypeFilling:    v1.OrderFillingType(v1.OrderFillingType_value[a.TypeFilling]),
		TypeTime:       v1.OrderTypeTime(v1.OrderTypeTime_value[a.TypeTime]),
		TimeExpiration: a.TimeExpiration.Time.UTC().Unix(),
		Comment:        a.Comment,
		PositionId:     positionId,
		PositionBy:     positionBy,
		Created:        a.Created.UTC().Unix(),
		Updated:        a.Updated.UTC().Unix(),
		Deleted:        deleted,
		Entry:          v1.DealEntry(v1.DealEntry_value[a.Entry]),
	}
}
