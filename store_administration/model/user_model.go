package model

import (
	"time"

	"github.com/vukovlevi/netstore/store_administration/db"
)

type User struct {
    Id int
    Firtname string
    Lastname string
    Username string
    Password string
    PasswordChanged bool
    PhoneNumber string
    Email string
    Role string
    DeletedAt *time.Time
}

func GetUserByUsername(username string) (User, error) {
    row := db.DB.QueryRow("SELECT user.id, firstname, lastname, username, password, password_changed, phone_number, email, role.name, deleted_at FROM user INNER JOIN role ON user.role_id = role.id WHERE username = ? AND deleted_at IS NULL", username)
    user := User{}
    err := row.Scan(&user.Id, &user.Firtname, &user.Lastname, &user.Username, &user.Password, &user.PasswordChanged, &user.PhoneNumber, &user.Email, &user.Role, &user.DeletedAt)
    return user, err
}
