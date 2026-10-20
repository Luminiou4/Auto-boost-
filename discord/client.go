package discord

import (
	"autoboost/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const defaultBaseURL = "https://discord.com/api/v10"

type Client struct {
	token      string
	baseURL    string
	userAgent  string
	httpClient *http.Client
	maxRetries int
}

type ClientOption func(*Client)

func WithProxy(proxyURL string) ClientOption {
	return func(c *Client) {
		if proxyURL == "" {
			return
		}
		parsed, err := url.Parse(proxyURL)
		if err != nil {
			return
		}
		transport := &http.Transport{Proxy: http.ProxyURL(parsed)}
		c.httpClient = &http.Client{Transport: transport}
	}
}

func WithBaseURL(base string) ClientOption {
	return func(c *Client) { c.baseURL = base }
}

func WithUserAgent(ua string) ClientOption {
	return func(c *Client) { c.userAgent = ua }
}

func WithMaxRetries(n int) ClientOption {
	return func(c *Client) { c.maxRetries = n }
}

func New(token string, opts ...ClientOption) *Client {
	c := &Client{
		token:      token,
		baseURL:    defaultBaseURL,
		userAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		httpClient: &http.Client{},
		maxRetries: 3,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) buildRequest(method, endpoint string, body []byte) (*http.Request, error) {
	fullURL := c.baseURL + endpoint
	var req *http.Request
	var err error
	if body != nil {
		req, err = http.NewRequest(method, fullURL, bytes.NewReader(body))
	} else {
		req, err = http.NewRequest(method, fullURL, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.token)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Super-Properties", "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiQ2hyb21lIn0=")
	return req, nil
}

func (c *Client) Get(endpoint string) ([]byte, int, error) {
	var statusCode int
	resp, err := utils.DoWithRetry(func() (*http.Response, error) {
		req, err := c.buildRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, err
		}
		return c.httpClient.Do(req)
	}, c.maxRetries)
	if err != nil {
		return nil, 0, err
	}
	statusCode = resp.StatusCode
	body, err := utils.ReadBody(resp)
	if err != nil {
		return nil, statusCode, err
	}
	if statusCode >= 400 {
		return nil, statusCode, parseAPIError(body, statusCode)
	}
	return body, statusCode, nil
}

func (c *Client) Post(endpoint string, payload interface{}) ([]byte, int, error) {
	return c.doWithBody(http.MethodPost, endpoint, payload)
}

func (c *Client) Put(endpoint string, payload interface{}) ([]byte, int, error) {
	return c.doWithBody(http.MethodPut, endpoint, payload)
}

func (c *Client) Delete(endpoint string, payload interface{}) ([]byte, int, error) {
	if payload == nil {
		return c.doNoBody(http.MethodDelete, endpoint)
	}
	return c.doWithBody(http.MethodDelete, endpoint, payload)
}

func (c *Client) doWithBody(method, endpoint string, payload interface{}) ([]byte, int, error) {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, fmt.Errorf("encode payload: %w", err)
	}
	var statusCode int
	resp, err := utils.DoWithRetry(func() (*http.Response, error) {
		req, err := c.buildRequest(method, endpoint, encoded)
		if err != nil {
			return nil, err
		}
		return c.httpClient.Do(req)
	}, c.maxRetries)
	if err != nil {
		return nil, 0, err
	}
	statusCode = resp.StatusCode
	body, err := utils.ReadBody(resp)
	if err != nil {
		return nil, statusCode, err
	}
	if statusCode >= 400 {
		return nil, statusCode, parseAPIError(body, statusCode)
	}
	return body, statusCode, nil
}

func (c *Client) doNoBody(method, endpoint string) ([]byte, int, error) {
	var statusCode int
	resp, err := utils.DoWithRetry(func() (*http.Response, error) {
		req, err := c.buildRequest(method, endpoint, nil)
		if err != nil {
			return nil, err
		}
		return c.httpClient.Do(req)
	}, c.maxRetries)
	if err != nil {
		return nil, 0, err
	}
	statusCode = resp.StatusCode
	body, err := utils.ReadBody(resp)
	if err != nil {
		return nil, statusCode, err
	}
	if statusCode >= 400 {
		return nil, statusCode, parseAPIError(body, statusCode)
	}
	return body, statusCode, nil
}

func parseAPIError(body []byte, status int) error {
	var apiErr APIError
	if err := json.Unmarshal(body, &apiErr); err != nil {
		return fmt.Errorf("HTTP %d: %s", status, string(body))
	}
	apiErr.Status = status
	return &apiErr
}
