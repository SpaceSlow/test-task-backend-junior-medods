package internal

import (
	"os"
	"sync"
	"time"
)

type ServerConfig struct {
	secretKey          string
	DSN                string
	tokenLifetime      time.Duration
	MaxTimeoutShutdown time.Duration
}

func (c *ServerConfig) SecretKey() string {
	return c.secretKey
}

func (c *ServerConfig) TokenLifetime() time.Duration {
	return c.tokenLifetime
}

var defaultConfig = &ServerConfig{
	secretKey:          os.Getenv("SECRET_KEY"),
	DSN:                os.Getenv("DSN"),
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
