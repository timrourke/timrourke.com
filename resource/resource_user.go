package resource

import (
	"errors"
	"fmt"
	"github.com/manyminds/api2go"
	"github.com/timrourke/timrourke.com/model"
	"github.com/timrourke/timrourke.com/storage"
	"net/http"
	"strconv"
	"strings"
)

// UserResource defines interface to storage layer
type UserResource struct {
	UserStorage *storage.UserStorage
}

const defaultPaginationLimit = 20

// ParseQueryParams parses request for query params
func (s UserResource) ParseQueryParams(r api2go.Request) (storage.QueryParams, error) {
	var (
		queryLimit   uint64 = defaultPaginationLimit
		queryOffset  uint64
		queryOrderBy string
	)

	requestParams := r.QueryParams

	offset, hasOffset := requestParams["page[offset]"]
	limit, hasLimit := requestParams["page[limit]"]
	pageNum, hasPageNum := requestParams["page[number]"]
	pageSize, hasPageSize := requestParams["page[size]"]
	sorts, hasSorts := requestParams["sort"]

	if hasLimit {
		parsedLimit, err := strconv.ParseUint(limit[0], 10, 64)
		if err != nil {
			return storage.QueryParams{}, err
		}

		queryLimit = parsedLimit
	} else if hasPageSize {
		parsedPageSize, err := strconv.ParseUint(pageSize[0], 10, 64)
		if err != nil {
			return storage.QueryParams{}, err
		}

		queryLimit = parsedPageSize
	}

	if queryLimit > 100 {
		queryLimit = 100
	}

	if hasPageNum {
		parsedPageNum, err := strconv.ParseUint(pageNum[0], 10, 64)
		if err != nil {
			return storage.QueryParams{}, nil
		}

		queryOffset = (parsedPageNum - 1) * queryLimit
	} else if hasOffset {
		parsedOffset, err := strconv.ParseUint(offset[0], 10, 64)
		if err != nil {
			return storage.QueryParams{}, nil
		}

		queryOffset = parsedOffset
	}

	if hasSorts {
		dir := "ASC"
		numSorts := len(sorts)
		queryOrderBy = ""

		for i, v := range sorts {
			if strings.HasPrefix(v, "-") {
				dir = "DESC"
				v = strings.TrimPrefix(v, "-")
			} else {
				dir = "ASC"
			}

			if (i + 1) < numSorts {
				dir = fmt.Sprintf("%s, ", dir)
			}

			queryOrderBy = fmt.Sprintf("%s%s %s", queryOrderBy, v, dir)
		}
	}

	params := storage.QueryParams{
		Limit:   queryLimit,
		Offset:  queryOffset,
		OrderBy: queryOrderBy,
	}

	return params, nil
}

// FindAll to satisfy api2go data source interface
func (s UserResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	params, _ := s.ParseQueryParams(r)
	_, result, err := s.UserStorage.GetAll(params)
	return &Response{Res: result}, err
}

// PaginatedFindAll can be used to load users in chunks
func (s UserResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	params, _ := s.ParseQueryParams(r)
	count, result, err := s.UserStorage.GetAll(params)
	if err != nil {
		return 0, &Response{}, err
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
