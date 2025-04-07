package config

import (
	"time"

	"github.com/joeshaw/envdecode"
)

type ConfServer struct {
	Port         int           `env:"SERVER_PORT,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	Debug        bool          `env:"SERVER_DEBUG,required"`
}

func NewConfServer() *ConfServer {
	var cfg ConfServer
	if err := envdecode.StrictDecode(&cfg); err != nil {
		panic("Failed to load server config: " + err.Error())
	}
	return &cfg
}
