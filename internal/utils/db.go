package utils

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewDB() *sqlx.DB {
	db, err := sqlx.Connect("mysql", buildDSN())
	if err != nil {
		panic(err.Error())
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(400 * time.Millisecond)

	return db
}

func buildDSN() string {
	var config = loadEnv()

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
}
