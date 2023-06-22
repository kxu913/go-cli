package route
import (
	"{{.ProjectName}}/graphql"

	"net/http"

	"github.com/labstack/echo/v4"
)

func Get{{.QueryName}}Graphql(c echo.Context) error {
	request := c.Request()
	query := request.URL.Query().Get("query")

	results := graphql.ExecuteQuery(query, graphql.{{.QueryName}}Schema)

	return c.JSON(http.StatusOK, results)
}

func {{.QueryName}}GroupGraphql(e *echo.Echo) {
	t := e.Group("/graphql/v1")
	t.GET("/{{.RouterName}}", Get{{.QueryName}}Graphql)


}