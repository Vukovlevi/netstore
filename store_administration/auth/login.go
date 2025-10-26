package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/vukovlevi/netstore/store_administration/model"
)

var (
    ErrUserNotFound = errors.New("user not found")
    ErrBadPassword = errors.New("password mismatch")
)

func LoginUser(username, password string, ctx *context.Context) error {
    user, err := model.GetUserByUsername(username)
    if err != nil {
        if err == sql.ErrNoRows {
            return ErrUserNotFound
        }
        return err
    }
    *ctx = context.WithValue(*ctx, "user", user)

    if ok := model.CheckPassword(password, user.Password); !ok {
        return ErrBadPassword
    }

    return CreateOrUpdateSessionForUser(user.Id, ctx)
}
