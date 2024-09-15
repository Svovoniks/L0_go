package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	user = "user"
	password = "password"
	host = "localhost"
	port = "55432"

	conntectStr := fmt.Sprintf(
		"user='%s' password='%s' host='%s' port='%s' sslmode=disable",
		user,
		password,
		host,
		port,
	)

	db, err := sql.Open("postgres", conntectStr)
	if err != nil {
		log.Fatal(fmt.Sprintf("Fail to connect to db:\n%s", err))
	}

	defer db.Close()

    rows, err := db.Query(`SELECT * FROM "order"`)
	if err != nil {
		log.Fatal(fmt.Sprintf("Fail to read data:\n%s", err))
	}

	defer rows.Close()

	for rows.Next() {
		var id string
		var json string

		if err := rows.Scan(&id, &json); err != nil {
			fmt.Printf("Fail to read row:\n%s", err)
		}

        print(id)
	}

}
