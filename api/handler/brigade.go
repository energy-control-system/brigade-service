package handler

import (
	"brigade-service/service/brigade"
	"fmt"
	"net/http"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
)

// CreateBrigade godoc
// @Summary Create brigade
// @Description Creates a brigade from inspector identifiers.
// @Tags brigades
// @Accept json
// @Produce json
// @Param request body brigade.CreateBrigadeRequest true "Brigade creation payload"
// @Success 200 {object} brigade.Brigade
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /brigades [post]
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

// GetBrigadeByID godoc
// @Summary Get brigade by ID
// @Description Returns a brigade with assigned inspectors.
// @Tags brigades
// @Produce json
// @Param id path int true "Brigade ID"
// @Success 200 {object} brigade.Brigade
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 404 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /brigades/{id} [get]
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

// GetAllBrigades godoc
// @Summary List brigades
// @Description Returns all brigades.
// @Tags brigades
// @Produce json
// @Success 200 {array} brigade.Brigade
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /brigades [get]
func GetAllBrigades(s *brigade.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		response, err := s.GetAllBrigades(c.Ctx())
		if err != nil {
			return fmt.Errorf("failed to get all brigades: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}
