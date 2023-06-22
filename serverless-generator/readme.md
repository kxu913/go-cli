## Create  serverless-generator project

### Start in local environment
#### Prepare VS code and Golang environment
Open the folder in VS Code.

#### Init serverless-generator
```./init.cmd```

### Debug application

#### Start the application
```go run .```
#### Access 
http://localhost:9003/demo/noauth

### Containerize and Deploy it to k8s
#### Containerize serverless-generator, NEED docker Start
```docker build -t kevin_913/serverless-generator```

#### Deploy it on k8s, NEED k8s start and kubectl installed
- Create resource ```./deploy.cmd```
- Destroy resource ```./deploy.cmd destroy```


