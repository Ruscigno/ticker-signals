package tradetransaction

import (
	"database/sql"
	"time"
)

type ValidTradeRequest byte

const (
	Undefined ValidTradeRequest = iota
	Valid     ValidTradeRequest = iota
	Invalid   ValidTradeRequest = iota
)

// TradeTransaction: Structure of a Trade Transaction
type TradeTransaction struct {
	InternalID      int32           `json:"internal_id"`       // Internal ID
	AccountID       int64           `json:"account_id"`        // Account ID
	OrderID         int64           `json:"order_id"`          // Order ticket
	CreationOrder   int64           `json:"creation_order"`    // Creation order
	DealID          sql.NullInt64   `json:"deal_id"`           // Deal ticket
	Symbol          sql.NullString  `json:"symbol"`            // Trade symbol name
	TradeType       string          `json:"trade_type"`        // Trade transaction type
	OrderType       string          `json:"order_type"`        // Order type
	OrderState      string          `json:"order_state"`       // Order state
	DealType        string          `json:"deal_type"`         // Deal type
	TimeType        string          `json:"time_type"`         // Order type by action period
	TimeExpiration  sql.NullTime    `json:"time_expiration"`   // Order expiration time
	Price           sql.NullFloat64 `json:"price"`             // Price
	PriceTrigger    sql.NullFloat64 `json:"price_trigger"`     // Stop limit order activation price
	PriceStopLoss   sql.NullFloat64 `json:"price_stop_loss"`   // Stop Loss level
	PriceTakeProfit sql.NullFloat64 `json:"price_take_profit"` // Take Profit level
	Volume          sql.NullFloat64 `json:"volume"`            // Volume in lots
	PositionID      sql.NullInt64   `json:"position_id"`       // Position ticket
	PositionBy      sql.NullInt64   `json:"position_by"`       // Ticket of an opposite position
	Created         time.Time       `json:"created"`
	Updated         time.Time       `json:"updated"`
	Deleted         sql.NullTime    `json:"deleted"`
}

// TradeRequest:The Structure of a Trade Request Result
type TradeRequest struct {
	AccountID      int64             `json:"account_id"`      // Account ID
	OrderID        int64             `json:"order_id"`        // Order ticket
	CreationOrder  int64             `json:"creation_order"`  // Creation order
	Action         string            `json:"action"`          // Trade operation type
	Magic          sql.NullInt64     `json:"magic"`           // Expert Advisor ID (magic number)
	Symbol         sql.NullString    `json:"symbol"`          // Trade symbol
	Volume         float64           `json:"volume"`          // Requested volume for a deal in lots
	Price          sql.NullFloat64   `json:"price"`           // Price
	StopLimit      sql.NullFloat64   `json:"stop_limit"`      // Stop Limit level of the order
	StopLoss       sql.NullFloat64   `json:"stop_loss"`       // Stop Loss level of the order
	TakeProfit     sql.NullFloat64   `json:"take_profit"`     // Take Profit level of the order
	Deviation      sql.NullInt64     `json:"deviation"`       // Maximal possible deviation from the requested price
	OrderType      string            `json:"order_type"`      // Order type
	TypeFilling    string            `json:"type_filling"`    // Order execution type
	TypeTime       string            `json:"type_time"`       // Order expiration type
	TimeExpiration sql.NullTime      `json:"time_expiration"` // Order expiration time (for the orders of ORDER_TIME_SPECIFIED type)
	Comment        string            `json:"comment"`         // Order comment
	PositionID     sql.NullInt64     `json:"position_id"`     // TradeTransaction ticket
	PositionBy     sql.NullInt64     `json:"position_by"`     // The ticket of an opposite tradetransaction
	Created        time.Time         `json:"created"`         // created time
	Updated        time.Time         `json:"updated"`         // last updated time
	Deleted        sql.NullTime      `json:"deleted"`         // deleted time
	Entry          string            `json:"entry"`           // only for local work
	Status         ValidTradeRequest // internal use only
}

// TradeTransaction: Execution of trade operations results in the opening of a tradetransaction
type TradeResult struct {
	AccountID       int64           `json:"account_id"`       // Account ID
	OrderID         int64           `json:"order_id"`         // Order ticket, if it is placed
	CreationOrder   int64           `json:"creation_order"`   // Creation order
	RetCode         uint32          `json:"retcode"`          // Operation return code
	DealID          sql.NullInt64   `json:"deal_id"`          // Deal ticket, if it is performed
	Volume          sql.NullFloat64 `json:"volume"`           // Deal volume, confirmed by broker
	Price           sql.NullFloat64 `json:"price"`            // Deal price, confirmed by broker
	Bid             sql.NullFloat64 `json:"bid"`              // Current Bid price
	Ask             sql.NullFloat64 `json:"ask"`              // Current Ask price
	Comment         string          `json:"comment"`          // Broker comment to operation (by default it is filled by description of trade server return code)
	RequestID       uint32          `json:"request_id"`       // Request ID set by the terminal during the dispatch
	RetcodeExternal sql.NullInt64   `json:"retcode_external"` // Return code of an external trading system
	Created         time.Time       `json:"created"`
	Updated         time.Time       `json:"updated"`
	Deleted         sql.NullTime    `json:"deleted"`
}
