package gateway_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kujilabo/redstart/user/gateway"
	"github.com/kujilabo/redstart/user/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func outputCasbinRule(t *testing.T, db *gorm.DB) {
	type Result struct {
		ID    int
		Ptype string
		V0    string
		V1    string
		V2    string
		V3    string
		V4    string
		V5    string
	}
	var results []Result
	if result := db.Raw("SELECT * FROM casbin_rule").Scan(&results); result.Error != nil {
		assert.Fail(t, result.Error.Error())
	}
	var s string
	s += "\n   id,ptype,                  v0,                  v1,         v2,         v3,         v4,         v5"
	for i := range results {
		result := &results[i]
		s += fmt.Sprintf("\n%5d,%5s,%20s,%20s, %10s, %10s, %10s, %10s", result.ID, result.Ptype, result.V0, result.V1, result.V2, result.V3, result.V4, result.V5)
	}
	t.Log(s)
}

func Test_AddPairOfUserAndGroup(t *testing.T) {
	t.Parallel()
	for i := 0; i < 1; i++ {
		fn := func(t *testing.T, ctx context.Context, ts testService) {
			orgID, _, owner := setupOrganization(ctx, t, ts)
			defer teardownOrganization(t, ts, orgID)

			outputCasbinRule(t, ts.db)
			assert.True(t, false)

			// given
			user1 := testAddAppUser(t, ctx, ts, owner, "LOGIN_ID_1", "USERNAME_1", "PASSWORD_1")
			user2 := testAddAppUser(t, ctx, ts, owner, "LOGIN_ID_2", "USERNAME_2", "PASSWORD_2")

			authorizationManager := gateway.NewAuthorizationManager(ctx, ts.dialect, ts.db, ts.rf)
			userGroupRepo := gateway.NewUserGroupRepository(ctx, ts.dialect, ts.db)
			ownerGroup, err := userGroupRepo.FindUserGroupByKey(ctx, owner, service.OwnerGroupKey)
			require.NoError(t, err)

			rbacRoleObject := service.NewRBACUserRoleObject(orgID, ownerGroup.UserGroupID())

			// when
			ok, err := authorizationManager.Authorize(ctx, owner, service.RBACSetAction, rbacRoleObject)
			assert.NoError(t, err)
			// then
			assert.True(t, ok, "owner should be able to make anyone belong to owner-group")
			if !ok {
				outputCasbinRule(t, ts.db)
			}

			// when
			ok, err = authorizationManager.Authorize(ctx, user2, service.RBACSetAction, rbacRoleObject)
			assert.NoError(t, err)
			// then
			assert.False(t, ok, "standard-user should not be able to make other users belong to owner-group")

			// given
			// - add user1 to owner-group
			err = authorizationManager.AddUserToGroup(ctx, owner, user1.AppUserID(), ownerGroup.UserGroupID())
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
}
