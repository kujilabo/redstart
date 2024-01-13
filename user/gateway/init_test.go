package gateway_test

import (
	"context"
	"os"
	"strconv"

	"gorm.io/gorm"

	"github.com/kujilabo/redstart/sqls"
	testlibgateway "github.com/kujilabo/redstart/testlib/gateway"
	"github.com/kujilabo/redstart/user/domain"
	"github.com/kujilabo/redstart/user/gateway"
)

var (
	invalidOrgID     *domain.OrganizationID
	invalidAppUserID *domain.AppUserID
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func init() {
	invalidOrgIDTmp, err := domain.NewOrganizationID(99999)
	if err != nil {
		panic(err)
	}
	invalidOrgID = invalidOrgIDTmp

	invalidAppUserIDTmp, err := domain.NewAppUserID(99999)
	if err != nil {
		panic(err)
	}
	invalidAppUserID = invalidAppUserIDTmp

	mysqlHost := getEnv("MYSQL_HOST", "127.0.0.1")
	mysqlPortS := getEnv("MYSQL_PORT", "3307")
	mysqlPort, err := strconv.Atoi(mysqlPortS)
	if err != nil {
		panic(err)
	}
	postgresHost := getEnv("POSTGRES_HOST", "127.0.0.1")
	postgresPortS := getEnv("POSTGRES_PORT", "5433")
	postgresPort, err := strconv.Atoi(postgresPortS)
	if err != nil {
		panic(err)
	}
	fns := []func() (*gorm.DB, error){
		func() (*gorm.DB, error) {
			return testlibgateway.InitMySQL(sqls.SQL, mysqlHost, mysqlPort)
		},
		func() (*gorm.DB, error) {
			return testlibgateway.InitPostgres(sqls.SQL, postgresHost, postgresPort)
		},
		// func() (*gorm.DB, error) {
		// 	return testlibgateway.InitSQLiteInFile(sqls.SQL)
		// },
	}

	for _, fn := range fns {
		db, err := fn()
		if err != nil {
			panic(err)
		}
		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}
		sqlDB.Close()
	}

	ctx := context.Background()
	for dialect, db := range testlibgateway.ListDB() {
		dialect := dialect
		db := db
		rf, err := gateway.NewRepositoryFactory(ctx, dialect, dialect.Name(), db, loc)
		if err != nil {
			panic(err)
		}
		authorizationManager := rf.NewAuthorizationManager(ctx)
		if err := authorizationManager.Init(ctx); err != nil {
			panic(err)
		}
	}
}
