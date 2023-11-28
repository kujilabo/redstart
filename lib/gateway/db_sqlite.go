package gateway

import (
	"database/sql"
	"embed"
	"log/slog"

	"github.com/golang-migrate/migrate/v4/database"
	migrate_sqlite3 "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	slog_gorm "github.com/orandin/slog-gorm"
	gorm_sqlite "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenSQLite(filePath string, logger *slog.Logger) (*gorm.DB, error) {
	return gorm.Open(gorm_sqlite.Open(filePath), &gorm.Config{
		Logger: slog_gorm.New(
			slog_gorm.WithLogger(logger), // Optional, use slog.Default() by default
			slog_gorm.WithTraceAll(),     // trace all messages
		),
	})
}

func MigrateSQLiteDB(db *gorm.DB, sqlFS embed.FS) error {
	driverName := "sqlite3"
	sourceDriver, err := iofs.New(sqlFS, driverName)
	if err != nil {
		return err
	}

	return migrateDB(db, "sqlite3", sourceDriver, func(sqlDB *sql.DB) (database.Driver, error) {
		return migrate_sqlite3.WithInstance(sqlDB, &migrate_sqlite3.Config{})
	})
}
