package domain

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	libdomain "github.com/kujilabo/redstart/lib/domain"
)

func TestNewSystemOwner(t *testing.T) {
	t.Parallel()
	model, err := libdomain.NewBaseModel(1, time.Now(), time.Now(), 1, 1)
	require.NoError(t, err)
	appUserID, err := NewAppUserID(1)
	require.NoError(t, err)
	organizationID, err := NewOrganizationID(1)
	require.NoError(t, err)
	appUser, err := NewAppUserModel(model, appUserID, organizationID, "LOGIN_ID", "USERNAME")
	assert.NoError(t, err)
	ower, err := NewOwnerModel(appUser)
	assert.NoError(t, err)
	systemOwner, err := NewSystemOwnerModel(ower)
	assert.NoError(t, err)
	log.Println(systemOwner)
}
