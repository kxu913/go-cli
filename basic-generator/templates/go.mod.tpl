module {{.ProjectName}}{{if .ModuleName}}/{{.ModuleName}}{{end}}

go 1.19

require (
    github.com/labstack/echo/v4 v4.10.2
    {{if .IncludeDB}}
    gorm.io/driver/postgres v1.5.0
	gorm.io/gorm v1.25.0 
    {{end}}
)
