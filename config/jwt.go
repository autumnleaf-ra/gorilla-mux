package config

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("8dsqa461a35cxzc023sad2513a5d")

type JWTClaim struct {
	Username string
	jwt.RegisteredClaims
}
