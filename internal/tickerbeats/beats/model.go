package tickerbeats

import (
	"github.com/Ruscigno/ticker-signals/internal/tickerbeats/signal"
	oo "github.com/Ruscigno/ticker-signals/internal/transaction/orders"
	pp "github.com/Ruscigno/ticker-signals/internal/transaction/positions"
	tmodel "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction"
)

type TickerBeats struct {
	Signal         *signal.Signal         `json:"signal"`
	ExpirationTime int32                  `json:"expiration_time"`
	GroupID        string                 `json:"group_id"`
	Positions      []*pp.Position         `json:"positions"`
	Orders         []*oo.Order            `json:"orders"`
	TradeRequests  []*tmodel.TradeRequest `json:"trade_requests"`
	Valid          bool                   `json:"valid"`
}
