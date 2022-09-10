package api

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"
	adapter "github.com/Ruscigno/ticker-signals/internal/api/adapters"
	bb "github.com/Ruscigno/ticker-signals/internal/tickerbeats/beats"
	accService "github.com/Ruscigno/ticker-signals/internal/transaction/accounts/service"
	infoS "github.com/Ruscigno/ticker-signals/internal/transaction/accountsinfo/service"
	deaService "github.com/Ruscigno/ticker-signals/internal/transaction/deals/service"
	ordService "github.com/Ruscigno/ticker-signals/internal/transaction/orders/service"
	posModel "github.com/Ruscigno/ticker-signals/internal/transaction/positions"
	posService "github.com/Ruscigno/ticker-signals/internal/transaction/positions/service"
	trService "github.com/Ruscigno/ticker-signals/internal/transaction/traderules/service"
	ttService "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction/service"

	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TransactionsServiceServer is used to implement TransactionsServiceServer
type TransactionsServiceServer struct {
	v1.UnimplementedTransactionsServiceServer
	AccSvc  accService.AccountsService
	InfoSvc infoS.AccountsInfoService
	DeaSvc  deaService.DealsService
	OrdSvc  ordService.OrdersService
	PosSvc  posService.PositionsService
	TtSvc   ttService.TradeTransactionService
	Beats   bb.TickerBeatsService
	TrSvc   trService.TradeRulesService
}

// NewTransactionsServiceServer creates a new API server handler
func NewTransactionsServiceServer(
	accSvc accService.AccountsService,
	infoSvc infoS.AccountsInfoService,
	deaSvc deaService.DealsService,
	ordSvc ordService.OrdersService,
	posSvc posService.PositionsService,
	ttSvc ttService.TradeTransactionService,
	beats bb.TickerBeatsService,
	trSvc trService.TradeRulesService,
) *TransactionsServiceServer {
	return &TransactionsServiceServer{
		AccSvc:  accSvc,
		InfoSvc: infoSvc,
		DeaSvc:  deaSvc,
		OrdSvc:  ordSvc,
		PosSvc:  posSvc,
		TtSvc:   ttSvc,
		Beats:   beats,
		TrSvc:   trSvc,
	}
}

// CreateAccount Create a new account
func (s *TransactionsServiceServer) CreateAccount(ctx context.Context, req *v1.CreateAccountRequest) (*empty.Empty, error) {
	// start := time.Now().UTC()
	// accountId := req.Account.GetAccountId()
	err := s.AccSvc.Insert(adapter.ProtoToAccount(req.GetAccount()), req.Account.GetTimeGMT())
	if err != nil {
		return &empty.Empty{}, err
	}
	err = s.InfoSvc.Insert(adapter.ProtoToAccountInfo(req.GetAccount()))
	if err != nil {
		return &empty.Empty{}, err
	}
	// zap.L().Info("CreateAccount processing finished", zap.Int64("accountID", accountId), zap.String("latency", fmt.Sprintf("%d", time.Since(start).Milliseconds())))
	return &empty.Empty{}, nil
}

// CreateDeals Create a new beats
func (s *TransactionsServiceServer) CreateDeals(ctx context.Context, req *v1.CreateDealsRequest) (*emptypb.Empty, error) {
	start := time.Now().UTC()
	var accountId int64 = 0
	for _, item := range req.GetDeals() {
		deal := adapter.ProtoToDeal(item)
		//First deal, first balance
		if item.DealType == v1.DealType_DEAL_TYPE_BALANCE {
			deal.Symbol.String = "BALANCE"
			deal.Symbol.Valid = true
		}
		if accountId != 0 && accountId != deal.AccountID {
			return &empty.Empty{}, fmt.Errorf("all records should be from the same accountID: %d", accountId)
		}
		accountId = deal.AccountID
		err := s.DeaSvc.Insert(deal)
		if err != nil {
			return &empty.Empty{}, err
		}
	}
	zap.L().Info("CreateDeals processing finished", zap.Int64("accountID", accountId), zap.String("latency", fmt.Sprintf("%d", time.Since(start).Milliseconds())))
	return &empty.Empty{}, nil
}

// CreateOrders a new order
func (s *TransactionsServiceServer) CreateOrders(ctx context.Context, req *v1.CreateOrdersRequest) (*emptypb.Empty, error) {
	start := time.Now().UTC()
	var accountId int64 = 0
	for _, item := range req.GetOrders() {
		ord := adapter.ProtoToOrder(item)
		if accountId != 0 && accountId != ord.AccountID {
			return &empty.Empty{}, fmt.Errorf("all records should be from the same accountID: %d", accountId)
		}
		accountId = ord.AccountID
		err := s.OrdSvc.Insert(ord)
		if err != nil {
			return &empty.Empty{}, err
		}
	}
	zap.L().Info("CreateOrders processing finished", zap.Int64("accountID", accountId), zap.String("latency", fmt.Sprintf("%d", time.Since(start).Milliseconds())))
	return &empty.Empty{}, nil
}

func (s *TransactionsServiceServer) CreatePositions(ctx context.Context, req *v1.CreatePositionsRequest) (*emptypb.Empty, error) {
	start := time.Now().UTC()
	var posList []*posModel.Position
	var accountId int64 = req.GetAccountId()
	// var TimeGMT time.Time = time.Unix(req.GetTimeGMT(), 0)
	for _, item := range req.GetPositions() {
		pos := adapter.ProtoToPosition(item)
		if accountId != 0 && accountId != pos.AccountID {
			return &empty.Empty{}, fmt.Errorf("all records should be from the same accountID: %d", accountId)
		}
		posList = append(posList, pos)
	}
	// var err error
	// if len(posList) == 0 {
	// 	err = s.posSvc.CloseAll(accountId, TimeGMT)
	// } else {
	inserted, err := s.PosSvc.InsertMulti(accountId, posList)
	// }
	if err != nil {
		return &empty.Empty{}, err
	}
	if inserted {
		zap.L().Info("CreatePositions processing finished", zap.Int64("accountID", accountId), zap.String("latency", fmt.Sprintf("%d", time.Since(start).Milliseconds())))
	}
	return &empty.Empty{}, nil
}

func (s *TransactionsServiceServer) CreateTradeTransaction(ctx context.Context, req *v1.CreateTradeTransactionRequest) (*emptypb.Empty, error) {
	start := time.Now().UTC()
	accountId := req.GetAccountId()
	timeGMT := req.GetTimeGMT()
	cOrder := req.GetCreationOrder()
	transaction := adapter.ProtoToTradeTransaction(accountId, cOrder, req.GetTradeTransaction())
	request := adapter.ProtoToTradeRequest(accountId, cOrder, req.GetTradeRequest())
	result := adapter.ProtoToTradeResult(accountId, cOrder, req.GetTradeResult())
	ra, err := s.TtSvc.Insert(transaction, request, result, timeGMT)
	if err != nil {
		return &empty.Empty{}, err
	}
	zap.L().Info("CreateTradeTransaction processing finished",
		zap.Int64("accountID", accountId),
		zap.Int64("creationOrder", cOrder),
		zap.String("latency", fmt.Sprintf("%d", time.Since(start).Milliseconds())),
		zap.Int64("RowsAffected", ra),
	)
	return &empty.Empty{}, nil
}
