package steps

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ONSdigital/dis-search-upstream-stub/models"
	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

//go:embed example_resources.json
var expectedResponseBody []byte

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	c.apiFeature.RegisterSteps(ctx)

	ctx.Step(`^I should receive a list of resources$`, c.iShouldReceiveAListOfResources)
}

func (c *Component) iShouldReceiveAListOfResources() error {
	var expectedResponse, actualResponse models.Resources

	responseBody := c.apiFeature.HTTPResponse.Body
	body, _ := io.ReadAll(responseBody)

	err := json.Unmarshal(expectedResponseBody, &expectedResponse)
	if err != nil {
		return fmt.Errorf("failed to unmarshal expected results from file - error: %v", err)
	}

	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		return fmt.Errorf("failed to unmarshal actual response from server - error: %v", err)
	}

	assert.Equal(c, expectedResponse, actualResponse)

	return c.StepError()
}
