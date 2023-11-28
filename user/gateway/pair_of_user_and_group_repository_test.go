package gateway_test

// func Test_pairOfUserAndGroupRepository_AddPairOfUserAndGroup(t *testing.T) {
// 	t.Parallel()
// 	fn := func(t *testing.T, ctx context.Context, ts testService) {
// 		orgID, _, owner := setupOrganization(ctx, t, ts)
// 		defer teardownOrganization(t, ts, orgID)

// 		user1 := testAddAppUser(t, ctx, ts, owner, "LOGIN_ID_1", "USERNAME_1", "PASSWORD_1")
// 		user2 := testAddAppUser(t, ctx, ts, owner, "LOGIN_ID_2", "USERNAME_2", "PASSWORD_2")
// 		role1 := testAddUserRole(t, ctx, ts, owner, "ROLE_KEY_1", "ROLE_NAME_1", "ROLE_DESC_1")
// 		role2 := testAddUserRole(t, ctx, ts, owner, "ROLE_KEY_2", "ROLE_NAME_2", "ROLE_DESC_2")
// 		role3 := testAddUserRole(t, ctx, ts, owner, "ROLE_KEY_3", "ROLE_NAME_3", "ROLE_DESC_3")

// 		pairOfUserAndRoleRepo := gateway.NewPairOfUserAndRoleRepository(ctx, ts.db)

// 		// user1 has role1, role2, role3
// 		for _, role := range []service.UserRole{role1, role2, role3} {
// 			err := pairOfUserAndRoleRepo.AddPairOfUserAndRole(ctx, owner, user1.GetAppUserID(), role.GetUerRoleID())
// 			require.NoError(t, err)
// 		}
// 		// user2 has role1
// 		for _, role := range []service.UserRole{role1} {
// 			err := pairOfUserAndRoleRepo.AddPairOfUserAndRole(ctx, owner, user2.GetAppUserID(), role.GetUerRoleID())
// 			require.NoError(t, err)
// 		}

// 		roles1, err := pairOfUserAndRoleRepo.FindUserRolesByUserID(ctx, owner, user1.GetAppUserID())
// 		require.NoError(t, err)
// 		assert.Len(t, roles1, 3)
// 	}
// 	testDB(t, fn)
// }
