package restapi

import (
	"net/http"
	"socialapp/internal/helper/common"
	httpHelper "socialapp/internal/helper/http"
	"socialapp/internal/model/request"
	"socialapp/internal/model/response"

	"github.com/labstack/echo/v4"
)

func (r *Restapi) CreatePost(c echo.Context) error {
	var request request.CreatePost
	err := c.Bind(&request)
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
	}
	userId := c.Get(string(common.EncodedUserJwtCtxKey)).(*response.User).ID

	request.UserID = int64(userId)

	code, err := r.service.CreatePost(c.Request().Context(), request)
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
}

func (r *Restapi) CreateComment(c echo.Context) error {
	var request request.CreateComment
	err := c.Bind(&request)
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
	}

	code, err := r.service.CreateComment(c.Request().Context(), request)
	return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
}

func (r *Restapi) FindAll(c echo.Context) error {
	var request request.FindAllPost
	err := c.Bind(&request)
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
	}
	if request.Limit == 0 {
		request.Limit = 10
	}
	if request.Offset == 0 {
		request.Offset = 0
	}

	userId := c.Get(string(common.EncodedUserJwtCtxKey)).(*response.User).ID

	request.UserID = int64(userId)

	ret, meta, code, err := r.service.FindAllPost(c.Request().Context(), request)
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "", ret, meta, err)
}
