package main

import (
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

var test_db *sqlx.DB

func init() {
	test_db = initDB()
}

type apiFeature struct {
	resp *httptest.ResponseRecorder
}

func (a *apiFeature) resetResponse(interface{}) {
	a.resp = httptest.NewRecorder()

	_ = test_db.MustExec("SET FOREIGN_KEY_CHECKS=0")
	_ = test_db.MustExec("TRUNCATE TABLE `users`")
	_ = test_db.MustExec("SET FOREIGN_KEY_CHECKS=1")
}

func (a *apiFeature) iSendRequestTo(method, endpoint string) error {
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return err
	}
	initRouter(test_db).ServeHTTP(a.resp, req)

	// handle panic
	defer func() {
		switch t := recover().(type) {
		case string:
			err = fmt.Errorf(t)
		case error:
			err = t
		}
	}()

	return err
}

func (a *apiFeature) theResponseCodeShouldBe(expectedStatus int) error {
	actual := a.resp.Code

	if actual != expectedStatus {
		return fmt.Errorf("expected status code to be %d, but it was %d",
			expectedStatus,
			actual)
	}

	return nil
}

func (a *apiFeature) theResponseShouldMatchText(expectedResponseText string) error {
	actual := strings.TrimSpace(a.resp.Body.String())

	if actual != expectedResponseText {
		return fmt.Errorf("expected response text to be %s, but it was %s",
			expectedResponseText,
			actual)
	}

	return nil
}

func (a *apiFeature) theResponseShouldMatchJson(expectedJson *gherkin.DocString) error {
	var (
		expected, actual []byte
		expectedData     interface{}
		actualData       interface{}
		err              error
	)

	// Attempt to unmarshal expected json to verify its validity from feature file
	if err = json.Unmarshal([]byte(expectedJson.Content), &expectedData); err != nil {
		return err
	}

	// Remarshal expected json so its format is standardized
	if expected, err = json.Marshal(expectedData); err != nil {
		return err
	}

	// Attempt to unmarshal actual json to verify its validity from feature file
	if err = json.Unmarshal([]byte(a.resp.Body.Bytes()), &actualData); err != nil {
		return err
	}

	// Remarshal actual json so its format is standardized
	if actual, err = json.Marshal(actualData); err != nil {
		return err
	}

	if string(actual) != string(expected) {
		return fmt.Errorf("expected json %s, does not match actual: %s",
			string(expected),
			string(actual))
	}

	return nil
}

func (a *apiFeature) thereAreUsers(users *gherkin.DataTable) error {
	var fields []string
	var marks []string
	head := users.Rows[0].Cells
	for _, cell := range head {
		fields = append(fields, cell.Value)
		marks = append(marks, "?")
	}

	stmt, err := test_db.Preparex("INSERT INTO users (" + strings.Join(fields, ", ") + ") VALUES(" + strings.Join(marks, ", ") + ")")

	if err != nil {
		return err
	}

	fmt.Printf("%+v", &users.Rows)

	for i := 1; i < len(users.Rows); i++ {
		var vals []interface{}
		for n, cell := range users.Rows[i].Cells {
			switch head[n].Value {
			case "id":
				vals = append(vals, cell.Value)
			case "username":
				vals = append(vals, cell.Value)
			case "email":
				vals = append(vals, cell.Value)
			case "created_at":
				createdAt, err := time.Parse(time.RFC3339, cell.Value)
				if err != nil {
					return err
				}

				vals = append(vals, createdAt)
			case "updated_at":
				updatedAt, err := time.Parse(time.RFC3339, cell.Value)
				if err != nil {
					return err
				}

				vals = append(vals, updatedAt)
			case "password_hash":
				vals = append(vals, cell.Value)
			default:
				return fmt.Errorf("unexpected column name: %s", head[n].Value)
			}
		}
		if _, err = stmt.Exec(vals...); err != nil {
			return err
		}
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	api := &apiFeature{}

	s.BeforeScenario(api.resetResponse)

	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`,
		api.iSendRequestTo)
	s.Step(`^the response code should be (\d+)$`,
		api.theResponseCodeShouldBe)
	s.Step(`^the response should match text "([^"]*)"$`,
		api.theResponseShouldMatchText)
	s.Step(`^the response should match json:$`,
		api.theResponseShouldMatchJson)
	s.Step(`^there are users:$`,
		api.thereAreUsers)
}
