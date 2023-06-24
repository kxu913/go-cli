package model

type ServiceRequest struct {
	Basic   BasicRequest
	DB      DBRequest
	Graphql GraphqlRequest
}

type BasicRequest struct {
	ProjectName string
	Prefix      string
	Port        int
	Modules     []string
}

type DBRequest struct {
	ProjectName string
	Prefix      string
	Host        string
	DBName      string
	DBPort      int
	User        string
	Pwd         string
	Table       string
}

func SetDBRequestFromBasicRequest(request *ServiceRequest) DBRequest {
	dbRequest := request.DB
	dbRequest.ProjectName = request.Basic.ProjectName
	dbRequest.Prefix = request.Basic.Prefix
	return dbRequest

}

type GraphqlRequest struct {
	ProjectName      string
	Prefix           string
	QueryName        string
	QueryDescription string
	SQL              string
}

func SetGraphqlRequestFromBasicRequest(request *ServiceRequest) GraphqlRequest {
	graphqlRequest := request.Graphql
	graphqlRequest.ProjectName = request.Basic.ProjectName
	graphqlRequest.Prefix = request.Basic.Prefix
	return graphqlRequest

}
