package storage

// QueryParams defines the struct of query params to use for SQL constraints
type QueryParams struct {
	Limit   uint64
	Offset  uint64
	OrderBy string
	Where   map[string]string
}
