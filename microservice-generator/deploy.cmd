@echo off
@setlocal
@REM Deploy microservice-generator} to k8s

echo "1 - Containize project"
docker build -t kevin_913/microservice-generator .
echo "Containize end..."
cd scripts
echo "2 - Create resources."
kubectl apply -f service.yaml
echo "All done"

