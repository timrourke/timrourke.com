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

// PostResource defines interface to storage layer
type PostResource struct {
	PostStorage *storage.PostStorage
}

// PostFilterableFields is a map of fields a post can sort or filter by, where
// the key is the jsonapi field name and the value is whether a filter should
// be performed using strict equality (true), or using a LIKE statement (false),
// in the SQL generated for the query
var PostFilterableFields = map[string]bool{
	"id":         true,
	"created-at": false,
	"updated-at": false,
	"permalink":  true,
}

func getPostsByUsersID(request api2go.Request, q *query.Query) *query.Query {
	usersID, ok := request.QueryParams["usersID"]

	if ok {
		q.Where("posts.user_id = :usersID")
		q.Bind("usersID", usersID[0])
	}

	return q
}

// PostRelationships defines the functions for modifying a Query to select...
var PostRelationships = map[string]RelationshipFunc{
	"usersID": getPostsByUsersID,
}

// FindAll to satisfy api2go data source interface
func (s PostResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	// 400
	params, err := ParseQueryParams(r, PostFilterableFields, PostRelationships)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(
			err,
			err.Error(),
			http.StatusBadRequest,
		)
	}

	// 500
	_, result, err := s.PostStorage.GetAll(params)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(
			err,
			"Internal Server Error",
			http.StatusInternalServerError)
	}

	return &Response{Res: result}, err
}

// PaginatedFindAll can be used to load posts in chunks
func (s PostResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	// 400
	params, err := ParseQueryParams(r, PostFilterableFields, PostRelationships)
	if err != nil {
		return 0, &Response{}, api2go.NewHTTPError(
			err,
			err.Error(),
			http.StatusBadRequest,
		)
	}

	// 500
	count, result, err := s.PostStorage.GetAll(params)
	if err != nil {
		return 0, &Response{}, api2go.NewHTTPError(
			err,
			err.Error(),
			http.StatusInternalServerError)
	}

	return count, &Response{Res: result}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the post with the given ID, otherwise an error
func (s PostResource) FindOne(id string, r api2go.Request) (api2go.Responder, error) {
	// 400
	_, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		errMessage := fmt.Sprintf("Post id must be integer: %s", id)

		return &Response{}, api2go.NewHTTPError(
			err,
			errMessage,
			http.StatusBadRequest)
	}

	// 404
	post, err := s.PostStorage.GetOne(id)
	if err == sql.ErrNoRows {
		errMessage := fmt.Sprintf("No post found with the id: %s", id)

		return &Response{}, api2go.NewHTTPError(
			err,
			errMessage,
			http.StatusNotFound,
		)
	}

	return &Response{Res: post}, err
}

// Create method to satisfy `api2go.DataSource` interface
func (s PostResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	// 400
	post, ok := obj.(model.Post)
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest)
	}

	// 500
	newPost, err := s.PostStorage.Insert(post)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Internal Server Error"),
			"Internal Server Error",
			http.StatusInternalServerError)
	}

	return &Response{Res: newPost, Code: http.StatusCreated}, err
}

// Delete to satisfy `api2go.DataSource` interface
func (s PostResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	// 400
	_, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		errMessage := fmt.Sprintf("Post id must be integer: %s", id)

		return &Response{}, api2go.NewHTTPError(
			err,
			errMessage,
			http.StatusBadRequest)
	}

	err = s.PostStorage.Delete(id)
	if err != nil {
		return &Response{Code: http.StatusInternalServerError}, err
	}
	return &Response{Code: http.StatusNoContent}, nil
}

// Update stores all changes on the post
func (s PostResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	post, ok := obj.(*model.Post)

	// 400
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest)
	}

	id := post.GetID()
	foundPost, err := s.PostStorage.GetOne(id)

	// 404
	if err == sql.ErrNoRows {
		errMessage := fmt.Sprintf("No post found with the id: %s", id)

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

	// Update fields in post
	foundPost.Title = post.Title
	foundPost.Excerpt = post.Excerpt
	foundPost.Content = post.Content
	foundPost.Permalink = post.Permalink
	// TODO: implement santization and validation

	// 500
	err = s.PostStorage.Update(foundPost)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(
			err,
			"Internal Server Error",
			http.StatusInternalServerError)
	}

	return &Response{Res: foundPost, Code: http.StatusNoContent}, err
}
