package util

import (
	"net/http"
	"socialapp/internal/helper/common"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func ResponseJSONHTTP(c echo.Context, code int, msg string, data interface{}, meta *common.Meta, err error) error {
	res := map[string]interface{}{
		"data":    data,
		"message": strings.ToLower(http.StatusText(code)),
	}
	if err != nil {
		res["message"] = errors.Cause(err).Error()
	} else {
		if msg != "" {
			res["message"] = msg
		}
	}

	if meta != nil {
		res["meta"] = meta
	}

	return c.JSON(code, res)
}

func GetJWTFromRequest(r *http.Request) string {
	// From query.
	query := r.URL.Query().Get("jwt")
	if query != "" {
		return query
	}

	// From header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}

	// From cookie.
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return ""
	}

	return cookie.Value
}
