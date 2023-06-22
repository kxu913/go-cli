package route

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"{{.ProjectName}}/middleware"
	"{{.ProjectName}}/db"

)

func GetTableList(c echo.Context) error {
	return c.JSON(http.StatusOK, db.TableList())
}

func DBGroup(e *echo.Echo) {
	t := e.Group("/{{.Prefix}}")
	t.GET("/db", GetTableList,middleware.BasicAuth(DemoBasicAuth))
}




