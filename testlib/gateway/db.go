package gateway

import (
	"database/sql"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source"
	"gorm.io/gorm"

	liberrors "github.com/kujilabo/redstart/lib/errors"
)

func ListDB() map[string]*gorm.DB {
	list := make(map[string]*gorm.DB)
	m, err := openMySQLForTest()
	if err != nil {
		panic(err)
	}
	list["mysql"] = m

	// s, err := openSQLiteForTest()
	// if err != nil {
	// 	panic(err)
	// }
	// list["sqlite3"] = s

	return list
}

func setupDB(db *gorm.DB, driverName string, sourceDriver source.Driver, getDatabaseDriver func(sqlDB *sql.DB) (database.Driver, error)) error {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
		return err
	}

	databaseDriver, err := getDatabaseDriver(sqlDB)
	if err != nil {
		log.Fatal(liberrors.Errorf("failed to WithInstance. err: %w", err))
		return err
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, driverName, databaseDriver)
	if err != nil {
		log.Fatal(liberrors.Errorf("failed to NewWithDatabaseInstance. err: %w", err))
		return err
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(liberrors.Errorf("failed to Up. driver:%s, err: %w", driverName, err))
			return err
		}
	}

	return nil
}
