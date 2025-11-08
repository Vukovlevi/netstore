package model

import (
	"time"

	"github.com/vukovlevi/netstore/store_administration/db"
)

const (
    HOURS_A_DAY = 24
    EXPIRES_IN_DAYS = 7
)

type Session struct {
    Id int
    UserId int
    Token string
    ExpiresAt time.Time
}

func GetSessionByUserId(userId int) (Session, error) {
    row := db.DB.QueryRow("SELECT id, user_id, token, expires_at FROM session WHERE user_id = ?", userId)
    session := Session{}
    err := row.Scan(&session.Id, &session.UserId, &session.Token, &session.ExpiresAt)
    return session, err
}

func GetSessionByToken(token string) (Session, error) {
    row := db.DB.QueryRow("SELECT id, user_id, token, expires_at FROM session WHERE token = ? AND expires_at > NOW()", token)
    session := Session{}
    err := row.Scan(&session.Id, &session.UserId, &session.Token, &session.ExpiresAt)
    return session, err
}

func (s *Session) UpdateExpiry() error {
    s.setNewExpiresAt()
    _, err := db.DB.Exec("UPDATE session SET expires_at = ? WHERE id = ?", s.ExpiresAt, s.Id)
    return err
}

func (s *Session) InsertNewSession() error {
    s.setNewExpiresAt()
    _, err := db.DB.Exec("INSERT INTO session VALUES (NULL, ?, ?, ?)", s.UserId, s.Token, s.ExpiresAt)
    return err
}

func (s *Session) ChangeExpiredToNew() error {
    s.setNewExpiresAt()
    _, err := db.DB.Exec("UPDATE session SET token = ?, expires_at = ? WHERE id = ?", s.Token, s.ExpiresAt, s.Id)
    return err
}

func (s *Session) setNewExpiresAt() {
    s.ExpiresAt = time.Now().Add(time.Hour * HOURS_A_DAY * EXPIRES_IN_DAYS)
}
