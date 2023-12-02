package gateway

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/kujilabo/redstart/user/domain"
)

type HasTableName interface {
	TableName() string
}

type wrappedDB struct {
	db             *gorm.DB
	organizationID domain.OrganizationID
}

func (x *wrappedDB) Table(name string, args ...interface{}) *wrappedDB {
	x.db = x.db.Table(name, args...)
	return x
}

func (x *wrappedDB) Select(query interface{}, args ...interface{}) *wrappedDB {
	x.db = x.db.Select(query, args...)
	return x
}

func (x *wrappedDB) Where(query interface{}, args ...interface{}) *wrappedDB {
	x.db = x.db.Where(query, args...)
	return x
}

func (x *wrappedDB) Joins(query string, args ...interface{}) *wrappedDB {
	x.db = x.db.Joins(query, args...)
	return x
}

func (x *wrappedDB) WhereOrganizationID(table HasTableName, organizationID domain.OrganizationID) *wrappedDB {
	x.db = x.db.Where(fmt.Sprintf("%s.`organization_id` = ?", table.TableName()), organizationID.Int())
	return x
}

func (x *wrappedDB) WhereNotRemoved(table HasTableName) *wrappedDB {
	x.db = x.db.Where(fmt.Sprintf("%s.`removed` = 0", table.TableName()))
	return x
}

func (x *wrappedDB) WhereAppUser() *wrappedDB {
	return x.WhereOrganizationID(&appUserEntity{}, x.organizationID).WhereNotRemoved(&appUserEntity{})
}

func (x *wrappedDB) WhereUserGroup() *wrappedDB {
	return x.WhereOrganizationID(&userGroupEntity{}, x.organizationID).WhereNotRemoved(&userGroupEntity{})
}

func (x *wrappedDB) WherePairOfUserAndGroup() *wrappedDB {
	return x.WhereOrganizationID(&pairOfUserAndGroupEntity{}, x.organizationID)
}
