package resource

import (
	"errors"
	"github.com/manyminds/api2go"
	"github.com/timrourke/timrourke.com/model"
	"github.com/timrourke/timrourke.com/storage"
	"net/http"
)

// UserResource defines interface to storage layer
type UserResource struct {
	UserStorage *storage.UserStorage
}

var UserFilterableFields = map[string]bool{
	"id":         true,
	"created-at": false,
	"updated-at": false,
	"email":      true,
	"username":   true,
}

// FindAll to satisfy api2go data source interface
func (s UserResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	params, err := ParseQueryParams(r, UserFilterableFields)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(
			err,
			err.Error(),
			http.StatusBadRequest,
		)
	}

	_, result, err := s.UserStorage.GetAll(params)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(
			err,
			err.Error(),
			http.StatusInternalServerError)
	}

	return &Response{Res: result}, err
}

// PaginatedFindAll can be used to load users in chunks
func (s UserResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	params, err := ParseQueryParams(r, UserFilterableFields)
	if err != nil {
		return 0, &Response{}, api2go.NewHTTPError(
			err,
			err.Error(),
			http.StatusBadRequest,
		)
	}

	count, result, err := s.UserStorage.GetAll(params)
	if err != nil {
		return 0, &Response{}, api2go.NewHTTPError(
			err,
			err.Error(),
			http.StatusInternalServerError)
	}

	return count, &Response{Res: result}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the user with the given ID, otherwise an error
func (s UserResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	user, err := s.UserStorage.GetOne(ID)
	return &Response{Res: user}, err
}

// Create method to satisfy `api2go.DataSource` interface
func (s UserResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	user, ok := obj.(model.User)
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest)
	}

	newUser, err := s.UserStorage.Insert(user)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Internal Server Error"),
			"Internal Server Error",
			http.StatusInternalServerError)
	}

	return &Response{Res: newUser, Code: http.StatusCreated}, err
}

// Delete to satisfy `api2go.DataSource` interface
func (s UserResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := s.UserStorage.Delete(id)
	if err != nil {
		return &Response{Code: http.StatusInternalServerError}, err
	}
	return &Response{Code: http.StatusNoContent}, nil
}

//Update stores all changes on the user
func (s UserResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	user, ok := obj.(model.User)
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest)
	}

	foundUser, err := s.UserStorage.GetOne(user.GetID())
	return &Response{Res: foundUser, Code: http.StatusNoContent}, err
}
