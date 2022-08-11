package config

import (
	"runtime"
	"strings"

	"github.com/blendle/zapdriver"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// setupViper sets up the viper configuration
func setupViper(fileName string) {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()                // read in environment variables that match
	viper.SetEnvPrefix("tickersignals") // will be uppercased automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigName(fileName)
	viper.AddConfigPath(".")
	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		zap.L().Warn("error reading config file", zap.Error(err), zapdriver.SourceLocation(runtime.Caller(0)))
	}
}
