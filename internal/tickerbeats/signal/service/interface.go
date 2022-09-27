package signalsvc

import (
	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"
	"github.com/Ruscigno/ticker-signals/internal/tickerbeats/signal"
	tm "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction"
)

// SignalService is a CRUD to the database
type SignalService interface {
	CloseDeadPositions(sourceAccountID, destinationAccountID int64, groupID string, minToExpire int32) ([]*tm.TradeRequest, error)
	ConfirmByExternalID(destinationAccountID, externalID, tickerBeatsID, positionID int64, status signal.SignalStatusEnum) error
	CreateTickerBeats(sourceAccountID, destinationAccountID int64, groupID string, minToExpire int32) (int64, error)
	GetTradeRequesByGroupID(destinationAccountID int64, groupID string, entry v1.DealEntry) ([]*tm.TradeRequest, error)
	GetSignalByDestination(accountID int64) ([]*signal.Signal, error)
	NeedToCloseAllPositions(accountID int64, groupID string, stop int32) ([]*tm.TradeRequest, error)
	RemoveDuplicatedSignals(destinationAccountID int64, groupID string) error
	UpdatePositionIdBeforeClose(sourceAccountID, destinationAccountID int64, groupID string) error
	UpdateStatus(sourceAccountID, destinationAccountID int64, status signal.SignalStatusEnum, groupID string) error
}
