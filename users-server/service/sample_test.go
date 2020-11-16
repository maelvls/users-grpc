// Putting some sample data.

package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadSampleUsers(t *testing.T) {
	db := NewDBOrPanic()
	txn := db.Txn(true)
	assert.NoError(t, LoadSampleUsers(txn))
	txn.Commit()
}
