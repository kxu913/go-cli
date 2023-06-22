# Go-cli
**`Go-cli` is a project that base on `golang`+`echo` to create a microservice，if you also have `Docker` and `Kubernetes+Istio` environemtn，you can use apis that in `serverless-generator` to deploy the service to Kubernetes and Istio.**

## Demo
<img src="public/cli-demo.gif" width="100%" height="100%">

## Quick Start
- Recommend use docker-compose to start.<br>
  Clone project and Run `docker-compose up -d` in scripts folder.
- Clone project<br>
  Start `basic-generator`、`db-api-generator`、`graphql-generator`、`microservice-generator`


## Module introduce.
*Why separate to different module? because you can determine how to combine or generate codes.*
### microservice-generator
Create a complete micro service，return a zip file. Module parameters see below：
- Basic：Detail introduce see `basic-generator`.
- DB：Detail introduce see `db-api-generator`.
- Graphql: Optional，if the section is null, the code will not generate for graphql query，Detail introduce see `graphql-generator`
- Eg.
```
curl --location 'http://localhost:1325/it/v1/create' \
--header 'Content-Type: application/json' \
--data '{
    "Basic": {
        "Prefix": "/api/v1",
        "ProjectName": "it",
        "Port": 8778,
        "Modules": [
            "DB"
        ]
    },
    "DB": {
        "Host": "172.22.192.1",
        "DBName": "workflow",
        "DBPort": 5432,
        "User": "postgres",
        "Pwd": "postgres",
        "Table": "task"
    },
    "Graphql": {
        "QueryName": "workflow",
        "QueryDescription": "Get Workflow",
        "SQL": "SELECT workflow.id AS workflow_id, workflow.status AS workflow_status FROM workflow"
    }
}'
```

### basic-generator
Create a basic project，you also can install it as a cli to run locally. Parameter as below：
- ProjectName：Project name.
- Prefix：Api prefix.
- Port: Run port.
- Modules：Generate modules，valid value ['JWT', 'DB', 'BASIC','ALL']
- Eg.
```
curl --location 'http://localhost:1323/cli/v1/init' \
--header 'Content-Type: application/json' \
--data '{
    "ProjectName": "it",
    "Prefix": "/api/v1",
    "Port": 8777,
    "Modules": [
        "DB"
    ]
}'
```

### db-api-generator
base on project which created by basic-generator，create CRUD apis for the table ，need provide valid parameters that used to connect database. Parameter as below：
- ProjectName：Same as above.
- Prefix：Same as above.
- Host: db host.
- DBname：db name.
- DBPort：db port.
- User: db user.
- Pwd: db password.
- Eg.
```
curl --location 'http://localhost:1324/cli/v1/db/workflow' \
--header 'Content-Type: application/json' \
--data '{
    "ProjectName": "it",
    "Prefix": "/api/v1",
    "Host": "172.22.192.1",
    "DBname": "workflow",
    "DBPort": 5432,
    "User": "postgres",
    "Pwd":"postgres"
}'
```

### graphql-generator
base on project which created by basic-generator，create query api that use graphql，need provide sql that use to query data., it will reuse db config that create by db-api-generator, Parameter as below：
- ProjectName：Same as above.
- QueryName：query name of graphql, unique.
- QueryDescription: query description.
- SQL：sql that used to query data from db.
- Eg.
```
curl --location 'http://localhost:9004/graphql/v1/sql' \
--header 'Content-Type: application/json' \
--data '{
    "ProjectName": "it",
    "QueryName": "task",
    "QueryDescription": "Get Task",
    "SQL":"SELECT task.id AS task_id,  workflow.id AS workflow_id, task.status AS task_status, task.started_time AS task_start_time, workflow.status AS workflow_status FROM task INNER JOIN workflow ON task.workflow_id=workflow.id"
}'
```

### serverless-generator
Use API to deploy docker image to Kubernetes，and create api in Istio：
#### Create namespace and Gateway.
- `http://localhost:9003/cli/ns/{ns}`
- Eg.
```
curl --location --request POST 'http://localhost:9003/cli/ns/demo'
```
#### Deploy service to Kubernetes，and inject api to gateway.
- `http://localhost:9003/cli/svc/{ns}`
- MetaData：
    - Name: Kubernetes service name.
    - Version: Kubernetes service version.
    - Prefix： Api prefix.
- Container：
    - Image：docker image.
    - Port：container port.
    - ForceUpdate：if `true`，Kubernetes will alway push image during deployment.
    - RunAsRoot：if `true`, Kubernetes will run the container as `root`.
    - Replicas：container replicas.
    - Environments：environments of container.
- Eg.
```
curl --location 'http://localhost:9003/cli/svc/demo' \
--header 'Content-Type: application/json' \
--data '{
    "MetaData": {
        "Name": "demo",
        "Version": "v1",
        "Prefix":"api/v1"
    },
    "Container": {
        "Image": "demo",
        "Port": 9999,
        "ForceUpdate": false,
        "RunAsRoot": false,
        "Environments": [
            {
                "name": "db_host",
                "value": "172.27.64.1"
            }
        ]
    },
    "Replicas": 1
}'
```

## Author
[Kevin Xu](http://kevin913.com.cn/about)<br />
<img src="public/my.jpg" width="200" height="200"><br/>




