package main

import (
	"embed"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed templates/*
var f embed.FS

var Output = GetEnvWithDefault("output", "d:/tmp/")

func createProjectFolder(pName string) {
	_, err := os.Stat(pName)
	if os.IsNotExist(err) {
		os.MkdirAll(pName, 0755)
	}
	_, e1 := os.Stat(pName + "/scripts")
	if os.IsNotExist(e1) {
		os.MkdirAll(pName+"/scripts", 0755)
	}
}

func createSRCFolder(apiModel *ApiModel) {
	pName := Output + apiModel.ProjectName
	_, e1 := os.Stat(pName + "/src")
	if os.IsNotExist(e1) {
		os.MkdirAll(pName, 0755)
	}
	_, e2 := os.Stat(pName + "/src/middleware")
	if os.IsNotExist(e2) {
		os.MkdirAll(pName+"/src/middleware", 0755)
	}

	_, e3 := os.Stat(pName + "/src/route")
	if os.IsNotExist(e3) {
		os.MkdirAll(pName+"/src/route", 0755)
	}
	_, e5 := os.Stat(pName + "/src/model")
	if os.IsNotExist(e5) {
		os.MkdirAll(pName+"/src/model", 0755)
	}
	if includeModule(apiModel, "DB") {
		_, e4 := os.Stat(pName + "/src/db")
		if os.IsNotExist(e4) {
			os.MkdirAll(pName+"/src/db", 0755)
		}
	}
}

var ModuleList = []string{"JWT", "DB", "BASIC"}

type Modules []string

func (modules *Modules) Set(val string) error {
	*modules = strings.Split(val, ",")
	return nil
}

func (modules *Modules) String() string {
	str := "["
	for _, v := range *modules {
		str += v
	}
	return str + "]"
}

var (
	pName             string
	prefix            string
	port              int
	local             bool
	web               bool
	modules           Modules
	rootFiles         []string = []string{"/main.go", "/readme.md", "/init.cmd", "/init.sh", "/deploy.cmd", "/deploy.sh", "/destroy.cmd", "/destroy.sh", "/Dockerfile", "/go.mod"}
	middlewareFiles   []string = []string{"/auth.go", "/go.mod"}
	modelFiles        []string = []string{"/go.mod"}
	routeFiles        []string = []string{"/router.go", "/go.mod"}
	dbFiles           []string = []string{"/db.go", "/config.go", "/go.mod"}
	k8sFiles          []string = []string{"/service.yaml"}
	jwtMiddlewareFile string   = "/jwtauth.go"
	jwtRouteFile      string   = "/jwtroute.go"
	dbRouteFile       string   = "/dbroute.go"

	allFiles map[string][]string = map[string][]string{"/": rootFiles, "/src/middleware": middlewareFiles, "/src/route": routeFiles, "/src/model": modelFiles, "/service": k8sFiles}

	sampleDateTime = "20230507164000"
)

func main() {

	flag.StringVar(&pName, "project", "demo", "项目名称")
	flag.StringVar(&prefix, "prefix", "demo", "API起始路径")
	flag.IntVar(&port, "port", 9000, "启动端口")
	flag.BoolVar(&local, "local", false, "debug in local")
	flag.BoolVar(&web, "web", false, "start as web")
	flag.Var(&modules, "modules", "Modules that need be created, current include 'JWT', 'DB', 'BASIC','ALL', default only include 'BASIC'")
	flag.Parse()
	if len(modules) == 0 {
		modules = append(modules, "BASIC")
	} else {
		for _, x := range modules {
			if strings.ToUpper(x) == "ALL" {
				modules = ModuleList
				break
			}
		}
	}

	if web {
		e := echo.New()
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{"*"},
		}))
		e.POST("/cli/v1/init", func(c echo.Context) error {
			apiModel := &ApiModel{}
			if err := c.Bind(apiModel); err != nil {
				fmt.Println(err)
				return c.JSON(http.StatusNotAcceptable, "invalid post body")
			}

			if apiModel.ProjectName == "" {
				return c.String(http.StatusNotAcceptable, "missing projectName")
			}
			if apiModel.Prefix == "" {
				apiModel.Prefix = "api/v1"
			}
			if apiModel.Port == 0 {
				apiModel.Port = 9000
			}
			if apiModel.Modules == nil && len(apiModel.Modules) <= 0 {
				apiModel.Modules = []string{"BASIC"}
			}
			for _, x := range apiModel.Modules {
				if strings.ToUpper(x) == "ALL" {
					apiModel.Modules = ModuleList
					break
				}
			}

			create(apiModel)
			// zipName := apiModel.ProjectName + "-" + time.Now().Format(sampleDateTime) + ".zip"
			// err := ZipSource(apiModel.ProjectName, zipName)
			// if err != nil {
			// 	return c.String(http.StatusInternalServerError, "zip files error.")
			// }
			// os.RemoveAll("./" + apiModel.ProjectName)
			return c.String(http.StatusCreated, "project created at "+Output+apiModel.ProjectName)
		})
		e.Logger.Fatal(e.Start(":1323"))
	} else {
		apiModel := &ApiModel{
			Port:        port,
			Prefix:      prefix,
			ProjectName: pName,
			Modules:     modules,
		}
		create(apiModel)
	}

}

type ApiModel struct {
	Prefix      string
	ProjectName string
	Port        int
	ModuleName  string
	Modules     []string
}

func create(apiModel *ApiModel) {

	createProjectFolder(Output + apiModel.ProjectName)
	createSRCFolder(apiModel)
	if includeModule(apiModel, "JWT") {
		allFiles["/src/middleware"] = append(allFiles["/src/middleware"], jwtMiddlewareFile)
		allFiles["/src/route"] = append(allFiles["/src/route"], jwtRouteFile)
	}
	if includeModule(apiModel, "DB") {
		allFiles["/src/route"] = append(allFiles["/src/route"], dbRouteFile)
		allFiles["/src/db"] = dbFiles
	}
	for route, xfiles := range allFiles {
		// hanlde data
		if strings.Contains(route, "src") {
			ss := strings.Split(route, "/")
			apiModel.ModuleName = ss[len(ss)-1]

		} else {
			apiModel.ModuleName = ""
		}
		var path = Output + apiModel.ProjectName + route
		if route == "/" {
			path = Output + apiModel.ProjectName
		}
		if route == "/service" {
			path = Output + apiModel.ProjectName + "/scripts"
		}

		for _, fileName := range xfiles {
			createFile(fileName, path, apiModel)

		}
	}
}

func createFile(fileName string, dest string, apiModel *ApiModel) {
	tpl := getTpl(fileName)
	os.OpenFile(dest+fileName, os.O_CREATE, 0o666)
	file, err := os.OpenFile(dest+fileName, os.O_RDWR, 0o666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data := map[string]any{
		"Prefix":      apiModel.Prefix,
		"ProjectName": apiModel.ProjectName,
		"Port":        apiModel.Port,
		"ModuleName":  apiModel.ModuleName,
		"IncludeJWT":  includeModule(apiModel, "JWT"),
		"IncludeDB":   includeModule(apiModel, "DB"),
	}
	e := tpl.Execute(file, data)
	if e != nil {
		panic(e)
	}
}

func includeModule(apiModel *ApiModel, module string) bool {
	for _, x := range apiModel.Modules {
		if strings.ToUpper(x) == strings.ToUpper(module) {
			return true
		}

	}
	return false
}

func getTpl(prefix string) *template.Template {
	tplPath := "templates"
	if local {

		tpl, _ := template.ParseFiles(tplPath + prefix + ".tpl")
		return tpl
	} else {

		tpl, err := template.ParseFS(f, tplPath+prefix+".tpl")
		if err != nil {
			panic(err)
		}
		return tpl
	}
}
