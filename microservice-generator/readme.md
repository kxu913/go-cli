## Create  microservice-generator project

### Start in local environment
#### Prepare VS code and Golang environment
Open the folder in VS Code.

#### Init microservice-generator
```./init.cmd```

### Debug application

#### Start the application
```go run .```
#### Access 
http://localhost:1325/it/v1/noauth

### Containerize and Deploy it to k8s
#### Containerize microservice-generator, NEED docker Start
```docker build -t kevin_913/microservice-generator```

#### Deploy it on k8s, NEED k8s start and kubectl installed
- Create resource ```./deploy.cmd```
- Destroy resource ```./deploy.cmd destroy```


