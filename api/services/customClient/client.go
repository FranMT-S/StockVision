package CustomClient

import (
	"net/http"
	"time"
)

type CustomClient struct {
	Client  *http.Client
	BaseURL string
	Headers map[string]string
}

func NewCustomClient(baseURL string) CustomClient {
	return CustomClient{
		Client: &http.Client{
			Timeout: time.Second * 60,
		},
		BaseURL: baseURL,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}
