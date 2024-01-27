package gateway_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kujilabo/redstart/user/domain"
	"github.com/kujilabo/redstart/user/gateway"
	"github.com/kujilabo/redstart/user/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func outputOrganization(t *testing.T, db *gorm.DB) {
	var results []gateway.OrganizationEntity
	if result := db.Find(&results); result.Error != nil {
		assert.Fail(t, result.Error.Error())
	}
	var s string
	s += "\n   id,version,           created_at,          updated_at,created_by,updated_by,      name,"
	for i := range results {
		result := &results[i]
		s += fmt.Sprintf("\n%5d,%8d,%20s,%20s,%10d,%10d,%10s", result.ID, result.Version, result.CreatedAt.Format(time.RFC3339), result.UpdatedAt.Format(time.RFC3339), result.CreatedBy, result.UpdatedBy, result.Name)
	}
	t.Log(s)
}

func Test_appUserRepository_FindSystemOwnerByOrganizationID_shouldReturnSystemOwner_whenExistingOrganizationIDIsSpecified(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService, orgID *domain.OrganizationID, sysOwner *service.SystemOwner, owner *service.Owner) {
		sysAdModel := domain.NewSystemAdminModel()
		sysAd := testNewSystemAdmin(sysAdModel)
		appUserRepo := gateway.NewAppUserRepository(ctx, ts.dialect, ts.db, ts.rf)

		// when
		testSysOwner, err := appUserRepo.FindSystemOwnerByOrganizationID(ctx, sysAd, orgID)

		// then
		require.NoError(t, err)
		assert.Equal(t, service.SystemOwnerLoginID, testSysOwner.LoginID())
	}
	testOrganization(t, fn)
}

func Test_appUserRepository_FindSystemOwnerByOrganizationID_shouldReturnError_whenInvalidOrganizationIDIsSpecified(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService, orgID *domain.OrganizationID, sysOwner *service.SystemOwner, owner *service.Owner) {
		sysAdModel := domain.NewSystemAdminModel()
		sysAd := testNewSystemAdmin(sysAdModel)
		appUserRepo := gateway.NewAppUserRepository(ctx, ts.dialect, ts.db, ts.rf)

		// when
		_, err := appUserRepo.FindSystemOwnerByOrganizationID(ctx, sysAd, invalidOrgID)

		// then
		assert.ErrorIs(t, err, service.ErrSystemOwnerNotFound)
	}
	testOrganization(t, fn)
}

func Test_appUserRepository_FindSystemOwnerByOrganizationName_shouldReturnSystemOwner_whenExistingNameIsSpecified(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService, orgID *domain.OrganizationID, sysOwner *service.SystemOwner, owner *service.Owner) {
		org := getOrganization(t, ctx, ts, orgID)
		sysAdModel := domain.NewSystemAdminModel()
		sysAd := testNewSystemAdmin(sysAdModel)

		appUserRepo := gateway.NewAppUserRepository(ctx, ts.dialect, ts.db, ts.rf)

		// when
		testSysOwner, err := appUserRepo.FindSystemOwnerByOrganizationName(ctx, sysAd, org.Name())

		// then
		require.NoError(t, err)
		assert.Equal(t, service.SystemOwnerLoginID, testSysOwner.LoginID())
	}
	testOrganization(t, fn)
}

func Test_appUserRepository_FindSystemOwnerByOrganizationName_shouldReturnError_whenInvalidNameIsSpecified(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService, orgID *domain.OrganizationID, sysOwner *service.SystemOwner, owner *service.Owner) {
		sysAdModel := domain.NewSystemAdminModel()
		sysAd := testNewSystemAdmin(sysAdModel)

		appUserRepo := gateway.NewAppUserRepository(ctx, ts.dialect, ts.db, ts.rf)

		// when
		_, err := appUserRepo.FindSystemOwnerByOrganizationName(ctx, sysAd, "NOT_FOUND")

		// then
		assert.ErrorIs(t, err, service.ErrSystemOwnerNotFound)
	}
	testOrganization(t, fn)
}

func Test_appUserRepository_FindAppUserByID_shouldReturnAppUser_whenExistingIDIsSpecified(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService, orgID *domain.OrganizationID, sysOwner *service.SystemOwner, owner *service.Owner) {
		appUserRepo := gateway.NewAppUserRepository(ctx, ts.dialect, ts.db, ts.rf)

		// given
		newAppUser := testAddAppUser(t, ctx, ts, owner, "LOGIN_ID", "USERNAME", "PASSWORD")

		// when
		appUser, err := appUserRepo.FindAppUserByID(ctx, owner, newAppUser.AppUserID())

		// then
		require.NoError(t, err)
		assert.Equal(t, "LOGIN_ID", appUser.LoginID(), "loginID should be 'LOGIN_ID'")
		assert.Equal(t, "USERNAME", appUser.Username(), "username should be 'USERNAME'")
	}
	testOrganization(t, fn)
}

func Test_appUserRepository_FindAppUserByID_shouldReturnError_whenInvaildIDIsSpecified(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService, orgID *domain.OrganizationID, sysOwner *service.SystemOwner, owner *service.Owner) {
		appUserRepo := gateway.NewAppUserRepository(ctx, ts.dialect, ts.db, ts.rf)

		_, err := appUserRepo.FindAppUserByID(ctx, owner, invalidAppUserID)
		assert.ErrorIs(t, err, service.ErrAppUserNotFound)
	}
	testOrganization(t, fn)
}

func Test_appUserRepository_FindAppUserByLoginID_shouldReturnAppUser_whenExistingLoginIDIsSpecified(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService, orgID *domain.OrganizationID, sysOwner *service.SystemOwner, owner *service.Owner) {
		appUserRepo := gateway.NewAppUserRepository(ctx, ts.dialect, ts.db, ts.rf)

		// given
		_ = testAddAppUser(t, ctx, ts, owner, "LOGIN_ID", "USERNAME", "PASSWORD")

		// when
		appUser, err := appUserRepo.FindAppUserByLoginID(ctx, owner, "LOGIN_ID")

		// then
		require.NoError(t, err)
		assert.Equal(t, "LOGIN_ID", appUser.LoginID(), "loginID should be 'LOGIN_ID'")
		assert.Equal(t, "USERNAME", appUser.Username(), "username should be 'USERNAME'")
	}
	testOrganization(t, fn)
}

func Test_appUserRepository_FindAppUserByLoginID_shouldReturnError_whenInvalidLoginIDIsSpecified(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService, orgID *domain.OrganizationID, sysOwner *service.SystemOwner, owner *service.Owner) {
		appUserRepo := gateway.NewAppUserRepository(ctx, ts.dialect, ts.db, ts.rf)

		// when
		_, err := appUserRepo.FindAppUserByLoginID(ctx, owner, "NOT_FOUND")

		// then
		assert.ErrorIs(t, err, service.ErrAppUserNotFound)
	}
	testOrganization(t, fn)
}

func Test_appUserRepository_FindOwnerByLoginID_shouldReturnOwner_whenExistingOwnerLoginIDIsSpecified(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService, orgID *domain.OrganizationID, sysOwner *service.SystemOwner, owner *service.Owner) {
		require.Equal(t, "OWNER_ID", owner.LoginID())
		require.Equal(t, "OWNER_NAME", owner.Username())

		appUserRepo := gateway.NewAppUserRepository(ctx, ts.dialect, ts.db, ts.rf)

		// when
		appUser, err := appUserRepo.FindOwnerByLoginID(ctx, sysOwner, owner.LoginID())

		// then
		require.NoError(t, err)
		assert.Equal(t, "OWNER_ID", appUser.LoginID())
		assert.Equal(t, "OWNER_NAME", appUser.Username())
	}
	testOrganization(t, fn)
}

func Test_appUserRepository_FindOwnerByLoginID_shouldReturnError_whenNotOwnerLoginIDIsSpecified(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService, orgID *domain.OrganizationID, sysOwner *service.SystemOwner, owner *service.Owner) {
		require.Equal(t, "OWNER_ID", owner.LoginID())
		require.Equal(t, "OWNER_NAME", owner.Username())

		appUserRepo := gateway.NewAppUserRepository(ctx, ts.dialect, ts.db, ts.rf)

		// given
		_ = testAddAppUser(t, ctx, ts, owner, "LOGIN_ID", "USERNAME", "PASSWORD")

		// when
		_, err := appUserRepo.FindOwnerByLoginID(ctx, sysOwner, "LOGIN_ID")

		// then
		assert.ErrorIs(t, err, service.ErrAppUserNotFound)
	}
	testOrganization(t, fn)
}

func Test_appUserRepository_VerifyPassword_shouldReturnTrue_whenCorrectPasswordIsSpecified(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService, orgID *domain.OrganizationID, sysOwner *service.SystemOwner, owner *service.Owner) {
		sysAdModel := domain.NewSystemAdminModel()
		sysAd := testNewSystemAdmin(sysAdModel)
		appUserRepo := gateway.NewAppUserRepository(ctx, ts.dialect, ts.db, ts.rf)

		// given
		_ = testAddAppUser(t, ctx, ts, owner, "LOGIN_ID", "USERNAME", "PASSWORD")

		// when
		verified, err := appUserRepo.VerifyPassword(ctx, sysAd, orgID, "LOGIN_ID", "PASSWORD")

		// then
		assert.True(t, verified)
		assert.NoError(t, err)
	}
	testOrganization(t, fn)
}

func Test_appUserRepository_VerifyPassword_shouldReturnFalse_whenWrongPasswordIsSpecified(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService, orgID *domain.OrganizationID, sysOwner *service.SystemOwner, owner *service.Owner) {
		sysAdModel := domain.NewSystemAdminModel()
		sysAd := testNewSystemAdmin(sysAdModel)
		appUserRepo := gateway.NewAppUserRepository(ctx, ts.dialect, ts.db, ts.rf)

		// given
		_ = testAddAppUser(t, ctx, ts, owner, "LOGIN_ID", "USERNAME", "PASSWORD")

		// when
		verified, err := appUserRepo.VerifyPassword(ctx, sysAd, orgID, "LOGIN_ID", "WRONG_PASSWORD")

		// then
		assert.False(t, verified)
		assert.NoError(t, err)
	}
	testOrganization(t, fn)
}
