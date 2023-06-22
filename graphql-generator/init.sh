#!/bin/bash
go work init
go work use -r .
go work use -r src/

# Init root mod
go mod tidy

# Init middleware mod
cd src/middleware
go mod tidy

# Init route mod
cd ../route
go mod tidy