package tradeRulessvc

import (
	"context"
	"runtime"
	"time"

	model "github.com/Ruscigno/ticker-signals/internal/transaction/traderules"
	repo "github.com/Ruscigno/ticker-signals/internal/transaction/traderules/repo"
	rt "github.com/Ruscigno/ticker-signals/internal/transaction/traderules/ruletypes"
	tt "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction"
	"github.com/blendle/zapdriver"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
)

// NewTradeRulesService creates a service to interact with PostgreSQL
func NewTradeRulesService(ctx context.Context, repo repo.TradeRulesRepository) TradeRulesService {
	return &tradeRulesService{
		ctx:         ctx,
		repo:        repo,
		records:     []*model.TradeRules{},
		expireCache: time.Now().UTC(),
	}
}

// TradeRulesService ...
type tradeRulesService struct {
	ctx         context.Context
	repo        repo.TradeRulesRepository
	records     []*model.TradeRules
	expireCache time.Time
}

// GetByAccount return all rules from an account
func (r *tradeRulesService) GetByAccount(accountID int64, ruleType int) ([]*model.TradeRules, error) {
	result := []*model.TradeRules{}
	for _, tr := range r.records {
		if !tr.Symbol.Valid {
			continue
		}
		if tr.AccountID == accountID && tr.RuleType == ruleType {
			result = append(result, tr)
		}
	}
	if len(result) > 0 {
		return result, nil
	}
	ac, err := r.repo.GetByAccount(accountID, ruleType)
	if err != nil {
		zap.L().Error("TradeRules: GetByAccount Error", zap.Error(err), zap.Int64("accountID", accountID), zapdriver.SourceLocation(runtime.Caller(0)))
		return nil, err
	}
	//TODO: should use a mutex here...
	r.records = append(r.records, ac...)
	return ac, nil
}

// GetBySymbol return all symbol's rules from an account
func (r *tradeRulesService) GetBySymbol(accountID int64, symbol string, ruleType int) ([]*model.TradeRules, error) {
	result := []*model.TradeRules{}
	for _, tr := range r.records {
		if !tr.Symbol.Valid {
			continue
		}
		if tr.AccountID == accountID && tr.Symbol.String == symbol && tr.RuleType == ruleType {
			result = append(result, tr)
		}
	}
	if len(result) > 0 {
		return result, nil
	}
	ac, err := r.repo.GetBySymbol(accountID, symbol, ruleType)
	if err != nil {
		zap.L().Error("TradeRules: GetBySymbol Error", zap.Error(err), zap.Int64("accountID", accountID), zap.String("symbol", symbol), zapdriver.SourceLocation(runtime.Caller(0)))
		return nil, err
	}
	//TODO: should use a mutex here...
	r.records = append(r.records, ac...)
	return ac, nil
}

func (r *tradeRulesService) GetAll() ([]*model.TradeRules, error) {
	var err error
	r.records, err = r.repo.GetAll()
	if err != nil {
		zap.L().Error("TradeRules: GetAll Error", zap.Error(err), zapdriver.SourceLocation(runtime.Caller(0)))
		return nil, err
	}
	r.expireCache = time.Now().Add(time.Hour * 1).UTC()
	return r.records, nil
}

func (r *tradeRulesService) ApplyRules(accountID int64, req []*tt.TradeRequest) error {
	var err error
	zap.L().Info("ApplyRules method", zap.Int64("account", accountID), zapdriver.SourceLocation(runtime.Caller(0)))
	if r.expireCache.After(time.Now().UTC()) {
		_, err = r.GetAll()
	}
	if err != nil {
		return err
	}
	for _, tr := range r.records {
		rule, ok := rt.RuleTypeList[tr.RuleType]
		if !ok {
			zap.L().Error("ApplyRules: unknown rule interface", zap.String("rule", spew.Sdump(rule)), zapdriver.SourceLocation(runtime.Caller(0)))
			continue
		}
		if accountID != tr.AccountID {
			continue
		}
		for _, trade := range req {
			if trade.Status == tt.Invalid {
				continue
			}
			err = rule.ApplyRules(tr.Rule, trade)
			zap.L().Info("Applying Rules", zap.Int64("account", tr.AccountID), zap.String("rule", tr.Rule), zapdriver.SourceLocation(runtime.Caller(0)))
			if err != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}
	return err
}
