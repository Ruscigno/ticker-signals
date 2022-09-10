package deals

import (
	"database/sql"
	"time"
)

//Deal Properties
type Deal struct {
	DealID     int64           `json:"deal_id"`
	AccountID  int64           `json:"account_id"`
	Ticket     int64           `json:"ticket"`
	Magic      sql.NullInt64   `json:"magic"`
	OrderID    sql.NullInt64   `json:"order_id"`
	Symbol     sql.NullString  `json:"symbol"`
	DealTime   time.Time       `json:"deal_time"`
	DealType   string          `json:"deal_type"`
	Entry      string          `json:"entry"`
	PositionId int64           `json:"position_id"`
	Volume     float64         `json:"volume"`
	Price      float64         `json:"price"`
	Commission sql.NullFloat64 `json:"commission"`
	Swap       sql.NullFloat64 `json:"swap"`
	Profit     sql.NullFloat64 `json:"profit"`
	Comment    string          `json:"comment"`
	ExternalId string          `json:"external_id"`
	Created    time.Time       `json:"created"`
	Updated    time.Time       `json:"updated"`
	Deleted    sql.NullTime    `json:"deleted"`
	Reason     int64           `json:"reason"`
	DealFee    sql.NullFloat64 `json:"deal_fee"`
}
