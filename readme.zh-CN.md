# Go-cli
**Go-cli是一个基于golang+echo的一个快速创建微服务的工具，如果你拥有Docker和Kubernetes+Istio的环境，也可以通过serverless-generator中包含的Api快速的生成基于Istio的Api.**

## 演示
<img src="public/cli-demo.gif" width="100%" height="100%">
<!-- <video id="video" controls="" preload="none" poster="封面">
      <source id="mp4" src="public/cli-demo.mp4" type="video/mp4">
</videos> -->

## 快速开始
- 推荐通过docker-compose启动<br>
  在script目录下运行 `docker-compose up -d`
- 克隆此项目到本地<br>
  分别启动`basic-generator`、`db-api-generator`、`graphql-generator`、`microservice-generator`


## 模块介绍
### microservice-generator
生成一个完整的微服务API，返回一个代码zip包。子模块参数参考下面：
- Basic：参数介绍参考`basic-generator`。
- DB：参数介绍参考`db-api-generator`。
- Graphql: 可选，不填就不会生成Graphql相应代码，参数介绍参考`graphql-generator`
- 示例
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
用来创建一个最基础的项目框架，可以生成cli在本地运行。包含参数：
- ProjectName：生成项目名称。
- Prefix：生成API前缀
- Port: 启动端口
- Modules：包含生成模块，可选参数 'JWT', 'DB', 'BASIC','ALL'
- 示例
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
基于basic-generator创建的项目，生成数据库代码，需要提供可连接的数据库参数。包含参数：
- ProjectName：项目名称，同上。
- Prefix：API前缀，同上。
- Host: 数据库服务器。
- DBname：数据库名称。
- DBPort：数据库端口。
- User: 用户名。
- Pwd: 密码。
- 示例
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
基于basic-generator创建的项目，生成Graphql代码，需要提供执行的数据库SQL。包含参数：
- ProjectName：项目名称，同上。
- QueryName：Graphql的查询名称。
- QueryDescription: Graphql的查询描述。
- SQL：查询数据需要的SQL。
- 示例
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
通过API部署docker image到Kubernetes，以及部署API到Istio：
#### 创建命名空间以及对应的Istio API规则。
- `http://localhost:9003/cli/ns/{ns}`
- 示例
```
curl --location --request POST 'http://localhost:9003/cli/ns/demo'
```
#### 部署服务到Kubernetes以及注入api到Istio。
- `http://localhost:9003/cli/svc/{ns}`
- MetaData：
    - Name: Kubernetes服务名称。
    - Version: Kubernetes服务版本。
    - Prefix： Api前缀。
- Container：
    - Image：部署镜像。
    - Port：容器端口。
    - ForceUpdate：强制更新，如果为`true`，每次部署都会重新拉取镜像。
    - RunAsRoot：是否以`root`启动容器。
    - Replicas：实例个数。
    - Environments：容器环境变量。
- 示例
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

## 作者
[Kevin Xu](http://kevin913.com.cn/about)<br />
<img src="public/my.jpg" width="200" height="200"><br/>




