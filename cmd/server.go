package cmd

import (
	"fmt"
	"os"
	"os/signal"
	database "socialapp/db"
	mw "socialapp/internal/delivery/middleware"
	"socialapp/internal/delivery/restapi"
	"socialapp/internal/repository"
	"socialapp/internal/service"
	"strconv"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	APP_PORT = "8080"
)

func Server() error {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	logger := zerolog.New(os.Stdout)
	db, err := database.NewDBDefaultSql()
	if err != nil {
		logger.Info().Msg(fmt.Sprintf("Postgres connection error: %s", err.Error()))
		return err
	}
	err = db.Ping()
	if err != nil {
		logger.Info().Msg(fmt.Sprintf("Postgres ping error: %s", err.Error()))
		return err
	}
	defer db.Close()

	// repository init
	userRepo := repository.NewUserRepository(logger, db)
	s3Repo := repository.NewS3Repository(logger)
	postRepo := repository.NewPostRepository(logger, db)
	friendshipRepo := repository.NewFriendshipRepository(logger, db)

	salt, err := strconv.Atoi(os.Getenv("BCRYPT_SALT"))
	if err != nil {
		salt = 8
	}
	// service registry
	service := service.New(
		service.Config{Salt: salt, JwtSecret: os.Getenv("JWT_SECRET")},
		logger,
		userRepo,
		s3Repo,
		postRepo,
		friendshipRepo,
	)

	// middleware init
	md := mw.New(logger, service)

	// restapi init
	rest := restapi.New(logger, md, service)

	// echo server
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Str("method", c.Request().Method).
				Int("status", v.Status).
				Msg("request")
			return nil
		},
	}))

	// add restapi route
	rest.MakeRoute(e)

	errs := make(chan error)
	go func() {
		logger.Log().Msg(fmt.Sprintf("start server on port %s", APP_PORT))
		errs <- e.Start(fmt.Sprintf(":%s", APP_PORT))
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	return <-errs
}
