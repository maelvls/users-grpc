// Putting some sample data.

package usersvc_test

import (
	"testing"

	usersvc "github.com/maelvls/users-grpc/users-server/usersvc"
	"github.com/stretchr/testify/assert"
)

func TestLoadSampleUsers(t *testing.T) {
	// Just a quick wiring test.
	db := usersvc.NewDBOrPanic()
	txn := db.Txn(true)
	assert.NoError(t, usersvc.LoadSampleUsers(txn))
	txn.Commit()
}
