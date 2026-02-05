package amundi

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

func (c *Client) doPost(ctx context.Context, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range c.config.DefaultHeaders {
		req.Header.Set(k, v)
	}

	resp, err := c.httpConfig.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}
