package model

import (
	"errors"
	"github.com/manyminds/api2go/jsonapi"
	"strconv"
	"time"
)

type Post struct {
	ID int64 `json:"-"`

	CreatedAt time.Time `json:"created-at" db:"created_at"`
	UpdatedAt time.Time `json:"updated-at" db:"updated_at"`
	Title     string    `json:"title" db:"title"`
	Excerpt   string    `json:"excerpt" db:"excerpt"`
	Content   string    `json:"content" db:"content"`
	Permalink string    `json:"permalink" db:"permalink"`
	User      *User     `json:"-"`
	UserId    string    `json:"-" db:"user_id"`
}

func (m Post) GetID() string {
	return strconv.FormatInt(m.ID, 10)
}

func (m *Post) SetID(id string) error {
	var err error
	m.ID, err = strconv.ParseInt(id, 10, 64)
	return err
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (m Post) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "users",
			Name: "users",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (m Post) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}

	result = append(result, jsonapi.ReferenceID{
		ID:   m.UserId,
		Type: "users",
		Name: "users",
	})

	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (m Post) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}
	// for key := range u.Chocolates {
	// 	result = append(result, u.Chocolates[key])
	// }

	return result
}

// SetToManyReferenceIDs sets the sweets reference IDs and satisfies the jsonapi.UnmarshalToManyRelations interface
func (m *Post) SetToManyReferenceIDs(name string, IDs []string) error {
	// if name == "sweets" {
	// 	u.ChocolatesIDs = IDs
	// 	return nil
	// }

	return errors.New("There is no to-many relationship with the name " + name)
}

// AddToManyIDs adds some new sweets that a users loves so much
func (m *Post) AddToManyIDs(name string, IDs []string) error {
	// if name == "sweets" {
	// 	u.ChocolatesIDs = append(u.ChocolatesIDs, IDs...)
	// 	return nil
	// }

	return errors.New("There is no to-many relationship with the name " + name)
}

// DeleteToManyIDs removes some sweets from a users because they made him very sick
func (m *Post) DeleteToManyIDs(name string, IDs []string) error {
	// if name == "sweets" {
	// 	for _, ID := range IDs {
	// 		for pos, oldID := range u.ChocolatesIDs {
	// 			if ID == oldID {
	// 				// match, this ID must be removed
	// 				u.ChocolatesIDs = append(u.ChocolatesIDs[:pos], u.ChocolatesIDs[pos+1:]...)
	// 			}
	// 		}
	// 	}
	// }

	return errors.New("There is no to-many relationship with the name " + name)
}
