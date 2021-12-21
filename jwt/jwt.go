package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/EvisuXiao/andrews-common/constants"
	"github.com/EvisuXiao/andrews-common/utils"
)

type Claims struct {
	constants.UserBrief
	jwt.StandardClaims
	secret []byte
}

func NewJwt(secret string) Claims {
	return Claims{secret: []byte(secret)}
}

func (c *Claims) SetExpired(expired time.Duration) {
	c.ExpiresAt = utils.LocalTime().Add(expired).Unix()
}

func (c Claims) GenerateTokenWithUsername(uid int, username string) (string, error) {
	c.UserBrief.Uid = uid
	c.UserBrief.Username = username
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(c.secret)
}

func (c Claims) GetUserBriefFromToken(tokenString string) (int, string, error) {
	_, err := jwt.ParseWithClaims(tokenString, &c, func(token *jwt.Token) (interface{}, error) {
		return c.secret, nil
	})
	return c.UserBrief.Uid, c.UserBrief.Username, err
}
