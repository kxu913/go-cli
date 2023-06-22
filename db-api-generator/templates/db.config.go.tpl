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
	dsn  = "host={{.Host}} port={{.Port}} user={{.UserName}} password={{.Password}} dbname={{.DBName}} sslmode=disable"
	
)

const (
	LIST_TABLES = "select tablename from pg_tables where schemaname='public'"
)

func Query(sql string) []map[string]any {
	g := OpenDBByGorm()
	defer Close(g)
	results := []map[string]any{}
	g.Raw(sql).Scan(&results)
	return results

}
