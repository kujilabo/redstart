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

func setupOrganization(ctx context.Context, t *testing.T, ts testService) (*domain.OrganizationID, *service.SystemOwner, *service.Owner) {
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
	authorizationManager := gateway.NewAuthorizationManager(ctx, ts.db, ts.rf)

	// 1. add organization
	orgID, err := orgRepo.AddOrganization(bg, sysAd, orgAddParam)
	require.NoError(t, err)
	assert.Greater(t, orgID.Int(), 0)

	// 2. add "system-owner" user
	sysOwnerID, err := appUserRepo.AddSystemOwner(bg, sysAd, orgID)
	require.NoError(t, err)
	require.Greater(t, sysOwnerID.Int(), 0)

	// TODO
	sysOwner, err := appUserRepo.FindSystemOwnerByOrganizationName(ctx, sysAd, orgName, service.IncludeGroups)
	require.NoError(t, err)

	// 3. add policy to "system-owner" userct(orgID)
	rbacSysOwner := service.NewRBACAppUser(orgID, sysOwnerID)
	rbacAllUserRolesObject := service.NewRBACAllUserRolesObject(orgID)
	// - "system-owner" "can" "set" "all-user-roles"
	err = authorizationManager.AddPolicyToUserBySystemAdmin(ctx, sysAd, orgID, rbacSysOwner, service.RBACSetAction, rbacAllUserRolesObject, service.RBACAllowEffect)
	require.NoError(t, err)

	// - "system-owner" "can" "unset" "all-user-roles"
	err = authorizationManager.AddPolicyToUserBySystemAdmin(ctx, sysAd, orgID, rbacSysOwner, service.RBACUnsetAction, rbacAllUserRolesObject, service.RBACAllowEffect)
	require.NoError(t, err)

	// 4. add owner-group
	ownerGroupID, err := userGorupRepo.AddOwnerGroup(ctx, sysOwner, orgID)
	require.NoError(t, err)

	// 5. add policty to "owner" group
	rbacOwnerGroup := service.NewRBACUserRole(orgID, ownerGroupID)
	// - "owner" group "can" "set" "all-user-roles"
	err = authorizationManager.AddPolicyToGroupBySystemAdmin(ctx, sysAd, orgID, rbacOwnerGroup, service.RBACSetAction, rbacAllUserRolesObject, service.RBACAllowEffect)
	require.NoError(t, err)
	// - "owner" group "can" "unset" "all-user-roles"
	err = authorizationManager.AddPolicyToGroupBySystemAdmin(ctx, sysAd, orgID, rbacOwnerGroup, service.RBACUnsetAction, rbacAllUserRolesObject, service.RBACAllowEffect)
	require.NoError(t, err)

	// 6. add first owner
	ownerID, err := appUserRepo.AddAppUser(ctx, sysOwner, firstOwnerAddParam)
	require.NoError(t, err)
	require.Greater(t, ownerID.Int(), 0)

	// - owner belongs to owner-group
	err = authorizationManager.AddUserToGroup(ctx, sysOwner, ownerID, ownerGroupID)
	require.NoError(t, err)

	owner, err := appUserRepo.FindOwnerByLoginID(ctx, sysOwner, firstOwnerAddParam.GetLoginID())
	require.NoError(t, err)

	return orgID, sysOwner, owner
}

func teardownOrganization(t *testing.T, ts testService, orgID *domain.OrganizationID) {
	// delete all organizations
	// ts.db.Exec("delete from space where organization_id = ?", orgID.Int())
	ts.db.Exec("delete from app_user where organization_id = ?", orgID.Int())
	ts.db.Exec("delete from organization where id = ?", orgID.Int())
	// db.Where("true").Delete(&spaceEntity{})
	// db.Where("true").Delete(&appUserEntity{})
	// db.Where("true").Delete(&organizationEntity{})
}

func testAddAppUser(t *testing.T, ctx context.Context, ts testService, owner service.OwnerModelInterface, loginID, username, password string) *service.AppUser {
	appUserRepo := ts.rf.NewAppUserRepository(ctx)
	userID1, err := appUserRepo.AddAppUser(ctx, owner, testNewAppUserAddParameter(t, loginID, username, password))
	require.NoError(t, err)
	user1, err := appUserRepo.FindAppUserByID(ctx, owner, userID1)
	require.NoError(t, err)
	require.Equal(t, loginID, user1.LoginID())

	return user1
}

func testAddUserGroup(t *testing.T, ctx context.Context, ts testService, owner service.OwnerModelInterface, key, name, description string) *service.UserGroup {
	userGorupRepo := ts.rf.NewUserGroupRepository(ctx)
	groupID1, err := userGorupRepo.AddUserGroup(ctx, owner, testNewUserGroupAddParameter(t, key, name, description))
	require.NoError(t, err)
	group1, err := userGorupRepo.FindUserGroupByID(ctx, owner, groupID1)
	require.NoError(t, err)
	require.Equal(t, key, group1.Key())
	require.Equal(t, name, group1.Name())
	require.Equal(t, description, group1.Description())

	return group1
}

type testSystemAdmin struct {
	*domain.SystemAdminModel
}

func (m *testSystemAdmin) AppUserID() *domain.AppUserID {
	return m.SystemAdminModel.AppUserID
}
func (m *testSystemAdmin) IsSystemAdmin() bool {
	return true
}
func testNewSystemAdmin(systemAdminModel *domain.SystemAdminModel) *testSystemAdmin {
	return &testSystemAdmin{
		systemAdminModel,
	}
}

type testAppUserModel struct {
	*domain.AppUserModel
}

func (m *testAppUserModel) AppUserID() *domain.AppUserID {
	return m.AppUserModel.AppUserID
}
func (m *testAppUserModel) OrganizationID() *domain.OrganizationID {
	return m.AppUserModel.OrganizationID
}
func (m *testAppUserModel) LoginID() string {
	return m.AppUserModel.LoginID
}
func (m *testAppUserModel) Username() string {
	return m.AppUserModel.Username
}
func testNewAppUser(appUserModel *domain.AppUserModel) *testAppUserModel {
	return &testAppUserModel{
		appUserModel,
	}
}

type testUserGroupModel struct {
	*domain.UserGroupModel
}

func (m *testUserGroupModel) Key() string {
	return m.UserGroupModel.Key
}
func (m *testUserGroupModel) Name() string {
	return m.UserGroupModel.Key
}
func (m *testUserGroupModel) Description() string {
	return m.UserGroupModel.Description
}
func testNewUserGroup(userGroupModel *domain.UserGroupModel) *testUserGroupModel {
	return &testUserGroupModel{
		userGroupModel,
	}
}
func testNewUserGroups(userGroupModels []*domain.UserGroupModel) []*testUserGroupModel {
	groups := make([]*testUserGroupModel, len(userGroupModels))
	for i, groupModel := range userGroupModels {
		groups[i] = testNewUserGroup(groupModel)
	}
	return groups
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

func getOrganization(t *testing.T, ctx context.Context, ts testService, orgID *domain.OrganizationID) *service.Organization {
	orgRepo := gateway.NewOrganizationRepository(ctx, ts.db)

	baseModel, err := libdomain.NewBaseModel(1, time.Now(), time.Now(), 1, 1)
	require.NoError(t, err)
	appUserID, _ := domain.NewAppUserID(1)
	appUserModel, err := domain.NewAppUserModel(baseModel, appUserID, orgID, "login_id", "username", nil)
	require.NoError(t, err)
	appUser, err := service.NewAppUser(ctx, ts.rf, appUserModel)
	require.NoError(t, err)

	org, err := orgRepo.GetOrganization(ctx, appUser)
	require.NoError(t, err)
	require.Equal(t, orgNameLength, len(org.Name()))

	return org
}
