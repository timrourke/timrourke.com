package query

import (
	"fmt"
	"strings"
)

type Query struct {
	Selects  []string
	Froms    []string
	Joins    map[string]string
	Conds    []string
	OrderBys []string
	Values   map[string]interface{}
}

func New() *Query {
	return &Query{}
}

func (q *Query) Compile() (string, map[string]interface{}) {
	selects := strings.Join(q.Selects, ", ")
	if len(selects) == 0 {
		selects = "*"
	}

	froms := strings.Join(q.Froms, ", ")

	joinsSlice := make([]string, 0)
	for join, on := range q.Joins {
		joinsSlice = append(joinsSlice, fmt.Sprintf("%s ON %s", join, on))
	}
	joins := strings.Join(joinsSlice, ", ")

	conds := strings.Join(q.Conds, " AND ")
	if len(conds) > 0 {
		conds = fmt.Sprintf("AND %s", conds)
	}

	orders := strings.Join(q.OrderBys, ", ")
	if len(orders) == 0 {
		orders = "id ASC"
	}

	sql := "SELECT %s FROM %s %s WHERE 1 %s ORDER BY %s LIMIT :limit,:offset"
	return fmt.Sprintf(sql,
		selects,
		froms,
		joins,
		conds,
		orders), q.Values
}

func (q *Query) Select(selection string) *Query {
	q.Selects = append(q.Selects, selection)
	return q
}

func (q *Query) Limit(offset uint64, limit uint64) *Query {
	q.Values["offset"] = offset
	q.Values["limit"] = limit
	return q
}

func (q *Query) Where(where string) *Query {
	q.Conds = append(q.Conds, where)
	return q
}

func (q *Query) Bind(key string, value interface{}) *Query {
	q.Values[key] = value
	return q
}

func (q *Query) OrderBy(orderBy string) *Query {
	q.OrderBys = append(q.OrderBys, orderBy)
	return q
}

func (q *Query) From(from string) *Query {
	q.Froms = append(q.Froms, from)
	return q
}
