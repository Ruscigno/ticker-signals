package orders

import (
	"database/sql"
	"time"
)

const TableName = "tickerbeats.orders"

//Order Properties
type Order struct {
	AccountID      int64           `json:"account_id"`
	OrderID        int64           `json:"order_id"`
	Ticket         int64           `json:"ticket"`
	Symbol         sql.NullString  `json:"symbol"`
	TimeSetup      sql.NullTime    `json:"time_setup"`
	OrderType      string          `json:"order_type"`
	State          string          `json:"state"`
	TimeExpiration sql.NullTime    `json:"time_expiration"`
	TimeDone       sql.NullTime    `json:"time_done"`
	TypeFilling    string          `json:"type_filling"`
	TypeTime       string          `json:"type_time"`
	Magic          sql.NullInt64   `json:"magic"`
	PositionId     int64           `json:"position_id"`
	VolumeInitial  float64         `json:"volume_initial"`
	VolumeCurrent  sql.NullFloat64 `json:"volume_current"`
	PriceOpen      float64         `json:"price_open"`
	StopLoss       sql.NullFloat64 `json:"stop_loss"`
	TakeProfit     sql.NullFloat64 `json:"take_profit"`
	PriceCurrent   float64         `json:"price_current"`
	PriceStopLimit sql.NullFloat64 `json:"price_stop_limit"`
	Comment        string          `json:"comment"`
	ExternalID     string          `json:"external_id"`
	Reason         string          `json:"reason"`
	Created        time.Time       `json:"created"`
	Updated        time.Time       `json:"updated"`
	Deleted        sql.NullTime    `json:"deleted"`
	PositionByID   sql.NullInt64   `json:"position_by_id"`
	// Direction      bool          `json:"direction"`
}
