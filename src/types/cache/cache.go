package cache

import (
	"errors"
	"fmt"
	_db "l0/types/db"
	"l0/types/logger"
	_order "l0/types/order"
	"reflect"
	"sync"
)

type Cache struct {
	dataStore sync.Map
}

func GetCache(db *_db.DB) *Cache {
	cache := new(Cache)

	count := 0
	allOrders, err := db.GetAll()
	if err == nil {
		for _, entry := range allOrders {
			cache.Put(&entry)
			count++
		}
	} else {
		logger.Logger.Info().
			Msg("Couldn't get orders from database")
	}

	logger.Logger.Info().
		Msg(fmt.Sprintf("Orders recovered from db: %v", count))

	return cache
}

func (c *Cache) Get(id string) (*string, error) {
	val, ok := c.dataStore.Load(id)
	if !ok {
		logger.Logger.Warn().
			Str("order_uid", id).
			Msg("Attempted to get non existing order")
		return nil, errors.New("No such order")
	}

	logger.Logger.Info().
		Str("order_uid", id).
		Msg("Order request successful")

	if stVal, okS := val.(string); okS {
		return &stVal, nil
	}

	logger.Logger.Info().
		Str("order_uid", id).
		Msg(fmt.Sprintf("Expecte value to be string but got: '%s'", reflect.TypeOf(val)))

	return nil, errors.New("failed to convert cached value to string")

}

func (c *Cache) GetAll() (allOrders []_order.Order)  {
	c.dataStore.Range(func(key, value any) bool {
		stVal, ok := value.(string)
		stKey, ok2 := key.(string)
		if ok && ok2 {
			allOrders = append(allOrders, _order.Order{
                Id: stKey,
                JsonStr: stVal,
            })
		}
		return true
	})
	return allOrders
}

func (c *Cache) Put(order *_order.Order) {
	c.dataStore.Store(order.Id, order.JsonStr)

	logger.Logger.Info().
		Str("order_uid", order.Id).
		Msg("Added order to cache")
}
