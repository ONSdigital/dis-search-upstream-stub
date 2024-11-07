package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dis-search-upstream-stub/config"
	"github.com/ONSdigital/dis-search-upstream-stub/data"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		r := mux.NewRouter()
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		api := Setup(r, cfg, &data.ResourceStore{})

		Convey("When created the following routes should have been added", func() {
			So(hasRoute(api.Router, "/resource", "GET"), ShouldBeTrue)
		})
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, http.NoBody)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}
