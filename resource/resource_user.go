package resource

import (
	"errors"
	"fmt"
	"github.com/manyminds/api2go"
	"github.com/timrourke/timrourke.com/model"
	"github.com/timrourke/timrourke.com/storage"
	"net/http"
)

type UserResource struct {
	UserStorage *storage.UserStorage
}

// FindAll to satisfy api2go data source interface
func (s UserResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	result, err := s.UserStorage.GetAll()
	return &Response{Res: result}, err
}

// PaginatedFindAll can be used to load users in chunks
//func (s UserResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
//	return uint(len(users)), &Response{Res: result}, nil
//}

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
	return &Response{Res: newUser, Code: http.StatusCreated}, err
}

// Delete to satisfy `api2go.DataSource` interface
func (s UserResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	fmt.Println("id to delete", id)
	err := s.UserStorage.Delete(id)
	if err != nil {
		return &Response{Code: http.StatusNotFound}, err
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
