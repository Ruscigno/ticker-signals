package api

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"
	"github.com/Ruscigno/ticker-signals/internal/api/adapters"
	svc "github.com/Ruscigno/ticker-signals/internal/tickerbeats/beats"
	trService "github.com/Ruscigno/ticker-signals/internal/transaction/traderules/service"
	"go.uber.org/zap"
)

// TransactionsServiceServer is used to implement TransactionsServiceServer
type TickerBeatsServiceServer struct {
	v1.UnimplementedTickerBeatsServiceServer
	ticker svc.TickerBeatsService
	trSvc  trService.TradeRulesService
}

// NewTickerBeatsServiceServer creates a new API server handler
func NewTickerBeatsServiceServer(ticker svc.TickerBeatsService, trSvc trService.TradeRulesService) *TickerBeatsServiceServer {
	return &TickerBeatsServiceServer{
		ticker: ticker,
		trSvc:  trSvc,
	}
}

func (t *TickerBeatsServiceServer) GetTickerBeats(ctx context.Context, req *v1.TickerBeatsRequest) (*v1.TickerBeatsResponse, error) {
	start := time.Now().UTC()
	from := time.Unix(req.GetServerTime(), 0)
	accID := req.GetAccountId()
	tickersBeats, err := t.ticker.GetTickerBeats(accID, from, req.GetSignalType(), req.GetLastOrderId())
	if err != nil {
		return nil, err
	}
	var beats []*v1.TradeBeats
	for _, tb := range tickersBeats {
		if !tb.Valid {
			continue
		}
		beat := adapters.BeatsToProto(tb)
		beats = append(beats, beat)
	}
	if len(beats) <= 0 || len(beats[0].TradeRequest) <= 0 {
		return &v1.TickerBeatsResponse{}, nil
	}
	err = t.ticker.BeatsSent(tickersBeats)
	if err != nil {
		return nil, err
	}
	zap.L().Info("GetTickerBeats processing finished",
		zap.Int64("accountID", accID),
		zap.String("latency", fmt.Sprintf("%d", time.Since(start).Milliseconds())),
	)
	return &v1.TickerBeatsResponse{Beats: beats}, nil
}
