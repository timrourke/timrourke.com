package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/manyminds/api2go"
	"github.com/timrourke/timrourke.com/model"
	"net/http"
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
func (s *UserStorage) GetAll(params QueryParams) (uint, []model.User, error) {
	var (
		users []model.User
		count uint
	)

	sql := fmt.Sprintf("SELECT * FROM users ORDER BY %s LIMIT ?,?", params.OrderBy)
	err := s.DB.Select(&users,
		sql,
		params.Offset,
		params.Limit,
	)

	// Get count of all users for pagination
	errCount := s.DB.Get(&count, "SELECT COUNT(*) FROM users")

	if err != nil || errCount != nil {
		err := errors.New("Server error retrieving all users")

		return 0, nil, err
	}

	return count, users, nil
}

// GetOne selects a single user
func (s *UserStorage) GetOne(id string) (*model.User, error) {
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		errMessage := fmt.Sprintf("User id must be integer: %s", id)

		return &model.User{}, api2go.NewHTTPError(
			errors.New(errMessage),
			errMessage,
			http.StatusBadRequest)
	}

	var user model.User

	err = s.DB.Get(&user, "SELECT * FROM users WHERE id=?", intID)
	if err == sql.ErrNoRows {
		errMessage := fmt.Sprintf("No user found with the id %d", intID)

		return &user, api2go.NewHTTPError(
			errors.New(errMessage),
			errMessage,
			http.StatusNotFound,
		)
	}

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
func (s *UserStorage) Update(c model.User) error {
	return nil
}
