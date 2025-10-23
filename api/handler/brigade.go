package handler

import (
	"brigade-service/service/brigade"
	"fmt"
	"net/http"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
)

func CreateBrigade(s *brigade.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var request brigade.CreateBrigadeRequest
		if err := c.ReadJson(&request); err != nil {
			return fmt.Errorf("failed to read create brigade request: %w", err)
		}

		response, err := s.CreateBrigade(c.Ctx(), request)
		if err != nil {
			return fmt.Errorf("failed to create brigade: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

type brigadeIDVars struct {
	ID int `path:"id"`
}

func GetBrigadeByID(s *brigade.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars brigadeIDVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read brigade id: %w", err)
		}

		response, err := s.GetBrigadeByID(c.Ctx(), vars.ID)
		if err != nil {
			return fmt.Errorf("failed to get brigade by id: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

func GetAllBrigades(s *brigade.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		response, err := s.GetAllBrigades(c.Ctx())
		if err != nil {
			return fmt.Errorf("failed to get all brigades: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}
