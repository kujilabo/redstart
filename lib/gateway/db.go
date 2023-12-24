package gateway

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"

	liberrors "github.com/kujilabo/redstart/lib/errors"
)

const MYSQL_ER_DUP_ENTRY = 1062
const MYSQL_ER_NO_REFERENCED_ROW_2 = 1452

const SQLITE_CONSTRAINT_PRIMARYKEY = 1555
const SQLITE_CONSTRAINT_UNIQUE = 2067

func ConvertDuplicatedError(err error, newErr error) error {
	var mysqlErr *mysql.MySQLError
	if ok := errors.As(err, &mysqlErr); ok && mysqlErr.Number == MYSQL_ER_DUP_ENTRY {
		return newErr
	}

	// var sqlite3Err sqlite3.Error
	// if ok := errors.As(err, &sqlite3Err); ok {
	// 	if int(sqlite3Err.ExtendedCode) == SQLITE_CONSTRAINT_PRIMARYKEY {
	// 		return newErr
	// 	} else if int(sqlite3Err.ExtendedCode) == SQLITE_CONSTRAINT_UNIQUE {
	// 		return newErr
	// 	}
	// }

	return err
}

func ConvertRelationError(err error, newErr error) error {
	var mysqlErr *mysql.MySQLError
	if ok := errors.As(err, &mysqlErr); ok && mysqlErr.Number == MYSQL_ER_NO_REFERENCED_ROW_2 {
		return newErr
	}

	return err
}

func migrateDB(db *gorm.DB, driverName string, sourceDriver source.Driver, getDatabaseDriver func(sqlDB *sql.DB) (database.Driver, error)) error {
	sqlDB, err := db.DB()
	if err != nil {
		return liberrors.Errorf("db.DB in gateway.migrateDB. err: %w", err)
	}

	databaseDriver, err := getDatabaseDriver(sqlDB)
	if err != nil {
		return liberrors.Errorf("getDatabaseDriver in gateway.migrateDB. err: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, driverName, databaseDriver)
	if err != nil {
		return liberrors.Errorf("NewWithInstance in gateway.migrateDB. err: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return liberrors.Errorf("failed to m.Up in gateway.migrateDB. err: %w", err)
	}

	return nil
}
