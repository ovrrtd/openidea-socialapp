package restapi

import (
	"net/http"
	httpHelper "socialapp/internal/helper/http"

	"github.com/labstack/echo/v4"
)

func (r *Restapi) UploadImage(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, http.StatusBadRequest, "", nil, nil, err)
	}

	imgUrl, code, err := r.service.UploadImage(c.Request().Context(), file)
	r.debugError(err)
	if err != nil {
		return httpHelper.ResponseJSONHTTP(c, code, "", nil, nil, err)
	}
	return httpHelper.ResponseJSONHTTP(c, code, "File uploaded sucessfully", map[string]string{"imageUrl": imgUrl}, nil, err)
}
