package accounts

import (
	"database/sql"
	"time"
)

//Account Properties
type Account struct {
	AccountID    int64        `json:"account_id"`
	TradeMode    string       `json:"trade_mode"`
	Leverage     int64        `json:"leverage"`
	MarginMode   string       `json:"margin_mode"`
	StopoutMode  string       `json:"stopout_mode"`
	TradeAllowed bool         `json:"trade_allowed"`
	TradeExpert  bool         `json:"trade_expert"`
	LimitOrders  int64        `json:"limit_orders"`
	Name         string       `json:"name"`
	Server       string       `json:"server"`
	Currency     string       `json:"currency"`
	Company      string       `json:"company"`
	Created      time.Time    `json:"created"`
	Updated      time.Time    `json:"updated"`
	Deleted      sql.NullTime `json:"deleted"`
}
