package order

import "testing"

func TestIsValidOrder(t *testing.T) {
	invalidOrder := map[string]any{"id": "34"}
	if IsValidOrder(invalidOrder) {
		t.Error("Faild to reject an invalid order")
	}

	validOrder := map[string]any{
		"order_uid":    "id",
		"track_number": "",
		"entry":        "",
		"delivery":     "",
		"payment":      "",
		"items":        "",
	}
	if !IsValidOrder(validOrder) {
		t.Error("Failed to validate a valid order")
	}
}

func TestRandomValidOrder(t *testing.T) {
	if !IsValidOrder(GenerateRandomOrder()) {
		t.Error("Generated order is invalid")
	}
}

var validMessage = `{"customer_id":"OfvqnGQgOZ","date_created":"2024-05-14T13:54:17+03:00","delivery":{"address":"LWR35dxF 62","city":"New York","email":"A9tPJinf@gmail.com","name":"hwflm8IV wLU7LQ","phone":"+970143226305","region":"Quebec","zip":"880840"},"delivery_service":"alpha","entry":"7GEv","internal_signature":"T6Awb7aG","items":[{"brand":"Apple","chrt_id":939560,"name":"BvY4yX5aJ5","nm_id":354935,"price":169,"rid":"Y1QGcPXal79xPlqj","sale":10,"size":"5","status":265,"total_price":33,"track_number":"Oi9qdg2V2JEC"}],"locale":"en","oof_shard":"8","order_uid":"8iIAICP7QPTn082Z","payment":{"amount":1358,"bank":"paypal","currency":"USD","custom_fee":40,"delivery_cost":376,"goods_total":63,"payment_dt":1727002457,"provider":"wbpay","request_id":"VXlvf2lq55Pc","transaction":"euslvFnewz0yKj7k"},"shardkey":"4","sm_id":90,"track_number":"2Y1YsCqSxXZw"}`
var invalidMessage = `{"customer_id":"OfvqnGQgOZ","date_created":"2024-05-14T13:54:17+03:00","delivery":{"address":"LWR35dxF 62","city":"New York","email":"A9tPJinf@gmail.com","name":"hwflm8IV wLU7LQ","phone":"+970143226305","region":"Quebec","zip":"880840"},"delivery_service":"alpha","entry":"7GEv","internal_signature":"T6Awb7aG","items":[{"brand":"Apple","chrt_id":939560,"name":"BvY4yX5aJ5","nm_id":354935,"price":169,"rid":"Y1QGcPXal79xPlqj","sale":10,"size":"5","status":265,"total_price":33,"track_number":"Oi9qdg2V2JEC"}],"locale":"en","oof_shard":"8","order_uid":89,"payment":{"amount":1358,"bank":"paypal","currency":"USD","custom_fee":40,"delivery_cost":376,"goods_total":63,"payment_dt":1727002457,"provider":"wbpay","request_id":"VXlvf2lq55Pc","transaction":"euslvFnewz0yKj7k"},"shardkey":"4","sm_id":90,"track_number":"2Y1YsCqSxXZw"}`
var invalidJsonMessage = `invalid Json`

func TestOrderFromMessage(t *testing.T) {
	order, err := OrderFromMessage([]byte(validMessage))

	if err != nil {
		t.Error("Failed to parse order from message")
	}

	if order.Id != "8iIAICP7QPTn082Z" {
		t.Error("Parsed order id is incorrect")
	}

	order, err = OrderFromMessage([]byte(invalidMessage))

	if err == nil {
		t.Error("Managed to parse an invalid message")
	}

	order, err = OrderFromMessage([]byte(invalidMessage))

	if err == nil {
		t.Error("Parsed invalid json")
	}
}
