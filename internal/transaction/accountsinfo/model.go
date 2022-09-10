package accountsInfo

import (
	"database/sql"
)

//AccountInfo Properties
type AccountInfo struct {
	AccountID           int64         `json:"account_id"`
	InfoID              int64         `json:"account_info_id"`
	Balance             float64       `json:"balance"`
	Credit              float64       `json:"credit"`
	Profit              float64       `json:"profit"`
	Equity              float64       `json:"equity"`
	Margin              float64       `json:"margin"`
	FreeMargin          float64       `json:"free_margin"`
	MarginLevel         float64       `json:"margin_level"`
	MarginCall          float64       `json:"margin_call"`
	MarginStopout       float64       `json:"margin_stopout"`
	TimeTradeServer     sql.NullTime  `json:"TimeTradeServer"`
	TimeCurrent         sql.NullTime  `json:"TimeCurrent"`
	TimeLocal           sql.NullTime  `json:"TimeLocal"`
	TimeGMT             sql.NullTime  `json:"TimeGMT"`
	LocalTimeGMTOffset  sql.NullInt64 `json:"LocalTimeGMTOffset"`
	ServerTimeGMTOffset sql.NullInt64 `json:"ServerTimeGMTOffset"`
}
