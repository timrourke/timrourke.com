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

func NewUserStorage(DB *sqlx.DB) *UserStorage {
	return &UserStorage{DB}
}

type UserStorage struct {
	DB *sqlx.DB
}

func (s *UserStorage) GetAll() ([]model.User, error) {
	var users []model.User

	err := s.DB.Select(&users, "SELECT * FROM users")
	if err != nil {
		errMessage := "Server error retrieving all users"

		return nil, api2go.NewHTTPError(
			errors.New(errMessage),
			errMessage,
			http.StatusInternalServerError)
	}

	return users, nil
}

func (s *UserStorage) GetOne(id string) (*model.User, error) {
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		errMessage := fmt.Sprintf("User id must be integer: %s", id)

		return &model.User{}, api2go.NewHTTPError(
			errors.New(errMessage),
			errMessage,
			http.StatusBadRequest)
	}
	var user model.User
	err = s.DB.Get(&user, "SELECT * FROM users WHERE id=?", intId)
	if err == sql.ErrNoRows {
		errMessage := fmt.Sprintf("No user found with the id %d", intId)

		return &user, api2go.NewHTTPError(
			errors.New(errMessage),
			errMessage,
			http.StatusNotFound,
		)
	}
	return &user, err
}

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

	insertId, err := result.LastInsertId()
	if err != nil {
		return &model.User{}, err
	}

	c.SetID(fmt.Sprintf("%d", insertId))
	return s.GetOne(c.GetID())
}

func (s *UserStorage) Delete(id string) error {
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("User id must be integer: %s", id)
	}

	result, err := s.DB.Exec("DELETE FROM users WHERE id=? LIMIT 1", id)
	if err != nil {
		return err
	}

	numDeleted, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if numDeleted == 0 {
		return fmt.Errorf("No user found with the id %d", intId)
	}
	return nil
}

func (s *UserStorage) Update(c model.User) error {
	return nil
}
