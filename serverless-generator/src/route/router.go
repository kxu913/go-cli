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
	operator.DestroyNS(c.Param("ns"))
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
	operator.CreateVirtualService(ns, &requestModel.MetaData, requestModel.Container.Port)
	return nil

}

func buildImage(project string, provider string, w *echo.Response) (string, int) {
	tag, port := operator.BuildImage(project, provider)
	writeLine(w, fmt.Sprintf("%s built.", project))
	operator.PushImageToRemote(provider, tag)
	writeLine(w, fmt.Sprintf("%s uploaded.", tag))
	return tag, port

}

func Deploy(c echo.Context) error {
	requestModel := &RequestModel{}
	if err := c.Bind(requestModel); err != nil {
		return c.String(http.StatusBadRequest, "Invalid Request Body.")
	}
	project := requestModel.MetaData.Name
	cloudProvider := requestModel.MetaData.CloudProvider
	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	buildImage(project, cloudProvider, w)
	tag, port := buildImage(project, cloudProvider, w)
	ns := operator.CreateNS(project)
	writeLine(w, fmt.Sprintf("%s created.", ns.GetName()))
	gw := operator.CreateGateway(project)
	writeLine(w, fmt.Sprintf("%s created.", gw.GetName()))
	requestModel.Container.Image = tag
	requestModel.Container.Port = int32(port)
	acct := operator.CreateServiceAccount(project, &requestModel.MetaData)
	writeLine(w, fmt.Sprintf("%s created.", acct.GetName()))
	svc := operator.CreateService(project, &requestModel.MetaData, requestModel.Container.Port)
	writeLine(w, fmt.Sprintf("%s created.", svc.GetName()))
	deployment := operator.DeployService(project, &requestModel.MetaData, requestModel.Replicas, &requestModel.Container)
	writeLine(w, fmt.Sprintf("%s deployed.", deployment.GetName()))
	operator.CreateVirtualService(project, &requestModel.MetaData, requestModel.Container.Port)
	writeLine(w, "Istio VirtualService created")

	return nil
}

func Destroy(c echo.Context) error {
	operator.DestroyNS(c.Param("ns"))
	return c.String(http.StatusOK, fmt.Sprintf("%s destroyed.", c.Param("ns")))
}

func Group(e *echo.Echo) {
	t := e.Group("/cli")
	t.POST("/ns/:ns", CreateNS)
	t.POST("/svc/:ns", CreateService)
	t.POST("/deploy", Deploy)
	t.DELETE("/destory/:ns", Destroy)
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

func writeLine(w *echo.Response, msg string) {
	fmt.Fprintf(w, "data: %s\n\n", msg)
	w.Flush()
}
