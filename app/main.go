package main

import (
	"fmt"
	"net/http"
	"os"
	"task-api/app/config"
	"task-api/foundation/database"
	"task-api/handlers"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

var validate = validator.New()

type CustomValidator struct {
	validator *validator.Validate
}

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		fmt.Printf("Error While parsing Config %v", err)
		return
	}

	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.999Z07:00"
	zerolog.LevelFieldName = "level_name"
	logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
	if cfg.LogLevel == "debug" {
		logger = logger.Level(zerolog.DebugLevel)
	} else if cfg.LogLevel == "error" {
		logger = logger.Level(zerolog.ErrorLevel)
	} else if cfg.LogLevel == "fatal" {
		logger = logger.Level(zerolog.FatalLevel)
	} else if cfg.LogLevel == "info" {
		logger = logger.Level(zerolog.InfoLevel)
	}

	db, err := database.Connect(cfg.Dns)
	if err != nil {
		fmt.Printf("unable to connect ")
		return
	}
	defer db.Close()

	e := echo.New()
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.AddTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: fmt.Sprintf(`{"level_name": "%s", "time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",`+
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",`+
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}",`+
			`"bytes_in":${bytes_in},"bytes_out":${bytes_out}}`+"\n", cfg.LogLevel),
	}))
	e.Validator = &CustomValidator{validator: validate}

	handlers.API(e, cfg, db, logger)

	e.Start(fmt.Sprintf(":%s", cfg.Port))
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
