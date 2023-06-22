## Create  graphql-generator project

### Start in local environment
#### Prepare VS code and Golang environment
Open the folder in VS Code.

#### Init graphql-generator
```./init.cmd```

### Debug application

#### Start the application
```go run .```
#### Access 
http://localhost:9004/demo/noauth

### Containerize and Deploy it to k8s
#### Containerize graphql-generator, NEED docker Start
```docker build -t kevin_913/graphql-generator```

#### Deploy it on k8s, NEED k8s start and kubectl installed
- Create resource ```./deploy.cmd```
- Destroy resource ```./deploy.cmd destroy```


