package config

import (
	"github.com/joeshaw/envdecode"
)

type ConfRedis struct {
	Host     string `env:"REDIS_HOST,required"`
	Port     string `env:"REDIS_PORT,required"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB,default=0"`
}

func (r *ConfRedis) Addr() string {
	return r.Host + ":" + r.Port
}

func NewConfRedis() *ConfRedis {
	var cfg ConfRedis
	if err := envdecode.StrictDecode(&cfg); err != nil {
		panic("Failed to load redis config: " + err.Error())
	}

	return &cfg
}
