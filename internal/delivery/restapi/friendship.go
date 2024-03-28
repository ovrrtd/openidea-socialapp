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

func (r *Restapi) CreateFriendship(c echo.Context) error {
	var request request.CreateFriendship
	err := c.Bind(&request)
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
	}
	userID := c.Get(string(common.EncodedUserJwtCtxKey)).(*response.User).ID
	request.AddedBy = int64(userID)
	code, err := r.service.CreateFriendship(c.Request().Context(), request)
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
}

func (r *Restapi) DeleteFriendship(c echo.Context) error {
	var request request.DeleteFriendship
	err := c.Bind(&request)
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
	}
	userID := c.Get(string(common.EncodedUserJwtCtxKey)).(*response.User).ID
	request.Friend2 = int64(userID)
	code, err := r.service.DeleteFriendship(c.Request().Context(), request)
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
}

func (r *Restapi) FindAllFriend(c echo.Context) error {
	var request request.FindAllFriendships
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

	if urlValues.Has("sortBy") {
		request.SortBy = urlValues.Get("sortBy")
		if request.SortBy != "createdAt" && request.SortBy != "friendCount" {
			return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
		}

	} else {
		if request.SortBy == "" {
			request.SortBy = "createdAt"
		}
	}

	if urlValues.Has("orderBy") {
		request.OrderBy = urlValues.Get("orderBy")
		if request.OrderBy != "asc" && request.OrderBy != "desc" {
			return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
		}
	} else {
		request.OrderBy = "desc"
	}

	if urlValues.Has("onlyFriend") {
		of := urlValues.Get("onlyFriend")
		if of != "true" && of != "false" {
			return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, nil)
		}
		if of == "true" {
			request.OnlyFriend = true
		}

		if of == "false" {
			request.OnlyFriend = false
		}
	} else {
		request.OnlyFriend = false
	}

	if urlValues.Has("search") {
		request.Search = urlValues.Get("search")
	}

	userId := c.Get(string(common.EncodedUserJwtCtxKey)).(*response.User).ID

	request.UserID = int64(userId)
	r.log.Debug().Msgf("request: %+v", request)

	ret, meta, code, err := r.service.FindAllFriendships(c.Request().Context(), request)
	r.debugError(err)
	return httpHelper.ResponseJSONHTTP(c, code, "", ret, meta, err)
}
