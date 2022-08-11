package utils

import (
	"errors"
	"os"
	"runtime"

	"github.com/blendle/zapdriver"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// custom error for app configuration load error
var ErrConfigurationLoad = errors.New("unable to load configuration")

const (
	DatabaseName             = "DATABSE_NAME"
	DatabaseConnectionString = "DATABASE_CONNECTION_STRING"
	DatabaseUserName         = "DATABASE_USER_NAME"
	DatabasePassword         = "DATABASE_PASSWORD"
	ServerPort               = "SERVER_PORT"
)

var (
	// defaults for the configuration
	defaults map[string]string = map[string]string{
		DatabaseName:             os.Getenv(DatabaseName),
		DatabaseConnectionString: os.Getenv(DatabaseConnectionString),
		DatabaseUserName:         os.Getenv(DatabaseUserName),
		DatabasePassword:         os.Getenv(DatabasePassword),
		ServerPort:               "31000",
	}
)

// AppConfig is the application configuration
type AppConfig struct {
	DatabaseName             string
	DatabaseConnectionString string
	DatabaseUserName         string
	DatabasePassword         string
	ServerPort               string
}

// setupViper sets up the viper configuration
func setupViper(fileName string) {
	viper.SetConfigName(fileName)
	viper.AddConfigPath(".")
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		zap.L().Warn("error reading config file", zap.Error(err), zapdriver.SourceLocation(runtime.Caller(0)))
	}
}

// GetAppConfig returns the application configuration
func GetAppConfig(fileName string) (*AppConfig, error) {
	setupViper(fileName)
	cfg := &AppConfig{
		DatabaseName:             viper.GetString(DatabaseName),
		DatabaseConnectionString: viper.GetString(DatabaseConnectionString),
		DatabaseUserName:         viper.GetString(DatabaseUserName),
		DatabasePassword:         viper.GetString(DatabasePassword),
		ServerPort:               viper.GetString(ServerPort),
	}
	if err := checkConfiguration(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// checkConfiguration checks the configuration for errors
func checkConfiguration(cfg *AppConfig) error {
	if cfg.DatabaseConnectionString == "" {
		return ErrConfigurationLoad
	}
	if cfg.DatabaseName == "" {
		return ErrConfigurationLoad
	}
	if cfg.ServerPort == "" {
		return ErrConfigurationLoad
	}
	if cfg.DatabaseUserName == "" {
		return ErrConfigurationLoad
	}
	if cfg.DatabasePassword == "" {
		return ErrConfigurationLoad
	}
	return nil
}
