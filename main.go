package main

import (
	"context"

	"time"

	"github.com/Ruscigno/ticker-signals/utils"

	"go.uber.org/zap"
)

func main() {
	logger := utils.SetupLogger("ticker-signals.log")
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	zap.L().Info("starting Ticker Signals Server...")
	defer zap.L().Info("See ya later, alligator!")

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		zap.L().Fatal("unable to set timezone properly", zap.Error(err))
		return
	}
	time.Local = loc // -> this is setting the global timezone

	cfg, err := utils.GetAppConfig("./config")
	if err != nil {
		zap.L().Fatal(utils.ErrConfigurationLoad.Error())
		return
	}

	ctx := context.Background()

	db, err := utils.NewMongoDBConnection(ctx, cfg)
	if err != nil {
		zap.L().Fatal(err.Error())
		return
	}
	defer db.MongoDbDisconnect()

	zap.L().Info("connected to MongoDB")

	// carRepo := repo.NewCarRepository(&ctx, db)
	// carService := service.NewCarService(&ctx, carRepo)
	// handler := handlers.NewCarHandler(carService)
	zap.L().Info("Ticker Signals Server started")
	// handler.Start(cfg.ServerPort)
}
