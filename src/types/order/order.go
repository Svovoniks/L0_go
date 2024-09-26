package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"l0/types/logger"
	"l0/types/random"
	"math/rand/v2"
	"reflect"
	"time"
)

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

func RandomValidOrder() string {
	orderMap := GenerateRandomOrder()

	enc, err := json.Marshal(orderMap)

	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("Shoud never happen")
	}

	return string(enc)

}

func GenerateRandomOrder() map[string]interface{} {
	order := map[string]interface{}{
		"order_uid":    random.RandomString(16),
		"track_number": random.RandomString(12),
		"entry":        random.RandomString(4),
		"delivery": map[string]interface{}{
			"name":    random.RandomString(8) + " " + random.RandomString(6),
			"phone":   random.RandomPhone(),
			"zip":     random.RandomZip(),
			"city":    random.RandomCity(),
			"address": fmt.Sprintf("%s %d", random.RandomString(8), rand.IntN(100)),
			"region":  random.RandomRegion(),
			"email":   random.RandomEmail(),
		},
		"payment": map[string]interface{}{
			"transaction":   random.RandomString(16),
			"request_id":    random.RandomString(12),
			"currency":      "USD",
			"provider":      random.RandomProvider(),
			"amount":        rand.IntN(5000),
			"payment_dt":    time.Now().Unix(),
			"bank":          random.RandomProvider(),
			"delivery_cost": rand.IntN(500),
			"goods_total":   rand.IntN(1000),
			"custom_fee":    rand.IntN(100),
		},
		"items": []map[string]interface{}{
			{
				"chrt_id":      rand.IntN(10000000),
				"track_number": random.RandomString(12),
				"price":        rand.IntN(1000),
				"rid":          random.RandomString(16),
				"name":         random.RandomString(10),
				"sale":         rand.IntN(50),
				"size":         fmt.Sprintf("%d", rand.IntN(10)),
				"total_price":  rand.IntN(1000),
				"nm_id":        rand.IntN(1000000),
				"brand":        random.RandomBrand(),
				"status":       rand.IntN(500),
			},
		},
		"locale":            "en",
		"internal_signature": random.RandomString(8),
		"customer_id":        random.RandomString(10),
		"delivery_service":   random.RandomProvider(),
		"shardkey":           fmt.Sprintf("%d", rand.IntN(10)),
		"sm_id":              rand.IntN(100),
		"date_created":       random.RandomDate(),
		"oof_shard":          fmt.Sprintf("%d", rand.IntN(10)),
	}

	return order
}
