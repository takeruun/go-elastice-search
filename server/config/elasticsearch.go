package config

import (
	"net/http"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticSearch struct {
	*elasticsearch.Client
}

func NewElasticSearch() *ElasticSearch {
	retryBackoff := backoff.NewExponentialBackOff()
	config := elasticsearch.Config{
		RetryOnStatus: []int{
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
			http.StatusTooManyRequests,
		},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		MaxRetries: 5,
	}

	config.Addresses = []string{"http://elasticsearch:9200"}

	es, err := elasticsearch.NewClient(config)
	if err != nil {
		panic(err)
	}

	return &ElasticSearch{
		es,
	}
}
