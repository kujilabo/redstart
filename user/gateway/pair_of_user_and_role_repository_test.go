package gateway_test

import (
	"context"
	"testing"

	libdomain "github.com/kujilabo/redstart/lib/domain"
	"github.com/kujilabo/redstart/user/gateway"
	"github.com/kujilabo/redstart/user/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_pairOfUserAndGroupRepository_FindUserGroupsByUserID(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService) {
		orgID, _, owner := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)

		// given
		user1 := testAddAppUser(t, ctx, ts, owner, "LOGIN_ID_1", "USERNAME_1", "PASSWORD_1")
		user2 := testAddAppUser(t, ctx, ts, owner, "LOGIN_ID_2", "USERNAME_2", "PASSWORD_2")
		group1 := testAddUserGroup(t, ctx, ts, owner, "GROUP_KEY_1", "GROUP_NAME_1", "GROUP_DESC_1")
		group2 := testAddUserGroup(t, ctx, ts, owner, "GROUP_KEY_2", "GROUP_NAME_2", "GROUP_DESC_2")
		group3 := testAddUserGroup(t, ctx, ts, owner, "GROUP_KEY_3", "GROUP_NAME_3", "GROUP_DESC_3")

		pairOfUserAndGroupRepo := gateway.NewPairOfUserAndGroupRepository(ctx, ts.db, ts.rf)

		// - user1 has group1, group2, group3
		for _, group := range []service.UserGroup{group1, group2, group3} {
			err := pairOfUserAndGroupRepo.AddPairOfUserAndGroup(ctx, owner, user1.GetAppUserID(), group.GetUerGroupID())
			require.NoError(t, err)
		}
		// - user2 has group1
		for _, group := range []service.UserGroup{group1} {
			err := pairOfUserAndGroupRepo.AddPairOfUserAndGroup(ctx, owner, user2.GetAppUserID(), group.GetUerGroupID())
			require.NoError(t, err)
		}

		// when
		groups1, err := pairOfUserAndGroupRepo.FindUserGroupsByUserID(ctx, owner, user1.GetAppUserID())
		require.NoError(t, err)
		groups2, err := pairOfUserAndGroupRepo.FindUserGroupsByUserID(ctx, owner, user2.GetAppUserID())
		require.NoError(t, err)

		// then
		// - user1 has group1, group2, group3
		assert.Len(t, groups1, 3)
		assert.Equal(t, "GROUP_KEY_1", groups1[0].GetKey())
		assert.Equal(t, "GROUP_KEY_2", groups1[1].GetKey())
		assert.Equal(t, "GROUP_KEY_3", groups1[2].GetKey())
		// - user2 has group1
		assert.Len(t, groups2, 1)
		assert.Equal(t, "GROUP_KEY_1", groups2[0].GetKey())
	}
	testDB(t, fn)
}

func Test_pairOfUserAndGroupRepository_AddPairOfUserAndGroup(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService) {
		orgID, _, owner := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)

		// given
		user1 := testAddAppUser(t, ctx, ts, owner, "LOGIN_ID_1", "USERNAME_1", "PASSWORD_1")
		user2 := testAddAppUser(t, ctx, ts, owner, "LOGIN_ID_2", "USERNAME_2", "PASSWORD_2")
		user3 := testAddAppUser(t, ctx, ts, owner, "LOGIN_ID_3", "USERNAME_3", "PASSWORD_2")
		group1 := testAddUserGroup(t, ctx, ts, owner, "GROUP_KEY_1", "GROUP_NAME_1", "GROUP_DESC_1")

		pairOfUserAndGroupRepo := gateway.NewPairOfUserAndGroupRepository(ctx, ts.db, ts.rf)
		userGroupRepo := gateway.NewUserGroupRepository(ctx, ts.db)
		ownerGroup, err := userGroupRepo.FindUserGroupByKey(ctx, owner, gateway.OwnerGroupKey)
		require.NoError(t, err)

		// when
		err = pairOfUserAndGroupRepo.AddPairOfUserAndGroup(ctx, owner, user1.GetAppUserID(), ownerGroup.GetUerGroupID())
		require.NoError(t, err)
		// then
		// - owner can set owner-role to user1
		{
			userGroups1, err := pairOfUserAndGroupRepo.FindUserGroupsByUserID(ctx, user1, user1.GetAppUserID())
			require.NoError(t, err)
			assert.Len(t, userGroups1, 1)
		}

		// when
		err = pairOfUserAndGroupRepo.AddPairOfUserAndGroup(ctx, user2, user3.GetAppUserID(), ownerGroup.GetUerGroupID())
		// then
		// - user2 can NOT set owner-role to user3
		assert.ErrorIs(t, err, libdomain.ErrPermissionDenied)

		// when
		err = pairOfUserAndGroupRepo.AddPairOfUserAndGroup(ctx, user1, user3.GetAppUserID(), group1.GetUerGroupID())
		// - user1 can set group1 to user3 because user1 has owner-role
		assert.NoError(t, err)
		{
			userGroups3, err := pairOfUserAndGroupRepo.FindUserGroupsByUserID(ctx, user1, user3.GetAppUserID())
			require.NoError(t, err)
			assert.Len(t, userGroups3, 1)
		}

	}
	testDB(t, fn)
}
