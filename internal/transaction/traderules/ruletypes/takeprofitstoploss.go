package tradeRules

import (
	"encoding/json"
	"fmt"
	"math"

	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"
	tt "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction"
)

type ExitRules struct {
	TakeProfitInPercentage float64 `json:"TakeProfitInPercentage"`
	TakeProfitInPoints     float64 `json:"TakeProfitInPoints"`
	StopLossInPercentage   float64 `json:"StopLossInPercentage"`
	StopLossInPoints       float64 `json:"StopLossInPoints"`
}

type accountExitRules struct {
	Rules *ExitRules
}

func (aa *accountExitRules) ApplyRules(rule string, tr *tt.TradeRequest) error {
	return nil
}

type SignalExit struct {
	SourceAccountID int64      `json:"SourceAccountID"`
	Rules           *ExitRules `json:"Rules"`
}

type signalExitRules struct {
	Rules *SignalExit
}

func (aa *signalExitRules) ApplyRules(rule string, trade *tt.TradeRequest) error {
	var r SignalExit
	if err := json.Unmarshal([]byte(rule), &r); err != nil {
		return err
	}
	err := calculateExits(rule, r.Rules, trade)
	if err != nil {
		return err
	}
	return nil
}

func calculateExits(rawRule string, r *ExitRules, t *tt.TradeRequest) error {
	sellType := v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_SELL)]
	buyType := v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_BUY)]
	entryIN := v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)]
	t.Status = tt.Invalid
	if t.Entry != entryIN {
		return nil
	}
	price := t.Price.Float64
	var takeProfit, stopLoss float64
	if t.OrderType == buyType {
		takeProfit = math.Min(price*(1+r.TakeProfitInPercentage), price+r.TakeProfitInPoints)
		stopLoss = math.Max(price*(1-r.StopLossInPercentage), price-r.StopLossInPoints)
		if takeProfit < price {
			t.TakeProfit.Valid = false
			t.StopLoss.Valid = false
			return fmt.Errorf("buy: invalid take profit, price: %.5f, take profit: %.5f", price, takeProfit)
		}
		if stopLoss > price {
			t.TakeProfit.Valid = false
			t.StopLoss.Valid = false
			return fmt.Errorf("buy: invalid stop loss, price: %.5f, stop loss: %.5f", price, stopLoss)
		}
		if t.TakeProfit.Valid && takeProfit > t.TakeProfit.Float64 {
			takeProfit = t.TakeProfit.Float64
		}
		if t.StopLoss.Valid && stopLoss < t.StopLoss.Float64 {
			stopLoss = t.StopLoss.Float64
		}

	} else if t.OrderType == sellType {
		takeProfit = math.Max(price*(1-r.TakeProfitInPercentage), price-r.TakeProfitInPoints)
		stopLoss = math.Min(price*(1+r.StopLossInPercentage), price+r.StopLossInPoints)
		if takeProfit > price {
			t.TakeProfit.Valid = false
			t.StopLoss.Valid = false
			return fmt.Errorf("sell: invalid take profit, price: %.5f, take profit: %.5f", price, takeProfit)
		}
		if stopLoss < price {
			t.TakeProfit.Valid = false
			t.StopLoss.Valid = false
			return fmt.Errorf("sell: invalid stop loss, price: %.5f, stop loss: %.5f", price, stopLoss)
		}
		if t.TakeProfit.Valid && takeProfit < t.TakeProfit.Float64 {
			takeProfit = t.TakeProfit.Float64
		}
		if t.StopLoss.Valid && stopLoss > t.StopLoss.Float64 {
			stopLoss = t.StopLoss.Float64
		}
	} else {
		return nil
	}
	takeProfit = math.Round(takeProfit*100000) / 100000
	stopLoss = math.Round(stopLoss*100000) / 100000
	t.StopLoss.Float64 = stopLoss
	t.TakeProfit.Float64 = takeProfit
	t.TakeProfit.Valid = true
	t.StopLoss.Valid = true
	t.Status = tt.Valid
	return nil
}
