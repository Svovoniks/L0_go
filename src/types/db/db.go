package db

import (
	"database/sql"
	"errors"
	"fmt"
	"l0/types/logger"
	_order "l0/types/order"

	_ "github.com/lib/pq"
)

type DB struct {
	Db *sql.DB
    Table string
}

type DBContext struct {
	Password string
	Host     string
	Port     string
	User     string
}

func GetDB(ctx DBContext) (*DB, error) {
	conntectStr := fmt.Sprintf(
		"user='%s' password='%s' host='%s' port='%s' sslmode=disable dbname=orders",
		ctx.User,
		ctx.Password,
		ctx.Host,
		ctx.Port,
	)


	db, err := sql.Open("postgres", conntectStr)
	if err != nil {
		logger.Logger.Warn().Msg(fmt.Sprintf("Fail to connect to db:\n%s", err))
		return nil, err
	}
	db.SetMaxOpenConns(20)

	db.SetMaxIdleConns(10)

	return &DB{
		Db: db,
        Table: "order",
	}, nil
}

func (d *DB) Put(order *_order.Order) bool {
	// not QueryRow because it spawans active connections
	// it overwhelms db and i couldn't figure out how to close them

	row, err := d.Db.Query(`insert into "order"(id, json_data) values($1, $2) on conflict (id) do update set json_data=$2`, order.Id, order.JsonStr)

	if err != nil {
		logger.Logger.Warn().
			Err(err).
			Msg("Couldn't write to database")
	} else {
		logger.Logger.Info().
			Str("order_uid", order.Id).
			Msg("Added order to database")
	}

	if row != nil {
		row.Close()
	}

	return err == nil
}

func (d *DB) Get(id string) (*_order.Order, error) {
	query := `select json_data from "order" where id=$1`
	row, err := d.Db.Query(query, id)

	if err != nil {
		logger.Logger.Warn().
			Str("order_uid", id).
			Str("query", query).
			Err(err).
			Msg("SQL query failed")
		return nil, err
	}

	defer row.Close()

	var json_data string

    if !row.Next(){
        return nil, errors.New("Not in database")
    }

	err = row.Scan(&json_data)

	if err != nil {
		logger.Logger.Warn().
			Str("order_uid", id).
			Err(err).
			Msg("Couldn't scan the row")
		return nil, err
	}

	logger.Logger.Info().
		Str("order_uid", id).
		Msg("Retrieved order from database")

	return &_order.Order{
		Id:      id,
		JsonStr: json_data,
	}, nil
}

func (d *DB) GetAll() ([]_order.Order, error) {
	query := `SELECT * FROM "order"`
	rows, err := d.Db.Query(query)
	if err != nil {
		logger.Logger.Warn().
			Str("query", query).
			Err(err).
			Msg("SQL query failed")
		return nil, err
	}

	defer rows.Close()

	var orders []_order.Order

	for rows.Next() {
		var id string
		var json string

		if err := rows.Scan(&id, &json); err != nil {
			logger.Logger.Warn().
				Err(err).
				Msg("Couldn't scan the row")
			continue
		}

		orders = append(orders, _order.Order{
			Id:      id,
			JsonStr: json,
		})
	}

	logger.Logger.Info().
		Int("order_count", len(orders)).
		Msg("Fetched all orders")

	return orders, nil
}

func (db *DB) DBCleanup() {
	db.Db.Exec(`delete from "order"`)
}
