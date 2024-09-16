package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

type Order struct {
	id      string
	jsonStr string
}

func main() {
    // db := GetDB()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:29092"},
		Topic:     "orders",
		Partition: 0,
		MaxBytes:  10e6,
	})

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
            fmt.Print(err)
			break
		}

        fmt.Printf("%s %s", msg.Key, msg.Value)
	}

    defer reader.Close()
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

		print(id)
	}

	conn, err := kafka.Dial("tcp", "localhost:29092")

	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		panic(err.Error())
	}

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}
	for k := range m {
		fmt.Println(k)
	}
}
