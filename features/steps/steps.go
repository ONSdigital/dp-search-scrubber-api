package steps

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	c.apiFeature.RegisterSteps(ctx)

	ctx.Step(`^I should receive a scrubber search empty response$`, c.iShouldReceiveAScrubberSearchEmptyResponse)
	ctx.Step(`^I should receive a scrubber search response with OAC codes populated$`, c.iShouldReceiveAScrubberSearchOACResponse)
	ctx.Step(`^I should receive a scrubber search response with Industry codes populated$`, c.iShouldReceiveAScrubberSearchIndustryResponse)
	ctx.Step(`^I should receive a scrubber search response full response$`, c.iShouldReceiveAScrubberSearchFullResponse)
}

func (c *Component) iShouldReceiveAScrubberSearchFullResponse() error {
	responseBody := c.apiFeature.HttpResponse.Body
	body, _ := ioutil.ReadAll(responseBody)

	trimmedBody, err := removeTimeParameter(string(body))
	if err != nil {
		return c.StepError()
	}

	assert.Equal(c, `{"query":"26513 W00009754","results":{"areas":[{"codes":{"W00009754":"W00009754"},"name":"Cardiff","region":"Wales","region_code":"W92000004"}],"industries":[{"code":"26513","name":"Manufacture of non-electronic measuring, testing etc. equipment, not for industrial process control"}]}}`, strings.TrimSpace(trimmedBody))

	return c.StepError()
}

func (c *Component) iShouldReceiveAScrubberSearchIndustryResponse() error {
	responseBody := c.apiFeature.HttpResponse.Body
	body, _ := ioutil.ReadAll(responseBody)

	trimmedBody, err := removeTimeParameter(string(body))
	if err != nil {
		return c.StepError()
	}

	assert.Equal(c, `{"query":"26513","results":{"industries":[{"code":"26513","name":"Manufacture of non-electronic measuring, testing etc. equipment, not for industrial process control"}]}}`, strings.TrimSpace(trimmedBody))

	return c.StepError()
}

func (c *Component) iShouldReceiveAScrubberSearchOACResponse() error {
	responseBody := c.apiFeature.HttpResponse.Body
	body, _ := ioutil.ReadAll(responseBody)

	trimmedBody, err := removeTimeParameter(string(body))
	if err != nil {
		return c.StepError()
	}

	assert.Equal(c, `{"query":"W00009754","results":{"areas":[{"codes":{"W00009754":"W00009754"},"name":"Cardiff","region":"Wales","region_code":"W92000004"}]}}`, strings.TrimSpace(trimmedBody))

	return c.StepError()
}

func (c *Component) iShouldReceiveAScrubberSearchEmptyResponse() error {
	responseBody := c.apiFeature.HttpResponse.Body
	body, _ := ioutil.ReadAll(responseBody)

	trimmedBody, err := removeTimeParameter(string(body))
	if err != nil {
		return c.StepError()
	}

	assert.Equal(c, `{"query":"dentists","results":{}}`, strings.TrimSpace(trimmedBody))

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
