package test_utils_local_context

import (
	_cache "l0/types/cache"
	_db "l0/types/db"
	"l0/types/kafka_reader"
	local_context "l0/types/local_context"
	"sync"
)

func GetLocalContext(db *_db.DB, reader *kafka_reader.KafkaReader) local_context.LocalContext {
	cache := _cache.GetCache(db)

	return local_context.LocalContext{
		Db:     db,
		Cache:  cache,
		Reader: reader,
		WaitGroup: new(sync.WaitGroup),
	}
}
