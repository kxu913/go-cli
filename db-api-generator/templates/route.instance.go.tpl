package route

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
	"{{.ProjectName}}/db"
	"{{.ProjectName}}/model"

)

func Create{{.TableInfo.TableName}}(c echo.Context) error {
	m := &model.{{.TableInfo.TableName}}{}

	if err := c.Bind(m); err != nil {
		return c.JSON(http.StatusNotAcceptable, "invalid post body")
	}
	return c.JSON(http.StatusOK, db.Create{{.TableInfo.TableName}}(m))
}

func Get{{.TableInfo.TableName}}(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid {{.TableInfo.TableName}} Id")
	}
	return c.JSON(http.StatusOK, db.Get{{.TableInfo.TableName}}(id))
}

func Get{{.TableInfo.TableName}}s(c echo.Context) error {
	return c.JSON(http.StatusOK, db.Get{{.TableInfo.TableName}}s())
}

func {{.TableInfo.TableName}}Group(e *echo.Echo) {
	group:= strings.ToLower("{{.TableInfo.TableName}}")
	groups:= strings.ToLower("{{.TableInfo.TableName}}")+"s"
	t := e.Group("{{.Prefix}}/")
	t.POST(group, Create{{.TableInfo.TableName}})
	t.GET(group+"/:id", Get{{.TableInfo.TableName}})
	t.GET(groups, Get{{.TableInfo.TableName}}s)
}




