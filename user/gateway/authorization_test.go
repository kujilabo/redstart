package gateway_test

import (
	"context"
	"testing"

	"github.com/kujilabo/redstart/user/gateway"
	"github.com/kujilabo/redstart/user/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AddPairOfUserAndGroup(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService) {
		orgID, _, owner := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)

		// given
		user1 := testAddAppUser(t, ctx, ts, owner, "LOGIN_ID_1", "USERNAME_1", "PASSWORD_1")
		user2 := testAddAppUser(t, ctx, ts, owner, "LOGIN_ID_2", "USERNAME_2", "PASSWORD_2")

		authorizationManager := gateway.NewAuthorizationManager(ctx, ts.db, ts.rf)
		userGroupRepo := gateway.NewUserGroupRepository(ctx, ts.db)
		ownerGroup, err := userGroupRepo.FindUserGroupByKey(ctx, owner, service.OwnerGroupKey)
		require.NoError(t, err)

		rbacRoleObject := service.NewRBACUserRoleObject(orgID, ownerGroup.GetUerGroupID())

		// when
		ok, err := authorizationManager.Authorize(ctx, owner, service.RBACSetAction, rbacRoleObject)
		assert.NoError(t, err)
		// then
		// - owner can make anyone belong to owner-group
		assert.True(t, ok)

		// when
		ok, err = authorizationManager.Authorize(ctx, user2, service.RBACSetAction, rbacRoleObject)
		assert.NoError(t, err)
		// then
		// - user2 can NOT make user3 belong to owner-group
		assert.False(t, ok)

		// given
		// - add user1 to owner-group
		err = authorizationManager.AddUserToGroup(ctx, owner, user1.GetAppUserID(), ownerGroup.GetUerGroupID())
		require.NoError(t, err)
		// when
		ok, err = authorizationManager.Authorize(ctx, user1, service.RBACSetAction, rbacRoleObject)
		assert.NoError(t, err)
		// then
		// - user1 can make sure user3 belong to group1 because user1 belongs to owner-group
		assert.True(t, ok)
	}
	testDB(t, fn)
}