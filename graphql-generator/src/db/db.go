package db

import (
	"fmt"
	"graphql-generator/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDBByGorm() *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Minute)
	return db
}

func Close(db *gorm.DB) error {
	sqldb, err := db.DB()
	if err != nil {
		return gorm.ErrUnsupportedDriver
	}
	sqldb.Close()
	return nil
}
func Query(sql string) []map[string]any {
	g := OpenDBByGorm()
	defer Close(g)
	results := []map[string]any{}
	g.Raw(sql).Scan(&results)
	fmt.Println(results)
	return results

}

func TableInfo(tableName string) []model.Field {

	g := OpenDBByGorm()
	defer Close(g)

	rows, err := g.Debug().Model(&model.Field{}).Raw(fmt.Sprintf(`
	SELECT a.attname AS "Name", t.typname AS "DBType"
		FROM pg_class c, pg_attribute a 
		LEFT JOIN pg_description b ON a.attrelid = b.objoid AND a.attnum = b.objsubid, pg_type t 
		WHERE c.relname = '%s' AND a.attnum > 0 AND a.attrelid = c.oid AND a.atttypid = t.oid
	`, tableName)).Rows()
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	infos := []model.Field{}
	var info model.Field
	for rows.Next() {
		g.ScanRows(rows, &info)
		infos = append(infos, info)
	}
	return infos

}
