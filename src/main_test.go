package main

import (
	"fmt"
	lcU "l0/test_utils/local_context"
	"l0/types/config"
	_db "l0/types/db"
	"testing"
)

var validMessage = `{"customer_id":"OfvqnGQgOZ","date_created":"2024-05-14T13:54:17+03:00","delivery":{"address":"LWR35dxF 62","city":"New York","email":"A9tPJinf@gmail.com","name":"hwflm8IV wLU7LQ","phone":"+970143226305","region":"Quebec","zip":"880840"},"delivery_service":"alpha","entry":"7GEv","internal_signature":"T6Awb7aG","items":[{"brand":"Apple","chrt_id":939560,"name":"BvY4yX5aJ5","nm_id":354935,"price":169,"rid":"Y1QGcPXal79xPlqj","sale":10,"size":"5","status":265,"total_price":33,"track_number":"Oi9qdg2V2JEC"}],"locale":"en","oof_shard":"8","order_uid":"8iIAICP7QPTn082Z","payment":{"amount":1358,"bank":"paypal","currency":"USD","custom_fee":40,"delivery_cost":376,"goods_total":63,"payment_dt":1727002457,"provider":"wbpay","request_id":"VXlvf2lq55Pc","transaction":"euslvFnewz0yKj7k"},"shardkey":"4","sm_id":90,"track_number":"2Y1YsCqSxXZw"}`
var invalidMessage = `{"customer_id":"OfvqnGQgOZ","date_created":"2024-05-14T13:54:17+03:00","delivery":{"address":"LWR35dxF 62","city":"New York","email":"A9tPJinf@gmail.com","name":"hwflm8IV wLU7LQ","phone":"+970143226305","region":"Quebec","zip":"880840"},"delivery_service":"alpha","entry":"7GEv","internal_signature":"T6Awb7aG","items":[{"brand":"Apple","chrt_id":939560,"name":"BvY4yX5aJ5","nm_id":354935,"price":169,"rid":"Y1QGcPXal79xPlqj","sale":10,"size":"5","status":265,"total_price":33,"track_number":"Oi9qdg2V2JEC"}],"locale":"en","oof_shard":"8","order_uid":89,"payment":{"amount":1358,"bank":"paypal","currency":"USD","custom_fee":40,"delivery_cost":376,"goods_total":63,"payment_dt":1727002457,"provider":"wbpay","request_id":"VXlvf2lq55Pc","transaction":"euslvFnewz0yKj7k"},"shardkey":"4","sm_id":90,"track_number":"2Y1YsCqSxXZw"}`
var invalidJsonMessage = `invalid Json`

func TestProcessMessage(t *testing.T) {
	cfg, err := config.GetConfig()
	if err != nil {
		t.Error("No config")
		return
	}

	db, err := _db.GetDB(_db.DBContext{
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
    db.DBCleanup()

	ctx := lcU.GetLocalContext(db, nil, nil)

	ProcessMessage([]byte(validMessage), &ctx)
	ProcessMessage([]byte(invalidMessage), &ctx)
	ProcessMessage([]byte(invalidJsonMessage), &ctx)

	allOrders, err := db.GetAll()
	if err != nil {
		t.Error("Couldn't read orders from db")
		return
	}

	if len(allOrders) != 1 {
		t.Error(fmt.Sprintf("Invalid Orders were added to db, or valid were not (expected allOrders len = 1 but got %v)", len(allOrders)), allOrders)

		return
	}
	if allOrders[0].Id != "8iIAICP7QPTn082Z" {
		t.Error("Parsed order id is incorrect")
	}
}
