package resource

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/manyminds/api2go"
	"github.com/timrourke/timrourke.com/model"
	"github.com/timrourke/timrourke.com/query"
	"github.com/timrourke/timrourke.com/storage"
	"net/http"
	"strconv"
)

// UserResource defines interface to storage layer
type UserResource struct {
	UserStorage *storage.UserStorage
}

// UserFilterableFields is a map of fields a user can sort or filter by, where
// the key is the jsonapi field name and the value is whether a filter should
// be performed using strict equality (true), or using a LIKE statement (false),
// in the SQL generated for the query
var UserFilterableFields = map[string]bool{
	"id":         true,
	"created-at": false,
	"updated-at": false,
	"email":      true,
	"username":   true,
}

func getUsersByPostsID(request api2go.Request, q *query.Query) *query.Query {
	postsID, ok := request.QueryParams["postsID"]

	if ok {
		q.Join("LEFT JOIN posts posts", "posts.user_id = users.id")
		q.Where("posts.id = :postsID")
		q.Bind("postsID", postsID[0])
	}

	return q
}

// UserRelationshipsByParam defines a map where the key is the query param and
// the function is the RelationshipFunc for modifying the query to get the given
// relationship
var UserRelationshipsByParam = map[string]RelationshipFunc{
	"postsID": getUsersByPostsID,
}

// FindAll to satisfy api2go data source interface
func (s UserResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	// 400
	params, err := ParseQueryParams(r, UserFilterableFields, UserRelationshipsByParam)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(
			err,
			err.Error(),
			http.StatusBadRequest,
		)
	}

	// 500
	_, result, err := s.UserStorage.GetAll(params)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(
			err,
			"Internal Server Error",
			http.StatusInternalServerError)
	}

	return &Response{Res: result}, err
}

// PaginatedFindAll can be used to load users in chunks
func (s UserResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	// 400
	params, err := ParseQueryParams(r, UserFilterableFields, UserRelationshipsByParam)
	if err != nil {
		return 0, &Response{}, api2go.NewHTTPError(
			err,
			err.Error(),
			http.StatusBadRequest,
		)
	}

	// 500
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
func (s UserResource) FindOne(id string, r api2go.Request) (api2go.Responder, error) {
	// 400
	_, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		errMessage := fmt.Sprintf("User id must be integer: %s", id)

		return &Response{}, api2go.NewHTTPError(
			err,
			errMessage,
			http.StatusBadRequest)
	}

	// 404
	user, err := s.UserStorage.GetOne(id)
	if err == sql.ErrNoRows {
		errMessage := fmt.Sprintf("No user found with the id: %s", id)

		return &Response{}, api2go.NewHTTPError(
			err,
			errMessage,
			http.StatusNotFound,
		)
	}

	return &Response{Res: user}, err
}

// Create method to satisfy `api2go.DataSource` interface
func (s UserResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	// 400
	user, ok := obj.(model.User)
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest)
	}

	// 500
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
	// 400
	_, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		errMessage := fmt.Sprintf("User id must be integer: %s", id)

		return &Response{}, api2go.NewHTTPError(
			err,
			errMessage,
			http.StatusBadRequest)
	}

	err = s.UserStorage.Delete(id)
	if err != nil {
		return &Response{Code: http.StatusInternalServerError}, err
	}
	return &Response{Code: http.StatusNoContent}, nil
}

// Update stores all changes on the user
func (s UserResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	user, ok := obj.(*model.User)

	// 400
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest)
	}

	id := user.GetID()
	foundUser, err := s.UserStorage.GetOne(id)

	// 404
	if err == sql.ErrNoRows {
		errMessage := fmt.Sprintf("No user found with the id: %s", id)

		return &Response{}, api2go.NewHTTPError(
			err,
			errMessage,
			http.StatusNotFound,
		)

		// 500
	} else if err != nil {
		return &Response{}, api2go.NewHTTPError(
			err,
			"Internal Server Error",
			http.StatusInternalServerError)
	}

	// Update fields in user
	foundUser.Email = user.Email
	foundUser.Username = user.Username
	// TODO: implement password hashing

	// 500
	err = s.UserStorage.Update(foundUser)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(
			err,
			"Internal Server Error",
			http.StatusInternalServerError)
	}

	return &Response{Res: foundUser, Code: http.StatusNoContent}, err
}
