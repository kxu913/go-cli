package main

import (
	"embed"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//go:embed templates/*
var f embed.FS

var (
	local   bool
	project string
	prefix  string
	host    string
	dbname  string
	port    int
	user    string
	pwd     string
	table   string
	web     bool
)

type ApiModel struct {
	ProjectName string
	Prefix      string
	Host        string
	DBname      string
	DBPort      int
	User        string
	Pwd         string
	Table       string
}

var Output = GetEnvWithDefault("output", "d:/tmp/")

func main() {

	flag.StringVar(&project, "project", "demo", "项目名称")
	flag.StringVar(&prefix, "prefix", "demo", "API起始路径")
	flag.StringVar(&host, "host", "localhost", "数据库Host")
	flag.StringVar(&dbname, "dbname", "workflow", "数据库")
	flag.StringVar(&user, "user", "postgres", "数据库用户名")
	flag.StringVar(&pwd, "pwd", "postgres", "数据库密码")
	flag.IntVar(&port, "port", 5432, "数据库端口")
	flag.StringVar(&table, "table", "table", "debug in local")
	flag.BoolVar(&web, "web", false, "start as web")
	flag.Parse()

	var dbConfig DBConfig
	var p string
	var t string

	if web {
		e := echo.New()
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{"*"},
		}))

		e.POST("/cli/v1/db/:table", func(c echo.Context) error {
			request := &ApiModel{}
			if e := c.Bind(request); e != nil {
				return c.String(http.StatusNotAcceptable, "invalid post body")
			}
			t = c.Param("table")
			dbConfig = DBConfig{
				Host:     request.Host,
				Port:     request.DBPort,
				DBName:   request.DBname,
				UserName: request.User,
				Password: request.Pwd,
			}
			p = request.ProjectName
			createOutputFolder(request.ProjectName)
			CreateTableGoFiles(request.ProjectName, request.Prefix, &dbConfig, c.Param("table"))
			ReplaceMainGo(p, c.Param("table"))
			return c.String(http.StatusCreated, "created")
		})
		e.Logger.Fatal(e.Start(":1324"))
	} else {

		p = project
		t = table
		dbConfig = DBConfig{
			Host:     host,
			Port:     port,
			DBName:   dbname,
			UserName: user,
			Password: pwd,
		}
	}

	createOutputFolder(p)

	// if strings.ToLower(table) == "all" {
	// 	tables := TableList(dbConfig)
	// 	for _, t := range tables {
	// 		CreateTableGoFiles(dbConfig, t)
	// 	}
	// } else {
	CreateTableGoFiles(p, prefix, &dbConfig, t)
	ReplaceMainGo(p, t)

	// }

}

func CreateTableGoFiles(p string, pre string, dbConfig *DBConfig, t string) {
	tableInfo := GetTableInfo(dbConfig, t)
	CreateModelGoFiles(p, tableInfo)
	CreateDBConifgGoFiles(p, dbConfig)
	CreateDBGoFiles(p, tableInfo)
	CreateRouteGoFiles(p, pre, tableInfo)
}

func CreateModelGoFiles(p string, tableInfo *TableInfo) {
	tpl := getTpl("model.go")

	file := WriteFile(fmt.Sprintf("%s%s/src/model/%s.go", Output, p, strings.ToLower(tableInfo.TableName)))

	defer file.Close()
	tpl.Execute(file, &tableInfo)

}
func CreateDBConifgGoFiles(p string, dbConfig *DBConfig) {
	// create db file
	tpl := getTpl("db.config.go")
	file := WriteFile(fmt.Sprintf("%s%s/src/db/config.go", Output, p))
	defer file.Close()
	tpl.Execute(file, dbConfig)

}

func CreateDBGoFiles(p string, tableInfo *TableInfo) {
	// create db file
	tpl := getTpl("db.instance.go")
	file := WriteFile(fmt.Sprintf("%s%s/src/db/%s_db.go", Output, p, strings.ToLower(tableInfo.TableName)))
	defer file.Close()
	tpl.Execute(file, map[string]any{
		"TableInfo":   &tableInfo,
		"ProjectName": p,
	})

}

func CreateRouteGoFiles(p string, pre string, tableInfo *TableInfo) {
	tpl := getTpl("route.instance.go")
	file := WriteFile(fmt.Sprintf("%s%s/src/route/%s_route.go", Output, p, strings.ToLower(tableInfo.TableName)))
	defer file.Close()
	tpl.Execute(file, map[string]any{
		"TableInfo":   &tableInfo,
		"Prefix":      pre,
		"ProjectName": p,
	})

}

func ReplaceMainGo(p string, t string) {
	mainName := fmt.Sprintf("%s%s/main.go", Output, p)
	b, err := ioutil.ReadFile(mainName)
	if err != nil {
		panic(err)
	}
	lineBytes := string(b)
	lines := strings.Split(lineBytes, "\n")
	var newLines []string
	for _, l := range lines {
		newLines = append(newLines, l)
		if strings.Index(l, "//auto generate groups") >= 0 {
			newLines = append(newLines, fmt.Sprintf("	route.%sGroup(e)\n", Ucfirst(t)))
		}
	}
	file, err := os.OpenFile(mainName, os.O_WRONLY, 0666)
	defer file.Close()
	file.WriteString(strings.Join(newLines, "\n"))

}

func createOutputFolder(p string) {
	root := Output + p
	_, err := os.Stat(root)
	if os.IsNotExist(err) {
		os.MkdirAll(root, 0755)
	}
	_, err0 := os.Stat(root + "/src")
	if os.IsNotExist(err0) {
		os.MkdirAll(root+"/src", 0755)
	}
	_, e1 := os.Stat(root + "/src/model")
	if os.IsNotExist(e1) {
		os.MkdirAll(root+"/src/model", 0755)
	}
	_, e2 := os.Stat(root + "/src/route")
	if os.IsNotExist(e2) {
		os.MkdirAll(root+"/src/route", 0755)
	}
	_, e3 := os.Stat(root + "/src/db")
	if os.IsNotExist(e3) {
		os.MkdirAll(root+"/src/db", 0755)
	}

}

func getTpl(name string) *template.Template {
	tplPath := "templates/"
	if local {

		tpl, _ := template.ParseFiles(tplPath + name + ".tpl")
		return tpl
	} else {
		tpl, err := template.ParseFS(f, tplPath+name+".tpl")
		if err != nil {
			panic(err)
		}
		return tpl
	}
}
