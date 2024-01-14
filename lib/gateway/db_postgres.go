package gateway

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"
	"time"

	"github.com/golang-migrate/migrate/v4/database"
	migrate_postgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	slog_gorm "github.com/orandin/slog-gorm"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenPostgres(username, password, host string, port int, database string, logger *slog.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", host, username, password, database, port, "disable", time.UTC.String())

	return gorm.Open(gorm_postgres.New(gorm_postgres.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger: slog_gorm.New(
			slog_gorm.WithLogger(logger), // Optional, use slog.Default() by default
			slog_gorm.WithTraceAll(),     // trace all messages
		),
	})
}

func MigratePostgresDB(db *gorm.DB, sqlFS fs.FS) error {
	driverName := "postgres"
	sourceDriver, err := iofs.New(sqlFS, driverName)
	if err != nil {
		return err
	}

	return migrateDB(db, driverName, sourceDriver, func(sqlDB *sql.DB) (database.Driver, error) {
		return migrate_postgres.WithInstance(sqlDB, &migrate_postgres.Config{})
	})
}
