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

func iSelect(selection string) error {
	q.Select(selection)
	return nil
}

func iAddTheFROMClause(from string) error {
	q.From(from)
	return nil
}

func iSelectFrom(column, table string) error {
	q.Select(column).From(table)
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

func iAddTheWHEREClause(where string) error {
	q.Where(where)
	return nil
}

func iAddTheJoinOn(join, on string) error {
	q.Join(join, on)
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I create a new Query$`, iCreateANewQuery)
	s.Step(`^I select "([^"]*)" from "([^"]*)"$`, iSelectFrom)
	s.Step(`^I add the WHERE clause "([^"]*)"$`, iAddTheWHEREClause)
	s.Step(`^I compile the Query$`, iCompileTheQuery)
	s.Step(`^the SQL should match "([^"]*)"$`, theSQLShouldMatch)
	s.Step(`^I select "([^"]*)"$`, iSelect)
	s.Step(`^I add the FROM clause "([^"]*)"$`, iAddTheFROMClause)
	s.Step(`^I add the join "([^"]*)" on "([^"]*)"$`, iAddTheJoinOn)
}
