package gateway_test

import (
	"context"

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

	fns := []func() (*gorm.DB, error){
		func() (*gorm.DB, error) {
			return testlibgateway.InitMySQL(sqls.SQL, "127.0.0.1", 3307)
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
	for driverName, db := range testlibgateway.ListDB() {
		rf, err := gateway.NewRepositoryFactory(ctx, driverName, db, loc)
		if err != nil {
			panic(err)
		}
		authorizationManager := rf.NewAuthorizationManager(ctx)
		if err := authorizationManager.Init(ctx); err != nil {
			panic(err)
		}
	}
}
