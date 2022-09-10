package tickerbeats

import (
	"github.com/Ruscigno/ticker-signals/internal/tickerbeats/signal"
	tmodel "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction"
)

type TickerBeats struct {
	Signal         *signal.Signal         `json:"signal"`
	ExpirationTime int32                  `json:"expiration_time"`
	GroupID        string                 `json:"group_id"`
	TradeRequests  []*tmodel.TradeRequest `json:"trade_requests"`
	Valid          bool                   `json:"valid"`
}
