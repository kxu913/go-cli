package db

import "os"

func GetEnvWithDefault(key string, defValue string) string {
	val, err := os.LookupEnv(key)
	if !err {
		return defValue
	}
	return val

}

var (
	host = GetEnvWithDefault("db_host", "localhost")
	dsn  = "host=" + host + " port=5432 user=postgres password=postgres dbname=user sslmode=disable"
	
)

const (
	LIST_TABLES = "select tablename from pg_tables where schemaname='public'"
)
