package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

type Order struct {
	id      string
	jsonStr string
}

func main() {
	db := GetDB()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:29092"},
		Topic:     "orders",
		Partition: 0,
		MaxBytes:  10e6,
	})

	var counter = 0

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Print(err)
			break
		}

		fmt.Printf("%s %s\n", msg.Key, msg.Value)

		if err := db.QueryRow(`insert into "order"(id, json_data) values($1, $2)`, strconv.Itoa(counter), msg.Value).Err(); err != nil {
            fmt.Printf("couldn't write to database: %s\n", err)
		}

		counter++
		test()
	}

	if err := reader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

func GetDB() *sql.DB {
	user := "user"
	password := "password"
	host := "localhost"
	port := "55432"

	conntectStr := fmt.Sprintf(
		"user='%s' password='%s' host='%s' port='%s' sslmode=disable dbname=orders",
		user,
		password,
		host,
		port,
	)

	db, err := sql.Open("postgres", conntectStr)
	if err != nil {
		log.Fatal(fmt.Sprintf("Fail to connect to db:\n%s", err))
	}
	return db
}

func test() {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	user = "user"
	password = "password"
	host = "localhost"
	port = "55432"

	conntectStr := fmt.Sprintf(
		"user='%s' password='%s' host='%s' port='%s' sslmode=disable dbname=orders",
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

		fmt.Printf("%s %s\n", id, json)
	}
}
