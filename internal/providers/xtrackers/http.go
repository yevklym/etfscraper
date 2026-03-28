package xtrackers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
)

func (c *Client) doPostBrowser(ctx context.Context, url string, body []byte) ([]byte, error) {
	if c.httpConfig.Debug {
		c.httpConfig.Logger.Printf("xtrackers: bypassing Akamai via headless browser fetch...")
	}

	// Give the browser max 60s to launch and fetch if context has no deadline
	navCtx := ctx
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		navCtx, cancel = context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
	}

	browser := rod.New().Context(navCtx)
	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect headless browser: %w", err)
	}
	defer func() { _ = browser.Close() }()

	page, err := stealth.Page(browser)
	if err != nil {
		return nil, fmt.Errorf("failed to create stealth page: %w", err)
	}
	defer func() { _ = page.Close() }()

	// Start navigation, which presents the cookie and role popups
	navURL := c.config.DefaultHeaders["Referer"]
	if err := page.Navigate(navURL); err != nil {
		return nil, fmt.Errorf("failed to navigate to product finder: %w", err)
	}

	// Wait for the page DOM to render enough for the Entry Gate modals.
	// We cannot use WaitLoad() because DWS pages load infinite tracking pixels
	// which prevents the browser "load" event from ever firing.
	time.Sleep(3 * time.Second)

	// 1. Accept cookies — use JS eval to find and click the button by text content.
	// This approach reliably triggers React's synthetic event handlers.
	if _, err := page.Eval(`() => {
		const btns = Array.from(document.querySelectorAll('button'));
		const acceptBtn = btns.find(b =>
			b.innerText.includes('Accept all cookies') ||
			b.innerText.includes('Alle Cookies akzeptieren'));
		if (acceptBtn) acceptBtn.click();
	}`); err != nil && c.httpConfig.Debug {
		c.httpConfig.Logger.Printf("xtrackers: cookie accept eval: %v", err)
	}

	time.Sleep(1 * time.Second)

	// 2. Select Retail Investor role (Privat / Private)
	if _, err := page.Eval(`() => {
		const labels = Array.from(document.querySelectorAll('*'));
		const role = labels.find(el =>
			el.innerText &&
			(el.innerText.trim() === 'Privat' ||
			 el.innerText.trim() === 'Private' ||
			 el.innerText.trim() === 'Private Investor'));
		if (role) role.parentElement.click();
	}`); err != nil && c.httpConfig.Debug {
		c.httpConfig.Logger.Printf("xtrackers: role select eval: %v", err)
	}

	time.Sleep(1 * time.Second)

	// 3. Click the "Accept & continue" button to submit the Entry Gate form
	if btn, err := page.Element("button.d-entry-gate__continue"); err == nil {
		if err := btn.Click("left", 1); err != nil && c.httpConfig.Debug {
			c.httpConfig.Logger.Printf("xtrackers: continue button click: %v", err)
		}
		// Wait for the session cookies to settle after redirect
		time.Sleep(3 * time.Second)
	}

	// 4. Evaluate fetch inside the now fully authenticated browser context
	fetchJS := `(url, body) => {
		return fetch(url, {
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json',
				'Accept': 'application/json',
				'client-id': 'passive-frontend'
			},
			body: body
		}).then(res => res.text());
	}`

	result, err := page.Eval(fetchJS, url, string(body))
	if err != nil {
		return nil, fmt.Errorf("in-browser fetch failed: %w", err)
	}

	return []byte(result.Value.Str()), nil
}

// doPost performs a standard HTTP POST request.
// This is preserved for future use with APIs that do not require Akamai bypass.
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
