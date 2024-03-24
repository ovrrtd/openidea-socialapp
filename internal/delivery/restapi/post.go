package restapi

import (
	"net/http"
	"socialapp/internal/helper/common"
	httpHelper "socialapp/internal/helper/http"
	"socialapp/internal/model/request"
	"socialapp/internal/model/response"
	"strconv"

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
	r.log.Debug().Msgf("request: %+v", request)
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
	userId := c.Get(string(common.EncodedUserJwtCtxKey)).(*response.User).ID

	request.UserID = int64(userId)

	code, err := r.service.CreateComment(c.Request().Context(), request)
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
}

func (r *Restapi) FindAll(c echo.Context) error {
	var request request.FindAllPost

	urlValues := c.QueryParams()
	if urlValues.Has("limit") {
		limit, err := strconv.Atoi(urlValues.Get("limit"))
		if err != nil {
			return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
		}
		request.Limit = limit
		if request.Limit < 0 {
			return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
		}
	} else {
		request.Limit = 10
	}

	if urlValues.Has("offset") {
		offset, err := strconv.Atoi(urlValues.Get("offset"))
		if err != nil {
			return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
		}
		request.Offset = offset
		if request.Offset < 0 {
			return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
		}
	} else {
		request.Offset = 0
	}
	request.Tags = []string{}
	if urlValues.Has("searchTag") {
		st := urlValues["searchTag"]
		for _, t := range st {
			if t == "" {
				return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
			}
			request.Tags = append(request.Tags, t)
		}
	}

	userId := c.Get(string(common.EncodedUserJwtCtxKey)).(*response.User).ID

	request.UserID = int64(userId)

	ret, meta, code, err := r.service.FindAllPost(c.Request().Context(), request)
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "", ret, meta, err)
}
