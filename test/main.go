package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/blendle/zapdriver"
	"go.uber.org/zap"

	"github.com/Ruscigno/ticker-signals/internal/api"
	"github.com/Ruscigno/ticker-signals/internal/utils"
	"github.com/Ruscigno/ticker-signals/internal/utils/app"
	"github.com/Ruscigno/ticker-signals/test/insert"
)

func main() {
	start := time.Now().UTC()
	loc, err := time.LoadLocation("")
	if err != nil {
		zap.L().Fatal("unable to set timezone properly", zap.Error(err), zapdriver.SourceLocation(runtime.Caller(0)))
	}
	// handle err
	time.Local = loc // -> this is setting the global timezone

	logger := app.SetupLogger("ticker-signals-test.log")
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()
	zap.L().Info("Testing ticker-signals...")

	dbURL := strings.ReplaceAll(os.Getenv(utils.TickerDatabaseURL), ";", " ")

	if dbURL == "" {
		zap.L().Fatal("Ohh no! All env variables must be set before starting",
			zap.String(utils.TickerDatabaseURL, dbURL),
			// zap.String(utils.TickerServerAddress, address),
			zapdriver.SourceLocation(runtime.Caller(0)))
		os.Exit(-1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	db, err := app.InitDatabase(dbURL)
	if err != nil {
		os.Exit(-1)
	}
	defer db.Close()
	zap.L().Info("Uhuuu, successfully connected to the database")

	svc := app.InitControllers(ctx, db)

	apiServer := api.NewTransactionsServiceServer(*svc.AccSvc, *svc.InfoSvc, *svc.DeaSvc, *svc.OrdSvc, *svc.PosSvc, *svc.TtSvc, *svc.Beats, *svc.TrSvc)
	if !startTests(&ctx, svc, apiServer) {
		os.Exit(-1)
	}
	zap.L().Info("See ya later, alligator!", zap.String("execution time", fmt.Sprintf("%d", time.Since(start).Milliseconds())))
}

func startTests(ctx *context.Context, svc *app.Controllers, apiServer *api.TransactionsServiceServer) bool {
	return dropTables(ctx, svc) &&
		createTables(ctx, svc) &&
		insert.TestInsert(ctx, svc, apiServer)
}

func createTables(ctx *context.Context, svc *app.Controllers) bool {
	return executeScript(ctx, svc, "../database/postgres/create-database.sql")
}

func dropTables(ctx *context.Context, svc *app.Controllers) bool {
	return executeScript(ctx, svc, "../database/postgres/delete-entities.sql")
}

func executeScript(ctx *context.Context, svc *app.Controllers, fileName string) bool {
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		zap.L().Fatal("unable create database", zap.Error(err), zapdriver.SourceLocation(runtime.Caller(0)))
		return false
	}
	sRaw := strings.ReplaceAll(string(raw), "\n", "")
	rawLines := strings.Split(sRaw, ";")
	for _, line := range rawLines {
		if strings.HasPrefix(line, "--") {
			continue
		}
		tx := svc.DB.MustBegin()
		_, err = svc.DB.Exec(line)
		if err != nil {
			zap.L().Fatal("unable create database", zap.String("script", line), zap.Error(err), zapdriver.SourceLocation(runtime.Caller(0)))
			return false
		}
		err = tx.Commit()
		if err != nil {
			zap.L().Fatal("transaction error", zap.String("script", line), zap.Error(err), zapdriver.SourceLocation(runtime.Caller(0)))
			return false
		}
	}
	return true
}
