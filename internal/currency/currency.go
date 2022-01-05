package currency

type Currency struct {
	cache                  InMemoryCache
	ttl                    int64
	apiKey                 string
	GetCurrencyRateOfRuble func(string, string) (float64, error)
}

type InMemoryCache interface {
	Set(key string, value interface{}, ttl int64)
	Get(key string) (interface{}, bool)
}

func New(apiKey string, cache InMemoryCache, ttl int64) *Currency {
	return &Currency{
		cache:  cache,
		ttl:    ttl,
		apiKey: apiKey,
	}
}

func (c *Currency) GetCurrencyRate(currency string) (float64, error) {
	var rate float64

	value, ok := c.cache.Get(currency)
	if !ok {
		rate, err := c.GetCurrencyRateOfRuble(currency, c.apiKey)
		if err != nil {
			return 0, err
		}

		c.cache.Set(currency, rate, c.ttl)

		return rate, nil
	}

	rate = value.(float64)
	return rate, nil
}
