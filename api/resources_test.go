package api_test

import (
	//"bytes"
	"context"
	"encoding/json"
	"fmt"
	//"fmt"
	"github.com/ONSdigital/dis-search-upstream-stub/config"
	"github.com/ONSdigital/dis-search-upstream-stub/data"
	"github.com/ONSdigital/dis-search-upstream-stub/models"
	"io"
	"net/http"
	"net/http/httptest"
	//"strings"
	"testing"
	"time"

	"github.com/ONSdigital/dis-search-upstream-stub/api"
	apiMock "github.com/ONSdigital/dis-search-upstream-stub/api/mock"
	//"github.com/ONSdigital/dp-api-clients-go/v2/headers"
	//dpresponse "github.com/ONSdigital/dp-net/v2/handlers/response"
	//dpHTTP "github.com/ONSdigital/dp-net/v2/http"
	//dprequest "github.com/ONSdigital/dp-net/v2/request"
	//"github.com/ONSdigital/dp-search-reindex-api/apierrors"
	//"github.com/ONSdigital/dp-search-reindex-api/mongo"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
	//"go.mongodb.org/mongo-driver/bson"
)

// Constants for testing
//const (
//	eTagValidResourceID1    = `"dcb67563ce9964e281fd3c4b6b448551638531bc"`
//	validResourceID1        = "UUID1"
//	validResourceID2        = "UUID2"
//	validResourceID3        = "UUID3"
//	notFoundResourceID      = "UUID4"
//	unLockableResourceID    = "UUID5"
//	expectedServerErrorMsg  = "internal server error"
//	validCount              = "3"
//	countNotANumber         = "notANumber"
//	countNegativeInt        = "-3"
//	expectedOffsetErrorMsg  = "invalid offset query parameter"
//	expectedLimitErrorMsg   = "invalid limit query parameter"
//	expectedLimitOverMaxMsg = "limit query parameter is larger than the maximum allowed"
//)

//var (
//	zeroTime      = time.Time{}.UTC()
//	errUnexpected = errors.New("an unexpected error occurred")
//)

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
		apiInstance := api.Setup(context.Background(), mux.NewRouter(), cfg, dataStorerMock)

		Convey("When a request is made to get a list of all resources", func() {
			req := httptest.NewRequest("GET", "http://localhost:29600/resource", nil)
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

		apiInstance := api.Setup(context.Background(), mux.NewRouter(), cfg, customValidPaginationDataStore)

		Convey("When a request is made to get a list of resources", func() {
			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:29600/resource?offset=%d&limit=%d", validOffset, validLimit), nil)
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
	//
	//	Convey("Given offset is greater than total number of resources in the Data Store", t, func() {
	//		greaterOffset := 10
	//
	//		greaterOffsetDataStore := &apiMock.DataStorerMock{
	//			GetResourcesFunc: func(ctx context.Context, options mongo.Options) (*models.Resources, error) {
	//				resources := expectedResources(ctx, t, cfg, false, cfg.DefaultLimit, greaterOffset, false)
	//				return &resources, nil
	//			},
	//		}
	//
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), greaterOffsetDataStore, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to get a list of resources", func() {
	//			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:25700/search-reindex-resources?offset=%d", greaterOffset), nil)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then a list of resources is returned with status code 200", func() {
	//				So(resp.Code, ShouldEqual, http.StatusOK)
	//
	//				payload, err := io.ReadAll(resp.Body)
	//				if err != nil {
	//					t.Errorf("failed to read payload with io.ReadAll, error: %v", err)
	//				}
	//
	//				resourcesReturned := models.Resources{}
	//				err = json.Unmarshal(payload, &resourcesReturned)
	//				So(err, ShouldBeNil)
	//
	//				Convey("And the returned list should be empty", func() {
	//					returnedResourceList := resourcesReturned.ResourceList
	//					So(returnedResourceList, ShouldHaveLength, 0)
	//
	//					Convey("And the etag of the response resources resource should be returned via the ETag header", func() {
	//						So(resp.Header().Get(dpresponse.ETagHeader), ShouldNotBeEmpty)
	//					})
	//				})
	//			})
	//		})
	//	})
	//}
	//
	//func TestGetResourcesHandlerWithEmptyResourceStoreSuccess(t *testing.T) {
	//	t.Parallel()
	//
	//	cfg, err := config.Get()
	//	if err != nil {
	//		t.Errorf("failed to retrieve default configuration, error: %v", err)
	//	}
	//
	//	Convey("Given a Search Reindex Resource API that returns an empty list of resources", t, func() {
	//		dataStorerMock := &apiMock.DataStorerMock{
	//			GetResourcesFunc: func(ctx context.Context, options mongo.Options) (*models.Resources, error) {
	//				resources := models.Resources{}
	//				return &resources, nil
	//			},
	//		}
	//
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), dataStorerMock, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to get a list of all the resources that exist in the resources collection", func() {
	//			req := httptest.NewRequest("GET", "http://localhost:25700/search-reindex-resources", nil)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then a resources resource is returned with status code 200", func() {
	//				So(resp.Code, ShouldEqual, http.StatusOK)
	//
	//				payload, err := io.ReadAll(resp.Body)
	//				if err != nil {
	//					t.Errorf("failed to read payload with io.ReadAll, error: %v", err)
	//				}
	//
	//				resourcesReturned := models.Resources{}
	//				err = json.Unmarshal(payload, &resourcesReturned)
	//				So(err, ShouldBeNil)
	//
	//				Convey("And the returned resources list should be empty", func() {
	//					So(resourcesReturned.ResourceList, ShouldHaveLength, 0)
	//				})
	//			})
	//		})
	//	})
	//}
	//
	//func TestGetResourcesHandlerFail(t *testing.T) {
	//	t.Parallel()
	//
	//	cfg, err := config.Get()
	//	if err != nil {
	//		t.Errorf("failed to retrieve default configuration, error: %v", err)
	//	}
	//
	//	dataStorerMock := &apiMock.DataStorerMock{
	//		GetResourcesFunc: func(ctx context.Context, options mongo.Options) (*models.Resources, error) {
	//			resources := expectedResources(ctx, t, cfg, false, options.Limit, options.Offset, false)
	//			return &resources, err
	//		},
	//	}
	//
	//	Convey("Given an outdated or invalid etag set in the if-match header", t, func() {
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), dataStorerMock, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to get a list of resources", func() {
	//			req := httptest.NewRequest("GET", "http://localhost:25700/search-reindex-resources", nil)
	//			err := headers.SetIfMatch(req, "invalid")
	//			if err != nil {
	//				t.Errorf("failed to set if-match header in request, error: %v", err)
	//			}
	//
	//			resp := httptest.NewRecorder()
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then a conflict with etag error is returned with status code 409", func() {
	//				So(resp.Code, ShouldEqual, http.StatusConflict)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldEqual, apierrors.ErrConflictWithETag.Error())
	//
	//				Convey("And the response ETag header should be empty", func() {
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldBeEmpty)
	//				})
	//			})
	//		})
	//	})
	//
	//	Convey("Given offset is not numeric", t, func() {
	//		nonNumericOffset := "stringOffset"
	//
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), dataStorerMock, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to get a list of resources", func() {
	//			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:25700/search-reindex-resources?offset=%s", nonNumericOffset), nil)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then a bad request error is returned with status code 400", func() {
	//				So(resp.Code, ShouldEqual, http.StatusBadRequest)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldEqual, expectedOffsetErrorMsg)
	//
	//				Convey("And the response ETag header should be empty", func() {
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldBeEmpty)
	//				})
	//			})
	//		})
	//	})
	//
	//	Convey("Given offset is negative", t, func() {
	//		negativeOffset := -3
	//
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), dataStorerMock, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to get a list of resources", func() {
	//			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:25700/search-reindex-resources?offset=%d", negativeOffset), nil)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then a bad request error is returned with status code 400", func() {
	//				So(resp.Code, ShouldEqual, http.StatusBadRequest)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldEqual, expectedOffsetErrorMsg)
	//
	//				Convey("And the response ETag header should be empty", func() {
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldBeEmpty)
	//				})
	//			})
	//		})
	//	})
	//
	//	Convey("Given limit is not numeric", t, func() {
	//		nonNumericLimit := "stringLimit"
	//
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), dataStorerMock, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to get a list of resources", func() {
	//			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:25700/search-reindex-resources?limit=%s", nonNumericLimit), nil)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then a bad request error is returned with status code 400", func() {
	//				So(resp.Code, ShouldEqual, http.StatusBadRequest)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldEqual, expectedLimitErrorMsg)
	//
	//				Convey("And the response ETag header should be empty", func() {
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldBeEmpty)
	//				})
	//			})
	//		})
	//	})
	//
	//	Convey("Given limit is negative", t, func() {
	//		negativeLimit := -1
	//
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), dataStorerMock, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to get a list of resources", func() {
	//			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:25700/search-reindex-resources?offset=0&limit=%d", negativeLimit), nil)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then a bad request error is returned with status code 400", func() {
	//				So(resp.Code, ShouldEqual, http.StatusBadRequest)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldEqual, expectedLimitErrorMsg)
	//
	//				Convey("And the response ETag header should be empty", func() {
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldBeEmpty)
	//				})
	//			})
	//		})
	//	})
	//
	//	Convey("Given limit is greater than the maximum allowed", t, func() {
	//		greaterLimit := 1001
	//
	//		greaterLimitDataStore := &apiMock.DataStorerMock{
	//			GetResourcesFunc: func(ctx context.Context, options mongo.Options) (*models.Resources, error) {
	//				resources := expectedResources(ctx, t, cfg, false, greaterLimit, cfg.DefaultOffset, false)
	//				return &resources, nil
	//			},
	//		}
	//
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), greaterLimitDataStore, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to get a list of resources", func() {
	//			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:25700/search-reindex-resources?limit=%d", greaterLimit), nil)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then a bad request error is returned with status code 400", func() {
	//				So(resp.Code, ShouldEqual, http.StatusBadRequest)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldEqual, expectedLimitOverMaxMsg)
	//
	//				Convey("And the response ETag header should be empty", func() {
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldBeEmpty)
	//				})
	//			})
	//		})
	//	})
	//
	//	Convey("Given a Search Reindex Resource API that that failed to connect to the Data Store", t, func() {
	//		dataStorerMock := &apiMock.DataStorerMock{
	//			GetResourcesFunc: func(ctx context.Context, options mongo.Options) (*models.Resources, error) {
	//				return nil, errors.New("something went wrong in the server")
	//			},
	//		}
	//
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), dataStorerMock, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to get a list of all the resources that exist in the resources collection", func() {
	//			req := httptest.NewRequest("GET", "http://localhost:25700/search-reindex-resources", nil)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then an error with status code 500 is returned", func() {
	//				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldEqual, expectedServerErrorMsg)
	//
	//				Convey("And the response ETag header should be empty", func() {
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldBeEmpty)
	//				})
	//			})
	//		})
	//	})
	//}
	//
	//func TestPutNumTasksHandler(t *testing.T) {
	//	t.Parallel()
	//
	//	cfg, err := config.Get()
	//	if err != nil {
	//		t.Errorf("failed to retrieve default configuration, error: %v", err)
	//	}
	//
	//	resourceStoreMock := &apiMock.DataStorerMock{
	//		AcquireResourceLockFunc: func(ctx context.Context, id string) (string, error) {
	//			switch id {
	//			case unLockableResourceID:
	//				return "", errors.New("acquiring lock failed")
	//			default:
	//				return "", nil
	//			}
	//		},
	//		UnlockResourceFunc: func(ctx context.Context, lockID string) {
	//			// mock UnlockResource to be successful by doing nothing
	//		},
	//		GetResourceFunc: func(ctx context.Context, id string) (*models.Resource, error) {
	//			switch id {
	//			case notFoundResourceID:
	//				return nil, mongo.ErrResourceNotFound
	//			default:
	//				resources := expectedResource(ctx, t, cfg, false, id, "", 0, false)
	//				return &resources, nil
	//			}
	//		},
	//		UpdateResourceFunc: func(ctx context.Context, id string, updates bson.M) error {
	//			switch id {
	//			case validResourceID2:
	//				return nil
	//			default:
	//				return errUnexpected
	//			}
	//		},
	//	}
	//
	//	Convey("Given valid resource id, valid value for number of tasks and no If-Match header is set", t, func() {
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), resourceStoreMock, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to update the number of tasks of a specific resource", func() {
	//			req := httptest.NewRequest("PUT", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s/number-of-tasks/%s", validResourceID2, validCount), nil)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then a status code 204 is returned", func() {
	//				So(resp.Code, ShouldEqual, http.StatusNoContent)
	//
	//				Convey("And the etag of the response resource resource should be returned via the ETag header", func() {
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldNotBeEmpty)
	//				})
	//			})
	//		})
	//	})
	//
	//	Convey("Given If-Match header is set to *", t, func() {
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), resourceStoreMock, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to update the number of tasks of a specific resource", func() {
	//			req := httptest.NewRequest("PUT", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s/number-of-tasks/%s", validResourceID2, validCount), nil)
	//
	//			err := headers.SetIfMatch(req, "*")
	//			if err != nil {
	//				t.Errorf("failed to set if-match header, error: %v", err)
	//			}
	//
	//			resp := httptest.NewRecorder()
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then a status code 204 is returned", func() {
	//				So(resp.Code, ShouldEqual, http.StatusNoContent)
	//
	//				Convey("And the etag of the response resource resource should be returned via the ETag header", func() {
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldNotBeEmpty)
	//				})
	//			})
	//		})
	//	})
	//
	//	Convey("Given a valid etag is set in the If-Match header", t, func() {
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), resourceStoreMock, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to update the number of tasks of a specific resource", func() {
	//			req := httptest.NewRequest("PUT", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s/number-of-tasks/%s", validResourceID2, validCount), nil)
	//			currentETag := `"ff9022ead6121d5a216cf0112970606b2572910d"`
	//
	//			err := headers.SetIfMatch(req, currentETag)
	//			if err != nil {
	//				t.Errorf("failed to set if-match header, error: %v", err)
	//			}
	//
	//			resp := httptest.NewRecorder()
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then a status code 204 is returned", func() {
	//				So(resp.Code, ShouldEqual, http.StatusNoContent)
	//
	//				Convey("And the etag of the response resource resource should be returned via the ETag header", func() {
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldNotBeEmpty)
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldNotEqual, currentETag)
	//				})
	//			})
	//		})
	//	})
	//
	//	Convey("Given an empty etag is set in the If-Match header", t, func() {
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), resourceStoreMock, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to update the number of tasks of a specific resource", func() {
	//			req := httptest.NewRequest("PUT", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s/number-of-tasks/%s", validResourceID2, validCount), nil)
	//
	//			err := headers.SetIfMatch(req, "")
	//			if err != nil {
	//				t.Errorf("failed to set if-match header, error: %v", err)
	//			}
	//
	//			resp := httptest.NewRecorder()
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then a status code 204 is returned", func() {
	//				So(resp.Code, ShouldEqual, http.StatusNoContent)
	//
	//				Convey("And the etag of the response resource resource should be returned via the ETag header", func() {
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldNotBeEmpty)
	//				})
	//			})
	//		})
	//	})
	//
	//	Convey("Given an outdated or invalid etag is set in the If-Match header", t, func() {
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), resourceStoreMock, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to update the number of tasks of a specific resource", func() {
	//			req := httptest.NewRequest("PUT", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s/number-of-tasks/%s", validResourceID2, validCount), nil)
	//
	//			err := headers.SetIfMatch(req, "invalid")
	//			if err != nil {
	//				t.Errorf("failed to set if-match header, error: %v", err)
	//			}
	//
	//			resp := httptest.NewRecorder()
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then a conflict with etag error is returned with status code 409", func() {
	//				So(resp.Code, ShouldEqual, http.StatusConflict)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldEqual, apierrors.ErrConflictWithETag.Error())
	//
	//				Convey("And the response ETag header should be empty", func() {
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldBeEmpty)
	//				})
	//			})
	//		})
	//	})
	//
	//	Convey("Given a specific resource does not exist in the Data Store", t, func() {
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), resourceStoreMock, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to update the number of tasks of a specific resource", func() {
	//			req := httptest.NewRequest("PUT", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s/number-of-tasks/%s", notFoundResourceID, validCount), nil)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then resource resource was not found returning a status code of 404", func() {
	//				So(resp.Code, ShouldEqual, http.StatusNotFound)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldEqual, "failed to find the specified reindex resource")
	//
	//				Convey("And the response ETag header should be empty", func() {
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldBeEmpty)
	//				})
	//			})
	//		})
	//	})
	//
	//	Convey("Given the value of number of tasks is not an integer", t, func() {
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), resourceStoreMock, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to update the number of tasks of a specific resource", func() {
	//			req := httptest.NewRequest("PUT", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s/number-of-tasks/%s", validResourceID2, countNotANumber), nil)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then it is a bad request returning a status code of 400", func() {
	//				So(resp.Code, ShouldEqual, http.StatusBadRequest)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldEqual, "number of tasks must be a positive integer")
	//
	//				Convey("And the response ETag header should be empty", func() {
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldBeEmpty)
	//				})
	//			})
	//		})
	//	})
	//
	//	Convey("Given the value of number of tasks is a negative integer", t, func() {
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), resourceStoreMock, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to update the number of tasks of a specific resource", func() {
	//			req := httptest.NewRequest("PUT", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s/number-of-tasks/%s", validResourceID2, countNegativeInt), nil)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then it is a bad request returning a status code of 400", func() {
	//				So(resp.Code, ShouldEqual, http.StatusBadRequest)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldEqual, "number of tasks must be a positive integer")
	//
	//				Convey("And the response ETag header should be empty", func() {
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldBeEmpty)
	//				})
	//			})
	//		})
	//	})
	//
	//	Convey("Given the Data Store is unable to lock the resource id", t, func() {
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), resourceStoreMock, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to update the number of tasks of a specific resource", func() {
	//			req := httptest.NewRequest("PUT", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s/number-of-tasks/%s", unLockableResourceID, validCount), nil)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then an error with status code 500 is returned", func() {
	//				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldEqual, expectedServerErrorMsg)
	//
	//				Convey("And the response ETag header should be empty", func() {
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldBeEmpty)
	//				})
	//			})
	//		})
	//	})
	//
	//	Convey("Given the request results in no modifications to the resource resource", t, func() {
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), resourceStoreMock, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to update the number of tasks of a specific resource", func() {
	//			req := httptest.NewRequest("PUT", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s/number-of-tasks/%s", validResourceID2, "0"), nil)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then a status code 304 is returned", func() {
	//				So(resp.Code, ShouldEqual, http.StatusNotModified)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldEqual, apierrors.ErrNewETagSame.Error())
	//
	//				Convey("And the response ETag header should be empty", func() {
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldBeEmpty)
	//				})
	//			})
	//		})
	//	})
	//
	//	Convey("Given an unexpected error occurs in the Data Store", t, func() {
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), resourceStoreMock, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a request is made to update the number of tasks of a specific resource", func() {
	//			req := httptest.NewRequest("PUT", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s/number-of-tasks/%s", validResourceID3, validCount), nil)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then the response returns a status code of 500", func() {
	//				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldEqual, "internal server error")
	//
	//				Convey("And the response ETag header should be empty", func() {
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldBeEmpty)
	//				})
	//			})
	//		})
	//	})
	//}
	//
	//func TestPatchResourceStatusHandler(t *testing.T) {
	//	t.Parallel()
	//
	//	cfg, err := config.Get()
	//	if err != nil {
	//		t.Errorf("failed to retrieve default configuration, error: %v", err)
	//	}
	//
	//	var etag1, etag2 string
	//
	//	resourceStoreMock := &apiMock.DataStorerMock{
	//		GetResourceFunc: func(ctx context.Context, id string) (*models.Resource, error) {
	//			switch id {
	//			case validResourceID1:
	//				newResource := expectedResource(ctx, t, cfg, false, validResourceID1, "", 0, false)
	//				etag1 = newResource.ETag
	//				return &newResource, nil
	//			case validResourceID2:
	//				newResource := expectedResource(ctx, t, cfg, false, validResourceID2, "", 0, false)
	//				etag2 = newResource.ETag
	//				return &newResource, nil
	//			case unLockableResourceID:
	//				newResource := expectedResource(ctx, t, cfg, false, unLockableResourceID, "", 0, false)
	//				return &newResource, nil
	//			case notFoundResourceID:
	//				return nil, mongo.ErrResourceNotFound
	//			default:
	//				return nil, errUnexpected
	//			}
	//		},
	//		AcquireResourceLockFunc: func(ctx context.Context, id string) (string, error) {
	//			switch id {
	//			case unLockableResourceID:
	//				return "", errors.New("acquiring lock failed")
	//			default:
	//				return "", nil
	//			}
	//		},
	//		UnlockResourceFunc: func(ctx context.Context, lockID string) {
	//			// mock UnlockResource to be successful by doing nothing
	//		},
	//		UpdateResourceFunc: func(ctx context.Context, id string, updates bson.M) error {
	//			switch id {
	//			case validResourceID2:
	//				return errUnexpected
	//			default:
	//				return nil
	//			}
	//		},
	//	}
	//
	//	validPatchBody := `[
	//		{ "op": "replace", "path": "/state", "value": "created" },
	//		{ "op": "replace", "path": "/total_search_documents", "value": 100 }
	//	]`
	//
	//	Convey("Given a Search Reindex Resource API that updates state of a resource via patch request", t, func() {
	//		httpClient := dpHTTP.NewClient()
	//		apiInstance := api.Setup(mux.NewRouter(), resourceStoreMock, &apiMock.AuthHandlerMock{}, taskNames, cfg, httpClient, &apiMock.IndexerMock{}, &apiMock.ReindexRequestedProducerMock{})
	//
	//		Convey("When a patch request is made with valid resource ID and valid patch operations", func() {
	//			req := httptest.NewRequest("PATCH", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s", validResourceID1), bytes.NewBufferString(validPatchBody))
	//			headers.SetIfMatch(req, etag1)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then the response should return a 204 status code", func() {
	//				So(resp.Code, ShouldEqual, http.StatusNoContent)
	//
	//				Convey("And the new eTag of the resource is returned via ETag header", func() {
	//					So(resp.Header().Get(dpresponse.ETagHeader), ShouldNotBeEmpty)
	//				})
	//			})
	//		})
	//
	//		Convey("When a patch request is made with invalid resource ID", func() {
	//			req := httptest.NewRequest("PATCH", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s", notFoundResourceID), bytes.NewBufferString(validPatchBody))
	//			headers.SetIfMatch(req, etag1)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then the response should return a 404 status code", func() {
	//				So(resp.Code, ShouldEqual, http.StatusNotFound)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldNotBeEmpty)
	//			})
	//		})
	//
	//		Convey("When connection to datastore has failed", func() {
	//			req := httptest.NewRequest("PATCH", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s", "invalid"), bytes.NewBufferString(validPatchBody))
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then the response should return a 500 status code", func() {
	//				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldNotBeEmpty)
	//			})
	//		})
	//
	//		Convey("When a patch request is made with no If-Match header", func() {
	//			req := httptest.NewRequest("PATCH", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s", validResourceID1), bytes.NewBufferString(validPatchBody))
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then the response should return a 201 status code as eTag check is ignored", func() {
	//				So(resp.Code, ShouldEqual, http.StatusNoContent)
	//			})
	//		})
	//
	//		Convey("When a patch request is made with outdated or unknown eTag in If-Match header", func() {
	//			req := httptest.NewRequest("PATCH", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s", validResourceID1), bytes.NewBufferString(validPatchBody))
	//			headers.SetIfMatch(req, "invalidETag")
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then the response should return a 201 status code as eTag check is ignored", func() {
	//				So(resp.Code, ShouldEqual, http.StatusConflict)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldNotBeEmpty)
	//			})
	//		})
	//
	//		Convey("When a patch request is made with invalid patch body given", func() {
	//			req := httptest.NewRequest("PATCH", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s", validResourceID1), bytes.NewBufferString("{}"))
	//			headers.SetIfMatch(req, etag1)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then the response should return a 400 status code", func() {
	//				So(resp.Code, ShouldEqual, http.StatusBadRequest)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldNotBeEmpty)
	//			})
	//		})
	//
	//		Convey("When a patch request is made with no patches given in request body", func() {
	//			req := httptest.NewRequest("PATCH", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s", validResourceID1), bytes.NewBufferString("[]"))
	//			headers.SetIfMatch(req, eTagValidResourceID1)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then the response should return a 400 status code", func() {
	//				So(resp.Code, ShouldEqual, http.StatusBadRequest)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldNotBeEmpty)
	//			})
	//		})
	//
	//		Convey("When a patch request is made with patch body containing invalid information", func() {
	//			patchBodyWithInvalidData := `[
	//				{ "op": "replace", "path": "/total_search_documents", "value": "invalid" }
	//			]`
	//
	//			req := httptest.NewRequest("PATCH", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s", validResourceID1), bytes.NewBufferString(patchBodyWithInvalidData))
	//			headers.SetIfMatch(req, etag1)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then the response should return a 400 status code", func() {
	//				So(resp.Code, ShouldEqual, http.StatusBadRequest)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldNotBeEmpty)
	//			})
	//		})
	//
	//		Convey("When acquiring resource lock to update resource has failed", func() {
	//			req := httptest.NewRequest("PATCH", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s", unLockableResourceID), bytes.NewBufferString(validPatchBody))
	//			unLockableResourceETag := `"24decf55038de874bc6fa9cf0930adc219f15db1"`
	//			headers.SetIfMatch(req, unLockableResourceETag)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then the response should return a 500 status code", func() {
	//				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldNotBeEmpty)
	//			})
	//		})
	//
	//		Convey("When a patch request is made which results in no modification", func() {
	//			patchBodyWithNoModification := `[
	//				{ "op": "replace", "path": "/state", "value": "created" }
	//			]`
	//
	//			req := httptest.NewRequest("PATCH", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s", validResourceID1), bytes.NewBufferString(patchBodyWithNoModification))
	//			headers.SetIfMatch(req, etag1)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then the response should return a 304 status code", func() {
	//				So(resp.Code, ShouldEqual, http.StatusNotModified)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldNotBeEmpty)
	//			})
	//		})
	//
	//		Convey("When the update to resource with patches has failed due to failing on UpdateResourceWithPatches func", func() {
	//			req := httptest.NewRequest("PATCH", fmt.Sprintf("http://localhost:25700/search-reindex-resources/%s", validResourceID2), bytes.NewBufferString(validPatchBody))
	//			headers.SetIfMatch(req, etag2)
	//			resp := httptest.NewRecorder()
	//
	//			apiInstance.Router.ServeHTTP(resp, req)
	//
	//			Convey("Then the response should return a 500 status code", func() {
	//				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
	//				errMsg := strings.TrimSpace(resp.Body.String())
	//				So(errMsg, ShouldNotBeEmpty)
	//			})
	//		})
	//	})
	//}
	//
	//func TestPreparePatchUpdatesSuccess(t *testing.T) {
	//	t.Parallel()
	//
	//	testCtx := context.Background()
	//
	//	cfg, err := config.Get()
	//	if err != nil {
	//		t.Errorf("failed to retrieve default configuration, error: %v", err)
	//	}
	//
	//	currentResource := expectedResource(testCtx, t, cfg, false, validResourceID1, "", 0, false)
	//
	//	Convey("Given valid patches", t, func() {
	//		validPatches := []dprequest.Patch{
	//			{
	//				Op:    dprequest.OpReplace.String(),
	//				Path:  models.ResourceTotalSearchDocumentsPath,
	//				Value: float64(100),
	//			},
	//			{
	//				Op:    dprequest.OpReplace.String(),
	//				Path:  models.ResourceNoOfTasksPath,
	//				Value: float64(2),
	//			},
	//		}
	//
	//		Convey("When preparePatchUpdates is called", func() {
	//			updatedResource, bsonUpdates, err := api.GetUpdatesFromResourcePatches(testCtx, validPatches, &currentResource)
	//
	//			Convey("Then updatedResource should contain updates from the patch", func() {
	//				So(updatedResource.TotalSearchDocuments, ShouldEqual, 100)
	//				So(updatedResource.NumberOfTasks, ShouldEqual, 2)
	//
	//				Convey("And bsonUpdates should contain updates from the patch", func() {
	//					So(bsonUpdates[models.ResourceTotalSearchDocumentsKey], ShouldEqual, 100)
	//					So(bsonUpdates[models.ResourceNoOfTasksKey], ShouldEqual, 2)
	//
	//					Convey("And LastUpdated should be updated", func() {
	//						So(updatedResource.LastUpdated, ShouldNotEqual, currentResource.LastUpdated)
	//						So(bsonUpdates[models.ResourceLastUpdatedKey], ShouldNotBeEmpty)
	//
	//						Convey("And no error should be returned", func() {
	//							So(err, ShouldBeNil)
	//						})
	//					})
	//				})
	//			})
	//		})
	//	})
	//
	//	Convey("Given patches which changes state to in-progress", t, func() {
	//		inProgressStatePatches := []dprequest.Patch{
	//			{
	//				Op:    dprequest.OpReplace.String(),
	//				Path:  models.ResourceStatePath,
	//				Value: models.ResourceStateInProgress,
	//			},
	//		}
	//
	//		Convey("When preparePatchUpdates is called", func() {
	//			updatedResource, bsonUpdates, err := api.GetUpdatesFromResourcePatches(testCtx, inProgressStatePatches, &currentResource)
	//
	//			Convey("Then updatedResource and bsonUpdates should contain updates from the patch", func() {
	//				So(updatedResource.State, ShouldEqual, models.ResourceStateInProgress)
	//				So(bsonUpdates[models.ResourceStateKey], ShouldEqual, models.ResourceStateInProgress)
	//
	//				Convey("And reindex started should be updated", func() {
	//					So(updatedResource.ReindexStarted, ShouldNotEqual, currentResource.ReindexStarted)
	//					So(bsonUpdates[models.ResourceReindexStartedKey], ShouldNotBeEmpty)
	//
	//					Convey("And LastUpdated should be updated", func() {
	//						So(updatedResource.LastUpdated, ShouldNotEqual, currentResource.LastUpdated)
	//						So(bsonUpdates[models.ResourceLastUpdatedKey], ShouldNotBeEmpty)
	//
	//						Convey("And no error should be returned", func() {
	//							So(err, ShouldBeNil)
	//						})
	//					})
	//				})
	//			})
	//		})
	//	})
	//
	//	Convey("Given patches which changes state to failed", t, func() {
	//		failedStatePatches := []dprequest.Patch{
	//			{
	//				Op:    dprequest.OpReplace.String(),
	//				Path:  models.ResourceStatePath,
	//				Value: models.ResourceStateFailed,
	//			},
	//		}
	//
	//		Convey("When preparePatchUpdates is called", func() {
	//			updatedResource, bsonUpdates, err := api.GetUpdatesFromResourcePatches(testCtx, failedStatePatches, &currentResource)
	//
	//			Convey("Then updatedResource and bsonUpdates should contain updates from the patch", func() {
	//				So(updatedResource.State, ShouldEqual, models.ResourceStateFailed)
	//				So(bsonUpdates[models.ResourceStateKey], ShouldEqual, models.ResourceStateFailed)
	//
	//				Convey("And ReindexFailed should be updated", func() {
	//					So(updatedResource.ReindexFailed, ShouldNotEqual, currentResource.ReindexFailed)
	//					So(bsonUpdates[models.ResourceReindexFailedKey], ShouldNotBeEmpty)
	//
	//					Convey("And LastUpdated should be updated", func() {
	//						So(updatedResource.LastUpdated, ShouldNotEqual, currentResource.LastUpdated)
	//						So(bsonUpdates[models.ResourceLastUpdatedKey], ShouldNotBeEmpty)
	//
	//						Convey("And no error should be returned", func() {
	//							So(err, ShouldBeNil)
	//						})
	//					})
	//				})
	//			})
	//		})
	//	})
	//
	//	Convey("Given patches which changes state to completed", t, func() {
	//		completedStatePatches := []dprequest.Patch{
	//			{
	//				Op:    dprequest.OpReplace.String(),
	//				Path:  models.ResourceStatePath,
	//				Value: models.ResourceStateCompleted,
	//			},
	//		}
	//
	//		Convey("When preparePatchUpdates is called", func() {
	//			updatedResource, bsonUpdates, err := api.GetUpdatesFromResourcePatches(testCtx, completedStatePatches, &currentResource)
	//
	//			Convey("Then updatedResource and bsonUpdates should contain updates from the patch", func() {
	//				So(updatedResource.State, ShouldEqual, models.ResourceStateCompleted)
	//				So(bsonUpdates[models.ResourceStateKey], ShouldEqual, models.ResourceStateCompleted)
	//
	//				Convey("And ReindexCompleted should be updated", func() {
	//					So(updatedResource.ReindexCompleted, ShouldNotEqual, currentResource.ReindexCompleted)
	//					So(bsonUpdates[models.ResourceReindexCompletedKey], ShouldNotBeEmpty)
	//
	//					Convey("And LastUpdated should be updated", func() {
	//						So(updatedResource.LastUpdated, ShouldNotEqual, currentResource.LastUpdated)
	//						So(bsonUpdates[models.ResourceLastUpdatedKey], ShouldNotBeEmpty)
	//
	//						Convey("And no error should be returned", func() {
	//							So(err, ShouldBeNil)
	//						})
	//					})
	//				})
	//			})
	//		})
	//	})
	//}
	//
	//func TestPreparePatchUpdatesFail(t *testing.T) {
	//	t.Parallel()
	//
	//	testCtx := context.Background()
	//
	//	cfg, err := config.Get()
	//	if err != nil {
	//		t.Errorf("failed to retrieve default configuration, error: %v", err)
	//	}
	//
	//	currentResource := expectedResource(testCtx, t, cfg, false, validResourceID1, "", 0, false)
	//
	//	Convey("Given patches with unknown path", t, func() {
	//		unknownPathPatches := []dprequest.Patch{
	//			{
	//				Op:    dprequest.OpReplace.String(),
	//				Path:  "/unknown",
	//				Value: "unknown",
	//			},
	//		}
	//
	//		Convey("When preparePatchUpdates is called", func() {
	//			updatedResource, bsonUpdates, err := api.GetUpdatesFromResourcePatches(testCtx, unknownPathPatches, &currentResource)
	//
	//			Convey("Then an error should be returned", func() {
	//				So(err, ShouldNotBeNil)
	//				So(err.Error(), ShouldEqual, fmt.Sprintf("provided path '%s' not supported", unknownPathPatches[0].Path))
	//
	//				So(updatedResource, ShouldResemble, models.Resource{})
	//				So(bsonUpdates, ShouldBeEmpty)
	//			})
	//		})
	//	})
	//
	//	Convey("Given patches with invalid number of tasks", t, func() {
	//		invalidNoOfTasksPatches := []dprequest.Patch{
	//			{
	//				Op:    dprequest.OpReplace.String(),
	//				Path:  models.ResourceNoOfTasksPath,
	//				Value: "unknown",
	//			},
	//		}
	//
	//		Convey("When preparePatchUpdates is called", func() {
	//			updatedResource, bsonUpdates, err := api.GetUpdatesFromResourcePatches(testCtx, invalidNoOfTasksPatches, &currentResource)
	//
	//			Convey("Then an error should be returned", func() {
	//				So(err, ShouldNotBeNil)
	//				So(err.Error(), ShouldEqual, fmt.Sprintf("wrong value type `%s` for `%s`, expected an integer", api.GetValueType(invalidNoOfTasksPatches[0].Value), invalidNoOfTasksPatches[0].Path))
	//
	//				So(updatedResource, ShouldResemble, models.Resource{})
	//				So(bsonUpdates, ShouldBeEmpty)
	//			})
	//		})
	//	})
	//
	//	Convey("Given patches with unknown state", t, func() {
	//		unknownStatePatches := []dprequest.Patch{
	//			{
	//				Op:    dprequest.OpReplace.String(),
	//				Path:  models.ResourceStatePath,
	//				Value: "unknown",
	//			},
	//		}
	//
	//		Convey("When preparePatchUpdates is called", func() {
	//			updatedResource, bsonUpdates, err := api.GetUpdatesFromResourcePatches(testCtx, unknownStatePatches, &currentResource)
	//
	//			Convey("Then an error should be returned", func() {
	//				So(err, ShouldNotBeNil)
	//				So(err.Error(), ShouldEqual, fmt.Sprintf("invalid resource state `%s` for `%s` - expected %v", unknownStatePatches[0].Value, unknownStatePatches[0].Path, models.ValidResourceStates))
	//
	//				So(updatedResource, ShouldResemble, models.Resource{})
	//				So(bsonUpdates, ShouldBeEmpty)
	//			})
	//		})
	//	})
	//
	//	Convey("Given patches with invalid state", t, func() {
	//		invalidStatePatches := []dprequest.Patch{
	//			{
	//				Op:    dprequest.OpReplace.String(),
	//				Path:  models.ResourceStatePath,
	//				Value: 12,
	//			},
	//		}
	//
	//		Convey("When preparePatchUpdates is called", func() {
	//			updatedResource, bsonUpdates, err := api.GetUpdatesFromResourcePatches(testCtx, invalidStatePatches, &currentResource)
	//
	//			Convey("Then an error should be returned", func() {
	//				So(err, ShouldNotBeNil)
	//				So(err.Error(), ShouldEqual, fmt.Sprintf("wrong value type `%s` for `%s`, expected string", api.GetValueType(invalidStatePatches[0].Value), invalidStatePatches[0].Path))
	//
	//				So(updatedResource, ShouldResemble, models.Resource{})
	//				So(bsonUpdates, ShouldBeEmpty)
	//			})
	//		})
	//	})
	//
	//	Convey("Given patches with invalid total search documents", t, func() {
	//		invalidTotalSearchDocsPatches := []dprequest.Patch{
	//			{
	//				Op:    dprequest.OpReplace.String(),
	//				Path:  models.ResourceTotalSearchDocumentsPath,
	//				Value: "invalid",
	//			},
	//		}
	//
	//		Convey("When preparePatchUpdates is called", func() {
	//			updatedResource, bsonUpdates, err := api.GetUpdatesFromResourcePatches(testCtx, invalidTotalSearchDocsPatches, &currentResource)
	//
	//			Convey("Then an error should be returned", func() {
	//				So(err, ShouldNotBeNil)
	//				So(err.Error(), ShouldEqual, fmt.Sprintf("wrong value type `%s` for `%s`, expected an integer", api.GetValueType(invalidTotalSearchDocsPatches[0].Value), invalidTotalSearchDocsPatches[0].Path))
	//
	//				So(updatedResource, ShouldResemble, models.Resource{})
	//				So(bsonUpdates, ShouldBeEmpty)
	//			})
	//		})
	//	})
}
