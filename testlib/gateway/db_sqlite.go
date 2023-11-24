package gateway

import (
	"database/sql"
	"embed"
	"os"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
	gorm_logrus "github.com/onrik/gorm-logrus"
	gormSQLite "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testDBFile string

func openSQLiteForTest() (*gorm.DB, error) {
	return gorm.Open(gormSQLite.Open(testDBFile), &gorm.Config{
		Logger: gorm_logrus.New(),
	})
}

func OpenSQLiteInMemory(sqlFS embed.FS) (*gorm.DB, error) {
	db, err := gorm.Open(gormSQLite.Open("file:memdb1?mode=memory&cache=shared"), &gorm.Config{
		Logger: gorm_logrus.New(),
	})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if err := setupSQLite(sqlFS, db); err != nil {
		return nil, err
	}
	return db, nil
}

func setupSQLite(sqlFS embed.FS, db *gorm.DB) error {
	driverName := "sqlite3"
	sourceDriver, err := iofs.New(sqlFS, driverName)
	if err != nil {
		return err
	}
	return setupDB(db, driverName, sourceDriver, func(sqlDB *sql.DB) (database.Driver, error) {
		return sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	})
}

func InitSQLiteInFile(sqlFS embed.FS) (*gorm.DB, error) {
	testDBFile = "./test.db"
	os.Remove(testDBFile)
	db, err := openSQLiteForTest()
	if err != nil {
		return nil, err
	}
	if err := setupSQLite(sqlFS, db); err != nil {
		return nil, err
	}
	return db, nil
}
