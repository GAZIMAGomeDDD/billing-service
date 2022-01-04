package currencycache

import (
	"github.com/GAZIMAGomeDDD/billing-service/internal/storage/inmemory"
	"github.com/GAZIMAGomeDDD/billing-service/pkg/exchangeratesapi"
)

type Cache struct {
	cache  *inmemory.Heap
	ttl    int64
	apiKey string
}

func New(apiKey string, cache *inmemory.Heap, ttl int64) *Cache {
	return &Cache{
		cache:  cache,
		ttl:    ttl,
		apiKey: apiKey,
	}
}

func (c *Cache) GetCurrencyExchangeRate(currency string) (float64, error) {
	var rate float64

	value, ok := c.cache.Get(currency)
	if !ok {
		rate, err := exchangeratesapi.GetCurrencyExchangeRate(currency, c.apiKey)
		if err != nil {
			return 0, err
		}

		c.cache.Set(currency, rate, c.ttl)

		return rate, nil
	}

	rate = value.(float64)
	return rate, nil
}
