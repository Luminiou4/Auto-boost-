package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"
)

func DoWithRetry(fn func() (*http.Response, error), maxRetries int) (*http.Response, error) {
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, err := fn()
		if err != nil {
			lastErr = err
			backoff := time.Duration(math.Pow(2, float64(attempt))) * 500 * time.Millisecond
			time.Sleep(backoff)
			continue
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := parseRetryAfter(resp)
			resp.Body.Close()
			if attempt < maxRetries {
				time.Sleep(retryAfter)
				continue
			}
			return nil, fmt.Errorf("rate limited after %d attempts, retry after %s", maxRetries, retryAfter)
		}

		return resp, nil
	}
	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

func parseRetryAfter(resp *http.Response) time.Duration {

	if val := resp.Header.Get("X-RateLimit-Reset-After"); val != "" {
		if secs, err := strconv.ParseFloat(val, 64); err == nil {
			return time.Duration(secs*1000) * time.Millisecond
		}
	}

	if val := resp.Header.Get("Retry-After"); val != "" {
		if secs, err := strconv.Atoi(val); err == nil {
			return time.Duration(secs) * time.Second
		}
	}

	return 5 * time.Second
}

func ReadBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func EncodeJSON(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
