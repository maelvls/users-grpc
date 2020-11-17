// Putting some sample data.

package service_test

import (
	"testing"

	service "github.com/maelvls/users-grpc/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestLoadSampleUsers(t *testing.T) {
	// Just a quick wiring test.
	db := service.NewDBOrPanic()
	txn := db.Txn(true)
	assert.NoError(t, service.LoadSampleUsers(txn))
	txn.Commit()
}
