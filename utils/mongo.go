package utils

import (
	"context"
	"runtime"

	"github.com/blendle/zapdriver"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDbConnection struct {
	ctx    context.Context
	Client *mongo.Client
	Cars   *mongo.Collection
}

// Connects to mongoDB and returns a connection
func NewMongoDBConnection(ctx context.Context, cfg *AppConfig) (*MongoDbConnection, error) {
	client, col, err := connectToMongoDB(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &MongoDbConnection{ctx: ctx, Client: client, Cars: col}, nil
}

// connectToMongoDB connects to mongoDB and returns a connection
func connectToMongoDB(ctx context.Context, cfg *AppConfig) (*mongo.Client, *mongo.Collection, error) {
	opt := options.Client().ApplyURI(cfg.DatabaseConnectionString)
	if cfg.DatabaseUserName != "" || cfg.DatabasePassword != "" {
		opt.SetAuth(options.Credential{
			Username: cfg.DatabaseUserName, Password: cfg.DatabasePassword,
		})
	}
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		zap.L().Error("unable to ping mongoDB", zap.Error(err), zapdriver.SourceLocation(runtime.Caller(0)))
	}
	return client, client.Database(cfg.DatabaseName).Collection("cars", nil), nil
}

// MongoDbDisconnect disconnects from mongoDB
func (s *MongoDbConnection) MongoDbDisconnect() {
	err := s.Client.Disconnect(s.ctx)
	if err != nil {
		zap.L().Error("unable to disconnect from mongoDB", zap.Error(err), zapdriver.SourceLocation(runtime.Caller(0)))
	}
}
