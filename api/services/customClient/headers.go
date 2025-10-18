package CustomClient

import "fmt"

// SetHeader permite agregar o actualizar un header
func (c *CustomClient) SetHeader(key, value string) {
	c.Headers[key] = value
}

// RemoveHeader elimina un header
func (c *CustomClient) RemoveHeader(key string) {
	delete(c.Headers, key)
}

// SetAuthToken establece un token de autorización en los headers
func (c *CustomClient) SetAuthToken(token string) {
	c.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
}
