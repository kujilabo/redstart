package gateway_test

import (
	"gorm.io/gorm"

	"github.com/kujilabo/redstart/sqls"
	testlibgateway "github.com/kujilabo/redstart/testlib/gateway"
	"github.com/kujilabo/redstart/user/domain"
)

var (
	invalidOrgID domain.OrganizationID
)

func init() {
	invalidOrgIDTmp, err := domain.NewOrganizationID(99999)
	if err != nil {
		panic(err)
	}
	invalidOrgID = invalidOrgIDTmp

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
}
