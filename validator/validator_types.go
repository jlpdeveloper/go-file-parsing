package validator

import (
	"context"
	"go-file-parsing/cache"
	"go-file-parsing/config"
)

type RowError struct {
	Row   int64
	Id    string
	Error error
}

type CacheData struct {
	Id   string
	Data map[string]string
}
type ColValidator func(*RowValidatorContext, []string) (map[string]string, error)

type RowValidatorContext struct {
	Config *config.ParserConfig
	GetMap func() map[string]string
}

type RowValidator interface {
	Validate(row string) (string, error)
}

func New(conf *config.ParserConfig, cacheChan chan CacheData, colValidators []ColValidator) CsvRowValidator {
	return CsvRowValidator{
		config:        conf,
		cacheChan:     cacheChan,
		colValidators: colValidators,
	}
}

func NewCacheChannel(cache cache.DistributedCache) chan CacheData {
	ctx := context.Background()
	cacheChan := make(chan CacheData, 1000)
	cachePoolSize := 100
	cachePool := make(chan func(data CacheData), cachePoolSize)
	for i := 0; i < cachePoolSize; i++ {
		cachePool <- func(cacheItem CacheData) {
			for key, value := range cacheItem.Data {
				_ = cache.SetField(ctx, cacheItem.Id, key, value)
			}
			// Return the map to the pool after use
			PutMap(cacheItem.Data)
		}
	}

	//Spin off a new goroutine that will write to the cache as it processes from the channel
	go func() {
		for cacheItem := range cacheChan {
			worker := <-cachePool
			go func(ci CacheData) {
				worker(ci)
				cachePool <- worker
			}(cacheItem)
		}
	}()
	return cacheChan
}
