package xtrackers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
)

// bypassEntryGate navigates to the product finder page and handles the DWS
// Entry Gate popups (cookie consent, investor role selection, accept & continue).
// After this function returns, the browser session has the necessary cookies.
func (c *Client) bypassEntryGate(page *rod.Page) {
	navURL := c.config.DefaultHeaders["Referer"]
	if err := page.Navigate(navURL); err != nil {
		if c.httpConfig.Debug {
			c.httpConfig.Logger.Printf("xtrackers: entry gate navigate: %v", err)
		}
		return
	}

	// Wait for the page DOM to render enough for the Entry Gate modals.
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
}

// launchBrowser creates a stealth browser + page and bypasses the Entry Gate.
// Caller must close the returned browser and page.
func (c *Client) launchBrowser(ctx context.Context) (*rod.Browser, *rod.Page, error) {
	if c.httpConfig.Debug {
		c.httpConfig.Logger.Printf("xtrackers: bypassing Akamai via headless browser fetch...")
	}

	navCtx := ctx
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		navCtx, cancel = context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
	}

	browser := rod.New().Context(navCtx)
	if err := browser.Connect(); err != nil {
		return nil, nil, fmt.Errorf("failed to connect headless browser: %w", err)
	}

	page, err := stealth.Page(browser)
	if err != nil {
		_ = browser.Close()
		return nil, nil, fmt.Errorf("failed to create stealth page: %w", err)
	}

	c.bypassEntryGate(page)

	return browser, page, nil
}

// doPostBrowser performs a POST request via in-browser fetch after Entry Gate bypass.
func (c *Client) doPostBrowser(ctx context.Context, url string, body []byte) ([]byte, error) {
	browser, page, err := c.launchBrowser(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = browser.Close() }()
	defer func() { _ = page.Close() }()

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
		return nil, fmt.Errorf("in-browser POST fetch failed: %w", err)
	}

	return []byte(result.Value.Str()), nil
}

// doGetBrowser performs a GET request via in-browser fetch after Entry Gate bypass.
func (c *Client) doGetBrowser(ctx context.Context, url string) ([]byte, error) {
	browser, page, err := c.launchBrowser(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = browser.Close() }()
	defer func() { _ = page.Close() }()

	fetchJS := `(url) => {
		return fetch(url, {
			method: 'GET',
			credentials: 'include',
			headers: {
				'Accept': 'application/json',
				'client-id': 'passive-frontend'
			}
		}).then(res => res.text());
	}`

	result, err := page.Eval(fetchJS, url)
	if err != nil {
		return nil, fmt.Errorf("in-browser GET fetch failed: %w", err)
	}

	return []byte(result.Value.Str()), nil
}

// doPost performs a standard HTTP POST request.
func (c *Client) doPost(ctx context.Context, url string, body []byte) (*http.Response, error) {
	return c.doRequest(ctx, http.MethodPost, url, bytes.NewReader(body))
}

// doGet performs a standard HTTP GET request.
func (c *Client) doGet(ctx context.Context, url string) (*http.Response, error) {
	return c.doRequest(ctx, http.MethodGet, url, nil)
}

// doRequest builds and executes an HTTP request with the default headers.
func (c *Client) doRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s request: %w", method, err)
	}

	for k, v := range c.config.DefaultHeaders {
		req.Header.Set(k, v)
	}

	resp, err := c.httpConfig.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s request failed: %w", method, err)
	}

	return resp, nil
}

// readAll reads all bytes from a reader.
func readAll(r io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
