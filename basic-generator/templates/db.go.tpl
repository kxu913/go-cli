package db

import (

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
func TableList() []string {

	g := OpenDBByGorm()
	defer Close(g)

	rows, err := g.Raw(LIST_TABLES).Rows()
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	tables := []string{}
	for rows.Next() {
		var table string
		rows.Scan(&table)
		tables = append(tables, table)
	}
	return tables

}
