package gateway

import (
	"database/sql"
	"embed"

	"github.com/golang-migrate/migrate/v4/database"
	migrate_sqlite3 "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	gorm_logrus "github.com/onrik/gorm-logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenSQLite(filePath string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(filePath), &gorm.Config{
		Logger: gorm_logrus.New(),
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
