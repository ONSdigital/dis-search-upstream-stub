package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ONSdigital/dis-search-upstream-stub/api"
	apiMock "github.com/ONSdigital/dis-search-upstream-stub/api/mock"
	"github.com/ONSdigital/dis-search-upstream-stub/config"
	"github.com/ONSdigital/dis-search-upstream-stub/data"
	"github.com/ONSdigital/dis-search-upstream-stub/models"
	dpresponse "github.com/ONSdigital/dp-net/v2/handlers/response"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

// Constants for testing
const (
	expectedServerErrorMsg  = "internal server error"
	expectedOffsetErrorMsg  = "invalid offset query parameter"
	expectedLimitErrorMsg   = "invalid limit query parameter"
	expectedLimitOverMaxMsg = "limit query parameter is larger than the maximum allowed"
)

// expectedStandardResource returns a release resource that can be used to define and test expected values within it
func expectedStandardResource(uri string) models.Resource {
	standardResource := models.Resource{
		URI:             uri,
		URIOld:          "/an/old/uri",
		ContentType:     "api_dataset_landing_page",
		CDID:            "A321B",
		DatasetID:       "ASELECTIONOFNUMBERSANDLETTERS456",
		Edition:         "an edition",
		MetaDescription: "a description",
		Summary:         "a summary",
		Title:           "a title",
		Language:        "string",
		Survey:          "string",
		CanonicalTopic:  "string",
	}

	expectedResource := standardResource
	topics := []string{"a", "b", "c", "d"}
	expectedResource.Topics = topics

	return expectedResource
}

// expectedReleaseResource returns a release resource that can be used to define and test expected values within it
func expectedReleaseResource(uri string) models.Resource {
	releaseResource := models.Resource{
		URI:             uri,
		URIOld:          "/an/old/uri",
		ContentType:     "api_dataset_landing_page",
		CDID:            "A321B",
		DatasetID:       "ASELECTIONOFNUMBERSANDLETTERS456",
		Edition:         "an edition",
		MetaDescription: "a description",
		ReleaseDate:     "2024-11-21:20:14Z",
		Summary:         "a summary",
		Title:           "a title",
		Language:        "string",
		Survey:          "string",
		CanonicalTopic:  "string",
		Release: models.Release{
			Cancelled:       true,
			Finalised:       true,
			Published:       true,
			ProvisionalDate: "October-November 2024",
		},
	}

	expectedResource := releaseResource
	topics := []string{"a", "b", "c", "d"}
	expectedResource.Topics = topics
	dateChanges := []string{"a change_notice", "a previous_date"}
	expectedResource.Release.DateChanges = dateChanges

	return expectedResource
}

func expectedResources(limit, offset int) models.Resources {
	resources := models.Resources{
		Limit:      limit,
		Offset:     offset,
		TotalCount: 2,
	}

	firstResource := expectedStandardResource("/a/uri")
	secondResource := expectedReleaseResource("/another/uri")

	if (offset == 0) && (limit > 1) {
		resources.Count = 2
		resources.Items = []models.Resource{firstResource, secondResource}
	}

	if (offset == 1) && (limit > 0) {
		resources.Count = 1
		resources.Items = []models.Resource{secondResource}
	}

	return resources
}

func TestGetResourcesHandlerSuccess(t *testing.T) {
	t.Parallel()

	cfg, configErr := config.Get()
	if configErr != nil {
		t.Errorf("failed to retrieve default configuration, error: %v", configErr)
	}

	dataStorerMock := &apiMock.DataStorerMock{
		GetResourcesFunc: func(ctx context.Context, options data.Options) (*models.Resources, error) {
			resources := expectedResources(cfg.DefaultLimit, cfg.DefaultOffset)
			return &resources, nil
		},
	}

	Convey("Given a list of resources exists in the Data Store", t, func() {
		apiInstance := api.Setup(mux.NewRouter(), cfg, dataStorerMock)

		Convey("When a request is made to get a list of all resources", func() {
			req := httptest.NewRequest("GET", "http://localhost:29600/resources", http.NoBody)
			resp := httptest.NewRecorder()
			apiInstance.Router.ServeHTTP(resp, req)

			Convey("Then a list of resources is returned with status code 200", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)

				payload, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("failed to read payload with io.ReadAll, error: %v", err)
				}

				resourcesReturned := models.Resources{}
				err = json.Unmarshal(payload, &resourcesReturned)
				So(err, ShouldBeNil)

				expectedResource1 := expectedStandardResource("/a/uri")
				expectedResource2 := expectedReleaseResource("/another/uri")

				Convey("And the returned list should contain expected resources", func() {
					returnedResourceList := resourcesReturned.Items
					So(returnedResourceList, ShouldHaveLength, 2)
					returnedResource1 := returnedResourceList[0]
					So(returnedResource1.URI, ShouldEqual, expectedResource1.URI)
					So(returnedResource1.URIOld, ShouldEqual, expectedResource1.URIOld)
					So(returnedResource1.ContentType, ShouldEqual, expectedResource1.ContentType)
					So(returnedResource1.CDID, ShouldEqual, expectedResource1.CDID)
					So(returnedResource1.DatasetID, ShouldEqual, expectedResource1.DatasetID)
					So(returnedResource1.Edition, ShouldEqual, expectedResource1.Edition)
					So(returnedResource1.MetaDescription, ShouldEqual, expectedResource1.MetaDescription)
					So(returnedResource1.ReleaseDate, ShouldEqual, expectedResource1.ReleaseDate)
					So(returnedResource1.Summary, ShouldEqual, expectedResource1.Summary)
					So(returnedResource1.Title, ShouldEqual, expectedResource1.Title)
					So(returnedResource1.Topics, ShouldEqual, expectedResource1.Topics)
					So(returnedResource1.Language, ShouldEqual, expectedResource1.Language)
					So(returnedResource1.Survey, ShouldEqual, expectedResource1.Survey)
					So(returnedResource1.CanonicalTopic, ShouldEqual, expectedResource1.CanonicalTopic)
					returnedResource2 := returnedResourceList[1]
					So(returnedResource2.URIOld, ShouldEqual, expectedResource2.URIOld)
					So(returnedResource2.ContentType, ShouldEqual, expectedResource2.ContentType)
					So(returnedResource2.CDID, ShouldEqual, expectedResource2.CDID)
					So(returnedResource2.DatasetID, ShouldEqual, expectedResource2.DatasetID)
					So(returnedResource2.Edition, ShouldEqual, expectedResource2.Edition)
					So(returnedResource2.MetaDescription, ShouldEqual, expectedResource2.MetaDescription)
					So(returnedResource2.ReleaseDate, ShouldEqual, expectedResource2.ReleaseDate)
					So(returnedResource2.Summary, ShouldEqual, expectedResource2.Summary)
					So(returnedResource2.Title, ShouldEqual, expectedResource2.Title)
					So(returnedResource2.Topics, ShouldEqual, expectedResource2.Topics)
					So(returnedResource2.Language, ShouldEqual, expectedResource2.Language)
					So(returnedResource2.Survey, ShouldEqual, expectedResource2.Survey)
					So(returnedResource2.CanonicalTopic, ShouldEqual, expectedResource2.CanonicalTopic)
					So(returnedResource2.Release.Cancelled, ShouldEqual, expectedResource2.Release.Cancelled)
					So(returnedResource2.Release.Finalised, ShouldEqual, expectedResource2.Release.Finalised)
					So(returnedResource2.Release.Published, ShouldEqual, expectedResource2.Release.Published)
					So(returnedResource2.Release.DateChanges, ShouldEqual, expectedResource2.Release.DateChanges)
					So(returnedResource2.Release.ProvisionalDate, ShouldEqual, expectedResource2.Release.ProvisionalDate)
				})
			})
		})
	})

	Convey("Given valid pagination parameters", t, func() {
		validOffset := 1
		validLimit := 20

		customValidPaginationDataStore := &apiMock.DataStorerMock{
			GetResourcesFunc: func(ctx context.Context, options data.Options) (*models.Resources, error) {
				resources := expectedResources(validLimit, validOffset)
				return &resources, nil
			},
		}

		apiInstance := api.Setup(mux.NewRouter(), cfg, customValidPaginationDataStore)

		Convey("When a request is made to get a list of resources", func() {
			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:29600/resources?offset=%d&limit=%d", validOffset, validLimit), http.NoBody)
			resp := httptest.NewRecorder()

			apiInstance.Router.ServeHTTP(resp, req)

			Convey("Then a list of resources is returned with status code 200", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)

				payload, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("failed to read payload with io.ReadAll, error: %v", err)
				}

				resourcesReturned := models.Resources{}
				err = json.Unmarshal(payload, &resourcesReturned)
				So(err, ShouldBeNil)

				expectedResource := expectedReleaseResource("/another/uri")

				Convey("And the returned list should contain the expected resource", func() {
					returnedResourceList := resourcesReturned.Items
					So(returnedResourceList, ShouldHaveLength, 1)
					returnedResource := returnedResourceList[0]
					So(returnedResource.URI, ShouldEqual, expectedResource.URI)
					So(returnedResource.URIOld, ShouldEqual, expectedResource.URIOld)
					So(returnedResource.ContentType, ShouldEqual, expectedResource.ContentType)
					So(returnedResource.CDID, ShouldEqual, expectedResource.CDID)
					So(returnedResource.DatasetID, ShouldEqual, expectedResource.DatasetID)
					So(returnedResource.Edition, ShouldEqual, expectedResource.Edition)
					So(returnedResource.MetaDescription, ShouldEqual, expectedResource.MetaDescription)
					So(returnedResource.ReleaseDate, ShouldEqual, expectedResource.ReleaseDate)
					So(returnedResource.Summary, ShouldEqual, expectedResource.Summary)
					So(returnedResource.Title, ShouldEqual, expectedResource.Title)
					So(returnedResource.Topics, ShouldEqual, expectedResource.Topics)
					So(returnedResource.Language, ShouldEqual, expectedResource.Language)
					So(returnedResource.Survey, ShouldEqual, expectedResource.Survey)
					So(returnedResource.CanonicalTopic, ShouldEqual, expectedResource.CanonicalTopic)
					So(returnedResource.Release.Cancelled, ShouldEqual, expectedResource.Release.Cancelled)
					So(returnedResource.Release.Finalised, ShouldEqual, expectedResource.Release.Finalised)
					So(returnedResource.Release.Published, ShouldEqual, expectedResource.Release.Published)
					So(returnedResource.Release.DateChanges, ShouldEqual, expectedResource.Release.DateChanges)
					So(returnedResource.Release.ProvisionalDate, ShouldEqual, expectedResource.Release.ProvisionalDate)
				})
			})
		})
	})

	Convey("Given offset is greater than total number of resources in the Data Store", t, func() {
		greaterOffset := 10

		greaterOffsetDataStore := &apiMock.DataStorerMock{
			GetResourcesFunc: func(ctx context.Context, options data.Options) (*models.Resources, error) {
				resources := expectedResources(cfg.DefaultLimit, greaterOffset)
				return &resources, nil
			},
		}

		apiInstance := api.Setup(mux.NewRouter(), cfg, greaterOffsetDataStore)

		Convey("When a request is made to get a list of resources", func() {
			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:29600/resources?offset=%d", greaterOffset), http.NoBody)
			resp := httptest.NewRecorder()

			apiInstance.Router.ServeHTTP(resp, req)

			Convey("Then a list of resources is returned with status code 200", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)

				payload, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("failed to read payload with io.ReadAll, error: %v", err)
				}

				resourcesReturned := models.Resources{}
				err = json.Unmarshal(payload, &resourcesReturned)
				So(err, ShouldBeNil)

				Convey("And the returned list should be empty", func() {
					returnedResourceList := resourcesReturned.Items
					So(returnedResourceList, ShouldHaveLength, 0)
				})
			})
		})
	})
}

func TestGetResourcesHandlerWithEmptyResourceStoreSuccess(t *testing.T) {
	t.Parallel()

	cfg, err := config.Get()
	if err != nil {
		t.Errorf("failed to retrieve default configuration, error: %v", err)
	}

	Convey("Given a Search Upstream API that returns an empty list of resources", t, func() {
		dataStorerMock := &apiMock.DataStorerMock{
			GetResourcesFunc: func(ctx context.Context, options data.Options) (*models.Resources, error) {
				resources := models.Resources{}
				return &resources, nil
			},
		}

		apiInstance := api.Setup(mux.NewRouter(), cfg, dataStorerMock)

		Convey("When a request is made to get a list of all the resources that exist in the resources collection", func() {
			req := httptest.NewRequest("GET", "http://localhost:29600/resources", http.NoBody)
			resp := httptest.NewRecorder()

			apiInstance.Router.ServeHTTP(resp, req)

			Convey("Then a list of resources is returned with status code 200", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)

				payload, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("failed to read payload with io.ReadAll, error: %v", err)
				}

				resourcesReturned := models.Resources{}
				err = json.Unmarshal(payload, &resourcesReturned)
				So(err, ShouldBeNil)

				Convey("And the returned resources list should be empty", func() {
					So(resourcesReturned.Items, ShouldHaveLength, 0)
				})
			})
		})
	})
}

func TestGetResourcesHandlerFail(t *testing.T) {
	t.Parallel()

	cfg, err := config.Get()
	if err != nil {
		t.Errorf("failed to retrieve default configuration, error: %v", err)
	}

	dataStorerMock := &apiMock.DataStorerMock{
		GetResourcesFunc: func(ctx context.Context, options data.Options) (*models.Resources, error) {
			resources := expectedResources(options.Limit, options.Offset)
			return &resources, err
		},
	}

	Convey("Given offset is not numeric", t, func() {
		nonNumericOffset := "stringOffset"

		apiInstance := api.Setup(mux.NewRouter(), cfg, dataStorerMock)

		Convey("When a request is made to get a list of resources", func() {
			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:29600/resources?offset=%s", nonNumericOffset), http.NoBody)
			resp := httptest.NewRecorder()

			apiInstance.Router.ServeHTTP(resp, req)

			Convey("Then a bad request error is returned with status code 400", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
				errMsg := strings.TrimSpace(resp.Body.String())
				So(errMsg, ShouldEqual, expectedOffsetErrorMsg)

				Convey("And the response ETag header should be empty", func() {
					So(resp.Header().Get(dpresponse.ETagHeader), ShouldBeEmpty)
				})
			})
		})
	})

	Convey("Given offset is negative", t, func() {
		negativeOffset := -3

		apiInstance := api.Setup(mux.NewRouter(), cfg, dataStorerMock)

		Convey("When a request is made to get a list of resources", func() {
			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:29600/resources?offset=%d", negativeOffset), http.NoBody)
			resp := httptest.NewRecorder()

			apiInstance.Router.ServeHTTP(resp, req)

			Convey("Then a bad request error is returned with status code 400", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
				errMsg := strings.TrimSpace(resp.Body.String())
				So(errMsg, ShouldEqual, expectedOffsetErrorMsg)
			})
		})
	})

	Convey("Given limit is not numeric", t, func() {
		nonNumericLimit := "stringLimit"

		apiInstance := api.Setup(mux.NewRouter(), cfg, dataStorerMock)

		Convey("When a request is made to get a list of resources", func() {
			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:29600/resources?limit=%s", nonNumericLimit), http.NoBody)
			resp := httptest.NewRecorder()

			apiInstance.Router.ServeHTTP(resp, req)

			Convey("Then a bad request error is returned with status code 400", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
				errMsg := strings.TrimSpace(resp.Body.String())
				So(errMsg, ShouldEqual, expectedLimitErrorMsg)
			})
		})
	})

	Convey("Given limit is negative", t, func() {
		negativeLimit := -1

		apiInstance := api.Setup(mux.NewRouter(), cfg, dataStorerMock)

		Convey("When a request is made to get a list of resources", func() {
			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:29600/resources?offset=0&limit=%d", negativeLimit), http.NoBody)
			resp := httptest.NewRecorder()

			apiInstance.Router.ServeHTTP(resp, req)

			Convey("Then a bad request error is returned with status code 400", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
				errMsg := strings.TrimSpace(resp.Body.String())
				So(errMsg, ShouldEqual, expectedLimitErrorMsg)

				Convey("And the response ETag header should be empty", func() {
					So(resp.Header().Get(dpresponse.ETagHeader), ShouldBeEmpty)
				})
			})
		})
	})

	Convey("Given limit is greater than the maximum allowed", t, func() {
		greaterLimit := 1001

		greaterLimitDataStore := &apiMock.DataStorerMock{
			GetResourcesFunc: func(ctx context.Context, options data.Options) (*models.Resources, error) {
				resources := expectedResources(greaterLimit, cfg.DefaultOffset)
				return &resources, nil
			},
		}

		apiInstance := api.Setup(mux.NewRouter(), cfg, greaterLimitDataStore)

		Convey("When a request is made to get a list of resources", func() {
			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:29600/resources?limit=%d", greaterLimit), http.NoBody)
			resp := httptest.NewRecorder()

			apiInstance.Router.ServeHTTP(resp, req)

			Convey("Then a bad request error is returned with status code 400", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
				errMsg := strings.TrimSpace(resp.Body.String())
				So(errMsg, ShouldEqual, expectedLimitOverMaxMsg)
			})
		})
	})

	Convey("Given a Search Upstream API that failed to connect to the Data Store", t, func() {
		dataStorerMock := &apiMock.DataStorerMock{
			GetResourcesFunc: func(ctx context.Context, options data.Options) (*models.Resources, error) {
				return nil, errors.New("something went wrong in the server")
			},
		}

		apiInstance := api.Setup(mux.NewRouter(), cfg, dataStorerMock)

		Convey("When a request is made to get a list of all the resources that exist in the resources collection", func() {
			req := httptest.NewRequest("GET", "http://localhost:29600/resources", http.NoBody)
			resp := httptest.NewRecorder()

			apiInstance.Router.ServeHTTP(resp, req)

			Convey("Then an error with status code 500 is returned", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
				errMsg := strings.TrimSpace(resp.Body.String())
				So(errMsg, ShouldEqual, expectedServerErrorMsg)
			})
		})
	})
}
