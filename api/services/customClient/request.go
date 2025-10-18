package CustomClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// doRequest ejecuta la petición HTTP
func (c *CustomClient) doRequest(method, endpoint string, queryParams map[string]string, body interface{}, result interface{}) error {
	bodyBytes, err := c.doRequestRaw(method, endpoint, queryParams, body)
	if err != nil {
		return err
	}

	// parse response
	if result != nil {
		if err := json.Unmarshal(bodyBytes, result); err != nil {
			return fmt.Errorf("error parsing JSON: %w", err)
		}
	}

	return nil
}

func (c *CustomClient) doRequestRaw(method, endpoint string, queryParams map[string]string, body interface{}) ([]byte, error) {
	req, err := c.BuildRequest(method, endpoint, queryParams, body)
	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	// 2. execute request
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	return c.validateAndProcessBody(resp)
}

func (c *CustomClient) BuildRequest(method, endpoint string, queryParams map[string]string, body interface{}) (*http.Request, error) {
	// 1. build url
	fullURL, err := c.buildURL(endpoint, queryParams)
	if err != nil {
		return nil, fmt.Errorf("error building URL: %w", err)
	}

	// 2. prepare body
	bodyReader, err := c.prepareBody(body)
	if err != nil {
		return nil, fmt.Errorf("error preparing body: %w", err)
	}

	// 3. create request
	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// 4. set headers
	c.setHeaders(req)

	return req, nil
}

// buildURL construye la URL completa con query params
func (c *CustomClient) buildURL(endpoint string, queryParams map[string]string) (string, error) {
	fullURL := c.BaseURL + endpoint

	if len(queryParams) == 0 {
		return fullURL, nil
	}

	params := url.Values{}
	for key, value := range queryParams {
		params.Add(key, value)
	}

	return fullURL + "?" + params.Encode(), nil
}

// prepareBody serializes the body to JSON if it exists
func (c *CustomClient) prepareBody(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonBody), nil
}

// setHeaders sets the headers configured in the request
func (c *CustomClient) setHeaders(req *http.Request) {
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}
}

// validateAndProcessBody validates the HTTP response and reads the body
func (c *CustomClient) validateAndProcessBody(resp *http.Response) ([]byte, error) {
	// check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// read body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	return responseBody, nil
}
