package router

import (
	"net/http"

	"good-todo-go/internal/presentation/public/api"

	"github.com/labstack/echo/v4"
)

func (s *Server) GetTodos(c echo.Context, params api.GetTodosParams) error {
	// TODO: implement with todoController
	return echo.NewHTTPError(http.StatusNotImplemented, "not implemented")
}

func (s *Server) CreateTodo(c echo.Context) error {
	// TODO: implement with todoController
	return echo.NewHTTPError(http.StatusNotImplemented, "not implemented")
}

func (s *Server) DeleteTodo(c echo.Context, todoId string) error {
	// TODO: implement with todoController
	return echo.NewHTTPError(http.StatusNotImplemented, "not implemented")
}

func (s *Server) GetTodo(c echo.Context, todoId string) error {
	// TODO: implement with todoController
	return echo.NewHTTPError(http.StatusNotImplemented, "not implemented")
}

func (s *Server) UpdateTodo(c echo.Context, todoId string) error {
	// TODO: implement with todoController
	return echo.NewHTTPError(http.StatusNotImplemented, "not implemented")
}
