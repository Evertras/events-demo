package token

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var SignKey []byte

// TODO: make this config
const tokenDuration = time.Hour * 10

type Claim struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

func New(userID string) (string, error) {
	now := time.Now()

	claims := Claim{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(tokenDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString(SignKey)
}

func Parse(token string) (*Claim, error) {
	claims := &Claim{}

	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return SignKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(*Claim); ok && t.Valid {
		return claims, nil
	}

	return nil, err
}
