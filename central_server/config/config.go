package config

import (
	"fmt"
	"os"
)

type Config struct {
	Ip   string
	Port string
}

func LoadConfig() *Config {
	ip := os.Getenv("IP_ADDRESS")
	if ip == "" {
		ip = "0.0.0.0"
	}
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "42069"
	}

	return &Config{
		Ip: ip,
		Port: port,
	}
}

func (c *Config) ToAddress() string {
	return fmt.Sprintf("%s:%s", c.Ip, c.Port)
}