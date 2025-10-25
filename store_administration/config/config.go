package config

import (
	"errors"

	"github.com/vukovlevi/netstore/store_administration/db"
)

type Config struct {
    dbConfig db.DBConfig
}

func CreateApplicationConfig() Config {
    return Config{
        dbConfig: db.CreateDatabaseConfig(),
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
