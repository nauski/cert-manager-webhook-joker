package joker

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	apiURL = "https://svc.joker.com/nic/replace"
)

type client struct {
	httpClient *http.Client
	username   string
	password   string
}

func New(config Config) Client {
	return &client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		username: config.Username,
		password: config.Password,
	}
}

func (c *client) CreateTXTRecord(ctx context.Context, zone, label, value string) error {
	return c.updateRecord(ctx, zone, label, "TXT", value, false)
}

func (c *client) DeleteTXTRecord(ctx context.Context, zone, label string) error {
	return c.updateRecord(ctx, zone, label, "TXT", "", true)
}

func (c *client) updateRecord(ctx context.Context, zone, label, recordType, value string, clear bool) error {
	data := url.Values{}
	data.Set("username", c.username)
	data.Set("password", c.password)
	data.Set("zone", zone)
	data.Set("label", label)
	data.Set("type", recordType)
	data.Set("verbose", "true")

	if clear {
		data.Set("clear", "true")
	} else {
		data.Set("value", value)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return c.handleAPIError(resp, body)
	}

	responseText := strings.TrimSpace(string(body))
	if responseText != "OK" {
		return fmt.Errorf("unexpected API response: %s", responseText)
	}

	return nil
}

func (c *client) handleAPIError(resp *http.Response, body []byte) error {
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("authentication failed: invalid credentials")
	case http.StatusBadRequest:
		return fmt.Errorf("bad request: %s", string(body))
	case http.StatusTooManyRequests:
		return fmt.Errorf("rate limited: %s", string(body))
	default:
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}
}
