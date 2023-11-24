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

func TestGetOrganization(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService) {
		// logrus.SetLevel(logrus.DebugLevel)
		orgID, _ := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)

		orgRepo := gateway.NewOrganizationRepository(ctx, ts.db)

		// get organization registered
		baseModel, err := libdomain.NewBaseModel(1, time.Now(), time.Now(), 1, 1)
		assert.NoError(t, err)
		appUserID, _ := domain.NewAppUserID(1)
		userModel, err := domain.NewAppUserModel(baseModel, appUserID, orgID, "login_id", "username")
		assert.NoError(t, err)
		{
			org, err := orgRepo.GetOrganization(ctx, userModel)
			assert.NoError(t, err)
			assert.Equal(t, orgNameLength, len(org.GetName()))
		}

		// get organization unregistered
		otherUserModel, err := domain.NewAppUserModel(baseModel, appUserID, invalidOrgID, "login_id", "username")
		assert.NoError(t, err)
		{
			_, err := orgRepo.GetOrganization(ctx, otherUserModel)
			assert.Equal(t, service.ErrOrganizationNotFound, err)
		}
	}
	testDB(t, fn)
}

func TestFindOrganizationByName(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService) {
		// logrus.SetLevel(logrus.DebugLevel)
		orgID, _ := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)
		systemAdminModel := domain.NewSystemAdminModel()

		orgRepo := gateway.NewOrganizationRepository(ctx, ts.db)

		var orgName string

		// get organization registered
		baseModel, err := libdomain.NewBaseModel(1, time.Now(), time.Now(), 1, 1)
		assert.NoError(t, err)
		appUserID, err := domain.NewAppUserID(1)
		require.NoError(t, err)

		userModel, err := domain.NewAppUserModel(baseModel, appUserID, orgID, "login_id", "username")
		assert.NoError(t, err)
		{
			org, err := orgRepo.GetOrganization(ctx, userModel)
			assert.NoError(t, err)
			assert.Equal(t, orgNameLength, len(org.GetName()))
			orgName = org.GetName()
		}

		// find organization registered by name
		{
			org, err := orgRepo.FindOrganizationByName(ctx, systemAdminModel, orgName)
			assert.NoError(t, err)
			assert.Equal(t, orgName, org.GetName())
		}

		// find organization unregistered by name
		{
			_, err := orgRepo.FindOrganizationByName(ctx, systemAdminModel, "NOT_FOUND")
			assert.Equal(t, service.ErrOrganizationNotFound, err)
		}
	}
	testDB(t, fn)
}
