package steps

import (
	_ "embed"
	"io"
	"strings"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

//go:embed example_resources.json
var expectedResponseBody []byte

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	c.apiFeature.RegisterSteps(ctx)

	ctx.Step(`^I should receive a hello-world response$`, c.iShouldReceiveAHelloworldResponse)
	ctx.Step(`^I should receive a list of resources$`, c.iShouldReceiveAListOfResources)
}

func (c *Component) iShouldReceiveAHelloworldResponse() error {
	responseBody := c.apiFeature.HTTPResponse.Body
	body, _ := io.ReadAll(responseBody)

	assert.Equal(c, `{"message":"Hello, World!"}`, strings.TrimSpace(string(body)))

	return c.StepError()
}

func (c *Component) iShouldReceiveAListOfResources() error {
	responseBody := c.apiFeature.HTTPResponse.Body
	body, _ := io.ReadAll(responseBody)

	assert.Equal(c, strings.TrimSpace(string(expectedResponseBody)), strings.TrimSpace(string(body)))

	return c.StepError()
}
