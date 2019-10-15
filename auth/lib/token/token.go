package token

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// TODO: THIS IS HILARIOUSLY INSECURE, make this a secret somewhere
var signKey = []byte("super sekrit key")

// TODO: make this config
const tokenDuration = time.Hour * 10

type Claim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func New(username string) (string, error) {
	now := time.Now()

	claims := Claim{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(tokenDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString(signKey)
}

func Parse(token string) (*Claim, error) {
	claims := &Claim{}

	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return signKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(*Claim); ok && t.Valid {
		return claims, nil
	}

	return nil, err
}
