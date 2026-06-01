package api

import (
	"brigade-service/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sunshineOfficial/golib/golog"
)

func TestBrigadeRoutesAllowUnauthenticatedRequests(t *testing.T) {
	builder := NewServerBuilder(t.Context(), golog.NewLogger("test"), config.Settings{
		Port: 80,
	})
	builder.AddBrigades(nil)

	routes := []struct {
		method string
		path   string
	}{
		{method: http.MethodPost, path: "/brigades"},
		{method: http.MethodGet, path: "/brigades/1"},
		{method: http.MethodPatch, path: "/brigades/1/archive"},
		{method: http.MethodGet, path: "/brigades"},
	}

	for _, route := range routes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			response := httptest.NewRecorder()
			request := httptest.NewRequest(route.method, route.path, nil)

			builder.router.ServeHTTP(response, request)

			if response.Code == http.StatusUnauthorized {
				t.Fatalf("status = %d, route must be open without authorization", response.Code)
			}
		})
	}
}
