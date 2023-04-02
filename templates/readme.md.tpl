## Create  {{.ProjectName}} project

### Start in local environment
#### Prepare VS code and Golang environment
Open the folder in VS Code.

#### Init {{.ProjectName}}
```./init.cmd```

### Debug application

#### Start the application
```go run .```
#### Access 
http://localhost:{{.Port}}/{{.Prefix}}/noauth

### Containerize and Deploy it to k8s
#### Containerize {{.ProjectName}}, NEED docker Start
```docker build -t kevin_913/{{.ProjectName}}```

#### Deploy it on k8s, NEED k8s start and kubectl installed
- Create resource ```./deploy.cmd```
- Destroy resource ```./deploy.cmd destroy```


