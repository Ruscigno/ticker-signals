package main

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"

	"os"

	v1 "github.com/Ruscigno/ruscigno-gosdk/ticker-beats/v1"
	"github.com/Ruscigno/ticker-signals/internal/api"
	"github.com/Ruscigno/ticker-signals/internal/utils"
	"github.com/Ruscigno/ticker-signals/internal/utils/app"
	"github.com/blendle/zapdriver"
	"github.com/spf13/viper"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/features/proto/echo"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	_ "github.com/lib/pq"

	"go.uber.org/zap"
)

var (
	defaults map[string]interface{} = map[string]interface{}{
		utils.TickerDatabaseURL: os.Getenv(utils.TickerDatabaseURL),
		utils.TickerConfigLevel: os.Getenv(utils.TickerConfigLevel),
		utils.TickerPort:        os.Getenv(utils.TickerPort),
	}
)

type config struct {
	dbURL       string
	serverPort  string
	configLevel string
}

var (
	cfg    *config
	system = "" // empty string represents the health of the system
)

type echoServer struct {
	pb.UnimplementedEchoServer
}

func (e *echoServer) UnaryEcho(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{
		Message: fmt.Sprintf("hello from localhost:%s", cfg.serverPort),
	}, nil
}

func setupViper() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		zap.L().Warn("fatal error reading config file", zap.Error(err), zapdriver.SourceLocation(runtime.Caller(0)))
	}
}

func main() {
	logger := app.SetupLogger("ticker-heart.log")
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	zap.L().Info("Starting ticker-heart...")
	defer zap.L().Info("See ya later, alligator!")

	loc, err := time.LoadLocation("")
	if err != nil {
		zap.L().Fatal("unable to set timezone properly", zap.Error(err), zapdriver.SourceLocation(runtime.Caller(0)))
	}
	time.Local = loc // -> this is setting the global timezone

	setupViper()
	cfg := &config{
		serverPort:  viper.GetString(utils.TickerPort),
		dbURL:       strings.ReplaceAll(viper.GetString(utils.TickerDatabaseURL), ";", " "),
		configLevel: strings.ToUpper(viper.GetString(utils.TickerConfigLevel)),
	}
	if cfg.configLevel == "TEST" || cfg.configLevel == "" {
		zap.L().Fatal("prod environment", zap.Error(err), zapdriver.SourceLocation(runtime.Caller(0)))
	}

	ctx := context.Background()
	if cfg.serverPort == "" || cfg.dbURL == "" {
		zap.L().Fatal("Ohh no! All env variables must be set before starting",
			zap.String(utils.TickerDatabaseURL, "***hidden***"),
			zapdriver.SourceLocation(runtime.Caller(0)))
	}
	db, err := app.InitDatabase(cfg.dbURL)
	if err != nil {
		os.Exit(-1)
	}
	defer db.Close()
	zap.L().Info("Uhuuu, successfully connected to the database")

	svc := app.InitControllers(ctx, db)

	lis, err := net.Listen("tcp", ":"+cfg.serverPort)
	if err != nil {
		zap.L().Fatal("failed to start the health check server", zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
	}
	zap.L().Info("Listening on port", zap.String("port", cfg.serverPort))
	s := grpc.NewServer()
	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(s, healthcheck)
	pb.RegisterEchoServer(s, &echoServer{})

	// Servers
	ticServer := api.NewTickerBeatsServiceServer(*svc.Beats, *svc.TrSvc)
	v1.RegisterTickerBeatsServiceServer(s, ticServer)

	go func() {
		// asynchronously inspect dependencies and toggle serving status as needed
		next := healthpb.HealthCheckResponse_SERVING
		zap.L().Info("Health Check Up & Running")
		for {
			healthcheck.SetServingStatus(system, next)

			if next == healthpb.HealthCheckResponse_SERVING {
				next = healthpb.HealthCheckResponse_NOT_SERVING
			} else {
				next = healthpb.HealthCheckResponse_SERVING
			}

			time.Sleep(time.Second * 5)
		}
	}()

	// See ya later, alligator!
	if err := s.Serve(lis); err != nil {
		zap.L().Fatal("failed to serve", zap.Error(err),
			zapdriver.SourceLocation(runtime.Caller(0)))
	}
}
