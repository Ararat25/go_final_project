package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	PasswordChecksum string `json:"checksum"`
	jwt.RegisteredClaims
}
