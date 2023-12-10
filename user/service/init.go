package service

import "github.com/kujilabo/redstart/lib/domain"

const (
	UserServiceContextKey domain.ContextKey = "user_service"

	SystemAdminLoginID = "__system_admin"
	SystemOwnerLoginID = "__system_owner"

	SystemOwnerGroupKey = "__system_owner"
	OwnerGroupKey       = "__owner"

	SystemOwnerGroupName = "System Owner"
	OwnerGroupName       = "Owner"
)