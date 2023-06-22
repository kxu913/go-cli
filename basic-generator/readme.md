# go-cli
Base on echo to create a go web application, support basic authorization.

If you just want use it, download cli.exe, and use cli.exe --help to see usage.

## Usage
- project: define project name that you will create.
- prefix: define API root path, it will be used in istio config, also see >>service.yaml.tpl
- port: define port.
- local: used to debug in local.

## Description
- go.mod.tpl: define go module.
- init.*.tpl: Init go project using gowork to manage modules.
- Dockerfile.tpl: create default Dockerfile.
- deploy.*.tpl: used to deploy the application into k8s. 
  -- *it also include istion configuration that used to manage traffic*

## Contribute
- Add more tpl files.
- Update main.go.


