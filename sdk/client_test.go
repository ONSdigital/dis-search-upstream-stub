package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/ONSdigital/dis-search-upstream-stub/models"
	healthcheck "github.com/ONSdigital/dp-api-clients-go/v2/health"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	c "github.com/smartystreets/goconvey/convey"
)

const testHost = "http://localhost:23900"

var (
	initialTestState = healthcheck.CreateCheckState(service)
)

func getMockResponse() models.Resources {
	items := make([]models.Resource, 1)

	mockItem := models.Resource{
		URI:             "http://www.ons.gov.uk/economy",
		URIOld:          "http://www.ons.gov.uk/economy",
		ContentType:     "bulletin",
		CDID:            "MM23",
		DatasetID:       "MRET",
		Edition:         "January",
		MetaDescription: "Some meta text",
		ReleaseDate:     "2024-11-21:20:14Z",
		Summary:         "A summary",
		Title:           "Bulletin title",
		Topics:          []string{"2213"},
		Language:        "en",
		Survey:          "Some survey text",
		CanonicalTopic:  "2213",
	}

	items = append(items, mockItem)

	mockResourcesResponse := models.Resources{
		Count:      1,
		TotalCount: 1,
		Limit:      10,
		Offset:     0,
		Items:      items,
	}
	return mockResourcesResponse
}

func TestHealthCheckerClient(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	timePriorHealthCheck := time.Now().UTC()
	path := "/health"

	c.Convey("Given clienter.Do returns an error", t, func() {
		clientError := errors.New("unexpected error")
		httpClient := newMockHTTPClient(&http.Response{}, clientError)
		upstreamAPIClient := newUpstreamAPIClient(t, httpClient)
		check := initialTestState

		c.Convey("When upstream API client Checker is called", func() {
			err := upstreamAPIClient.Checker(ctx, &check)
			c.So(err, c.ShouldBeNil)

			c.Convey("Then the expected check is returned", func() {
				c.So(check.Name(), c.ShouldEqual, service)
				c.So(check.Status(), c.ShouldEqual, health.StatusCritical)
				c.So(check.StatusCode(), c.ShouldEqual, 0)
				c.So(check.Message(), c.ShouldEqual, clientError.Error())
				c.So(*check.LastChecked(), c.ShouldHappenAfter, timePriorHealthCheck)
				c.So(check.LastSuccess(), c.ShouldBeNil)
				c.So(*check.LastFailure(), c.ShouldHappenAfter, timePriorHealthCheck)
			})

			c.Convey("And client.Do should be called once with the expected parameters", func() {
				doCalls := httpClient.DoCalls()
				c.So(doCalls, c.ShouldHaveLength, 1)
				c.So(doCalls[0].Req.URL.Path, c.ShouldEqual, path)
			})
		})
	})

	c.Convey("Given a 500 response for health check", t, func() {
		httpClient := newMockHTTPClient(&http.Response{StatusCode: http.StatusInternalServerError}, nil)
		upstreamAPIClient := newUpstreamAPIClient(t, httpClient)
		check := initialTestState

		c.Convey("When upstream API client Checker is called", func() {
			err := upstreamAPIClient.Checker(ctx, &check)
			c.So(err, c.ShouldBeNil)

			c.Convey("Then the expected check is returned", func() {
				c.So(check.Name(), c.ShouldEqual, service)
				c.So(check.Status(), c.ShouldEqual, health.StatusCritical)
				c.So(check.StatusCode(), c.ShouldEqual, 500)
				c.So(check.Message(), c.ShouldEqual, service+healthcheck.StatusMessage[health.StatusCritical])
				c.So(*check.LastChecked(), c.ShouldHappenAfter, timePriorHealthCheck)
				c.So(check.LastSuccess(), c.ShouldBeNil)
				c.So(*check.LastFailure(), c.ShouldHappenAfter, timePriorHealthCheck)
			})

			c.Convey("And client.Do should be called once with the expected parameters", func() {
				doCalls := httpClient.DoCalls()
				c.So(doCalls, c.ShouldHaveLength, 1)
				c.So(doCalls[0].Req.URL.Path, c.ShouldEqual, path)
			})
		})
	})
}

func TestGetResources(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c.Convey("Given request to get resources", t, func() {
		body, err := json.Marshal(getMockResponse())
		if err != nil {
			t.Errorf("failed to setup test data, error: %v", err)
		}

		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			},
			nil)

		upstreamAPIClient := newUpstreamAPIClient(t, httpClient)

		c.Convey("When GetResources is called", func() {
			resp, err := upstreamAPIClient.GetResources(ctx, Options{})

			c.Convey("Then the expected response body is returned", func() {
				c.So(*resp, c.ShouldResemble, getMockResponse())

				c.Convey("And no error is returned", func() {
					c.So(err, c.ShouldBeNil)

					c.Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						c.So(doCalls, c.ShouldHaveLength, 1)
						c.So(doCalls[0].Req.Method, c.ShouldEqual, "GET")
						c.So(doCalls[0].Req.URL.Path, c.ShouldEqual, "/resources")
						c.So(doCalls[0].Req.Header["Authorization"], c.ShouldBeEmpty)
					})
				})
			})
		})
	})
}

func newMockHTTPClient(r *http.Response, err error) *dphttp.ClienterMock {
	return &dphttp.ClienterMock{
		SetPathsWithNoRetriesFunc: func(paths []string) {
			// Mocked function is called but do nothing
		},
		DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return r, err
		},
		GetPathsWithNoRetriesFunc: func() []string {
			return []string{"/healthcheck"}
		},
	}
}

func newUpstreamAPIClient(_ *testing.T, httpClient *dphttp.ClienterMock) *Client {
	healthClient := healthcheck.NewClientWithClienter(service, testHost, httpClient)
	return NewWithHealthClient(healthClient, "/resources")
}
