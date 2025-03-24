package data

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"slender/internal/logger"
	"slender/internal/model"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz"
const (
	letterIdxBits = 5                    // 5 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // all 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // of letter indices fitting in 63 bits
)

const charBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-=[],./!@#$%^&*()_+{}<>?"
const (
	charIdxBits = 7                  // 7 bits to represent a letter index
	charIdxMask = 1<<charIdxBits - 1 // all 1-bits, as many as letterIdxBits
	charIdxMax  = 63 / charIdxBits   // of letter indices fitting in 63 bits
)

// BoundaryGenerator generates a "boundary" string with a specified number of digits.
func BoundaryGenerator(digit uint8) string {
	b := make([]byte, digit)
	l := len(letterBytes)

	for i, cache, remain := int(digit)-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < l {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// SaltGenerator generates a "salt" string with a specified number of digits.
func SaltGenerator(digit uint8) string {
	b := make([]byte, digit)
	l := len(charBytes)

	for i, cache, remain := int(digit)-1, rand.Int63(), charIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), charIdxMax
		}
		if idx := int(cache & charIdxMask); idx < l {
			b[i] = charBytes[idx]
			i--
		}
		cache >>= charIdxBits
		remain--
	}

	return string(b)
}

// Sha256Generator generates a sha256 string of length 64.
func Sha256Generator(password, salt string) string {
	h := sha256.New()
	h.Write([]byte(password + salt))
	return hex.EncodeToString(h.Sum(nil))
}

// ClaimsGenerator generates a JWT claims.
func ClaimsGenerator(username, token string, iat, eat time.Time) model.JWTClaims {
	return model.JWTClaims{
		TokenClaims: model.TokenClaims{
			Token: token,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "slender",
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(iat),
			NotBefore: jwt.NewNumericDate(iat),
			ExpiresAt: jwt.NewNumericDate(eat),
		},
	}
}

// JWTGenerator generates a JWT string.
func JWTGenerator(secret string, claims model.JWTClaims) string {
	key := []byte(secret)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString(key)
	if err != nil {
		logger.Err("error generating JWT string", err)
	}

	return s
}

// ParseJWT parses JWT string.
func ParseJWT(secret string, tokenString string) (model.JWTClaims, error) {
	var resClaims model.JWTClaims
	token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		logger.Debug("error parsing JWT string", "err", err)
		return resClaims, err
	}

	if claims, ok := token.Claims.(*model.JWTClaims); ok {
		resClaims = *claims
	}

	return resClaims, nil
}
