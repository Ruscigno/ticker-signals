package adapters

import (
	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"
	tb "github.com/Ruscigno/ticker-signals/internal/tickerbeats/beats"
	sig "github.com/Ruscigno/ticker-signals/internal/tickerbeats/signal"
)

// BeatsToProto transform tb.TickerBeats int v1.TickerBeats
// TODO: try to use proto.Marshal and Unmarshal
func BeatsToProto(beats *tb.TickerBeats) *v1.TradeBeats {
	signal := SignalToProto(beats.Signal)
	var tr []*v1.TradeRequest
	ctr := len(beats.TradeRequests)
	if ctr > 0 {
		tr = make([]*v1.TradeRequest, ctr)
		for i, beat := range beats.TradeRequests {
			tr[i] = TradeRequestToProto(beat)
		}
	}
	return &v1.TradeBeats{
		Signal:         signal,
		ExpirationTime: beats.ExpirationTime,
		TradeRequest:   tr,
	}
}

// SignalToProto transform sig.Signal into v1.Signal
// TODO: try to use proto.Marshal and Unmarshal
func SignalToProto(s *sig.Signal) *v1.Signal {
	return &v1.Signal{
		SignalId:             s.SignalID,
		SourceAccountId:      s.SourceAccountID,
		DestinationAccountId: s.DestinationAccountID,
		Active:               s.Active,
		MaxDeposit:           s.MaxDepositPercent,
		StopIfLessThan:       s.StopIfLessThan,
		MaxSpread:            s.MaxSpread,
		MinutesToExpire:      s.MinutesToExpire,
	}
}
