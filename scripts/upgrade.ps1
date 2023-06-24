param(
    [string]$version
)

Set-Location "D:\workspace\tool\go_cli\basic-generator"
docker build -t basic-generator .
docker tag basic-generator:latest ccr.ccs.tencentyun.com/minicloud/basic-generator:$version
docker push ccr.ccs.tencentyun.com/minicloud/basic-generator:$version

Set-Location "D:\workspace\tool\go_cli\db-api-generator"
docker build -t db-api-generator .
docker tag db-api-generator:latest ccr.ccs.tencentyun.com/minicloud/db-api-generator:$version
docker push ccr.ccs.tencentyun.com/minicloud/db-api-generator:$version

Set-Location "D:\workspace\tool\go_cli\graphql-generator"
docker build -t graphql-generator .
docker tag graphql-generator:latest ccr.ccs.tencentyun.com/minicloud/graphql-generator:$version
docker push ccr.ccs.tencentyun.com/minicloud/graphql-generator:$version

Set-Location "D:\workspace\tool\go_cli\microservice-generator"
docker build -t microservice-generator .
docker tag microservice-generator:latest ccr.ccs.tencentyun.com/minicloud/microservice-generator:$version
docker push ccr.ccs.tencentyun.com/minicloud/microservice-generator:$version

Set-Location "D:\workspace\tool\go_cli\serverless-generator"
docker build -t serverless-generator .
docker tag serverless-generator:latest ccr.ccs.tencentyun.com/minicloud/serverless-generator:$version
docker push ccr.ccs.tencentyun.com/minicloud/serverless-generator:$version

Set-Location "D:\workspace\tool\go_cli\scripts"