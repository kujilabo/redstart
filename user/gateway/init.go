package gateway

import (
	"go.opentelemetry.io/otel"
)

var (
	tracer = otel.Tracer("github.com/kujilabo/redstart/user/gateway")

	AppUserTableName = "app_user"

	// SystemStudentLoginID = "system-student"
	// GuestLoginID         = "guest"

	// AdministratorRole = "Administrator"
	// ManagerRole       = "Manager"
	// UserRole          = "User"
	// GuestRole         = "Guest"
	// UnknownRole       = "Unknown"
)
