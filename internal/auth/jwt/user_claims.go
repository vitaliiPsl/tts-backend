package jwt

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	Id string `json:"id"`
	jwt.StandardClaims
}