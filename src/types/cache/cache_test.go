package cache

import (
	"l0/types/config"
	_db "l0/types/db"
	_order "l0/types/order"
	"l0/types/random"
	"testing"
)

func TestPutGet(t *testing.T) {
	cache := new(Cache)

	cache.Put(&_order.Order{
		Id:      "id",
		JsonStr: "json",
	})

	if str, err := cache.Get("id"); err != nil {
		t.Error("Failed to get value from cache")
		if *str != "json" {
			t.Error("Value read didn't match value written")
		}
	}
}

func TestGetCache(t *testing.T) {
	cfg, err := config.GetConfig()
	if err != nil {
		t.Error("No config", err)
		return
	}
	db, err := _db.GetDB(_db.DBContext{
		Password: cfg.DbPassword,
		Host:     cfg.DbHost,
		Port:     cfg.DbPort,
		User:     cfg.DbUser,
	})
	if err != nil {
		t.Error("Couldn't connect to database", err)
		return
	}
	defer db.Db.Close()
	defer db.DBCleanup()

	ord_ls := [5]_order.Order{}

	for i := range 5 {
		ord_ls[i] = _order.Order{
			Id:      random.RandomString(10),
			JsonStr: random.RandomString(10),
		}
		if !db.Put(&ord_ls[i]) {
			t.Error("failed to put order into database")
		}
	}

	cache := GetCache(db)

	for _, ord := range ord_ls {
		if _, err := cache.Get(ord.Id); err != nil {
			t.Error("cache does not contain all values from database")
		}
	}
}
