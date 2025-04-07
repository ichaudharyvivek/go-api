package config

import (
	"fmt"

	"github.com/joeshaw/envdecode"
)

type ConfDB struct {
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT,required"`
	Username string `env:"DB_USER,required"`
	Password string `env:"DB_PASS,required"`
	DBName   string `env:"DB_NAME,required"`
	Debug    bool   `env:"DB_DEBUG,required"`
	SSLMode  string `env:"DB_SSLMODE,required"`
}

func NewConfDB() *ConfDB {
	var cfg ConfDB
	if err := envdecode.StrictDecode(&cfg); err != nil {
		panic("Failed to load DB config: " + err.Error())
	}
	return &cfg
}

func (c *ConfDB) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.DBName, c.SSLMode)
}
