package api

import (
	"brigade-service/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sunshineOfficial/golib/golog"
)

func TestBrigadeAuthorizationPolicy(t *testing.T) {
	builder := NewServerBuilder(t.Context(), golog.NewLogger("test"), config.Settings{
		Port: 80,
	})
	builder.AddBrigades(nil)

	t.Run("creation requires authorization", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/brigades", nil)

		builder.router.ServeHTTP(response, request)

		if response.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
		}
	})

	t.Run("get by id allows internal calls without authorization", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/brigades/1", nil)

		builder.router.ServeHTTP(response, request)

		if response.Code == http.StatusUnauthorized {
			t.Fatalf("status = %d, route must stay open for internal service calls", response.Code)
		}
	})
}
