@echo off
@setlocal
@REM Deploy {{.ProjectName}}} to k8s

echo "1 - Containize project"
docker build -t kevin_913/{{.ProjectName}} .
echo "Containize end..."
cd scripts
echo "2 - Create resources."
kubectl apply -f service.yaml
echo "All done"

