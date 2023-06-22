@echo off
@setlocal
@REM Destroy serverless-generator} in k8s

cd scripts

kubectl delete -f service.yaml


