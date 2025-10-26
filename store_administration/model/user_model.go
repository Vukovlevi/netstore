package model

import (
	"database/sql"
	"errors"

	"github.com/vukovlevi/netstore/store_administration/db"
	"golang.org/x/crypto/bcrypt"
)

const (
    PASSWORD_HASH_COST = 12
)

type User struct {
    Id int `json:"id"`
    Firstname string `json:"firstname"`
    Lastname string `json:"lastname"`
    Username string `json:"username"`
    Password string `json:"password,omitempty,omitzero"`
    PasswordChanged bool `json:"passwordChanged,omitempty,omitzero"`
    PhoneNumber string `json:"phoneNumber"`
    Email string `json:"email"`
    Role string `json:"role"`
    RoleId int `json:"roleId,omitempty,omitzero"`
    DeletedAt sql.NullTime `json:"deletedAt"`
}

func GetUserByUsername(username string) (User, error) {
    row := db.DB.QueryRow("SELECT user.id, firstname, lastname, username, password, password_changed, phone_number, email, role.name, deleted_at FROM user INNER JOIN role ON user.role_id = role.id WHERE username = ? AND deleted_at IS NULL", username)
    user := User{}
    err := row.Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Username, &user.Password, &user.PasswordChanged, &user.PhoneNumber, &user.Email, &user.Role, &user.DeletedAt)
    return user, err
}

func GetUserByUserId(userId int) (User, error) {
    row := db.DB.QueryRow("SELECT user.id, firstname, lastname, username, password, password_changed, phone_number, email, role.name, deleted_at FROM user INNER JOIN role ON user.role_id = role.id WHERE user.id = ? AND deleted_at IS NULL", userId)
    user := User{}
    err := row.Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Username, &user.Password, &user.PasswordChanged, &user.PhoneNumber, &user.Email, &user.Role, &user.DeletedAt)
    user.Password = ""
    return user, err
}

func (u *User) InsertNewUser() error {
    _, err := db.DB.Exec("INSERT INTO user VALUES (NULL, ?, ?, ?, ?, ?, ?, ?, ?, NULL)", u.Firstname, u.Lastname, u.Username, u.Password, u.PasswordChanged, u.PhoneNumber, u.Email, u.RoleId)
    return err
}

func (u *User) UpdateUser() error {
    _, err := db.DB.Exec("UPDATE user SET firstname = ?, lastname = ?, username = ?, phone_number = ?, email = ?, role_id = ? WHERE id = ?", u.Firstname, u.Lastname, u.Username, u.PhoneNumber, u.Email, u.RoleId, u.Id)
    return err
}

func (u *User) DeleteUser() error {
    _, err := db.DB.Exec("UPDATE user SET deleted_at = NOW() WHERE id = ?", u.Id)
    return err
}

func (u *User) EncryptPassword() error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), PASSWORD_HASH_COST)
    if err != nil {
        return err
    }

    u.Password = string(hashedPassword)
    return nil
}

//Returns user-readable error
func (u *User) ValidateInsert() error {
    if u.Firstname == "" || u.Lastname == "" || u.Username == "" || u.Password == "" || u.RoleId == 0 {
        return errors.New("missing required parameter for saving new user") //TODO: hungarian error message
    }
    return nil
}

//Returns user-readable error
func (u *User) ValidateUpdate() error {
    if u.Id == 0 || u.Firstname == "" || u.Lastname == "" || u.Username == "" || u.RoleId == 0 {
        return errors.New("missing required paramter for updating user") //TODO: hungarian error message
    }
    return nil
}

//Returns user-readable error
func (u *User) ValidateDelete() error {
    if u.Id == 0 {
        return errors.New("missing id for deleting user") //TODO: hungarian error message
    }
    return nil
}
