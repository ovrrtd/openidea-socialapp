package restapi

import (
	"net/http"
	"socialapp/internal/helper/common"
	httpHelper "socialapp/internal/helper/http"
	"socialapp/internal/model/request"
	"socialapp/internal/model/response"

	"github.com/labstack/echo/v4"
)

func (r *Restapi) Register(c echo.Context) error {
	var request request.Register
	err := c.Bind(&request)
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
	}

	ret, code, err := r.service.Register(c.Request().Context(), request)
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "User registered successfully", ret, nil, err)
}

func (r *Restapi) Login(c echo.Context) error {
	var request request.Login
	err := c.Bind(&request)
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
	}
	ret, code, err := r.service.Login(c.Request().Context(), request)
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "User logged successfully", ret, nil, err)
}

func (r *Restapi) UpdateAccount(c echo.Context) error {
	var request request.UpdateAccount
	err := c.Bind(&request)
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
	}

	request.ID = c.Get(string(common.EncodedUserJwtCtxKey)).(*response.User).ID

	_, code, err := r.service.UpdateAccount(c.Request().Context(), request)
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
}

func (r *Restapi) LinkEmail(c echo.Context) error {
	var request request.LinkEmail
	err := c.Bind(&request)
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
	}

	request.ID = c.Get(string(common.EncodedUserJwtCtxKey)).(*response.User).ID

	_, code, err := r.service.LinkEmail(c.Request().Context(), request)
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
}

func (r *Restapi) LinkPhone(c echo.Context) error {
	var request request.LinkPhone
	err := c.Bind(&request)
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
	}

	request.ID = c.Get(string(common.EncodedUserJwtCtxKey)).(*response.User).ID

	_, code, err := r.service.LinkPhone(c.Request().Context(), request)
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
}
