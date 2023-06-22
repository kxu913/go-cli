@echo off
@setlocal
@REM Destroy graphql-generator} in k8s

cd scripts

kubectl delete -f service.yaml


