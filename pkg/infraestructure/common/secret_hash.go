package common

import "golang.org/x/crypto/bcrypt"

// HashSecret applies bcrypt hash on the secret.
func HashSecret(s string) string {
	data, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(data)
}

// MatchHashAndSecret returns true if the secret matches the hash.
func MatchHashAndSecret(hash, secret string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret)) == nil
}
