@echo off
@setlocal
@REM Destroy microservice-generator} in k8s

cd scripts

kubectl delete -f service.yaml


