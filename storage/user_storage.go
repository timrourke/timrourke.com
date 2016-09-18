package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/timrourke/timrourke.com/model"
	"github.com/timrourke/timrourke.com/query"
	"strconv"
)

// NewUserStorage returns a new instance of UserStorage
func NewUserStorage(DB *sqlx.DB) *UserStorage {
	return &UserStorage{DB}
}

// UserStorage forms SQL queries for users
type UserStorage struct {
	DB *sqlx.DB
}

// GetAll selects a list of users
func (s *UserStorage) GetAll(q *query.Query) (uint, []model.User, error) {
	var (
		users []model.User
		count uint
	)

	q.Select("users.*").From("users users")

	sql, boundValues := q.Compile()

	rows, err := s.DB.NamedQuery(sql, boundValues)
	if err != nil {
		return 0, nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var m model.User
		err = rows.StructScan(&m)
		if err != nil {
			return 0, nil, err
		}

		users = append(users, m)
	}

	// Get count of all users for pagination
	errCount := s.DB.Get(&count, "SELECT COUNT(*) FROM users")

	if err != nil {
		return 0, nil, err
	} else if errCount != nil {
		return 0, nil, errCount
	}

	return count, users, nil
}

// GetOne selects a single user
func (s *UserStorage) GetOne(ID string) (*model.User, error) {
	var user model.User

	err := s.DB.Get(&user, "SELECT * FROM users WHERE id=?", ID)

	return &user, err
}

// Insert inserts a single user
func (s *UserStorage) Insert(c model.User) (*model.User, error) {
	result, err := s.DB.NamedExec(`INSERT INTO users (
		username,
		email,
		password_hash
	) VALUES (
		:username,
		:email,
		:password_hash
	)`, &c)

	if err != nil {
		return &model.User{}, err
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		return &model.User{}, err
	}

	// Set ID on return struct for rendering to json
	c.SetID(fmt.Sprintf("%d", insertID))

	return s.GetOne(c.GetID())
}

// Delete deletes a single user
func (s *UserStorage) Delete(id string) error {
	_, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("User id must be integer: %s", id)
	}

	_, err = s.DB.Exec("DELETE FROM users WHERE id=? LIMIT 1", id)
	if err != nil {
		return err
	}
	return nil
}

// Update updates a single user
func (s *UserStorage) Update(c *model.User) error {
	_, err := s.DB.NamedExec(`UPDATE users SET 
		username=:username,
		email=:email,
		password_hash=:password_hash
		WHERE id=:id`, &c)

	if err != nil {
		return err
	}

	return nil
}
