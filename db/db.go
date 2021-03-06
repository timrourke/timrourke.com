package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func ConnectToDB(username string, password string, dbname string) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True",
		username, password, dbname)

	db, err := sqlx.Connect("mysql", connectionString)
	if err != nil {
		fmt.Println("database error", err)
		return nil, errors.Wrap(err, "could not connect to database")
	}

	return db, nil
}
