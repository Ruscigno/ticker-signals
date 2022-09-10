package tradeRules

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"
	tmodel "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction"
	"github.com/stretchr/testify/assert"
)

func TestAccountExitRules(t *testing.T) {
	aer := &accountExitRules{
		Rules: &ExitRules{
			TakeProfitInPercentage: 0.10,
			TakeProfitInPoints:     100,
			StopLossInPercentage:   0.20,
			StopLossInPoints:       200,
		},
	}
	assert.Nil(t, aer.ApplyRules("", nil))
}

func TestNoPrevious_SL_TP(t *testing.T) {
	aer := &signalExitRules{
		Rules: &SignalExit{
			SourceAccountID: 10000,
			Rules: &ExitRules{
				TakeProfitInPercentage: 0.10,
				TakeProfitInPoints:     100,
				StopLossInPercentage:   0.20,
				StopLossInPoints:       200,
			},
		},
	}
	trades := []*tmodel.TradeRequest{
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_BUY)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
		},
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_SELL)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
		},
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_BUY)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_OUT)],
		},
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_SELL)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_OUT)],
		},
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_BUY_LIMIT)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
		},
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_SELL_LIMIT)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
		},
	}
	tradesResult := []*tmodel.TradeRequest{
		{
			Price:      sql.NullFloat64{Valid: true, Float64: 100},
			StopLoss:   sql.NullFloat64{Valid: true, Float64: 80},
			TakeProfit: sql.NullFloat64{Valid: true, Float64: 110},
			Status:     tmodel.Valid,
		},
		{
			Price:      sql.NullFloat64{Valid: true, Float64: 100},
			StopLoss:   sql.NullFloat64{Valid: true, Float64: 120},
			TakeProfit: sql.NullFloat64{Valid: true, Float64: 90},
			Status:     tmodel.Valid,
		},
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
	}
	runTest(t, aer, trades, tradesResult, nil)
}

func TestPrevious_SL_TP(t *testing.T) {
	aer := &signalExitRules{
		Rules: &SignalExit{
			SourceAccountID: 10000,
			Rules: &ExitRules{
				TakeProfitInPercentage: 0.10,
				TakeProfitInPoints:     100,
				StopLossInPercentage:   0.20,
				StopLossInPoints:       200,
			},
		},
	}
	trades := []*tmodel.TradeRequest{
		{
			AccountID:  10000,
			Price:      sql.NullFloat64{Valid: true, Float64: 100},
			OrderType:  v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_BUY)],
			Entry:      v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
			StopLoss:   sql.NullFloat64{Valid: true, Float64: 99},
			TakeProfit: sql.NullFloat64{Valid: true, Float64: 101},
		},
		{
			AccountID:  10000,
			Price:      sql.NullFloat64{Valid: true, Float64: 100},
			OrderType:  v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_SELL)],
			Entry:      v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
			StopLoss:   sql.NullFloat64{Valid: true, Float64: 102},
			TakeProfit: sql.NullFloat64{Valid: true, Float64: 98},
		},
		{
			AccountID:  10000,
			Price:      sql.NullFloat64{Valid: true, Float64: 100},
			OrderType:  v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_BUY)],
			Entry:      v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
			StopLoss:   sql.NullFloat64{Valid: true, Float64: 10},
			TakeProfit: sql.NullFloat64{Valid: true, Float64: 200},
		},
		{
			AccountID:  10000,
			Price:      sql.NullFloat64{Valid: true, Float64: 100},
			OrderType:  v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_SELL)],
			Entry:      v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
			StopLoss:   sql.NullFloat64{Valid: true, Float64: 200},
			TakeProfit: sql.NullFloat64{Valid: true, Float64: 10},
		},
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_BUY)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_OUT)],
		},
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_SELL)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_OUT)],
		},
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_BUY_LIMIT)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
		},
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_SELL_LIMIT)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
		},
	}
	tradesResult := []*tmodel.TradeRequest{
		{
			Price:      sql.NullFloat64{Valid: true, Float64: 100},
			StopLoss:   sql.NullFloat64{Valid: true, Float64: 99},
			TakeProfit: sql.NullFloat64{Valid: true, Float64: 101},
			Status:     tmodel.Valid,
		},
		{
			Price:      sql.NullFloat64{Valid: true, Float64: 100},
			StopLoss:   sql.NullFloat64{Valid: true, Float64: 102},
			TakeProfit: sql.NullFloat64{Valid: true, Float64: 98},
			Status:     tmodel.Valid,
		},
		{
			Price:      sql.NullFloat64{Valid: true, Float64: 100},
			StopLoss:   sql.NullFloat64{Valid: true, Float64: 80},
			TakeProfit: sql.NullFloat64{Valid: true, Float64: 110},
			Status:     tmodel.Valid,
		},
		{
			Price:      sql.NullFloat64{Valid: true, Float64: 100},
			StopLoss:   sql.NullFloat64{Valid: true, Float64: 120},
			TakeProfit: sql.NullFloat64{Valid: true, Float64: 90},
			Status:     tmodel.Valid,
		},
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
	}
	runTest(t, aer, trades, tradesResult, nil)
}

func TestNegative_TP(t *testing.T) {
	aer := &signalExitRules{
		Rules: &SignalExit{
			SourceAccountID: 10000,
			Rules: &ExitRules{
				TakeProfitInPercentage: -0.10,
				TakeProfitInPoints:     -100,
				StopLossInPercentage:   0.20,
				StopLossInPoints:       200,
			},
		},
	}
	trades := []*tmodel.TradeRequest{
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_BUY)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
		},
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_SELL)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
		},
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_BUY)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
		},
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_SELL)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
		},
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_BUY)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_OUT)],
		},
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_SELL)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_OUT)],
		},
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_BUY_LIMIT)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
		},
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_SELL_LIMIT)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
		},
	}
	tradesResult := []*tmodel.TradeRequest{
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
	}
	errors := []error{
		fmt.Errorf("buy: invalid take profit, price: %.5f, take profit: %.5f", 100.0, 0.0),
		fmt.Errorf("sell: invalid take profit, price: %.5f, take profit: %.5f", 100.0, 200.0),
		fmt.Errorf("buy: invalid take profit, price: %.5f, take profit: %.5f", 100.0, 0.0),
		fmt.Errorf("sell: invalid take profit, price: %.5f, take profit: %.5f", 100.0, 200.0),
		fmt.Errorf("buy: invalid take profit, price: %.5f, take profit: %.5f", 100.0, 0.0),
		fmt.Errorf("buy: invalid take profit, price: %.5f, take profit: %.5f", 100.0, 0.0),
		fmt.Errorf("buy: invalid take profit, price: %.5f, take profit: %.5f", 100.0, 0.0),
		fmt.Errorf("buy: invalid take profit, price: %.5f, take profit: %.5f", 100.0, 0.0),
	}
	runTest(t, aer, trades, tradesResult, errors)
}

func TestNegative_SL(t *testing.T) {
	aer := &signalExitRules{
		Rules: &SignalExit{
			SourceAccountID: 10000,
			Rules: &ExitRules{
				TakeProfitInPercentage: 0.10,
				TakeProfitInPoints:     1,
				StopLossInPercentage:   -0.20,
				StopLossInPoints:       -2,
			},
		},
	}
	trades := []*tmodel.TradeRequest{
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_BUY)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
		},
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_SELL)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
		},
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_BUY)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
		},
		{
			AccountID: 10000,
			Price:     sql.NullFloat64{Valid: true, Float64: 100},
			OrderType: v1.OrderType_name[int32(v1.OrderType_ORDER_TYPE_SELL)],
			Entry:     v1.DealEntry_name[int32(v1.DealEntry_DEAL_ENTRY_IN)],
		},
	}
	tradesResult := []*tmodel.TradeRequest{
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
		{
			Price:  sql.NullFloat64{Valid: true, Float64: 100},
			Status: tmodel.Invalid,
		},
	}
	errors := []error{
		fmt.Errorf("buy: invalid stop loss, price: %.5f, stop loss: %.5f", 100.0, 120.0),
		fmt.Errorf("sell: invalid stop loss, price: %.5f, stop loss: %.5f", 100.0, 80.0),
		fmt.Errorf("buy: invalid stop loss, price: %.5f, stop loss: %.5f", 100.0, 120.0),
		fmt.Errorf("sell: invalid stop loss, price: %.5f, stop loss: %.5f", 100.0, 80.0),
	}
	runTest(t, aer, trades, tradesResult, errors)
}

func runTest(t *testing.T, aer *signalExitRules, trades []*tmodel.TradeRequest, tradesResult []*tmodel.TradeRequest, testError []error) {
	for i, tr := range trades {
		rule, err := json.Marshal(aer.Rules)
		assert.NoError(t, err)
		err = aer.ApplyRules(string(rule), tr)
		if err != nil {
			assert.Equal(t, testError[i], err, "interaction: %d", i)
			assert.Equal(t, tradesResult[i].Status, tr.Status, "valid. interaction: %d", i)
			continue
		}
		assert.Equal(t, tradesResult[i].StopLoss.Float64, tr.StopLoss.Float64, fmt.Sprintf("StopLoss.Float64, interaction: %d", i))
		assert.Equal(t, tradesResult[i].StopLoss.Valid, tr.StopLoss.Valid, fmt.Sprintf("StopLoss.Valid, interaction: %d", i))
		assert.Equal(t, tradesResult[i].TakeProfit.Float64, tr.TakeProfit.Float64, fmt.Sprintf("TakeProfit.Float64, interaction: %d", i))
		assert.Equal(t, tradesResult[i].TakeProfit.Valid, tr.TakeProfit.Valid, fmt.Sprintf("TakeProfit.Valid, interaction: %d", i))
		assert.Equal(t, tradesResult[i].Status, tr.Status, "valid. interaction: %d", i)
	}
}
