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
    CentralServerAddress string
    CentralServerPort string
    Psk string
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

    centralServerAddress := os.Getenv("CENTRAL_SERVER_ADDRESS")
    if centralServerAddress == "" {
        centralServerAddress = "localhost"
    }

    centralServerPort := os.Getenv("CENTRAL_SERVER_PORT")
    if centralServerPort == "" {
        centralServerPort = "42069"
    }

    psk := os.Getenv("PSK")
    return Config{
        dbConfig: db.CreateDatabaseConfig(),
        Ip: ip,
        Port: port,
        CentralServerAddress: centralServerAddress,
        CentralServerPort: centralServerPort,
        Psk: psk,
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
