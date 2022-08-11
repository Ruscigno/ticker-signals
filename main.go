package main

import (
	"context"

	"time"

	"github.com/Ruscigno/ticker-signals/internal/config"
	s "github.com/Ruscigno/ticker-signals/internal/server"
	"github.com/Ruscigno/ticker-signals/utils"

	"go.uber.org/zap"
)

func main() {
	start := time.Now()
	logger := utils.SetupLogger("ticker-signals.log")
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		zap.L().Fatal("unable to set timezone properly", zap.Error(err))
		return
	}
	time.Local = loc // -> this is setting the global timezone

	cfg, err := config.GetAppConfig("./config")
	if err != nil {
		zap.L().Fatal(err.Error())
		return
	}
	ctx := context.Background()

	// db, err := utils.NewMongoDBConnection(ctx, cfg)
	// if err != nil {
	// 	zap.L().Fatal(err.Error())
	// 	return
	// }
	// defer db.MongoDbDisconnect()
	// zap.L().Info("connected to MongoDB")

	s.StartHttpServer(ctx, cfg, start)
	// carRepo := repo.NewCarRepository(&ctx, db)
	// carService := service.NewCarService(&ctx, carRepo)
	// handler := handlers.NewCarHandler(carService)
	// handler.Start(cfg.ServerPort)

	zap.L().Info("Ticker Signals Server", zap.String("started", "in "+time.Since(start).String()))
	defer zap.L().Info("See ya later, alligator!")
}
