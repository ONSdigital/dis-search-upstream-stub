package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ONSdigital/dis-search-upstream-stub/api"
	apiMock "github.com/ONSdigital/dis-search-upstream-stub/api/mock"
	"github.com/ONSdigital/dis-search-upstream-stub/config"
	"github.com/ONSdigital/dis-search-upstream-stub/data"
	"github.com/ONSdigital/dis-search-upstream-stub/models"
	dpresponse "github.com/ONSdigital/dp-net/v2/handlers/response"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Constants for testing
const (
	expectedServerErrorMsg  = "internal server error"
	expectedOffsetErrorMsg  = "invalid offset query parameter"
	expectedLimitErrorMsg   = "invalid limit query parameter"
	expectedLimitOverMaxMsg = "limit query parameter is larger than the maximum allowed"
)

// expectedResource returns a Resource that can be used to define and test expected values within it
func expectedResource(uri string) models.Resource {
	topics := []string{"a", "b", "c", "d"}
	dateChanges := []string{"a change_notice", "a previous_date"}
	resource := models.Resource{
		Uri:             uri,
		UriOld:          "/an/old/uri",
		ContentType:     "api_dataset_landing_page",
		CDID:            "ASELECTIONOFNUMBERSANDLETTERS123",
		DatasetID:       "ASELECTIONOFNUMBERSANDLETTERS456",
		Edition:         "an edition",
		MetaDescription: "a description",
		ReleaseDate:     time.Time{}.UTC(),
		Summary:         "a summary",
		Title:           "a title",
		Topics:          topics,
		Language:        "string",
		Survey:          "string",
		CanonicalTopic:  "string",
		Cancelled:       true,
		Finalised:       true,
		Published:       true,
		DateChanges:     dateChanges,
		ProvisionalDate: "October-November 2024",
	}

	return resource
}

func expectedResources(limit, offset int) models.Resources {
	resources := models.Resources{
		Limit:      limit,
		Offset:     offset,
		TotalCount: 2,
	}

	firstResource := expectedResource("/a/uri")
	secondResource := expectedResource("/another/uri")

	if (offset == 0) && (limit > 1) {
		resources.Count = 2
		resources.ResourceList = []models.Resource{firstResource, secondResource}
	}

	if (offset == 1) && (limit > 0) {
		resources.Count = 1
		resources.ResourceList = []models.Resource{secondResource}
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
			req := httptest.NewRequest("GET", "http://localhost:29600/resource", http.NoBody)
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

				expectedResource1 := expectedResource("/a/uri")
				expectedResource2 := expectedResource("/another/uri")

				Convey("And the returned list should contain expected resources", func() {
					returnedResourceList := resourcesReturned.ResourceList
					So(returnedResourceList, ShouldHaveLength, 2)
					returnedResource1 := returnedResourceList[0]
					So(returnedResource1.Uri, ShouldEqual, expectedResource1.Uri)
					So(returnedResource1.UriOld, ShouldEqual, expectedResource1.UriOld)
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
					So(returnedResource1.Cancelled, ShouldEqual, expectedResource1.Cancelled)
					So(returnedResource1.Finalised, ShouldEqual, expectedResource1.Finalised)
					So(returnedResource1.Published, ShouldEqual, expectedResource1.Published)
					So(returnedResource1.DateChanges, ShouldEqual, expectedResource1.DateChanges)
					So(returnedResource1.ProvisionalDate, ShouldEqual, expectedResource1.ProvisionalDate)
					returnedResource2 := returnedResourceList[1]
					So(returnedResource2.UriOld, ShouldEqual, expectedResource2.UriOld)
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
					So(returnedResource2.Cancelled, ShouldEqual, expectedResource2.Cancelled)
					So(returnedResource2.Finalised, ShouldEqual, expectedResource2.Finalised)
					So(returnedResource2.Published, ShouldEqual, expectedResource2.Published)
					So(returnedResource2.DateChanges, ShouldEqual, expectedResource2.DateChanges)
					So(returnedResource2.ProvisionalDate, ShouldEqual, expectedResource2.ProvisionalDate)
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
			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:29600/resource?offset=%d&limit=%d", validOffset, validLimit), http.NoBody)
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

				expectedResource := expectedResource("/another/uri")

				Convey("And the returned list should contain the expected resource", func() {
					returnedResourceList := resourcesReturned.ResourceList
					So(returnedResourceList, ShouldHaveLength, 1)
					returnedResource := returnedResourceList[0]
					So(returnedResource.Uri, ShouldEqual, expectedResource.Uri)
					So(returnedResource.UriOld, ShouldEqual, expectedResource.UriOld)
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
					So(returnedResource.Cancelled, ShouldEqual, expectedResource.Cancelled)
					So(returnedResource.Finalised, ShouldEqual, expectedResource.Finalised)
					So(returnedResource.Published, ShouldEqual, expectedResource.Published)
					So(returnedResource.DateChanges, ShouldEqual, expectedResource.DateChanges)
					So(returnedResource.ProvisionalDate, ShouldEqual, expectedResource.ProvisionalDate)
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
			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:29600/resource?offset=%d", greaterOffset), http.NoBody)
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
					returnedResourceList := resourcesReturned.ResourceList
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

	Convey("Given a Search Reindex Resource API that returns an empty list of resources", t, func() {
		dataStorerMock := &apiMock.DataStorerMock{
			GetResourcesFunc: func(ctx context.Context, options data.Options) (*models.Resources, error) {
				resources := models.Resources{}
				return &resources, nil
			},
		}

		apiInstance := api.Setup(mux.NewRouter(), cfg, dataStorerMock)

		Convey("When a request is made to get a list of all the resources that exist in the resources collection", func() {
			req := httptest.NewRequest("GET", "http://localhost:29600/resource", http.NoBody)
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
					So(resourcesReturned.ResourceList, ShouldHaveLength, 0)
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
			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:29600/resource?offset=%s", nonNumericOffset), http.NoBody)
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
			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:29600/resource?offset=%d", negativeOffset), http.NoBody)
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
			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:29600/resource?limit=%s", nonNumericLimit), http.NoBody)
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
			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:29600/resource?offset=0&limit=%d", negativeLimit), http.NoBody)
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
			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:29600/resource?limit=%d", greaterLimit), http.NoBody)
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
			req := httptest.NewRequest("GET", "http://localhost:29600/resource", http.NoBody)
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
