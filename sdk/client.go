package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ONSdigital/dis-search-upstream-stub/models"
	apiError "github.com/ONSdigital/dis-search-upstream-stub/sdk/errors"
	healthcheck "github.com/ONSdigital/dp-api-clients-go/v2/health"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
)

const (
	service = "dis-search-upstream-stub"
)

type Client struct {
	hcCli             *healthcheck.Client
	resourcesEndpoint string
}

// New creates a new instance of Client with a given upstream api url
func New(upstreamAPIURL, resourcesEndpoint string) *Client {
	return &Client{
		hcCli:             healthcheck.NewClient(service, upstreamAPIURL),
		resourcesEndpoint: resourcesEndpoint,
	}
}

// NewWithHealthClient creates a new instance of upstream API Client,
// reusing the URL and Clienter from the provided healthcheck client
func NewWithHealthClient(hcCli *healthcheck.Client, resourcesEndpoint string) *Client {
	return &Client{
		hcCli:             healthcheck.NewClientWithClienter(service, hcCli.URL, hcCli.Client),
		resourcesEndpoint: resourcesEndpoint,
	}
}

// URL returns the URL used by this client
func (cli *Client) URL() string {
	return cli.hcCli.URL
}

func (cli *Client) ResourcesEndpoint() string {
	return cli.resourcesEndpoint
}

// Health returns the underlying Healthcheck Client for this upstream API client
func (cli *Client) Health() *healthcheck.Client {
	return cli.hcCli
}

// Checker calls upstream api health endpoint and returns a check object to the caller
func (cli *Client) Checker(ctx context.Context, check *health.CheckState) error {
	return cli.hcCli.Checker(ctx, check)
}

// GetResources gets a list of upstream resources
func (cli *Client) GetResources(ctx context.Context, options Options) (*models.Resources, apiError.Error) {
	path := fmt.Sprintf("%s%s", cli.hcCli.URL, cli.resourcesEndpoint)
	if options.Query != nil {
		path = path + "?" + options.Query.Encode()
	}

	respInfo, apiErr := cli.callUpstreamAPI(ctx, path, http.MethodGet, options.Headers, nil)
	if apiErr != nil {
		return nil, apiErr
	}

	var resourcesResponse models.Resources

	if err := json.Unmarshal(respInfo.Body, &resourcesResponse); err != nil {
		return nil, apiError.StatusError{
			Err: fmt.Errorf("failed to unmarshal upstream resources response - error is: %v", err),
		}
	}

	fmt.Println(&resourcesResponse)

	return &resourcesResponse, nil
}

type ResponseInfo struct {
	Body    []byte
	Headers http.Header
	Status  int
}

// callUpstreamAPI calls the Upstream API endpoint given by path for the provided REST method, request headers, and body payload.
// It returns the response body and any error that occurred.
func (cli *Client) callUpstreamAPI(ctx context.Context, path, method string, headers http.Header, payload []byte) (*ResponseInfo, apiError.Error) {
	URL, err := url.Parse(path)
	if err != nil {
		return nil, apiError.StatusError{
			Err: fmt.Errorf("failed to parse path: \"%v\" error is: %v", path, err),
		}
	}

	path = URL.String()

	var req *http.Request

	if payload != nil {
		req, err = http.NewRequest(method, path, bytes.NewReader(payload))
	} else {
		req, err = http.NewRequest(method, path, http.NoBody)
	}

	// check req, above, didn't error
	if err != nil {
		return nil, apiError.StatusError{
			Err: fmt.Errorf("failed to create request for call to upstream api, error is: %v", err),
		}
	}

	// set any headers against request
	setHeaders(req, headers)

	if payload != nil {
		req.Header.Add("Content-type", "application/json")
	}

	resp, err := cli.hcCli.Client.Do(ctx, req)
	if err != nil {
		return nil, apiError.StatusError{
			Err:  fmt.Errorf("failed to call upstream api, error is: %v", err),
			Code: http.StatusInternalServerError,
		}
	}
	defer func() {
		err = closeResponseBody(resp)
	}()

	respInfo := &ResponseInfo{
		Headers: resp.Header.Clone(),
		Status:  resp.StatusCode,
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 400 {
		return respInfo, apiError.StatusError{
			Err:  fmt.Errorf("failed as unexpected code from upstream api: %v", resp.StatusCode),
			Code: resp.StatusCode,
		}
	}

	if resp.Body == nil {
		return respInfo, nil
	}

	respInfo.Body, err = io.ReadAll(resp.Body)
	if err != nil {
		return respInfo, apiError.StatusError{
			Err:  fmt.Errorf("failed to read response body from call to upstream api, error is: %v", err),
			Code: resp.StatusCode,
		}
	}
	return respInfo, nil
}

// closeResponseBody closes the response body and logs an error if unsuccessful
func closeResponseBody(resp *http.Response) apiError.Error {
	if resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			return apiError.StatusError{
				Err:  fmt.Errorf("error closing http response body from call to upstream api, error is: %v", err),
				Code: http.StatusInternalServerError,
			}
		}
	}

	return nil
}
