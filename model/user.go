package model

import (
	"github.com/manyminds/api2go/jsonapi"
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
}

type NewUser struct {
	Username             string `json:"username"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password-confirmation"`
}

func (m User) GetID() string {
	return strconv.FormatInt(m.ID, 10)
}

func (m *User) SetID(id string) error {
	var err error
	m.ID, err = strconv.ParseInt(id, 10, 64)
	return err
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (m User) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type:        "posts",
			Name:        "posts",
			IsNotLoaded: true,
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (m User) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}

	return result
}
