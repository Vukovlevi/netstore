package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/vukovlevi/netstore/store_administration/model"
)

const (
    TOKEN_LENGTH = 32
)

func CreateOrUpdateSessionForUser(userId int, ctx *context.Context) error {
    session, err := model.GetSessionByUserId(userId)
    if err != nil {
        if err == sql.ErrNoRows {
            return createSessionForUser(userId, ctx)
        }
        return err
    }
    *ctx = context.WithValue(*ctx, "session", &session)

    if session.ExpiresAt.Before(time.Now()) {
        session.Token = generateToken(TOKEN_LENGTH)
        return session.ChangeExpiredToNew()
    }

    return session.UpdateExpiry()
}

func createSessionForUser(userId int, ctx *context.Context) error {
    session := model.Session{
        UserId: userId,
        Token: generateToken(TOKEN_LENGTH),
    }
    *ctx = context.WithValue(*ctx, "session", session)
    return session.InsertNewSession()
}
