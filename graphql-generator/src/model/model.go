package model

type Query struct {
	ProjectName      string
	QueryName        string
	QueryDescription string
	SQL              string
	Fields           []Field
}

type Field struct {
	Name        string
	Alias       string
	DBType      string
	GraphqlType string
}
