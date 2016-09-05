package model

import (
	"strconv"
	"time"
)

type User struct {
	ID int64 `json:"-"`

	CreatedAt    time.Time `json:"created-at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated-at" db:"updated_at"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	exists       bool      `sql:"-"`
}

func (m User) GetID() string {
	return strconv.FormatInt(m.ID, 10)
}

func (m *User) SetID(id string) error {
	var err error
	m.ID, err = strconv.ParseInt(id, 10, 64)
	return err
}
