# db_cli
***Just for Postgresql, if you want to use for another database, you can update type mapping in `db_generater.go`***
Used to create some useful go files according by db tables, include below:
- model files, according by table columns.
- db files, basic CRUD functions and a constants file that describe select fileds when you want use SQL to query data.
- router files, a route group that include the resources CRUD operations.

If you just want use it, download db_cli.exe, and use db_cli.exe --help to see usage.

## Usage
- project: define project name that you want to create files, don't worry the value, you can change it after you copy this files to your workspace.
- prefix: define prefix name that route group will be behind, don't worry the value, you can change it after you copy this files to your workspace
- host: db host, default is `localhost`.
- db: db you want to generate files, `required`
- port: db port, default is `5432`.
- user: db user, default is `postgres`.
- pwd: db password, default is `postgres`.
- table: table that you want to creat this files, if you want to speify it, the cli will create retreive all tables in DB and generate files, default is `all`.

## Contribute
- Add more tpl files.
- Update main.go.


