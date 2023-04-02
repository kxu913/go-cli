package main

import (
	"embed"
	"flag"
	"os"
	"strings"
	"text/template"
)

//go:embed templates/*
var f embed.FS

func createProjectFolder() {
	_, err := os.Stat(pName)
	if os.IsNotExist(err) {
		os.MkdirAll(pName, 0755)
	}
	_, e1 := os.Stat(pName + "/scripts")
	if os.IsNotExist(e1) {
		os.MkdirAll(pName+"/scripts", 0755)
	}
}

func createSRCFolder() {
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
}

var (
	pName           string
	prefix          string
	port            int
	local           bool
	rootFiles       []string = []string{"/main.go", "/readme.md", "/init.cmd", "/init.sh", "/deploy.cmd", "/deploy.sh", "/destroy.cmd", "/destroy.sh", "/Dockerfile", "/go.mod"}
	middlewareFiles []string = []string{"/auth.go", "/go.mod"}
	routeFiles      []string = []string{"/router.go", "/go.mod"}
	k8sFiles        []string = []string{"/service.yaml"}

	allFiles map[string][]string = map[string][]string{"/": rootFiles, "/src/middleware": middlewareFiles, "/src/route": routeFiles, "/service": k8sFiles}
)

func main() {

	flag.StringVar(&pName, "project", "demo", "项目名称")
	flag.StringVar(&prefix, "prefix", "demo", "API起始路径")
	flag.IntVar(&port, "port", 9000, "启动端口")
	flag.BoolVar(&local, "local", false, "debug in local")
	flag.Parse()
	data := map[string]any{
		"Port":        port,
		"Prefix":      prefix,
		"ProjectName": pName,
	}
	createProjectFolder()
	createSRCFolder()
	for route, xfiles := range allFiles {
		// hanlde data
		if strings.Contains(route, "src") {
			ss := strings.Split(route, "/")
			data["ModuleName"] = ss[len(ss)-1]

		} else {
			data["ModuleName"] = ""
		}
		var path = pName + route
		if route == "/" {
			path = pName
		}
		if route == "/service" {
			path = pName + "/scripts"
		}

		for _, fileName := range xfiles {
			createFile(fileName, path, data)

		}

	}

}

func createFile(fileName string, dest string, data map[string]any) {
	tpl := getTpl(fileName)
	file, err := os.OpenFile(dest+fileName, os.O_CREATE, 0o666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	tpl.Execute(file, data)
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
