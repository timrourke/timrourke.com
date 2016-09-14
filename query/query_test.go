package query

import (
	"fmt"
	"github.com/DATA-DOG/godog"
	//"strings"
)

var (
	q      *Query
	sql    string
	values map[string]interface{}
)

func iCreateANewQuery() error {
	q = New()
	return nil
}

func iSelectFrom(column, table string) error {
	q.Select("t.*").From("dinosaurs t")
	return nil
}

func iCompileTheQuery() error {
	sql, values = q.Compile()
	return nil
}

func theSQLShouldMatch(expectedSql string) error {
	if expectedSql == sql {
		return nil
	}
	return fmt.Errorf("expected SQL '%s' did not match actual '%s'",
		expectedSql,
		sql)
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I create a new Query$`, iCreateANewQuery)
	s.Step(`^I select "([^"]*)" from "([^"]*)"$`, iSelectFrom)
	s.Step(`^I compile the Query$`, iCompileTheQuery)
	s.Step(`^the SQL should match "([^"]*)"$`, theSQLShouldMatch)
}
