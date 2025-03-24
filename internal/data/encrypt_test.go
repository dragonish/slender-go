package data

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBoundaryGenerator(t *testing.T) {
	res1 := BoundaryGenerator(16)
	res2 := BoundaryGenerator(16)

	assert.Len(t, res1, 16)
	assert.Len(t, res2, 16)
	assert.NotEqual(t, res1, res2)
}

func TestSaltGenerator(t *testing.T) {
	res1 := SaltGenerator(160)
	res2 := SaltGenerator(160)

	assert.Len(t, res1, 160)
	assert.Len(t, res2, 160)
	assert.NotEqual(t, res1, res2)
}

func TestSha256Generator(t *testing.T) {
	password := "abcdef"

	res1 := Sha256Generator(password, SaltGenerator(6))
	res2 := Sha256Generator(password, SaltGenerator(6))
	res3 := Sha256Generator(password, "123456")
	res4 := Sha256Generator(password, "123456")

	assert.Len(t, res1, 64)
	assert.Len(t, res2, 64)
	assert.NotEqual(t, res1, res2)
	assert.Len(t, res3, 64)
	assert.Len(t, res4, 64)
	assert.Equal(t, res3, res4)
}

func TestJWTGenerator(t *testing.T) {
	now := time.Now()
	expired := now.AddDate(0, 1, 0)
	unexpected := now.AddDate(0, -1, 0)
	user := "user"
	token := "token"
	secret := "secret"

	claims := ClaimsGenerator(user, token, now, expired)

	res1 := JWTGenerator(secret, claims)
	res2 := JWTGenerator(secret, claims)
	res3 := JWTGenerator("my_secret", claims)

	uClaims := ClaimsGenerator(user, token, now, unexpected)
	res4 := JWTGenerator(secret, uClaims)

	assert.NotEmpty(t, res1)
	assert.NotEmpty(t, res2)
	assert.Equal(t, res1, res2)
	assert.NotEmpty(t, res3)
	assert.NotEqual(t, res1, res3)
	assert.NotEmpty(t, res4)
}

func TestParseJWT(t *testing.T) {
	user := "user"
	token := "token"
	secret := "secret"
	now := time.Now()

	claims1 := ClaimsGenerator(user, token, now, now.AddDate(0, 1, 0))
	jwt1 := JWTGenerator(secret, claims1)
	res1, _ := ParseJWT(secret, jwt1)

	iat2 := now.AddDate(0, -1, 0)
	expired2 := now.AddDate(0, 0, -15)
	claims2 := ClaimsGenerator(user, token, iat2, expired2)
	jwt2 := JWTGenerator(secret, claims2)
	res2, _ := ParseJWT(secret, jwt2)

	iat3 := now.AddDate(0, 0, 1)
	expired3 := now.AddDate(0, 1, 0)
	claims3 := ClaimsGenerator(user, token, iat3, expired3)
	jwt3 := JWTGenerator(secret, claims3)
	res3, _ := ParseJWT(secret, jwt3)

	res4, _ := ParseJWT("my_secret", jwt1)
	res5, _ := ParseJWT(secret, "abcdef")

	assert.Equal(t, token, res1.Token)
	assert.Equal(t, user, res1.Subject)
	assert.Empty(t, res2.Token)
	assert.Empty(t, res3.Token)
	assert.Empty(t, res4.Token)
	assert.Empty(t, res5.Token)
}
