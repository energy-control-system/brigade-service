package handler

import (
	"brigade-service/service/brigade"
	"fmt"
	"net/http"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
	"github.com/sunshineOfficial/golib/pagination"
)

// CreateBrigade godoc
// @Summary Create brigade
// @Description Creates a brigade from inspector identifiers.
// @Tags brigades
// @Produce json
// @Param request body brigade.CreateBrigadeRequest true "Brigade creation payload"
// @Success 200 {object} brigade.Brigade
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Security bearer
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

// ArchiveBrigade godoc
// @Summary Archive brigade
// @Description Marks a brigade as archived.
// @Tags brigades
// @Produce json
// @Param id path int true "Brigade ID"
// @Success 204
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 404 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Security bearer
// @Router /brigades/{id}/archive [patch]
func ArchiveBrigade(s *brigade.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars brigadeIDVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read brigade id: %w", err)
		}

		if err := s.ArchiveBrigade(c.Ctx(), vars.ID); err != nil {
			return fmt.Errorf("failed to archive brigade: %w", err)
		}

		return c.WriteJson(http.StatusNoContent, nil)
	}
}

// GetAllBrigades godoc
// @Summary List brigades
// @Description Returns all brigades.
// @Tags brigades
// @Produce json
// @Param limit query int false "Maximum number of items to return; 0 means no limit"
// @Param offset query int false "Number of items to skip"
// @Success 200 {array} brigade.Brigade
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Security bearer
// @Router /brigades [get]
func GetAllBrigades(s *brigade.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars pagination.Pagination
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read pagination: %w", err)
		}

		response, err := s.GetAllBrigades(c.Ctx(), vars)
		if err != nil {
			return fmt.Errorf("failed to get all brigades: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}
