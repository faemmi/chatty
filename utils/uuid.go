package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"github.com/google/uuid"
)

func IsValidUUID(value string) bool {
    _, err := uuid.Parse(value)
    return err == nil
}


// CommutativeUUIDHash calculates a commutative hash for two UUIDs.
func CommutativeUUIDHashFromString(s1, s2 string) string {
	// Sort the strings to ensure order doesn't matter
	uuids := []string{s1, s2}
	sort.Strings(uuids)

	// Concatenate the sorted UUIDs
	data := []byte(uuids[0] + uuids[1])

	// Calculate the SHA-256 hash
	hash := sha256.Sum256(data)

	// Return the hash as a hex string
	return hex.EncodeToString(hash[:])

}

// CommutativeUUIDHash calculates a commutative hash for two UUIDs.
func CommutativeUUIDHash(u1, u2 uuid.UUID) string {
	return CommutativeUUIDHashFromString(u1.String(), u2.String())
}

