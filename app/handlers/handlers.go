package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"task-api/app/config"
	"task-api/buisness/task"
	"task-api/buisness/user"
	"task-api/mid"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func API(e *echo.Echo, cfg config.Config, db *sqlx.DB, log zerolog.Logger) {
	m := mid.Mid{
		TokenSeceret: cfg.JwtTokenSeceret,
	}

	taskHandlers := taskHandlers{
		task: task.New(db, log),
	}
	api := e.Group("/api/v1")
	taskG := api.Group("/task")
	taskG.POST("", taskHandlers.create, m.Authentication)
	taskG.GET("", taskHandlers.query, m.Authentication)
	taskG.GET("/:task_id", taskHandlers.queryByID, m.Authentication)
	taskG.PUT("/:task_id", taskHandlers.update, m.Authentication)
	taskG.DELETE("/:task_id", taskHandlers.delete, m.Authentication)
	taskG.PATCH("/mark-as-done", taskHandlers.updateStatus, m.Authentication)

	userHandlers := userHandlers{
		userService: user.New(db, cfg.JwtTokenSeceret, log),
	}
	userG := api.Group("/user")
	userG.POST("/login", userHandlers.login)
	userG.POST("", userHandlers.create)
	userG.GET("", userHandlers.queryByID, m.Authentication)
	userG.PUT("/:user_id", userHandlers.update, m.Authentication)
	userG.DELETE("/:user_id", userHandlers.delete, m.Authentication)
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func badRequest(c echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, ErrorResponse{Message: message})
}

func Error(c echo.Context, err error, resourceName string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: fmt.Sprintf("Request %s Not Found", resourceName)})
	}
	c.Logger().Error(err)
	return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Internal Server Error"})
}

func getUserID(c echo.Context) int {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return 0
	}

	return userID
}
