package tradeRules

import (
	tt "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction"
)

// *** RULE TYPES ****************************************
type AccountMaxOpenPositions struct {
	Max int `json:"Max"`
}

type SymbolMaxOpenPositions struct {
	Symbol string `json:"Symbol"`
	Max    int    `json:"Max"`
}

type SignalMaxOpenPositions struct {
	SourceAccountID int64 `json:"SourceAccountID"`
	Max             int   `json:"Max"`
}

//****************************************************

var RuleTypeList = map[int]Rules{
	1: &accountExitRules{Rules: &ExitRules{}},
	2: &signalExitRules{Rules: &SignalExit{Rules: &ExitRules{}}},
}

type Rules interface {
	ApplyRules(rule string, tr *tt.TradeRequest) error
}
