package model

import (
	"database/sql"
	"errors"
	"log/slog"

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
    row := db.DB.QueryRow("SELECT user.id, firstname, lastname, username, password_changed, phone_number, email, role.name, deleted_at FROM user INNER JOIN role ON user.role_id = role.id WHERE user.id = ? AND deleted_at IS NULL", userId)
    user := User{}
    err := row.Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Username, &user.PasswordChanged, &user.PhoneNumber, &user.Email, &user.Role, &user.DeletedAt)
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

func GetAllUser() ([]User, error) {
    rows, err := db.DB.Query("SELECT user.id, firstname, lastname, username, phone_number, email, role.name, deleted_at FROM user INNER JOIN role ON user.role_id = role.id WHERE deleted_at IS NULL")
    if err != nil {
        return []User{}, err
    }

    users := make([]User, 0)
    for rows.Next() {
        user := User{}
        err = rows.Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Username, &user.PhoneNumber, &user.Email, &user.Role, &user.DeletedAt)
        if err != nil {
            return []User{}, err
        }
        users = append(users, user)
    }

    return users, nil
}

func (u *User) UpdatePassword() error {
    _, err := db.DB.Exec("UPDATE user SET password = ?, password_changed = TRUE WHERE id = ?", u.Password, u.Id)
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
func (u *User) ValidateDelete(deleteBy User, storeLeaderRole string) error {
    if u.Id == 0 {
        return errors.New("missing id for deleting user") //TODO: hungarian error message
    }

    user, err := GetUserByUserId(u.Id)
    if err != nil {
        slog.Error("could not get user to be deleted with given id", "error", err, "id", u.Id)
        return errors.New("could not get user with given id") //TODO: hungarian error message
    }

    if user.Role != storeLeaderRole {
        return nil
    }

    if deleteBy.Role != storeLeaderRole {
        return errors.New("an HR person can not delete store leader account") //TODO: hungarian error message
    }

    numOfStoreLeaders := 0
    row := db.DB.QueryRow("SELECT COUNT(*) FROM user INNER JOIN role ON user.role_id = role.id WHERE role.name = ? AND user.deleted_at IS NULL", storeLeaderRole)
    err = row.Scan(&numOfStoreLeaders)
    if err != nil {
        slog.Error("could not count store leaders in db", "error", err)
        return errors.New("could not count store leaders in db") //TODO: hungarian error message
    }
    if numOfStoreLeaders < 2 {
        return errors.New("could not delete store leader, because there has to remain at least 1 store leader in the db") //TODO: hungarian error message
    }

    return nil
}

//Returns user-readable error
func (u *User) ValidateUpdatePassword(oldPassword, newPassword string) error {
    if len(newPassword) < 8 {
        return errors.New("the new password is not long enough (at least 8 characters)") //TODO: user-readable error message
    }

    realOldPassword := ""
    row := db.DB.QueryRow("SELECT password FROM user WHERE id = ?", u.Id)
    err := row.Scan(&realOldPassword)
    if err != nil {
        slog.Error("could not scan real old password of user during password change", "error", err)
        return errors.New("something went wrong") //TODO: user-readable error message
    }

    if !CheckPassword(oldPassword, realOldPassword) {
        return errors.New("password does not match your password") //TODO: user-readable error message
    }

    u.Password = newPassword
    return nil
}

func CheckPassword(tryPassword, realPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(realPassword), []byte(tryPassword))
    return err == nil
}
