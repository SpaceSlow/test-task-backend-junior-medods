package users

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type accessTokenClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	IP    net.IP `json:"ip"`
}

type AccessToken struct {
	jwt   string
	email string
	ip    net.IP
}

func NewAccessToken(jwt string) *AccessToken {
	return &AccessToken{jwt: jwt}
}

func GenerateAccessToken(email string, ip net.IP, tokenLifetime time.Duration, secretKey string) (*AccessToken, error) {
	const method = "NewAccessToken"

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, accessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifetime)),
		},
		Email: email,
		IP:    ip,
	})

	signedJWT, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", method, err)
	}
	access := AccessToken{jwt: signedJWT, email: email, ip: ip}

	return &access, nil
}

func (t AccessToken) JWT() string {
	return t.jwt
}

func (t AccessToken) Email() string {
	return t.email
}

func (t AccessToken) IP() net.IP {
	return t.ip
}

func (t *AccessToken) Parse(secretKey string) error {
	claims := &accessTokenClaims{}
	token, err := jwt.ParseWithClaims(
		t.JWT(),
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
		jwt.WithValidMethods([]string{"HS512"}),
	)
	if err != nil || !token.Valid {
		return ErrInvalidAccessToken
	}

	t.email = claims.Email
	t.ip = claims.IP
	return nil
}

type RefreshToken []byte

func NewRefreshToken() (*RefreshToken, error) {
	const method = "NewRefreshToken"

	randomBytes := make([]byte, 72)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", method, err)
	}
	refresh := RefreshToken(randomBytes)

	return &refresh, nil
}

func ParseRefreshToken(refreshBase64 string) (*RefreshToken, error) {
	token, err := base64.StdEncoding.DecodeString(refreshBase64)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}
	refresh := RefreshToken(token)

	return &refresh, nil
}

func (t RefreshToken) GenerateHash() ([]byte, error) {
	return bcrypt.GenerateFromPassword(t, bcrypt.DefaultCost)
}

func (t RefreshToken) Valid(hash string) bool {
	return errors.Is(nil, bcrypt.CompareHashAndPassword([]byte(hash), t))
}

func (t RefreshToken) String() string {
	return base64.StdEncoding.EncodeToString(t)
}
