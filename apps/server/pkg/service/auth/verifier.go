package auth

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/chat-system/server/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

const (
	defaultValidDuration = 6 * time.Hour
)

type Verifier struct {
	issuer string
	secret string
}

func NewVerifier(conf *config.Config) *Verifier {
	return &Verifier{
		issuer: conf.Auth.JWTIssuer,
		secret: conf.Auth.JWTSecret,
	}
}

func (v *Verifier) CreateToken(grants *Grants) (*AccessToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":    v.issuer,
		"grants": grants,
		"exp":    json.Number(strconv.FormatInt(time.Now().Add(defaultValidDuration).Unix(), 10)),
	})

	jwtToken, err := token.SignedString([]byte(v.secret))

	if err != nil {
		return nil, err
	}

	return &AccessToken{
		Grants: grants,
		Jwt:    jwtToken,
	}, nil
}

func (v *Verifier) ParseToken(rawJWT string) (*AccessToken, error) {
	token, err := jwt.Parse(rawJWT, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return "", nil
		}

		return []byte(v.secret), nil
	}, jwt.WithIssuer(v.issuer))

	if err != nil {
		return nil, ErrInvalidAccessToken
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	// parse grants
	if grants, ok := claims["grants"].(map[string]any); ok {
		return &AccessToken{
			Grants: &Grants{
				Id:             grants["id"].(string),
				Email:          grants["email"].(string),
				CanSendMessage: grants["canSendMessage"].(bool),
			},
			Jwt: rawJWT,
		}, nil
	} else {
		return nil, ErrInvalidAccessToken
	}

}
