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

	libdomain "github.com/kujilabo/redstart/lib/domain"
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

func setupOrganization(ctx context.Context, t *testing.T, ts testService) (domain.OrganizationID, service.SystemOwner, service.Owner) {
	bg := context.Background()
	orgName := RandString(orgNameLength)
	sysAd, err := service.NewSystemAdmin(ctx, ts.rf)
	require.NoError(t, err)

	firstOwnerAddParam, err := service.NewFirstOwnerAddParameter("OWNER_ID", "OWNER_NAME", "OWNER_PASSWORD")
	require.NoError(t, err)
	orgAddParam, err := service.NewOrganizationAddParameter(orgName, firstOwnerAddParam)
	require.NoError(t, err)

	orgRepo := gateway.NewOrganizationRepository(ctx, ts.db)

	appUserRepo := gateway.NewAppUserRepository(ctx, ts.driverName, ts.db, ts.rf)
	userGorupRepo := gateway.NewUserGroupRepository(ctx, ts.db)
	pairOfUserAndRole := gateway.NewPairOfUserAndGroupRepository(ctx, ts.db, ts.rf)
	rbacRepo := gateway.NewRBACRepository(ctx, ts.db)

	// add organization
	orgID, err := orgRepo.AddOrganization(bg, sysAd, orgAddParam)
	require.NoError(t, err)
	assert.Greater(t, orgID.Int(), 0)

	rbacDomain := service.NewRBACOrganization(orgID)

	// add system-owner-group
	sysOwnerGroupID, err := userGorupRepo.AddSystemOwnerGroup(ctx, sysAd, orgID)
	require.NoError(t, err)

	// add owner-group
	ownerGroupID, err := userGorupRepo.AddOwnerGroup(ctx, sysAd, orgID)
	require.NoError(t, err)

	// add system-owner
	sysOwnerID, err := appUserRepo.AddSystemOwner(bg, sysAd, orgID)
	require.NoError(t, err)
	require.Greater(t, sysOwnerID.Int(), 0)

	sysOwner, err := appUserRepo.FindSystemOwnerByOrganizationName(bg, sysAd, orgName, service.IncludeGroups)
	require.NoError(t, err)

	rbacAllUserRolesObject := service.NewRBACAllUserRolesObject(orgID)

	rbacSysOwnerRole := service.NewRBACUserRole(orgID, sysOwnerGroupID)

	// those who belong to system-owner-group can set "all user roles"
	err = rbacRepo.AddPolicy(rbacDomain, rbacSysOwnerRole, service.RBACSetAction, rbacAllUserRolesObject, service.RBACAllowEffect)
	require.NoError(t, err)

	// those who belong to system-owner-group can unset "all user roles"
	err = rbacRepo.AddPolicy(rbacDomain, rbacSysOwnerRole, service.RBACUnsetAction, rbacAllUserRolesObject, service.RBACAllowEffect)
	require.NoError(t, err)

	// systen-owner belongs to system-owner-group
	err = pairOfUserAndRole.AddPairOfUserAndGroupBySystemAdmin(ctx, sysAd, orgID, sysOwnerID, sysOwnerGroupID)
	require.NoError(t, err)

	// add owner
	ownerID, err := appUserRepo.AddAppUser(ctx, sysOwner, firstOwnerAddParam)
	require.NoError(t, err)
	require.Greater(t, ownerID.Int(), 0)

	rbacOwnerRole := service.NewRBACUserRole(orgID, ownerGroupID)

	// those who belong to owner-group can set all roles
	err = rbacRepo.AddPolicy(rbacDomain, rbacOwnerRole, service.RBACSetAction, rbacAllUserRolesObject, service.RBACAllowEffect)
	require.NoError(t, err)

	// those who belong to owner-group can set all roles
	err = rbacRepo.AddPolicy(rbacDomain, rbacOwnerRole, service.RBACUnsetAction, rbacAllUserRolesObject, service.RBACAllowEffect)
	require.NoError(t, err)

	// owner belongs to owner-group
	// err = pairOfUserAndRole.AddPairOfUserAndGroup(ctx, sysOwner, ownerID, ownerGroupID)

	require.NoError(t, err)

	owner, err := appUserRepo.FindOwnerByLoginID(ctx, sysOwner, firstOwnerAddParam.GetLoginID())
	require.NoError(t, err)

	return orgID, sysOwner, owner
}

func teardownOrganization(t *testing.T, ts testService, orgID domain.OrganizationID) {
	// delete all organizations
	// ts.db.Exec("delete from space where organization_id = ?", orgID.Int())
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

func testAddAppUser(t *testing.T, ctx context.Context, ts testService, owner domain.OwnerModel, loginID, username, password string) service.AppUser {
	appUserRepo := ts.rf.NewAppUserRepository(ctx)
	userID1, err := appUserRepo.AddAppUser(ctx, owner, testNewAppUserAddParameter(t, loginID, username, password))
	require.NoError(t, err)
	user1, err := appUserRepo.FindAppUserByID(ctx, owner, userID1)
	require.NoError(t, err)
	require.Equal(t, loginID, user1.GetLoginID())

	return user1
}

func testAddUserGroup(t *testing.T, ctx context.Context, ts testService, owner domain.OwnerModel, key, name, description string) service.UserGroup {
	userGorupRepo := ts.rf.NewUserGroupRepository(ctx)
	groupID1, err := userGorupRepo.AddUserGroup(ctx, owner, testNewUserGroupAddParameter(t, key, name, description))
	require.NoError(t, err)
	group1, err := userGorupRepo.FindUserGroupByID(ctx, owner, groupID1)
	require.NoError(t, err)
	require.Equal(t, key, group1.GetKey())
	require.Equal(t, name, group1.GetName())
	require.Equal(t, description, group1.GetDescription())

	return group1
}

func testNewAppUserAddParameter(t *testing.T, loginID, username, password string) service.AppUserAddParameter {
	p, err := service.NewAppUserAddParameter(loginID, username, password)
	require.NoError(t, err)
	return p
}

func testNewUserGroupAddParameter(t *testing.T, key, name, description string) service.UserGroupAddParameter {
	p, err := service.NewUserGroupAddParameter(key, name, description)
	require.NoError(t, err)
	return p
}

func getOrganization(t *testing.T, ctx context.Context, ts testService, orgID domain.OrganizationID) service.Organization {
	orgRepo := gateway.NewOrganizationRepository(ctx, ts.db)

	baseModel, err := libdomain.NewBaseModel(1, time.Now(), time.Now(), 1, 1)
	require.NoError(t, err)
	appUserID, _ := domain.NewAppUserID(1)
	userModel, err := domain.NewAppUserModel(baseModel, appUserID, orgID, "login_id", "username", nil)
	require.NoError(t, err)

	org, err := orgRepo.GetOrganization(ctx, userModel)
	require.NoError(t, err)
	require.Equal(t, orgNameLength, len(org.GetName()))

	return org
}
