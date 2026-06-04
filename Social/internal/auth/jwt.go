package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secret   string
	exp      time.Duration
	issue    string
	audience string
}

// this function creates JWTAuthenticator
func NewJWTAuthenticator(secret string, exp time.Duration, iss string, aud string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret:   secret,
		exp:      exp,
		issue:    iss,
		audience: aud,
	}
}

// function to generate jwt token
func (a *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	//return the jwt token
	return tokenString, nil
}

// function to validate jwt token
func (a *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any,error){
		//Check whether the algorithm is HMAC(HS256)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(a.secret), nil 
	},
	jwt.WithExpirationRequired(),// token MUST have an exp field — reject if missing
	jwt.WithAudience(a.audience),
	jwt.WithIssuer(a.issue),
	jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}), //only accept HS256 — double safety on top of the manual check
)
}
