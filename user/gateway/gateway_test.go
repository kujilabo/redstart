package gateway_test

import (
	"context"
	"crypto/rand"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	testlibgateway "github.com/kujilabo/redstart/testlib/gateway"
	"github.com/kujilabo/redstart/user/domain"
	"github.com/kujilabo/redstart/user/gateway"
	"github.com/kujilabo/redstart/user/service"
)

const orgNameLength = 8

type testService struct {
	driverName string
	db         *gorm.DB
	rf         service.RepositoryFactory
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var (
	loc = time.UTC
)

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		val, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
		if err != nil {
			panic(err)
		}
		b[i] = letterRunes[val.Int64()]
	}
	return string(b)
}

func testDB(t *testing.T, fn func(t *testing.T, ctx context.Context, ts testService)) {
	ctx := context.Background()
	for driverName, db := range testlibgateway.ListDB() {
		driverName := driverName
		db := db
		t.Run(driverName, func(t *testing.T) {
			// t.Parallel()
			sqlDB, err := db.DB()
			require.NoError(t, err)
			defer sqlDB.Close()

			rf, err := gateway.NewRepositoryFactory(ctx, driverName, db, loc)
			require.NoError(t, err)
			testService := testService{driverName: driverName, db: db, rf: rf}

			fn(t, ctx, testService)
		})
	}
}

func setupOrganization(ctx context.Context, t *testing.T, ts testService) (domain.OrganizationID, service.Owner) {
	bg := context.Background()
	orgName := RandString(orgNameLength)
	sysAd, err := service.NewSystemAdmin(ctx, ts.rf)
	assert.NoError(t, err)

	firstOwnerAddParam, err := service.NewFirstOwnerAddParameter("OWNER_ID", "OWNER_NAME", "")
	assert.NoError(t, err)
	orgAddParam, err := service.NewOrganizationAddParameter(orgName, firstOwnerAddParam)
	assert.NoError(t, err)

	orgRepo := gateway.NewOrganizationRepository(ctx, ts.db)

	// register new organization
	orgID, err := orgRepo.AddOrganization(bg, sysAd, orgAddParam)
	assert.NoError(t, err)
	assert.Greater(t, orgID.Int(), 0)

	appUserRepo := gateway.NewAppUserRepository(ctx, ts.driverName, ts.db, ts.rf)
	sysOwnerID, err := appUserRepo.AddSystemOwner(bg, sysAd, orgID)
	assert.NoError(t, err)
	assert.Greater(t, sysOwnerID.Int(), 0)

	sysOwner, err := appUserRepo.FindSystemOwnerByOrganizationName(bg, sysAd, orgName)
	assert.NoError(t, err)

	firstOwnerID, err := appUserRepo.AddFirstOwner(bg, sysOwner, firstOwnerAddParam)
	assert.NoError(t, err)
	assert.Greater(t, firstOwnerID.Int(), 0)

	firstOwner, err := appUserRepo.FindOwnerByLoginID(bg, sysOwner, "OWNER_ID")
	assert.NoError(t, err)

	// spaceRepo := gateway.NewSpaceRepository(ctx, ts.db)
	// _, err = spaceRepo.AddDefaultSpace(bg, sysOwner)
	// assert.NoError(t, err)
	// _, err = spaceRepo.AddPersonalSpace(bg, sysOwner, firstOwner)
	// assert.NoError(t, err)

	return orgID, firstOwner
}

func teardownOrganization(t *testing.T, ts testService, orgID domain.OrganizationID) {
	// delete all organizations
	ts.db.Exec("delete from space where organization_id = ?", orgID.Int())
	ts.db.Exec("delete from app_user where organization_id = ?", orgID.Int())
	ts.db.Exec("delete from organization where id = ?", orgID.Int())
	// db.Where("true").Delete(&spaceEntity{})
	// db.Where("true").Delete(&appUserEntity{})
	// db.Where("true").Delete(&organizationEntity{})
}

// func setupOrganization(ctx context.Context, t *testing.T, ts testService) (domain.OrganizationID, service.SystemOwner, service.Owner) {
// 	orgName := RandString(orgNameLength)
// 	userRf, err := ts.rf.NewUserRepositoryFactory(ctx)
// 	require.NoError(t, err)
// 	sysAd, err := service.NewSystemAdmin(ctx, userRf)
// 	require.NoError(t, err)

// 	firstOwnerAddParam, err := service.NewFirstOwnerAddParameter("OWNER_ID", "OWNER_NAME", "")
// 	require.NoError(t, err)
// 	orgAddParam, err := service.NewOrganizationAddParameter(orgName, firstOwnerAddParam)
// 	require.NoError(t, err)
// 	orgRepo := gateway.NewOrganizationRepository(ctx, ts.db)

// 	// register new organization
// 	orgID, err := orgRepo.AddOrganization(ctx, sysAd, orgAddParam)
// 	require.NoError(t, err)
// 	require.Greater(t, orgID.Int(), 0)
// 	logrus.Debugf("OrgID: %d \n", orgID.Int())
// 	org, err := orgRepo.FindOrganizationByID(ctx, sysAd, orgID)
// 	require.NoError(t, err)
// 	logrus.Debugf("OrgID: %d \n", org.GetID().Int())

// 	appUserRepo := gateway.NewAppUserRepository(ctx, ts.driverName, ts.db, ts.rf)
// 	sysOwnerID, err := appUserRepo.AddSystemOwner(ctx, sysAd, orgID)
// 	require.NoError(t, err)
// 	require.Greater(t, sysOwnerID.Int(), 0)

// 	sysOwner, err := appUserRepo.FindSystemOwnerByOrganizationName(ctx, sysAd, orgName)
// 	require.NoError(t, err)
// 	require.Greater(t, int(uint(sysOwnerID)), 0)

// 	firstOwnerID, err := appUserRepo.AddFirstOwner(ctx, sysOwner, firstOwnerAddParam)
// 	require.NoError(t, err)
// 	require.Greater(t, int(uint(firstOwnerID)), 0)

// 	firstOwner, err := appUserRepo.FindOwnerByLoginID(ctx, sysOwner, "OWNER_ID")
// 	require.NoError(t, err)

// 	spaceRepo := userG.NewSpaceRepository(ctx, ts.db)
// 	_, err = spaceRepo.AddDefaultSpace(ctx, sysOwner)
// 	require.NoError(t, err)
// 	_, err = spaceRepo.AddPersonalSpace(ctx, sysOwner, firstOwner)
// 	require.NoError(t, err)

// 	return orgID, sysOwner, firstOwner
// }

func testNewAppUser(t *testing.T, ctx context.Context, ts testService, sysOwner service.SystemOwner, owner service.Owner, loginID, username, password string) service.AppUser {
	appUserRepo := ts.rf.NewAppUserRepository(ctx)
	userID1, err := appUserRepo.AddAppUser(ctx, owner, testNewAppUserAddParameter(t, loginID, username, password))
	require.NoError(t, err)
	user1, err := appUserRepo.FindAppUserByID(ctx, owner, userID1)
	require.NoError(t, err)
	require.Equal(t, loginID, user1.GetLoginID())

	return user1
}

func testNewAppUserAddParameter(t *testing.T, loginID, username, password string) service.AppUserAddParameter {
	p, err := service.NewAppUserAddParameter(loginID, username, password)
	require.NoError(t, err)
	return p
}
