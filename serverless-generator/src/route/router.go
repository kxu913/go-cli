package route

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"serverless-generator/model"
	"serverless-generator/operator"
)

type RequestModel struct {
	Replicas  int
	MetaData  model.MetaData
	Container model.Container
}

func CreateNS(c echo.Context) error {
	ns := operator.CreateNS(c.Param("ns"))
	gw := operator.CreateGateway(c.Param("ns"))

	return c.JSON(http.StatusOK, map[string]any{
		"Namesapce": ns,
		"Gateway":   gw,
	})

}

func CreateService(c echo.Context) error {
	requestModel := &RequestModel{}
	ns := c.Param("ns")
	if err := c.Bind(requestModel); err != nil {
		return c.String(http.StatusBadRequest, "Invalid Request Body.")
	}
	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	acct := operator.CreateServiceAccount(ns, &requestModel.MetaData)
	fmt.Fprintf(w, "data: ServiceAccount %s.\n\n", acct)
	w.Flush()
	svc := operator.CreateService(ns, &requestModel.MetaData, requestModel.Container.Port)
	fmt.Fprintf(w, "data: Service status is %s.\n\n", &svc.Status)
	w.Flush()
	deployment := operator.DeployService(ns, &requestModel.MetaData, requestModel.Replicas, &requestModel.Container)
	fmt.Fprintf(w, "data: Deployment %s.\n\n", deployment)
	w.Flush()
	vs := operator.CreateVirtualService(ns, &requestModel.MetaData, requestModel.Container.Port)
	fmt.Fprintf(w, "data: VirtualService status is %s.\n\n", &vs.Status)
	w.Flush()
	return nil

}
func Group(e *echo.Echo) {
	t := e.Group("/cli")
	t.POST("/ns/:ns", CreateNS)
	t.POST("/svc/:ns", CreateService)
}

func DemoBasicAuth(email string, password string, ctx echo.Context) (bool, error) {

	//Implement using your basic auth
	if email == "kevin" && password == "888888" {
		ctx.Request().Header.Add("email", email)
		ctx.Request().Header.Add("password", password)
		return true, nil
	}

	return false, nil
}
