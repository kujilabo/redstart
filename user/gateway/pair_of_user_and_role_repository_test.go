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

func Test_pairOfUserAndRoleRepository_FindUserRolesByUserID(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService) {
		orgID, _, owner := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)

		// given
		user1 := testAddAppUser(t, ctx, ts, owner, "LOGIN_ID_1", "USERNAME_1", "PASSWORD_1")
		user2 := testAddAppUser(t, ctx, ts, owner, "LOGIN_ID_2", "USERNAME_2", "PASSWORD_2")
		role1 := testAddUserRole(t, ctx, ts, owner, "ROLE_KEY_1", "ROLE_NAME_1", "ROLE_DESC_1")
		role2 := testAddUserRole(t, ctx, ts, owner, "ROLE_KEY_2", "ROLE_NAME_2", "ROLE_DESC_2")
		role3 := testAddUserRole(t, ctx, ts, owner, "ROLE_KEY_3", "ROLE_NAME_3", "ROLE_DESC_3")

		pairOfUserAndRoleRepo := gateway.NewPairOfUserAndRoleRepository(ctx, ts.db, ts.rf)

		// - user1 has role1, role2, role3
		for _, role := range []service.UserRole{role1, role2, role3} {
			err := pairOfUserAndRoleRepo.AddPairOfUserAndRole(ctx, owner, user1.GetAppUserID(), role.GetUerRoleID())
			require.NoError(t, err)
		}
		// - user2 has role1
		for _, role := range []service.UserRole{role1} {
			err := pairOfUserAndRoleRepo.AddPairOfUserAndRole(ctx, owner, user2.GetAppUserID(), role.GetUerRoleID())
			require.NoError(t, err)
		}

		// when
		roles1, err := pairOfUserAndRoleRepo.FindUserRolesByUserID(ctx, owner, user1.GetAppUserID())
		require.NoError(t, err)
		roles2, err := pairOfUserAndRoleRepo.FindUserRolesByUserID(ctx, owner, user2.GetAppUserID())
		require.NoError(t, err)

		// then
		// - user1 has role1, role2, role3
		assert.Len(t, roles1, 3)
		assert.Equal(t, "ROLE_KEY_1", roles1[0].GetKey())
		assert.Equal(t, "ROLE_KEY_2", roles1[1].GetKey())
		assert.Equal(t, "ROLE_KEY_3", roles1[2].GetKey())
		// - user2 has role1
		assert.Len(t, roles2, 1)
		assert.Equal(t, "ROLE_KEY_1", roles2[0].GetKey())
	}
	testDB(t, fn)
}

func Test_pairOfUserAndRoleRepository_AddPairOfUserAndRole(t *testing.T) {
	t.Parallel()
	fn := func(t *testing.T, ctx context.Context, ts testService) {
		orgID, _, owner := setupOrganization(ctx, t, ts)
		defer teardownOrganization(t, ts, orgID)

		// given
		user1 := testAddAppUser(t, ctx, ts, owner, "LOGIN_ID_1", "USERNAME_1", "PASSWORD_1")
		user2 := testAddAppUser(t, ctx, ts, owner, "LOGIN_ID_2", "USERNAME_2", "PASSWORD_2")
		user3 := testAddAppUser(t, ctx, ts, owner, "LOGIN_ID_3", "USERNAME_3", "PASSWORD_2")
		role1 := testAddUserRole(t, ctx, ts, owner, "ROLE_KEY_1", "ROLE_NAME_1", "ROLE_DESC_1")

		pairOfUserAndRoleRepo := gateway.NewPairOfUserAndRoleRepository(ctx, ts.db, ts.rf)
		userRoleRepo := gateway.NewUserRoleRepository(ctx, ts.db)
		ownerRole, err := userRoleRepo.FindUserRoleByKey(ctx, owner, gateway.OwnerRoleKey)
		require.NoError(t, err)

		// when
		err = pairOfUserAndRoleRepo.AddPairOfUserAndRole(ctx, owner, user1.GetAppUserID(), ownerRole.GetUerRoleID())
		require.NoError(t, err)
		// then
		// - owner can set owner-role to user1
		{
			userRoles1, err := pairOfUserAndRoleRepo.FindUserRolesByUserID(ctx, user1, user1.GetAppUserID())
			require.NoError(t, err)
			assert.Len(t, userRoles1, 1)
		}

		// when
		err = pairOfUserAndRoleRepo.AddPairOfUserAndRole(ctx, user2, user3.GetAppUserID(), ownerRole.GetUerRoleID())
		// then
		// - user2 can NOT set owner-role to user3
		assert.ErrorIs(t, err, libdomain.ErrPermissionDenied)

		// when
		err = pairOfUserAndRoleRepo.AddPairOfUserAndRole(ctx, user1, user3.GetAppUserID(), role1.GetUerRoleID())
		// - user1 can set role1 to user3 because user1 has owner-role
		assert.NoError(t, err)
		{
			userRoles3, err := pairOfUserAndRoleRepo.FindUserRolesByUserID(ctx, user1, user3.GetAppUserID())
			require.NoError(t, err)
			assert.Len(t, userRoles3, 1)
		}

	}
	testDB(t, fn)
}
