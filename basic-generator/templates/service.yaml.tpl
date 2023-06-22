apiVersion: v1
kind: Service
metadata:
  name: {{.ProjectName}}
  labels:
    app: {{.ProjectName}}
    service: {{.ProjectName}}
spec:
  ports:
  - port: {{.Port}}
    name: http
  selector:
    app: {{.ProjectName}}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{.ProjectName}}-svc
  labels:
    account: {{.ProjectName}}-svc
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.ProjectName}}-v1
  labels:
    app: {{.ProjectName}}
    version: v1
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{.ProjectName}}
      version: v1
  template:
    metadata:
      labels:
        app: {{.ProjectName}}
        version: v1
    spec:
      serviceAccountName: {{.ProjectName}}-svc
      containers:
      - name: {{.ProjectName}}
        image: kevin_913/{{.ProjectName}}
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: {{.Port}}
        securityContext:
          runAsUser: 1000
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{.ProjectName}}
spec:
  hosts:
  - "*"
  gateways:
  - kevin-gateway
  http:
  - match:
    - uri:
        prefix: /{{.Prefix}}

    route:
    - destination:
        host: {{.ProjectName}}
        port:
          number: {{.Port}}