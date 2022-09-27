package tradeRules

import (
	"database/sql"
)

type TradeRules struct {
	AccountID   int64          `json:"account_id"`
	RuleID      int            `json:"rule_id"`
	Active      bool           `json:"active"`
	Description sql.NullString `json:"description"`
	Symbol      sql.NullString `json:"symbol"`
	RuleType    int            `json:"rule_type"`
	RuleVersion int            `json:"rule_version"`
	Rule        string         `json:"rule"`
}
