package users

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
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
	const method = "NewAccessToken"

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifetime)),
		},
		IP: ip,
	})

	signedJWT, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", method, err)
	}
	access := AccessToken(signedJWT)

	return &access, nil
}

func (t AccessToken) String() string {
	return string(t)
}

func (t AccessToken) IP(secretKey string) (net.IP, error) {
	const method = "AccessToken.IP"

	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(
		t.String(),
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
		jwt.WithValidMethods([]string{"HS512"}),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", method, err)
	}

	if !token.Valid {
		return nil, ErrInvalidAccessToken
	}

	return claims.IP, nil
}

type RefreshToken string

func NewRefreshToken() (*RefreshToken, error) {
	const method = "NewRefreshToken"

	randomBytes := make([]byte, 128)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", method, err)
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
