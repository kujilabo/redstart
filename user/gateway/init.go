package gateway

import (
	"go.opentelemetry.io/otel"
)

var (
	tracer = otel.Tracer("github.com/kujilabo/redstart/user/gateway")

	AppUserTableName = "app_user"

	SystemAdminLoginID = "__system_admin"
	SystemOwnerLoginID = "__system_owner"

	SystemOwnerRoleKey = "__system_owner"
	OwnerRoleKey       = "__owner"

	SystemOwnerRoleName = "System Owner"
	OwnerRoleName       = "Owner"
	// SystemStudentLoginID = "system-student"
	// GuestLoginID         = "guest"

	// AdministratorRole = "Administrator"
	// ManagerRole       = "Manager"
	// UserRole          = "User"
	// GuestRole         = "Guest"
	// UnknownRole       = "Unknown"
)
