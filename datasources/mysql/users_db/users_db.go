package users_db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

//7th layer
//here we initiate the database
var (
	Client *sql.DB
	err    error
)

func init() {
	datasourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		,
	)
	Client, err = sql.Open("mysql", datasourceName)
	if err != nil {
		panic(err)
	}
	if err = Client.Ping(); err != nil {
		panic(err)
	}
	log.Println("database successfully configured")
}
