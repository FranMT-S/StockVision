package CustomClient

import (
	"net/http"
)

func (c *CustomClient) GetRaw(endpoint string, queryParams map[string]string) ([]byte, error) {
	// create a new client with the same base url and headers, and remove content type to return raw response
	client := NewCustomClient(c.BaseURL)
	client.Headers = c.Headers
	client.Headers["Content-Type"] = ""
	return client.doRequestRaw(http.MethodGet, endpoint, queryParams, nil)
}

func (c *CustomClient) Get(endpoint string, queryParams map[string]string, result interface{}) error {
	return c.doRequest(http.MethodGet, endpoint, queryParams, nil, result)
}

func (c *CustomClient) Post(endpoint string, queryParams map[string]string, body interface{}, result interface{}) error {
	return c.doRequest(http.MethodPost, endpoint, queryParams, body, result)
}

func (c *CustomClient) Put(endpoint string, queryParams map[string]string, body interface{}, result interface{}) error {
	return c.doRequest(http.MethodPut, endpoint, queryParams, body, result)
}
func (c *CustomClient) Patch(endpoint string, queryParams map[string]string, body interface{}, result interface{}) error {
	return c.doRequest(http.MethodPatch, endpoint, queryParams, body, result)
}

func (c *CustomClient) Delete(endpoint string, queryParams map[string]string, result interface{}) error {
	return c.doRequest(http.MethodDelete, endpoint, queryParams, nil, result)
}
