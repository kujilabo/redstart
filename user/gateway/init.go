package gateway

import (
	"go.opentelemetry.io/otel"

	libdomain "github.com/kujilabo/redstart/lib/domain"
)

var (
	UserGatewayContextKey libdomain.ContextKey = "user_gateway"

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
