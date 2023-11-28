package gateway_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	libdomain "github.com/kujilabo/redstart/lib/domain"
	"github.com/kujilabo/redstart/user/domain"
	"github.com/kujilabo/redstart/user/gateway"
	"github.com/kujilabo/redstart/user/service"
)

func Test_organizationRepository_GetOrganization(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService) {
		orgID, _, _ := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)

		orgRepo := gateway.NewOrganizationRepository(ctx, ts.db)

		// get organization registered
		baseModel, err := libdomain.NewBaseModel(1, time.Now(), time.Now(), 1, 1)
		assert.NoError(t, err)
		appUserID, _ := domain.NewAppUserID(1)
		userModel, err := domain.NewAppUserModel(baseModel, appUserID, orgID, "login_id", "username", nil)
		assert.NoError(t, err)
		{
			org, err := orgRepo.GetOrganization(ctx, userModel)
			assert.NoError(t, err)
			assert.Equal(t, orgNameLength, len(org.GetName()))
		}

		// get organization unregistered
		otherUserModel, err := domain.NewAppUserModel(baseModel, appUserID, invalidOrgID, "login_id", "username", nil)
		assert.NoError(t, err)
		{
			_, err := orgRepo.GetOrganization(ctx, otherUserModel)
			assert.ErrorIs(t, err, service.ErrOrganizationNotFound)
		}
	}
	testDB(t, fn)
}

func Test_organizationRepository_FindOrganizationByName(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService) {
		orgID, _, _ := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)
		sysAdModel := domain.NewSystemAdminModel()

		orgRepo := gateway.NewOrganizationRepository(ctx, ts.db)

		var orgName string

		// get organization registered
		baseModel, err := libdomain.NewBaseModel(1, time.Now(), time.Now(), 1, 1)
		assert.NoError(t, err)
		appUserID, err := domain.NewAppUserID(1)
		require.NoError(t, err)

		userModel, err := domain.NewAppUserModel(baseModel, appUserID, orgID, "login_id", "username", nil)
		assert.NoError(t, err)
		{
			org, err := orgRepo.GetOrganization(ctx, userModel)
			assert.NoError(t, err)
			assert.Equal(t, orgNameLength, len(org.GetName()))
			orgName = org.GetName()
		}

		// find organization registered by name
		{
			org, err := orgRepo.FindOrganizationByName(ctx, sysAdModel, orgName)
			assert.NoError(t, err)
			assert.Equal(t, orgName, org.GetName())
		}

		// find organization unregistered by name
		{
			_, err := orgRepo.FindOrganizationByName(ctx, sysAdModel, "NOT_FOUND")
			assert.Equal(t, service.ErrOrganizationNotFound, err)
		}
	}
	testDB(t, fn)
}

func Test_organizationRepository_FindOrganizationByID(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService) {
		orgID, _, _ := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)
		sysAdModel := domain.NewSystemAdminModel()

		orgRepo := gateway.NewOrganizationRepository(ctx, ts.db)

		// get organization registered
		baseModel, err := libdomain.NewBaseModel(1, time.Now(), time.Now(), 1, 1)
		assert.NoError(t, err)
		appUserID, err := domain.NewAppUserID(1)
		require.NoError(t, err)

		userModel, err := domain.NewAppUserModel(baseModel, appUserID, orgID, "login_id", "username", nil)
		assert.NoError(t, err)
		{
			org, err := orgRepo.GetOrganization(ctx, userModel)
			assert.NoError(t, err)
			assert.Equal(t, orgNameLength, len(org.GetName()))
		}

		// find organization registered by ID
		{
			org, err := orgRepo.FindOrganizationByID(ctx, sysAdModel, orgID)
			assert.NoError(t, err)
			assert.Equal(t, orgID.Int(), org.GetID().Int())
		}

		// find organization unregistered by ID
		{
			_, err := orgRepo.FindOrganizationByID(ctx, sysAdModel, invalidOrgID)
			assert.Equal(t, service.ErrOrganizationNotFound, err)
		}
	}
	testDB(t, fn)
}
