package app

import (
	"context"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/blendle/zapdriver"
	"github.com/jmoiron/sqlx"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	accRepository "github.com/Ruscigno/ticker-signals/internal/transaction/accounts/repo"
	accService "github.com/Ruscigno/ticker-signals/internal/transaction/accounts/service"

	ttRepository "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction/repo"
	ttService "github.com/Ruscigno/ticker-signals/internal/transaction/tradetransaction/service"

	trRepository "github.com/Ruscigno/ticker-signals/internal/transaction/traderules/repo"
	trService "github.com/Ruscigno/ticker-signals/internal/transaction/traderules/service"

	bb "github.com/Ruscigno/ticker-signals/internal/tickerbeats/beats"
	sigR "github.com/Ruscigno/ticker-signals/internal/tickerbeats/signal/repo"
	sigS "github.com/Ruscigno/ticker-signals/internal/tickerbeats/signal/service"

	"github.com/Ruscigno/ticker-signals/internal/utils"

	_ "github.com/lib/pq"
)

type WriteSyncer struct {
	io.Writer
}

func (ws WriteSyncer) Sync() error {
	return nil
}

func GetWriteSyncer(logName string) zapcore.WriteSyncer {
	var ioWriter = &lumberjack.Logger{
		Filename:   logName,
		MaxSize:    20, // MB
		MaxBackups: 5,  // number of backups
		MaxAge:     28, //days
		LocalTime:  true,
		Compress:   false, // disabled by default
	}
	var sw = WriteSyncer{
		ioWriter,
	}
	return sw
}

func SetupLogger(fileName string) *zap.Logger {
	// The bundled Config struct only supports the most common configuration
	// options. More complex needs, like splitting logs between multiple files
	// or writing to non-file outputs, require use of the zapcore package.
	//
	// In this example, imagine we're both sending our logs to Kafka and writing
	// them to the console. We'd like to encode the console output and the Kafka
	// topics differently, and we'd also like special treatment for
	// high-priority logs.

	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// Assume that we have clients for two Kafka topics. The clients implement
	// zapcore.WriteSyncer and are safe for concurrent use. (If they only
	// implement io.Writer, we can use zapcore.AddSync to add a no-op Sync
	// method. If they're not safe for concurrent use, we can add a protecting
	// mutex with zapcore.Lock.)
	logFile := GetWriteSyncer(fileName)
	topicDebugging := zapcore.AddSync(logFile)
	topicErrors := zapcore.AddSync(logFile)

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	var config zap.Config
	if os.Getenv(utils.TickerConfigLevel) == "PROD" {
		config = zap.NewProductionConfig()
		config.EncoderConfig = zap.NewProductionEncoderConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	}
	config.EncoderConfig.CallerKey = zapdriver.SourceLocation(runtime.Caller(0)).String
	configConsole := config
	configConsole.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Optimize the Kafka output for machine consumption and the console output
	// for human operators.
	kafkaEncoder := zapcore.NewJSONEncoder(config.EncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(configConsole.EncoderConfig)

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(kafkaEncoder, topicErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(kafkaEncoder, topicDebugging, lowPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	// From a zapcore.Core, it's easy to construct a Logger.
	logger := zap.New(core)
	return logger
}

func InitDatabase(dbURL string) (*sqlx.DB, error) {
	zap.L().Info("Trying to connect to the database")
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		zap.L().Fatal("Error trying to connect to the database", zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute)
	return db, err
}

type Controllers struct {
	DB     *sqlx.DB
	TrSvc  *trService.TradeRulesService
	AccSvc *accService.AccountsService
	SigSvc *sigS.SignalService
	Beats  *bb.TickerBeatsService
	TtSvc  *ttService.TradeTransactionService
}

func InitControllers(ctx context.Context, db *sqlx.DB) *Controllers {
	// Repos
	trRepo := trRepository.NewTradeRulesRepo(ctx, db)
	accRepo := accRepository.NewAccountsRepo(ctx, db)
	sigRepo := sigR.NewSignalRepository(ctx, db)
	ttRepo := ttRepository.NewTradeTransactionRepo(ctx, db)

	//Controllers
	trSvc := trService.NewTradeRulesService(ctx, trRepo)
	accSvc := accService.NewAccountsService(ctx, accRepo)
	sigSvc := sigS.NewSignalService(ctx, sigRepo)
	beats := bb.NewTickerBeatsService(ctx, sigSvc, trSvc)
	ttSvc := ttService.NewTradeTransactionService(ctx, ttRepo, beats)

	return &Controllers{
		DB:     db,
		TrSvc:  &trSvc,
		AccSvc: &accSvc,
		SigSvc: &sigSvc,
		Beats:  &beats,
		TtSvc:  &ttSvc,
	}
}
