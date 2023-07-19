package utils

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func NewClient() *retryablehttp.Client {
	Client := retryablehttp.NewClient()
	Client.Logger = nil
	Client.RetryMax = 2
	Client.Backoff = retryablehttp.LinearJitterBackoff

	Client.HTTPClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10,
		},
		Timeout: 10 * time.Second,
	}

	return Client
}