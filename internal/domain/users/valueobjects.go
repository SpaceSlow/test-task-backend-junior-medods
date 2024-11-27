package users

import (
	"crypto/rand"
	"encoding/base64"
	"net"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserClaims struct {
	jwt.RegisteredClaims
	IP net.IP
}

type AccessToken string

func NewAccessToken(ip net.IP, tokenLifetime time.Duration, secretKey string) (*AccessToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifetime)),
		},
		IP: ip,
	})

	signedJWT, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}
	access := AccessToken(signedJWT)

	return &access, nil
}

func (t AccessToken) String() string {
	return string(t)
}

type RefreshToken string

func NewRefreshToken() (*RefreshToken, error) {
	randomBytes := make([]byte, 128)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	refresh := RefreshToken(base64.StdEncoding.EncodeToString(randomBytes))

	return &refresh, nil
}

func (t RefreshToken) Hash() ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(t.String()), bcrypt.DefaultCost)
}

func (t RefreshToken) String() string {
	return string(t)
}
