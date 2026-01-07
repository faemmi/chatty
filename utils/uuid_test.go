package utils

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCommutativeUUIDHash(t *testing.T) {
	id1 := uuid.Must(uuid.NewUUID())
	id2 := uuid.Must(uuid.NewUUID())

	hash1 := CommutativeUUIDHash(id1, id2)
	hash2 := CommutativeUUIDHash(id2, id1)

	assert.Equal(t, hash1, hash2, "Hashes neq")
}
