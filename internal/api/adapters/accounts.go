package adapters

import (
	"time"

	acc "github.com/Ruscigno/ticker-signals/internal/transaction/accounts"

	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"
)

// ProtoToAccount transform an proto to an account model
// TODO: try to use proto.Marshal and Unmarshal
func ProtoToAccount(a *v1.Account) *acc.Account {
	return &acc.Account{
		AccountID:    a.AccountId,
		TradeMode:    v1.AccountTradeMode_name[int32(a.TradeMode)],
		Leverage:     a.Leverage,
		MarginMode:   v1.AccountMarginMode_name[int32(a.MarginMode)],
		StopoutMode:  v1.AccountStopoutMode_name[int32(a.StopoutMode)],
		TradeAllowed: a.TradeAllowed,
		TradeExpert:  a.TradeExpert,
		LimitOrders:  a.LimitOrders,
		Name:         a.Name,
		Server:       a.Server,
		Currency:     a.Currency,
		Company:      a.Company,
		Created:      time.Unix(a.TimeGMT, 0),
	}
}
