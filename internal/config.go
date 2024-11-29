package internal

import (
	"os"
	"sync"
	"time"
)

type ServerConfig struct {
	secretKey          string
	DSN                string
	smtpAddress        string
	smtpSender         string
	smtpPassword       string
	tokenLifetime      time.Duration
	MaxTimeoutShutdown time.Duration
}

func (c *ServerConfig) SecretKey() string {
	return c.secretKey
}

func (c *ServerConfig) TokenLifetime() time.Duration {
	return c.tokenLifetime
}

func (c *ServerConfig) SMTPAddress() string {
	return c.smtpAddress
}

func (c *ServerConfig) SMTPSender() string {
	return c.smtpSender
}

func (c *ServerConfig) SMTPPassword() string {
	return c.smtpPassword
}

var defaultConfig = &ServerConfig{
	secretKey:          os.Getenv("SECRET_KEY"),
	DSN:                os.Getenv("DSN"),
	smtpAddress:        os.Getenv("SMTP_ADDRESS"),
	smtpSender:         os.Getenv("SMTP_SENDER"),
	smtpPassword:       os.Getenv("SMTP_PASSWORD"),
	tokenLifetime:      time.Hour,
	MaxTimeoutShutdown: 5 * time.Second,
}

var serverConfig *ServerConfig = nil
var once sync.Once

func LoadServerConfig() *ServerConfig {
	once.Do(func() {
		serverConfig = defaultConfig
	})
	return serverConfig
}
