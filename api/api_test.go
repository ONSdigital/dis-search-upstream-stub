package api

import (
	"context"
	"github.com/ONSdigital/dis-search-upstream-stub/config"
	"github.com/ONSdigital/dis-search-upstream-stub/data"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		r := mux.NewRouter()
		ctx := context.Background()
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		api := Setup(ctx, r, cfg, &data.ResourceStore{})

		// TODO: remove hello world example handler route test case
		Convey("When created the following routes should have been added", func() {
			// Replace the check below with any newly added api endpoints
			So(hasRoute(api.Router, "/hello", "GET"), ShouldBeTrue)
			So(hasRoute(api.Router, "/resource", "GET"), ShouldBeTrue)
		})
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, http.NoBody)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}
