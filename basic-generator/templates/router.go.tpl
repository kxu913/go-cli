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
	t := e.Group("{{.Prefix}}")
	t.GET("/auth", AuthApi, middleware.BasicAuth(DemoBasicAuth))
	t.GET("/noauth", NoAuthApi)

}

func DemoBasicAuth(email string, password string, ctx echo.Context) (bool, error) {

	//Implement using your basic auth
	if(email=="kevin" && password=="888888"){
		ctx.Request().Header.Add("email", email)
		ctx.Request().Header.Add("password", password)
		return true,nil
	}

	return false,nil
}



