package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/vukovlevi/netstore/store_administration/db"
)

type Config struct {
    dbConfig db.DBConfig
    Ip string
    Port string
}

func CreateApplicationConfig() Config {
    ip := os.Getenv("IP_ADDRESS")
    if ip == "" {
        ip = "0.0.0.0"
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }
    return Config{
        dbConfig: db.CreateDatabaseConfig(),
        Ip: ip,
        Port: port,
    }
}

func (c *Config) Apply() error {
    err := c.setupDB()
    if err != nil {
        return err
    }
    return nil
}

func (c *Config) setupDB() error {
    err := db.ConnectToConfig(c.dbConfig)
    if err != nil {
        return errors.Join(errors.New("connection to database failed"), err)
    }
    return nil
}

func (c *Config) ToAddress() string {
    return fmt.Sprintf("%s:%s", c.Ip, c.Port)
}