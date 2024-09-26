package localcontext

import (
	"l0/types/cache"
	"l0/types/config"
	_db "l0/types/db"
	"l0/types/kafka_reader"
	"l0/types/kafka_writer"
	"sync"
)

type LocalContext struct {
	Db        *_db.DB
	Cache     *cache.Cache
	Reader    *kafka_reader.KafkaReader
	Writer    *kafka_writer.KafkaWriter
	WaitGroup *sync.WaitGroup
}

func GetLocalContext(cfg *config.Config, reader *kafka_reader.KafkaReader, writer *kafka_writer.KafkaWriter) (*LocalContext, error) {
	db, err := _db.GetDB(_db.DBContext{
        Password: cfg.DbPassword,
        Host: cfg.DbHost,
        Port: cfg.DbPort,
        User: cfg.DbUser,
    })
    if err != nil  {
        return nil, err
    }
	cache := cache.GetCache(db)


	return &LocalContext{
		Db:        db,
		Cache:     cache,
		Reader:    reader,
		Writer:    writer,
		WaitGroup: new(sync.WaitGroup),
	}, nil
}
