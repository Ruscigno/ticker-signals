package utils

import (
	"io"
	"os"
	"runtime"

	"github.com/blendle/zapdriver"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type writeSyncer struct {
	io.Writer
}

func (ws writeSyncer) Sync() error {
	return nil
}

func getWriteSyncer(logName string) zapcore.WriteSyncer {
	var ioWriter = &lumberjack.Logger{
		Filename:   logName,
		MaxSize:    20, // MB
		MaxBackups: 5,  // number of backups
		MaxAge:     28, //days
		LocalTime:  true,
		Compress:   false, // disabled by default
	}
	var sw = writeSyncer{
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
	logFile := getWriteSyncer(fileName)
	topicDebugging := zapcore.AddSync(logFile)
	topicErrors := zapcore.AddSync(logFile)

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	var config zap.Config
	//TODO: production config must be different from development config
	config = zap.NewDevelopmentConfig()
	config.EncoderConfig = zap.NewDevelopmentEncoderConfig()

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
