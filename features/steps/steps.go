package steps

import (
	"encoding/json"
	"io"
	"os"

	"github.com/ONSdigital/dp-search-scrubber-api/models"
	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	c.apiFeature.RegisterSteps(ctx)

	ctx.Step(`^the response body is the same as the json in "([^"]*)"$`, c.theResponseBodyIsTheSameAsTheJSONIn)
}

func (c *Component) theResponseBodyIsTheSameAsTheJSONIn(expectedFile string) error {
	responseBody := c.apiFeature.HTTPResponse.Body
	actualRawContent, _ := io.ReadAll(responseBody)

	var expected models.ScrubberResp
	var actual models.ScrubberResp

	expectedRawContent, err := os.ReadFile(expectedFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(expectedRawContent, &expected)
	if err != nil {
		return err
	}

	err = json.Unmarshal(actualRawContent, &actual)
	if err != nil {
		return err
	}

	expected.Time = actual.Time // Workaround for the time the request took
	assert.Equal(c, expected.Query, actual.Query)

	return c.StepError()
}
