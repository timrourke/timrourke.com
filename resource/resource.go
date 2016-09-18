package resource

import (
	"fmt"
	"github.com/manyminds/api2go"
	"github.com/timrourke/timrourke.com/query"
	"strconv"
	"strings"
)

const defaultPaginationLimit = 20

// RelationshipFunc defines the function type for modifying a Query to select
// a relationship
type RelationshipFunc func(api2go.Request, *query.Query) *query.Query

// ParseQueryParams parses request for query params
func ParseQueryParams(r api2go.Request, filterableFields map[string]bool, relationshipsByParam map[string]RelationshipFunc) (*query.Query, error) {
	var (
		err          error
		q            *query.Query
		queryLimit   uint64 = defaultPaginationLimit
		queryOffset  uint64
		queryOrderBy = "id ASC"
	)

	requestParams := r.QueryParams

	q = query.New()

	offset, hasOffset := requestParams["page[offset]"]
	limit, hasLimit := requestParams["page[limit]"]
	pageNum, hasPageNum := requestParams["page[number]"]
	pageSize, hasPageSize := requestParams["page[size]"]
	sorts, hasSorts := requestParams["sort"]

	if hasLimit {
		parsedLimit, err := strconv.ParseUint(limit[0], 10, 64)
		if err != nil {
			return q, err
		}

		queryLimit = parsedLimit
	} else if hasPageSize {
		parsedPageSize, err := strconv.ParseUint(pageSize[0], 10, 64)
		if err != nil {
			return q, err
		}

		queryLimit = parsedPageSize
	}

	if queryLimit > 100 {
		queryLimit = 100
	}

	if hasPageNum {
		parsedPageNum, err := strconv.ParseUint(pageNum[0], 10, 64)
		if err != nil {
			return q, nil
		}

		queryOffset = (parsedPageNum - 1) * queryLimit
	} else if hasOffset {
		parsedOffset, err := strconv.ParseUint(offset[0], 10, 64)
		if err != nil {
			return q, nil
		}

		queryOffset = parsedOffset
	}

	// Modify query for filters
	handleFilters(requestParams, filterableFields, q)

	// Modify query for any relationships passed as query params. These params
	// are generally provided by api2go when fetching relationships, for example
	// GET /users/1/posts
	for _, relationship := range relationshipsByParam {
		relationship(r, q)
	}

	// Modify query for sorting
	if hasSorts && len(sorts[0]) > 0 {
		queryOrderBy, err = handleSorts(sorts, filterableFields)
		if err != nil {
			return q, err
		}
	}

	q.OrderBy(queryOrderBy)
	q.Limit(queryOffset, queryLimit)

	return q, nil
}

func handleFilters(params map[string][]string, filterableFields map[string]bool, q *query.Query) {
	for filterColumn, strictEquals := range filterableFields {
		filterValue, ok := params[fmt.Sprintf("filter[%s]", filterColumn)]

		if !ok || len(filterValue[0]) == 0 {
			continue
		}

		if strictEquals {
			q.Where(fmt.Sprintf("%s = :%s", filterColumn, filterColumn))
		} else {
			q.Where(fmt.Sprintf("%s LIKE \"%:%s%\"", filterColumn, filterColumn))
		}

		q.Bind(fmt.Sprintf("%s", filterColumn), filterValue[0])
	}
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
