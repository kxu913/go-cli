package main

import (
	"{{.ProjectName}}/route"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	route.Group(e)
	//auto generate groups

	{{if .IncludeJWT}}route.JWTGroup(e){{end}}
	{{if .IncludeDB}}route.DBGroup(e){{end}}
	e.Logger.Fatal(e.Start(":{{.Port}}"))

}
