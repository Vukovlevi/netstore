package db

import (
	"fmt"
	"os"
)

const (
    ENV_DB_USERNAME="DB_USERNAME"
    ENV_DB_PASSWORD="DB_PASSWORD"
    ENV_DB_PROTOCOL="DB_PROTOCOL"
    ENV_DB_HOST="DB_HOST"
    ENV_DB_PORT="DB_PORT"
    ENV_DB_NAME="DB_NAME"
)

type DBConfig struct {
    username string
    password string
    protocol string
    host string
    port string
    dbName string
}

func CreateDatabaseConfig() DBConfig {
    return DBConfig{
        username: os.Getenv(ENV_DB_USERNAME),
        password: os.Getenv(ENV_DB_PASSWORD),
        protocol: os.Getenv(ENV_DB_PROTOCOL),
        host: os.Getenv(ENV_DB_HOST),
        port: os.Getenv(ENV_DB_PORT),
        dbName: os.Getenv(ENV_DB_NAME),
    }
}

func (c *DBConfig) ToConnectionString() string {
    return fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=true&loc=Local", c.username, c.password, c.protocol, c.host, c.port, c.dbName)
}
