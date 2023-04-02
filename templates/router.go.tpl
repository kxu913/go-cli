package route

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"{{.ProjectName}}/middleware"

)



func NoAuthApi(c echo.Context) error {
	return c.String(http.StatusOK, "Well done.")

}

func AuthApi(c echo.Context) error {

	return c.String(http.StatusOK, "Well done.")

}

func Group(e *echo.Echo) {
	t := e.Group("/{{.Prefix}}")

	t.GET("/auth", AuthApi, middleware.CustomAuth())
	t.GET("/noauth", NoAuthApi)

}
