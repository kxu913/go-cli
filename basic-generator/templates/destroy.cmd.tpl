@echo off
@setlocal
@REM Destroy {{.ProjectName}}} in k8s

cd scripts

kubectl delete -f service.yaml


