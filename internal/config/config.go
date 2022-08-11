package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

// custom error for app configuration load error
var ErrConfigurationLoad = errors.New("unable to load configuration")

const (
	ApiUri                   = "api.uri"
	SiteUrl                  = "site.url"
	DatabaseDriver           = "database.driver"
	DatabaseDsn              = "database.dsn"
	DatabaseServer           = "database.server"
	DatabaseName             = "database.name"
	DatabaseConnectionString = "database.connection_string"
	DatabaseConns            = "database.conns"
	DatabaseConnsIdle        = "database.conns_idle"
	DatabaseUser             = "database.user"
	DatabasePassword         = "database.password"
	HttpHost                 = "http.host"
	HttpPort                 = "http.port"
)

var (
	// defaults for the configuration
	defaults map[string]string = map[string]string{
		ApiUri:                   os.Getenv(ApiUri),
		SiteUrl:                  os.Getenv(SiteUrl),
		DatabaseDriver:           os.Getenv(DatabaseDriver),
		DatabaseDsn:              os.Getenv(DatabaseDsn),
		DatabaseServer:           os.Getenv(DatabaseServer),
		DatabaseName:             os.Getenv(DatabaseName),
		DatabaseConnectionString: os.Getenv(DatabaseConnectionString),
		DatabaseUser:             os.Getenv(DatabaseUser),
		DatabasePassword:         os.Getenv(DatabasePassword),
		HttpHost:                 os.Getenv(HttpHost),
		HttpPort:                 os.Getenv(HttpPort),
	}
)

type Settings struct {
	ApiUri                   string `yaml:"ApiUri" json:"-" flag:"api-uri"`
	SiteUrl                  string `yaml:"SiteUrl" json:"-" flag:"site-url"`
	DatabaseDriver           string `yaml:"DatabaseDriver" json:"-" flag:"database-driver"`
	DatabaseDsn              string `yaml:"DatabaseDsn" json:"-" flag:"database-dsn"`
	DatabaseServer           string `yaml:"DatabaseServer" json:"-" flag:"database-server"`
	DatabaseName             string `yaml:"DatabaseName" json:"-" flag:"database-name"`
	DatabaseConnectionString string `yaml:"DatabaseConnectionString" json:"-" flag:"database-connection-string"`
	DatabaseConns            int    `yaml:"DatabaseConns" json:"-" flag:"database-conns"`
	DatabaseConnsIdle        int    `yaml:"DatabaseConnsIdle" json:"-" flag:"database-conns-idle"`
	DatabaseUser             string `yaml:"DatabaseUser" json:"-" flag:"database-user"`
	DatabasePassword         string `yaml:"DatabasePassword" json:"-" flag:"database-password"`
	HttpHost                 string `yaml:"HttpHost" json:"-" flag:"http-host"`
	HttpPort                 int    `yaml:"HttpPort" json:"-" flag:"http-port"`
}

// AppConfig is the application configuration
type AppConfig struct {
	once     sync.Once
	db       *gorm.DB
	settings *Settings
}

// GetAppConfig returns the application configuration
func GetAppConfig(fileName string) (*AppConfig, error) {
	setupViper(fileName)
	cfg := &AppConfig{
		settings: &Settings{
			ApiUri:                   viper.GetString(ApiUri),
			SiteUrl:                  viper.GetString(SiteUrl),
			DatabaseDriver:           viper.GetString(DatabaseDriver),
			DatabaseDsn:              viper.GetString(DatabaseDsn),
			DatabaseServer:           viper.GetString(DatabaseServer),
			DatabaseName:             viper.GetString(DatabaseName),
			DatabaseConnectionString: viper.GetString(DatabaseConnectionString),
			DatabaseConns:            viper.GetInt(DatabaseConns),
			DatabaseConnsIdle:        viper.GetInt(DatabaseConnsIdle),
			DatabaseUser:             viper.GetString(DatabaseUser),
			DatabasePassword:         viper.GetString(DatabasePassword),
			HttpHost:                 viper.GetString(HttpHost),
			HttpPort:                 viper.GetInt(HttpPort),
		},
	}
	return cfg, nil
}

// SiteUrl returns the public server URL (default is "http://localhost:2342/").
func (c *AppConfig) SiteUrl() string {
	if c.settings.SiteUrl == "" {
		return fmt.Sprintf("http://localhost:%d/", c.HttpPort())
	}

	return strings.TrimRight(c.settings.SiteUrl, "/") + "/"
}

// BaseUri returns the site base URI for a given resource.
func (c *AppConfig) BaseUri(res string) string {
	if c.SiteUrl() == "" {
		return res
	}

	u, err := url.Parse(c.SiteUrl())

	if err != nil {
		return res
	}

	return strings.TrimRight(u.EscapedPath(), "/") + res
}

// ApiUri returns the api URI.
func (c *AppConfig) ApiUri() string {
	return c.BaseUri(ApiUri)
}

func (c *AppConfig) HttpHost() string {
	if c.settings.HttpHost == "" {
		return "0.0.0.0"
	}

	return c.settings.HttpHost
}

// HttpPort returns the built-in HTTP server port.
func (c *AppConfig) HttpPort() int {
	if c.settings.HttpPort == 0 {
		return 31010
	}

	return c.settings.HttpPort
}
