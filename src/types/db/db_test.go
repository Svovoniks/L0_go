package db

import (
	"l0/types/config"
	"l0/types/order"
	"testing"
)

func TestGet(t *testing.T) {
	cfg, err := config.GetConfig()
	if err != nil {
		t.Error("No config")
		return
	}

	db, err := GetDB(DBContext{
		Password: cfg.DbPassword,
		Host:     cfg.DbHost,
		Port:     cfg.DbPort,
		User:     cfg.DbUser,
	})
	if err != nil {
		t.Error(err)
		return
	}
	defer db.DBCleanup()

	if !db.Put(&order.Order{Id: "id", JsonStr: "val"}) {
		t.Error("Couldn't write to db")
	}

	val, errD := db.Get("id")

	if errD != nil {
		t.Error("Couldn't read from db", errD)
		return
	}

	if val.JsonStr != "val" {
		t.Error("Read incorrect value")
	}
}
