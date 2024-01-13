package gateway_test

import (
	"context"
	"testing"

	"github.com/kujilabo/redstart/user/domain"
	"github.com/kujilabo/redstart/user/gateway"
	"github.com/kujilabo/redstart/user/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_appUserRepository_FindSystemOwnerByOrganizationID(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService) {
		orgID, _, _ := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)

		sysAdModel := domain.NewSystemAdminModel()
		sysAd := testNewSystemAdmin(sysAdModel)

		appUserRepo := gateway.NewAppUserRepository(ctx, ts.dialect, ts.db, ts.rf)

		{
			sysOwner, err := appUserRepo.FindSystemOwnerByOrganizationID(ctx, sysAd, orgID)
			require.NoError(t, err)
			assert.Equal(t, service.SystemOwnerLoginID, sysOwner.LoginID())
		}

		{
			_, err := appUserRepo.FindSystemOwnerByOrganizationID(ctx, sysAd, invalidOrgID)
			assert.ErrorIs(t, err, service.ErrSystemOwnerNotFound)
		}
	}
	testDB(t, fn)
}

func Test_appUserRepository_FindSystemOwnerByOrganizationName(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService) {
		orgID, _, _ := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)

		org := getOrganization(t, ctx, ts, orgID)
		sysAdModel := domain.NewSystemAdminModel()
		sysAd := testNewSystemAdmin(sysAdModel)

		appUserRepo := gateway.NewAppUserRepository(ctx, ts.dialect, ts.db, ts.rf)

		{
			sysOwner, err := appUserRepo.FindSystemOwnerByOrganizationName(ctx, sysAd, org.Name())
			require.NoError(t, err)
			assert.Equal(t, service.SystemOwnerLoginID, sysOwner.LoginID())
		}

		{
			_, err := appUserRepo.FindSystemOwnerByOrganizationName(ctx, sysAd, "NOT_FOUND")
			assert.ErrorIs(t, err, service.ErrSystemOwnerNotFound)
		}
	}
	testDB(t, fn)
}

func Test_appUserRepository_FindAppUserByID(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService) {
		orgID, _, owner := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)

		appUserAddParam, err := service.NewAppUserAddParameter("LOGIN_ID", "USERNAME", "PASSWORD", "", "", "", "")
		require.NoError(t, err)

		appUserRepo := gateway.NewAppUserRepository(ctx, ts.dialect, ts.db, ts.rf)

		appUserID, err := appUserRepo.AddAppUser(ctx, owner, appUserAddParam)
		require.NoError(t, err)
		require.Greater(t, appUserID.Int(), 0)

		{
			appUser, err := appUserRepo.FindAppUserByID(ctx, owner, appUserID)
			require.NoError(t, err)
			assert.Equal(t, appUserAddParam.LoginID(), appUser.LoginID())
			assert.Equal(t, appUserAddParam.Username(), appUser.Username())
		}

		{
			_, err := appUserRepo.FindAppUserByID(ctx, owner, invalidAppUserID)
			assert.ErrorIs(t, err, service.ErrAppUserNotFound)
		}
	}
	testDB(t, fn)
}

func Test_appUserRepository_FindAppUserByLoginID(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService) {
		orgID, _, owner := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)

		appUserAddParam, err := service.NewAppUserAddParameter("LOGIN_ID", "USERNAME", "PASSWORD", "", "", "", "")
		require.NoError(t, err)

		appUserRepo := gateway.NewAppUserRepository(ctx, ts.dialect, ts.db, ts.rf)

		appUserID, err := appUserRepo.AddAppUser(ctx, owner, appUserAddParam)
		require.NoError(t, err)
		require.Greater(t, appUserID.Int(), 0)

		{
			appUser, err := appUserRepo.FindAppUserByLoginID(ctx, owner, appUserAddParam.LoginID())
			require.NoError(t, err)
			assert.Equal(t, appUserAddParam.LoginID(), appUser.LoginID())
			assert.Equal(t, appUserAddParam.Username(), appUser.Username())
		}

		{
			_, err := appUserRepo.FindAppUserByLoginID(ctx, owner, "NOT_FOUND")
			assert.ErrorIs(t, err, service.ErrAppUserNotFound)
		}
	}
	testDB(t, fn)
}

func Test_appUserRepository_FindOwnerByLoginID(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService) {
		orgID, sysOwner, owner := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)

		appUserAddParam, err := service.NewAppUserAddParameter("LOGIN_ID", "USERNAME", "PASSWORD", "", "", "", "")
		require.NoError(t, err)

		appUserRepo := gateway.NewAppUserRepository(ctx, ts.dialect, ts.db, ts.rf)

		appUserID, err := appUserRepo.AddAppUser(ctx, owner, appUserAddParam)
		require.NoError(t, err)
		require.Greater(t, appUserID.Int(), 0)

		{
			appUser, err := appUserRepo.FindOwnerByLoginID(ctx, sysOwner, owner.LoginID())
			require.NoError(t, err)
			assert.Equal(t, owner.LoginID(), appUser.LoginID())
			assert.Equal(t, owner.Username(), appUser.Username())
		}

		{
			_, err := appUserRepo.FindOwnerByLoginID(ctx, sysOwner, appUserAddParam.LoginID())
			assert.ErrorIs(t, err, service.ErrAppUserNotFound)
		}
	}
	testDB(t, fn)
}
