package middleware

import (
	"net/http"
	"os"
	"socialapp/internal/helper/common"
	"socialapp/internal/helper/errorer"
	httpHelper "socialapp/internal/helper/http"
	"socialapp/internal/helper/jwt"
	"socialapp/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type middleware struct {
	logger  zerolog.Logger
	service service.Service
}

type Middleware interface {
	Authentication(isThrowError bool) func(next echo.HandlerFunc) echo.HandlerFunc
}

func New(logger zerolog.Logger, service service.Service) Middleware {
	return &middleware{
		logger:  logger,
		service: service,
	}
}

func (m *middleware) Authentication(isThrowError bool) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			m.logger.Info().Msg("Authentication")
			token := httpHelper.GetJWTFromRequest(c.Request())

			if token == "" && isThrowError {
				return httpHelper.ResponseJSONHTTP(c, http.StatusUnauthorized, "", nil, nil, errorer.ErrUnauthorized)
			}

			if token != "" {
				claims := &common.UserClaims{}
				err := jwt.VerifyJwt(token, claims, os.Getenv("JWT_SECRET"))
				if err != nil {
					if err == errorer.ErrUnauthorized {
						return httpHelper.ResponseJSONHTTP(c, http.StatusUnauthorized, "", nil, nil, errorer.ErrUnauthorized)
					}
					return httpHelper.ResponseJSONHTTP(c, http.StatusForbidden, "", nil, nil, errorer.ErrForbidden)
				}

				usr, code, err := m.service.GetUserByID(c.Request().Context(), claims.Id)
				if err != nil {
					return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
				}
				c.Set(common.EncodedUserJwtCtxKey.ToString(), usr)
			}

			return next(c)
		}
	}
}
