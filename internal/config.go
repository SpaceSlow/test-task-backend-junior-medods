package internal

import (
	"sync"
	"time"
)

type ServerConfig struct {
	MaxTimeoutShutdown time.Duration
}

var defaultConfig = &ServerConfig{
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
