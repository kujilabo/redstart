package config

import (
	"database/sql"
	"embed"
	"log/slog"
	"os"

	"gorm.io/gorm"

	libdomain "github.com/kujilabo/redstart/lib/domain"
	liberrors "github.com/kujilabo/redstart/lib/errors"
	libgateway "github.com/kujilabo/redstart/lib/gateway"
)

type SQLite3Config struct {
	File string `yaml:"file" validate:"required"`
}

type MySQLConfig struct {
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"required"`
	Database string `yaml:"database" validate:"required"`
}

type PostgresConfig struct {
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"required"`
	Database string `yaml:"database" validate:"required"`
}

type DBConfig struct {
	DriverName string          `yaml:"driverName"`
	SQLite3    *SQLite3Config  `yaml:"sqlite3"`
	MySQL      *MySQLConfig    `yaml:"mysql"`
	Postgres   *PostgresConfig `yaml:"postgres"`
	Migration  bool            `yaml:"migration"`
}

func InitDB(cfg *DBConfig, sqlFS embed.FS) (libgateway.DialectRDBMS, *gorm.DB, *sql.DB, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	switch cfg.DriverName {
	// case "sqlite3":
	// 	db, err := libgateway.OpenSQLite("./"+cfg.SQLite3.File, logger)
	// 	if err != nil {
	// 		return nil, nil, liberrors.Errorf("OpenSQLite. err: %w", err)
	// 	}

	// 	sqlDB, err := db.DB()
	// 	if err != nil {
	// 		return nil, nil, err
	// 	}

	// 	if err := sqlDB.Ping(); err != nil {
	// 		return nil, nil, err
	// 	}

	// 	if cfg.Migration {
	// 		if err := libgateway.MigrateSQLiteDB(db, sqlFS); err != nil {
	// 			return nil, nil, err
	// 		}
	// 	}

	// 	return db, sqlDB, nil
	case "mysql":
		db, err := libgateway.OpenMySQL(cfg.MySQL.Username, cfg.MySQL.Password, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database, logger)
		if err != nil {
			return nil, nil, nil, err
		}

		sqlDB, err := db.DB()
		if err != nil {
			return nil, nil, nil, err
		}

		if err := sqlDB.Ping(); err != nil {
			return nil, nil, nil, err
		}

		if err := libgateway.MigrateMySQLDB(db, sqlFS); err != nil {
			return nil, nil, nil, liberrors.Errorf("failed to MigrateMySQLDB. err: %w", err)
		}

		dialect := libgateway.DialectMySQL{}
		return &dialect, db, sqlDB, nil
	case "postgres":
		db, err := libgateway.OpenPostgres(cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Database, logger)
		if err != nil {
			return nil, nil, nil, err
		}

		sqlDB, err := db.DB()
		if err != nil {
			return nil, nil, nil, err
		}

		if err := sqlDB.Ping(); err != nil {
			return nil, nil, nil, err
		}

		if err := libgateway.MigratePostgresDB(db, sqlFS); err != nil {
			return nil, nil, nil, liberrors.Errorf("failed to MigrateMySQLDB. err: %w", err)
		}

		dialect := libgateway.DialectPostgres{}
		return &dialect, db, sqlDB, nil
	default:
		return nil, nil, nil, libdomain.ErrInvalidArgument
	}
}
