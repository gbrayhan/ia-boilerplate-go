package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cucumber/godog"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	resp *http.Response
	body []byte
	err  error
	base = "http://localhost:8080"
)

func theServiceIsInitialized() error {
	time.Sleep(2 * time.Second)
	return nil
}

func iSendAGETRequestTo(path string) error {
	resp, err = http.Get(base + path)
	return err
}

func iSendAPOSTRequestToWithBody(path string, payload *godog.DocString) error {
	resp, err = http.Post(base+path, "application/json", bytes.NewBufferString(payload.Content))
	return err
}

func theResponseCodeShouldBe(code int) error {
	if resp.StatusCode != code {
		return fmt.Errorf("expected status code %d but got %d", code, resp.StatusCode)
	}
	body, err = ioutil.ReadAll(resp.Body)
	return err
}

func theJSONResponseShouldContain(field, value string) error {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}
	if data[field] != value {
		return fmt.Errorf("expected %q = %q, but got %v", field, value, data[field])
	}
	return nil
}

func theJSONResponseShouldContainKey(key string) error {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}
	if _, ok := data[key]; !ok {
		return fmt.Errorf("expected JSON to contain key %q", key)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^the service is initialized$`, theServiceIsInitialized)
	ctx.Step(`^I send a GET request to "([^"]*)"$`, iSendAGETRequestTo)
	ctx.Step(`^I send a POST request to "([^"]*)" with body:$`, iSendAPOSTRequestToWithBody)
	ctx.Step(`^the response code should be (\d+)$`, theResponseCodeShouldBe)
	ctx.Step(`^the JSON response should contain "([^"]*)": "([^"]*)"$`, theJSONResponseShouldContain)
	ctx.Step(`^the JSON response should contain key "([^"]*)"$`, theJSONResponseShouldContainKey)
}
