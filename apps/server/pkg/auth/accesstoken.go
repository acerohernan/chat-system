package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	defaultValidDuration = 6 * time.Hour
)

type AccessToken struct {
	issuer string
	secret string
	grants *Grants
}

type Grants struct {
	Email          string `json:"email"`
	CanSendMessage bool   `json:"canSendMessage"`
}

func NewAccessToken(issuer string, secret string, grants *Grants) *AccessToken {
	return &AccessToken{
		issuer: issuer,
		secret: secret,
		grants: grants,
	}
}

func (t *AccessToken) toJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":    t.issuer,
		"exp":    time.Now().Add(defaultValidDuration),
		"grants": t.grants,
	})

	jwt, err := token.SignedString([]byte(t.secret))

	if err != nil {
		return "", err
	}

	return jwt, nil
}
