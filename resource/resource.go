package resource

import (
	"fmt"
	"github.com/manyminds/api2go"
	"github.com/timrourke/timrourke.com/storage"
	"strconv"
	"strings"
)

const defaultPaginationLimit = 20

// ParseQueryParams parses request for query params
func ParseQueryParams(r api2go.Request, filterableFields map[string]bool) (storage.QueryParams, error) {
	var (
		err          error
		queryLimit   uint64 = defaultPaginationLimit
		queryOffset  uint64
		queryOrderBy string = "id ASC"
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

	if hasSorts && len(sorts[0]) > 0 {
		queryOrderBy, err = handleSorts(sorts, filterableFields)
		if err != nil {
			return storage.QueryParams{}, err
		}
	}

	params := storage.QueryParams{
		Limit:   queryLimit,
		Offset:  queryOffset,
		OrderBy: queryOrderBy,
	}

	return params, nil
}

// handleSorts builds a SQL string for an order by statement
func handleSorts(sorts []string, filterableFields map[string]bool) (string, error) {
	dir := "ASC"
	numSorts := len(sorts)
	queryOrderBy := ""

	for i, fieldName := range sorts {
		if strings.HasPrefix(fieldName, "-") {
			dir = "DESC"
			fieldName = strings.TrimPrefix(fieldName, "-")
		} else {
			dir = "ASC"
		}

		_, isFilterableField := filterableFields[fieldName]
		if !isFilterableField {
			return "", fmt.Errorf("'%s' is not a valid sort field for this model", fieldName)
		}

		if (i + 1) < numSorts {
			dir = fmt.Sprintf("%s, ", dir)
		}

		// Convert dashes to underscores
		columnName := strings.Replace(fieldName, "-", "_", -1)

		queryOrderBy = fmt.Sprintf("%s%s %s", queryOrderBy, columnName, dir)
	}

	if queryOrderBy == "" {
		queryOrderBy = "id ASC"
	}

	return queryOrderBy, nil
}
