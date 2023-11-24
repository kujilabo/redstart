package gateway

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4/database"
	migrate_mysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	slog_gorm "github.com/orandin/slog-gorm"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenMySQL(username, password, host string, port int, database string, logger *slog.Logger) (*gorm.DB, error) {
	c := mysql.Config{
		DBName:               database,
		User:                 username,
		Passwd:               password,
		Addr:                 fmt.Sprintf("%s:%d", host, port),
		Net:                  "tcp",
		ParseTime:            true,
		MultiStatements:      true,
		Params:               map[string]string{"charset": "utf8"},
		Collation:            "utf8mb4_unicode_ci",
		AllowNativePasswords: true,
		Loc:                  time.UTC,
	}
	dsn := c.FormatDSN()
	return gorm.Open(gorm_mysql.Open(dsn), &gorm.Config{
		Logger: slog_gorm.New(
			slog_gorm.WithLogger(logger), // Optional, use slog.Default() by default
			slog_gorm.WithTraceAll(),     // trace all messages
		),
	})
}

func MigrateMySQLDB(db *gorm.DB, sqlFS embed.FS) error {
	driverName := "mysql"
	sourceDriver, err := iofs.New(sqlFS, driverName)
	if err != nil {
		return err
	}

	return migrateDB(db, driverName, sourceDriver, func(sqlDB *sql.DB) (database.Driver, error) {
		return migrate_mysql.WithInstance(sqlDB, &migrate_mysql.Config{})
	})
}
