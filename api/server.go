package api

import (
	"brigade-service/api/handler"
	"brigade-service/config"
	"brigade-service/service/brigade"
	"context"
	"fmt"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
	"github.com/sunshineOfficial/golib/gohttp/gorouter/middleware"
	"github.com/sunshineOfficial/golib/gohttp/gorouter/plugin"
	"github.com/sunshineOfficial/golib/gohttp/goserver"
	"github.com/sunshineOfficial/golib/golog"
)

type ServerBuilder struct {
	server goserver.Server
	router *gorouter.Router
}

func NewServerBuilder(ctx context.Context, log golog.Logger, settings config.Settings) *ServerBuilder {
	return &ServerBuilder{
		server: goserver.NewHTTPServer(ctx, log, fmt.Sprintf(":%d", settings.Port)),
		router: gorouter.NewRouter(log).Use(
			middleware.Metrics(),
			middleware.Recover,
			middleware.LogError,
		),
	}
}

func (s *ServerBuilder) AddDebug() {
	s.router.Install(plugin.NewPProf(), plugin.NewMetrics(), plugin.NewSwaggo("api/brigade-service"))
}

func (s *ServerBuilder) AddBrigades(service *brigade.Service) {
	r := s.router.SubRouter("/brigades")
	r.HandlePost("", handler.CreateBrigade(service))
	r.HandleGet("/{id}", handler.GetBrigadeByID(service))
	r.HandleGet("", handler.GetAllBrigades(service))
}

func (s *ServerBuilder) Build() goserver.Server {
	s.server.UseHandler(s.router)

	return s.server
}
