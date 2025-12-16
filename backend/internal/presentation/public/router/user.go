package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) GetMe(c echo.Context) error {
	// TODO: implement with userController
	return echo.NewHTTPError(http.StatusNotImplemented, "not implemented")
}

func (s *Server) UpdateMe(c echo.Context) error {
	// TODO: implement with userController
	return echo.NewHTTPError(http.StatusNotImplemented, "not implemented")
}
