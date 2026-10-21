package auth

import "github.com/golang-jwt/jwt/v5"

type Autheticator interface {
	GenerateToken(claims jwt.Claims) (string, error)
	ValidateToken(toke string) (*jwt.Token, error)
}
