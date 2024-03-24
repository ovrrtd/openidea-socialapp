package restapi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "socialapperce",
		Help:    "Histogram of socialapp server request duration.",
		Buckets: prometheus.LinearBuckets(1, 1, 10), // Adjust bucket sizes as needed
	}, []string{"path", "method", "status"})
)

func (r *Restapi) MakeRoute(e *echo.Echo) {
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// user
	NewRoute(e, http.MethodPatch, "/v1/user", r.UpdateAccount, r.middleware.Authentication(true))
	NewRoute(e, http.MethodPost, "/v1/user/link", r.LinkEmail, r.middleware.Authentication(true))
	NewRoute(e, http.MethodPost, "/v1/user/link/phone", r.LinkPhone, r.middleware.Authentication(true))
	NewRoute(e, http.MethodPost, "/v1/user/register", r.Register)
	NewRoute(e, http.MethodPost, "/v1/user/login", r.Login)
	// friendship
	NewRoute(e, http.MethodGet, "/v1/friend", r.FindAllFriend, r.middleware.Authentication(true))
	NewRoute(e, http.MethodPost, "/v1/friend", r.CreateFriendship, r.middleware.Authentication(true))
	NewRoute(e, http.MethodDelete, "/v1/friend", r.DeleteFriendship, r.middleware.Authentication(true))
	// image
	NewRoute(e, http.MethodPost, "/v1/image", r.UploadImage, r.middleware.Authentication(true))
	// post
	NewRoute(e, http.MethodPost, "/v1/post", r.CreatePost, r.middleware.Authentication(true))
	NewRoute(e, http.MethodGet, "/v1/post", r.FindAll, r.middleware.Authentication(true))
	NewRoute(e, http.MethodPost, "/v1/post/comment", r.CreateComment, r.middleware.Authentication(true))

}

func NewRoute(app *echo.Echo, method string, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) {
	app.Add(method, path, wrapHandlerWithMetrics(path, method, handler), middleware...)
}

func wrapHandlerWithMetrics(path string, method string, handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		startTime := time.Now()

		// Execute the actual handler and catch any errors
		err := handler(c)

		// Regardless of whether an error occurred, record the metrics
		duration := time.Since(startTime).Seconds()

		requestHistogram.WithLabelValues(path, method, strconv.Itoa(c.Response().Status)).Observe(duration)
		return err
	}
}
