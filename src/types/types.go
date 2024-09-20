package types

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sync"

	_ "github.com/lib/pq"
)

type LocalContext struct {
	Db         *DB
	Cache      *Cache
	Reader_ctx *context.Context
	Writer_ctx *context.Context
	WaitGroup  *sync.WaitGroup
}

func GetLocalContext(readerContext *context.Context, writerContext *context.Context) LocalContext {
	db := GetDB()
	cache := GetCache(db)

	return LocalContext{
		Db:         db,
		Cache:      cache,
		Reader_ctx: readerContext,
		Writer_ctx: writerContext,
		WaitGroup:  new(sync.WaitGroup),
	}
}

const OrderIdJsonKey string = "order_uid"

var RequiredOrderFields = []string{
	OrderIdJsonKey,
	"track_number",
	"entry",
	"delivery",
	"payment",
	"items",
}

type Order struct {
	Id      string
	JsonStr string
}

type Cache struct {
	lock      sync.RWMutex
	dataStore map[string]string
}

func GetCache(db *DB) *Cache {
	cache := Cache{
		dataStore: make(map[string]string),
	}

	for _, entry := range db.GetAll() {
		cache.Put(&entry)
	}

	return &cache
}

func (c *Cache) Get(id string) string {
	return c.dataStore[id]
}

func (c *Cache) Put(order *Order) {
	c.dataStore[order.Id] = order.JsonStr
}

func IsValidOrder(jsonMap map[string]any) bool {
	for _, field := range RequiredOrderFields {
		if _, ok := jsonMap[field]; !ok {
			return false
		}
	}

	return true
}

func OrderFromMessage(message []byte) (*Order, error) {
	var jsonMap map[string]any
	err := json.Unmarshal(message, &jsonMap)

	if err != nil {
		return nil, err
	}

	if !IsValidOrder(jsonMap) {
		return nil, errors.New("Not a valid order")
	}

	if _, ok := jsonMap[OrderIdJsonKey].(string); !ok {
		return nil, errors.New(fmt.Sprint("Expected order_uid to be string, but got:", reflect.TypeOf(jsonMap[OrderIdJsonKey])))
	}

	return &Order{
		Id:      jsonMap[OrderIdJsonKey].(string),
		JsonStr: string(message),
	}, nil

}

type DB struct {
	db *sql.DB
}

func GetDB() *DB {
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
	return &DB{
		db: db,
	}
}

func (d *DB) Put(order *Order) bool {
	if err := d.db.QueryRow(`insert into "order"(id, json_data) values($1, $2)`, order.Id, order.JsonStr).Err(); err != nil {
		fmt.Printf("couldn't write to database: %s\n", err)
		return false
	}
	return true
}

func (d *DB) Get(id string) (order *Order) {
	row, err := d.db.Query(`select json_data from "orders" where id=$1`, id)

	if err != nil {
		return nil
	}

	defer row.Close()

	var json_data string

	err = row.Scan(&json_data)

	if err != nil {
		return nil
	}

	return &Order{
		Id:      id,
		JsonStr: json_data,
	}
}

func (d *DB) GetAll() (orders []Order) {
	rows, err := d.db.Query(`SELECT * FROM "order"`)
	if err != nil {
		return nil
	}

	defer rows.Close()

	for rows.Next() {
		var id string
		var json string

		if err := rows.Scan(&id, &json); err != nil {
			fmt.Printf("Fail to read row:\n%s", err)
			continue
		}

		orders = append(orders, Order{
			Id:      id,
			JsonStr: json,
		})
	}

	return orders
}
