package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	host := "haproxy-svc"
	port := "5432"
	user := "postgres"
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s sslmode=disable",
		host, port, user)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("create database testdb")
	if err != nil {
		panic(err)
	}
}
