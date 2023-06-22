package route

import (
	"graphql-generator/graphql"
	"graphql-generator/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ParseSQL(c echo.Context) error {
	request := &model.Query{}
	if err := c.Bind(request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid Post Body")
	}
	r := graphql.ParseSQL(request)
	graphql.CreateGraphqlDefiniton(request.ProjectName, r)
	graphql.GenerateRouterFunc(request.ProjectName, r)
	graphql.ReplaceMain(request.ProjectName, request)

	return c.JSON(http.StatusOK, r)

}

func DeleteApi(c echo.Context) error {
	graphql.DeleteFiles(c.Param("project"), c.Param("router"))
	graphql.UpdateMain(c.Param("project"), c.Param("router"))
	return c.String(http.StatusOK, "Done")

}

func Group(e *echo.Echo) {
	t := e.Group("/graphql/v1")
	t.POST("/sql", ParseSQL)
	t.DELETE("/sql/:project/:router", DeleteApi)

}
