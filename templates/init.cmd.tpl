@echo off
@setlocal
@REM Init module using gowork
go work init
go work use -r .
go work use -r src/

@REM Init root mod
go mod tidy

@REM Init middleware mod
cd src/middleware
go mod tidy

@REM Init route mod
cd ../route
go mod tidy