package auth

import (
	"database/sql"
	"errors"

	"github.com/vukovlevi/netstore/store_administration/model"
	"golang.org/x/crypto/bcrypt"
)

var (
    ErrUserNotFound = errors.New("user not found")
    ErrBadPassword = errors.New("password mismatch")
)

func LoginUser(username, password string) error {
    user, err := model.GetUserByUsername(username)
    if err != nil {
        if err == sql.ErrNoRows {
            return ErrUserNotFound
        }
        return err
    }

    if ok := checkPassword(password, user.Password); !ok {
        return ErrBadPassword
    }

    return nil
}

func checkPassword(tryPassword, realPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(realPassword), []byte(tryPassword))
    return err == nil
}
