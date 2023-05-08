package steps

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	c.apiFeature.RegisterSteps(ctx)

	ctx.Step(`^the response body is the same as the json in "([^"]*)"$`, c.theResponseBodyIsTheSameAsTheJsonIn)
}

func (c *Component) theResponseBodyIsTheSameAsTheJsonIn(expected string) error {
	responseBody := c.apiFeature.HttpResponse.Body
	body, _ := io.ReadAll(responseBody)

	content, err := os.ReadFile(expected)
	if err != nil {
		return err
	}

	str := strings.ReplaceAll(string(content), "\n", "")
	str = strings.ReplaceAll(str, " ", "")

	trimmedBody, err := removeTimeParameter(string(body))
	if err != nil {
		return c.StepError()
	}

	assert.Equal(c, str, strings.ReplaceAll(trimmedBody, " ", ""))

	return c.StepError()
}

func removeTimeParameter(responseBody string) (string, error) {
	var data map[string]interface{}

	if err := json.Unmarshal([]byte(responseBody), &data); err != nil {
		return "", err
	}

	delete(data, "time")

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(jsonBytes)), nil
}
