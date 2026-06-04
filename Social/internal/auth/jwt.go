package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)


type JWTAuthenticator struct {
	secret string 
	exp time.Duration 
	issue string
	audience string 
}

//this function creates JWTAuthenticator 
func NewJWTAuthenticator(secret string, exp time.Duration, iss string, aud string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret: secret,
		exp: exp,
		issue: iss,
		audience:aud,
	 }
}



//function to generate jwt token 
func (a *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string,error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	//return the jwt token 
	return tokenString,nil 
}

//function to validate jwt token 
func (a *JWTAuthenticator) ValidateToken(token string) (*jwt.Token,error){
	return nil, nil 
}