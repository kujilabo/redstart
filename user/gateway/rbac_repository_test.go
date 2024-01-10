package gateway_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/kujilabo/redstart/user/domain"
	"github.com/kujilabo/redstart/user/gateway"
	"github.com/kujilabo/redstart/user/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type test_sdoa struct {
	subject string
	domain  string
	object  string
	action  string
	want    bool
}

func (t *test_sdoa) String() string {
	return fmt.Sprintf("%s,%s,%s,%s,%v", t.subject, t.domain, t.object, t.action, t.want)
}

func initRBACRepository(t *testing.T, db *gorm.DB, conf string) error {
	t.Helper()

	a, err := gormadapter.NewAdapterByDB(db)
	require.NoError(t, err)

	m, err := model.NewModelFromString(conf)
	require.NoError(t, err)

	err = a.SavePolicy(m)
	require.NoError(t, err)

	return nil
}

func initEnforcer(t *testing.T, db *gorm.DB, conf string) *casbin.Enforcer {
	t.Helper()

	a, err := gormadapter.NewAdapterByDB(db)
	require.NoError(t, err)

	m, err := model.NewModelFromString(conf)
	require.NoError(t, err)

	e, err := casbin.NewEnforcer(m, a)
	require.NoError(t, err)

	return e
}

func addPolicy(t *testing.T, rbacRepository service.RBACRepository, dom, sub, act, obj string) {
	t.Helper()
	err := rbacRepository.AddPolicy(domain.NewRBACDomain(dom), domain.NewRBACUser(sub), domain.NewRBACAction(act), domain.NewRBACObject(obj), service.RBACAllowEffect)
	require.NoError(t, err)
}

func TestA(t *testing.T) {
	t.Parallel()

	fn := func(t *testing.T, ctx context.Context, ts testService) {
		defer teardownCasbin(t, ts)
		rbacRepo := gateway.RBACRepository{
			DB:   ts.db,
			Conf: gateway.Conf,
		}

		err := initRBACRepository(t, ts.db, gateway.Conf)
		require.NoError(t, err)
		addPolicy(t, &rbacRepo, "domain1", "alice", "read", "domain:1_data:1")
		// rbacRepo.AddPolicy(domain.NewRBACDomain("domain1"), domain.NewRBACUser("alice"), domain.NewRBACAction("write"), domain.NewRBACObject("data1"), service.RBACAllowEffect)

		// const policy = `
		// p, alice, domain:1_data:1, read, allow, domain1
		// p, bob, domain:2_data:2, write, allow, domain2
		// p, bob, domain:1_data:2, write, allow, domain1
		// p, charlie, domain:1_data*, read, allow, domain1
		// p, domain:1_data2_admin, domain:1_data:2, read, allow, domain1
		// p, domain:1_data2_admin, domain:1_data:2, write, allow, domain1

		// g, alice, domain:1_data2_admin, domain1
		// g2, domain:1_data_child, domain:1_data_parent, domain1
		// g2, domain:2_data_child, domain:2_data_parent, domain2
		// `
		tests := []test_sdoa{
			{subject: "alice", domain: "domain1", object: "domain:1_data:1", action: "read", want: true},
			// {subject: "alice", domain: "domain1", object: "domain:1_data:1", action: "write", want: false},
			// {subject: "alice", domain: "domain1", object: "domain:1_data:2", action: "read", want: true},
			// {subject: "alice", domain: "domain1", object: "domain:1_data:2", action: "write", want: true},

			// {subject: "bob", domain: "domain1", object: "domain:1_data:1", action: "read", want: false},
			// {subject: "bob", domain: "domain1", object: "domain:1_data:1", action: "write", want: false},
			// {subject: "bob", domain: "domain1", object: "domain:1_data:2", action: "read", want: false},
			// {subject: "bob", domain: "domain1", object: "domain:1_data:2", action: "write", want: true},

			// {subject: "charlie", domain: "domain1", object: "domain:1_data:2", action: "read", want: true},
			// {subject: "charlie", domain: "domain1", object: "domain:1_data_parent", action: "read", want: true},
		}
		for _, tt := range tests {
			t.Run(tt.String(), func(t *testing.T) {
				e := initEnforcer(t, ts.db, gateway.Conf)
				ok, err := e.Enforce(tt.subject, tt.object, tt.action, tt.domain)
				require.NoError(t, err)
				assert.Equal(t, tt.want, ok)
			})
		}
	}
	testDB(t, fn)
}

func teardownCasbin(t *testing.T, ts testService) {
	// delete all organizations
	// ts.db.Exec("delete from space where organization_id = ?", orgID.Int())
	ts.db.Exec("delete from casbin_rule")
	// db.Where("true").Delete(&spaceEntity{})
	// db.Where("true").Delete(&appUserEntity{})
	// db.Where("true").Delete(&organizationEntity{})
}

// func Test_rbac_model(t *testing.T) {
// 	const conf = `
// 	[request_definition]
// 	r = sub, obj, act

// 	[policy_definition]
// 	p = sub, obj, act

// 	[role_definition]
// 	g = _, _

// 	[policy_effect]
// 	e = some(where (p.eft == allow))

// 	[matchers]
// 	m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
// 	`

// 	const policy = `
// 	p, alice, data1, read
// 	p, bob, data2, write
// 	p, data2_admin, data2, read
// 	p, data2_admin, data2, write

// 	g, alice, data2_admin
// 	`

// 	e := NewEnforcer(t, conf, policy)
// 	tests := []test_soa{
// 		{subject: "alice", object: "data1", action: "read", want: true},
// 		{subject: "alice", object: "data1", action: "write", want: false},
// 		{subject: "alice", object: "data2", action: "read", want: true},
// 		{subject: "alice", object: "data2", action: "write", want: true},

// 		{subject: "bob", object: "data1", action: "read", want: false},
// 		{subject: "bob", object: "data1", action: "write", want: false},
// 		{subject: "bob", object: "data2", action: "read", want: false},
// 		{subject: "bob", object: "data2", action: "write", want: true},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.String(), func(t *testing.T) {
// 			ok, err := e.Enforce(tt.subject, tt.object, tt.action)
// 			require.NoError(t, err)
// 			assert.Equal(t, tt.want, ok)
// 		})
// 	}
// }
