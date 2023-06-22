package route

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"{{.ProjectName}}/middleware"

)

func Login(c echo.Context) error {

	email := c.Request().Header.Get("email")
	pass := c.Request().Header.Get("password")

	token, expiredAt, err := middleware.CreateToken(email, pass)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Genreate Token Failed.")
	}
	return c.JSON(http.StatusOK, map[string]any{
		"Token":     token,
		"ExpiredAt": expiredAt,
	})
}


func JWTTokenApi(c echo.Context) error {
	email := c.Request().Header.Get("email")
	pass := c.Request().Header.Get("password")
	return c.JSON(http.StatusOK, map[string]any{
		"Email":    email,
		"Password": pass,
	})

}

func JWTGroup(e *echo.Echo) {
	t := e.Group("/{{.Prefix}}")
	t.POST("/login", Login, middleware.BasicAuth(DemoBasicAuth))
	t.GET("/info", JWTTokenApi,middleware.JwtAuth(DemoBasicAuth))
}




