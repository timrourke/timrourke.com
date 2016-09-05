package main

import (
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"net/http"
	"net/http/httptest"
	"strings"
)

type apiFeature struct {
	resp *httptest.ResponseRecorder
}

func (a *apiFeature) resetResponse(interface{}) {
	a.resp = httptest.NewRecorder()
}

func (a *apiFeature) iSendRequestTo(method, endpoint string) error {
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return err
	}
	initRouter(initDB()).ServeHTTP(a.resp, req)

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
		data             interface{}
		err              error
	)

	// Attempt to unmarshal json to verify its validity from feature file
	if err = json.Unmarshal([]byte(expectedJson.Content), &data); err != nil {
		return err
	}

	// Remarshal json so its format is standardized
	if expected, err = json.Marshal(data); err != nil {
		return err
	}

	actual = a.resp.Body.Bytes()
	if string(actual) != string(expected) {
		return fmt.Errorf("expected json %s, does not match actual: %s",
			string(expected),
			string(actual))
	}

	return nil
}

func FeatureContext(s *godog.Suite) {
	api := &apiFeature{}

	s.BeforeScenario(api.resetResponse)

	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`, api.iSendRequestTo)
	s.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	s.Step(`^the response should match text "([^"]*)"$`, api.theResponseShouldMatchText)
	s.Step(`^the response should match json:$`, api.theResponseShouldMatchJson)
}
