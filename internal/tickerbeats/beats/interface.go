package tickerbeats

import (
	"time"

	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"
	"github.com/Ruscigno/ticker-signals/internal/tickerbeats/signal"
)

// TickerBeatsService is a CRUD to the database
type TickerBeatsService interface {
	BeatsSent(beats []*TickerBeats) error
	ConfirmByExternalID(destinationAccountID, externalid, tickerBeatsID, positionID int64, status signal.SignalStatusEnum) error
	GetTickerBeats(accountID int64, from time.Time, sType v1.SignalType, lastOrderID int64) ([]*TickerBeats, error)
}
