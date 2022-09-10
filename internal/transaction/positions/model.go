package position

import (
	"database/sql"
	"time"
)

const TableName = "tickerbeats.positions"

//Position Properties
type Position struct {
	PositionID     int64          `json:"position_id"`
	AccountID      int64          `json:"account_id"`
	Ticket         int64          `json:"ticket"`
	Symbol         sql.NullString `json:"symbol"`
	PositionTime   sql.NullTime   `json:"position_time"`
	PositionType   string         `json:"position_type"`
	Volume         float64        `json:"volume"`
	PriceOpen      float64        `json:"price_open"`
	StopLoss       float64        `json:"stop_loss"`
	TakeProfit     float64        `json:"take_profit"`
	PriceCurrent   float64        `json:"price_current"`
	Commission     float64        `json:"commission"`
	Swap           float64        `json:"swap"`
	Profit         float64        `json:"profit"`
	Comment        string         `json:"comment"`
	Created        time.Time      `json:"created"`
	Updated        time.Time      `json:"updated"`
	Deleted        sql.NullTime   `json:"deleted"`
	PositionUpdate sql.NullTime   `json:"position_update"`
	Reason         string         `json:"reason"`
	ExternalID     string         `json:"external_id"`
	Magic          sql.NullInt64  `json:"magic"`
}
