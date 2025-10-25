package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectToConfig(conf DBConfig) error {
    db, err := sql.Open("mysql", conf.ToConnectionString())
    if err != nil {
        return err
    }

    err = db.Ping()
    if err != nil {
        return err
    }

    DB = db
    return nil
}

func Disconnect() {
    DB.Close()
}
