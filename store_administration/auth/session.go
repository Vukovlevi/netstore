package auth

import (
	"database/sql"

	"github.com/vukovlevi/netstore/store_administration/model"
)

const (
    TOKEN_LENGTH = 64
)

func CreateOrUpdateSessionForUser(userId int) error {
    session, err := model.GetSessionByUserId(userId)
    if err != nil {
        if err == sql.ErrNoRows {
            return createSessionForUser(userId)
        }
        return err
    }
    return session.UpdateExpiry()
}

func createSessionForUser(userId int) error {
    session := model.Session{
        UserId: userId,
        Token: generateToken(TOKEN_LENGTH),
    }
    return session.InsertNewSession()
}
