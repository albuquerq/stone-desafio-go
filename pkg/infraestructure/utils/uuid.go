package utils

import "github.com/google/uuid"

// GenUUID return a new UUID.
func GenUUID() string {
	return uuid.New().String()
}
