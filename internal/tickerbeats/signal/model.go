package signal

import (
	"database/sql"
	"math"
	"time"
)

type SignalStatusEnum int
type OrderBoostTypeEnum int32
type SignalStatusEnumUsed map[SignalStatusEnum]bool

const (
	Canceled             SignalStatusEnum = 1
	Created              SignalStatusEnum = 2
	Sent                 SignalStatusEnum = 4
	Expired              SignalStatusEnum = 8
	Failed               SignalStatusEnum = 16
	TransactionConfirmed SignalStatusEnum = 32
	OrderConfirmed       SignalStatusEnum = 64
	DealConfirmed        SignalStatusEnum = 128
	PositionConfirmed    SignalStatusEnum = 256
	PositionClose        SignalStatusEnum = 512
)

const (
	Undefined  OrderBoostTypeEnum = 1
	Multiplier OrderBoostTypeEnum = 2
	OrderLimit OrderBoostTypeEnum = 4
)

type Signal struct {
	SignalID             int64           `json:"SignalID"`
	SourceAccountID      int64           `json:"SourceAccountID"`
	DestinationAccountID int64           `json:"DestinationAccountID"`
	Active               bool            `json:"Active"`
	MaxDepositPercent    int32           `json:"MaxDepositPercent"`
	StopIfLessThan       int32           `json:"StopIfLessThan"`
	MaxSpread            int64           `json:"MaxSpread"`
	MinutesToExpire      int32           `json:"MinutesToExpire"`
	OrderBoost           sql.NullFloat64 `json:"OrderBoost"`
	OrderBoostType       sql.NullInt32   `json:"OrderBoostType"`
}

type SignalResult struct {
	SourceAccountID       int64         `json:"SourceAccountID"`
	DestinationAccountID  int64         `json:"DestinationAccountID"`
	SourceBeatsID         sql.NullInt64 `json:"SourceBeatsID"`
	DestinationBeatsID    sql.NullInt64 `json:"DestinationBeatsID"`
	SourcePositionID      int64         `json:"SourcePositionID"`
	DestinationPositionID sql.NullInt64 `json:"DestinationPositionID"`
	SignalType            string        `json:"SignalType"`
	SignalStatus          int16         `json:"SignalStatus"`
	ExternalID            int64         `json:"ExternalID"`
	GroupId               string        `json:"GroupId"`
	Entry                 string        `json:"Entry"`
	SentTime              sql.NullTime  `json:"SentTime"`
	ConfirmationTime      sql.NullTime  `json:"ConfirmationTime"`
	ExpireAt              sql.NullTime  `json:"ExpireAt"`
	Created               time.Time     `json:"Created"`
	Updated               time.Time     `json:"Updated"`
}

func CalcSignalStatusEnumEnabled(status SignalStatusEnum) SignalStatusEnumUsed {
	// Enum value maps for DealEntry.
	var result SignalStatusEnumUsed = SignalStatusEnumUsed{
		Canceled:             false,
		Created:              false,
		Sent:                 false,
		Expired:              false,
		Failed:               false,
		TransactionConfirmed: false,
		OrderConfirmed:       false,
		DealConfirmed:        false,
		PositionConfirmed:    false,
		PositionClose:        false,
	}
	values := make(map[int]SignalStatusEnum)
	for i := 0; i < len(result); i++ {
		if i == 0 {
			values[i] = SignalStatusEnum(1)
		} else {
			values[i] = SignalStatusEnum(math.Exp2(float64(i)))
		}
	}
	vStatus := int32(status)
	for i := len(result) - 1; i > 0; i-- {
		v := values[i]
		if vStatus >= int32(v) {
			result[v] = true
			vStatus -= int32(v)
		}
	}
	return result
}
