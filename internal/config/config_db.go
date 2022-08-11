package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"go.uber.org/zap"
	"moul.io/zapgorm2"

	"github.com/Ruscigno/ticker-signals/internal/mutex"
)

// SQL Databases.
// TODO: PostgresSQL support requires upgrading GORM, so generic column data types can be used.
const (
	MySQL    = "mysql"
	MariaDB  = "mariadb"
	Postgres = "postgres"
)

// SQLite default DSNs.
const (
	SQLiteTestDB    = ".test.db"
	SQLiteMemoryDSN = ":memory:"
)

// dsnPattern is a regular expression matching a database DSN string.
var dsnPattern = regexp.MustCompile(
	`^(?:(?P<user>.*?)(?::(?P<password>.*))?@)?` +
		`(?:(?P<net>[^\(]*)(?:\((?P<server>[^\)]*)\))?)?` +
		`\/(?P<name>.*?)` +
		`(?:\?(?P<params>[^\?]*))?$`)

// DatabaseDriver returns the database driver name.
func (c *AppConfig) DatabaseDriver() string {
	switch strings.ToLower(c.settings.DatabaseDriver) {
	case MySQL, MariaDB:
		c.settings.DatabaseDriver = MySQL
	case Postgres:
		c.settings.DatabaseDriver = Postgres
	default:
		zap.L().Warn("config: unsupported database driver, using postgres", zap.String("driver", c.settings.DatabaseDriver))
		c.settings.DatabaseDriver = Postgres
		c.settings.DatabaseDsn = ""
	}

	return c.settings.DatabaseDriver
}

// DatabaseDsn returns the database data source name (DSN).
func (c *AppConfig) DatabaseDsn() string {
	if c.settings.DatabaseDsn == "" {
		switch c.DatabaseDriver() {
		case MySQL, MariaDB:
			address := c.DatabaseServer()

			// Connect via TCP or Unix Domain Socket?
			if strings.HasPrefix(address, "/") {
				zap.L().Debug("mariadb: connecting via Unix domain socket")
				address = fmt.Sprintf("unix(%s)", address)
			} else {
				address = fmt.Sprintf("tcp(%s)", address)
			}

			return fmt.Sprintf(
				"%s:%s@%s/%s?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true",
				c.DatabaseUser(),
				c.DatabasePassword(),
				address,
				c.DatabaseName(),
			)
		case Postgres:
			return fmt.Sprintf(
				"user=%s password=%s dbname=%s host=%s port=%d sslmode=disable TimeZone=UTC",
				c.DatabaseUser(),
				c.DatabasePassword(),
				c.DatabaseName(),
				c.DatabaseHost(),
				c.DatabasePort(),
			)
		default:
			zap.L().Error("config: empty database dsn")
			return ""
		}
	}

	return c.settings.DatabaseDsn
}

// ParseDatabaseDsn parses the database dsn and extracts user, password, database server, and name.
func (c *AppConfig) ParseDatabaseDsn() {
	if c.settings.DatabaseDsn == "" || c.settings.DatabaseServer != "" {
		return
	}

	matches := dsnPattern.FindStringSubmatch(c.settings.DatabaseDsn)
	names := dsnPattern.SubexpNames()

	for i, match := range matches {
		switch names[i] {
		case "user":
			c.settings.DatabaseUser = match
		case "password":
			c.settings.DatabasePassword = match
		case "server":
			c.settings.DatabaseServer = match
		case "name":
			c.settings.DatabaseName = match
		}
	}
}

// DatabaseServer the database server.
func (c *AppConfig) DatabaseServer() string {
	c.ParseDatabaseDsn()

	if c.settings.DatabaseServer == "" {
		return "localhost"
	}

	return c.settings.DatabaseServer
}

// DatabaseHost the database server host.
func (c *AppConfig) DatabaseHost() string {
	if s := strings.Split(c.DatabaseServer(), ":"); len(s) > 0 {
		return s[0]
	}

	return c.settings.DatabaseServer
}

// DatabasePort the database server port.
func (c *AppConfig) DatabasePort() int {
	const defaultPort = 5432

	if s := strings.Split(c.DatabaseServer(), ":"); len(s) != 2 {
		return defaultPort
	} else if port, err := strconv.Atoi(s[1]); err != nil {
		return defaultPort
	} else if port < 1 || port > 65535 {
		return defaultPort
	} else {
		return port
	}
}

// DatabasePortString the database server port as string.
func (c *AppConfig) DatabasePortString() string {
	return strconv.Itoa(c.DatabasePort())
}

// DatabaseName the database schema name.
func (c *AppConfig) DatabaseName() string {
	c.ParseDatabaseDsn()

	if c.settings.DatabaseName == "" {
		return "tickerheart"
	}

	return c.settings.DatabaseName
}

// DatabaseUser returns the database user name.
func (c *AppConfig) DatabaseUser() string {
	c.ParseDatabaseDsn()

	if c.settings.DatabaseUser == "" {
		return "tickerheart"
	}

	return c.settings.DatabaseUser
}

// DatabasePassword returns the database user password.
func (c *AppConfig) DatabasePassword() string {
	c.ParseDatabaseDsn()

	return c.settings.DatabasePassword
}

// DatabaseConns returns the maximum number of open connections to the database.
func (c *AppConfig) DatabaseConns() int {
	limit := c.settings.DatabaseConns

	if limit <= 0 {
		limit = (runtime.NumCPU() * 2) + 16
	}

	if limit > 1024 {
		limit = 1024
	}

	return limit
}

// DatabaseConnsIdle returns the maximum number of idle connections to the database (equal or less than open).
func (c *AppConfig) DatabaseConnsIdle() int {
	limit := c.settings.DatabaseConnsIdle

	if limit <= 0 {
		limit = runtime.NumCPU() + 8
	}

	if limit > c.DatabaseConns() {
		limit = c.DatabaseConns()
	}

	return limit
}

// Db returns the db connection.
func (c *AppConfig) Db() *gorm.DB {
	if c.db == nil {
		zap.L().Fatal("config: database not connected")
	}

	return c.db
}

// CloseDb closes the db connection (if any).
func (c *AppConfig) CloseDb() error {
	if c.db != nil {
		if err := c.db.Close(); err == nil {
			c.db = nil
		} else {
			return err
		}
	}

	return nil
}

// SetDbOptions sets the database collation to unicode if supported.
func (c *AppConfig) SetDbOptions() {
	switch c.DatabaseDriver() {
	case MySQL, MariaDB:
		c.Db().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	case Postgres:
		// Ignore for now.
	}
}

// InitDb initializes the database without running previously failed migrations.
func (c *AppConfig) InitDb() {
	c.MigrateDb(false, nil)
}

// MigrateDb initializes the database and migrates the schema if needed.
func (c *AppConfig) MigrateDb(runFailed bool, ids []string) {
	// c.SetDbOptions()
	// entity.SetDbProvider(c)
	// entity.InitDb(true, runFailed, ids)

	// entity.Admin.InitPassword(c.AdminPassword())

	// go entity.SaveErrorMessages()
}

// InitTestDb drops all tables in the currently configured database and re-creates them.
func (c *AppConfig) InitTestDb() {
	// c.SetDbOptions()
	// entity.SetDbProvider(c)
	// entity.ResetTestFixtures()

	// entity.Admin.InitPassword(c.AdminPassword())

	// go entity.SaveErrorMessages()
}

// connectDb establishes a database connection.
func (c *AppConfig) connectDb() error {
	mutex.Db.Lock()
	defer mutex.Db.Unlock()

	dbDriver := c.DatabaseDriver()
	dbDsn := c.DatabaseDsn()

	if dbDriver == "" {
		return errors.New("config: database driver not specified")
	}

	if dbDsn == "" {
		return errors.New("config: database DSN not specified")
	}

	log := zapgorm2.New(zap.L())
	log.SetAsDefault()
	// TODO: open with the correct options.
	// db, err := gorm.Open(dbDriver, dbDsn, &gorm.Config{Logger: log})
	db, err := gorm.Open(dbDriver, dbDsn)
	if err != nil || db == nil {
		for i := 1; i <= 12; i++ {
			db, err = gorm.Open(dbDriver, dbDsn)

			if db != nil && err == nil {
				break
			}

			time.Sleep(5 * time.Second)
		}

		if err != nil || db == nil {
			zap.L().Fatal(err.Error())
		}
	}

	db.LogMode(false)
	// TODO: open with the correct options.
	// db.SetLogger(log)

	db.DB().SetMaxOpenConns(c.DatabaseConns())
	db.DB().SetMaxIdleConns(c.DatabaseConnsIdle())
	db.DB().SetConnMaxLifetime(10 * time.Minute)

	c.db = db

	return err
}

// ImportSQL imports a file to the currently configured database.
func (c *AppConfig) ImportSQL(filename string) {
	contents, err := os.ReadFile(filename)

	if err != nil {
		zap.L().Error(err.Error())
		return
	}

	statements := strings.Split(string(contents), ";\n")
	q := c.Db().Unscoped()

	for _, stmt := range statements {
		// Skip empty lines and comments
		if len(stmt) < 3 || stmt[0] == '#' || stmt[0] == ';' {
			continue
		}

		var result struct{}

		q.Raw(stmt).Scan(&result)
	}
}
