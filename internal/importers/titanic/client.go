package titanic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHttpTimeout  = 20 * time.Second
	defaultRetryDelay   = 60 * time.Second
	defaultMaxHttpTries = 3
)

type HttpStatusError struct {
	statusCode int
	url        string
	body       string
}

func (err *HttpStatusError) Error() string {
	if err.body == "" {
		return fmt.Sprintf("unexpected status %d from %s", err.statusCode, err.url)
	}
	return fmt.Sprintf("unexpected status %d from %s: %s", err.statusCode, err.url, err.body)
}

func (importer *TitanicImporter) GetJson(url string, target any) error {
	resp, err := importer.PerformRequest(func() (*http.Request, error) {
		return http.NewRequest(http.MethodGet, url, nil)
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	return nil
}

func (importer *TitanicImporter) PostJson(url string, payload []byte, target any) error {
	resp, err := importer.PerformRequest(func() (*http.Request, error) {
		req, reqErr := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
		if reqErr != nil {
			return nil, reqErr
		}
		req.Header.Set("Content-Type", "application/json")
		return req, nil
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	return nil
}

func (importer *TitanicImporter) GetBytes(url string) ([]byte, error) {
	resp, err := importer.PerformRequest(func() (*http.Request, error) {
		return http.NewRequest(http.MethodGet, url, nil)
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (importer *TitanicImporter) PerformRequest(requestFactory func() (*http.Request, error)) (*http.Response, error) {
	for attempt := range defaultMaxHttpTries {
		req, err := requestFactory()
		if err != nil {
			return nil, err
		}

		resp, err := importer.client.Do(req)
		if err != nil {
			if attempt == defaultMaxHttpTries-1 {
				return nil, err
			}
			time.Sleep(backoffDelay(attempt))
			continue
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			// Try to get retry delay from header, otherwise use default
			retryAfter := retryAfterDuration(resp.Header.Get("Retry-After"))
			resp.Body.Close()
			time.Sleep(retryAfter)
			continue
		}

		if resp.StatusCode >= http.StatusInternalServerError {
			if attempt == defaultMaxHttpTries-1 {
				body := readErrorBody(resp.Body)
				resp.Body.Close()
				return nil, &HttpStatusError{statusCode: resp.StatusCode, url: req.URL.String(), body: body}
			}
			resp.Body.Close()
			time.Sleep(backoffDelay(attempt))
			continue
		}

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			body := readErrorBody(resp.Body)
			resp.Body.Close()
			return nil, &HttpStatusError{statusCode: resp.StatusCode, url: req.URL.String(), body: body}
		}

		return resp, nil
	}

	return nil, fmt.Errorf("request failed after %d attempts", defaultMaxHttpTries)
}

func readErrorBody(body io.Reader) string {
	data, err := io.ReadAll(io.LimitReader(body, 4*1024))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func retryAfterDuration(value string) time.Duration {
	if value == "" {
		return defaultRetryDelay
	}
	if seconds, err := strconv.Atoi(value); err == nil {
		return time.Duration(seconds) * time.Second
	}
	if parsed, err := http.ParseTime(value); err == nil {
		delta := time.Until(parsed)
		if delta > 0 {
			return delta
		}
	}
	return defaultRetryDelay
}

func backoffDelay(attempt int) time.Duration {
	base := time.Second * 2
	delay := base * time.Duration(1<<attempt)
	if delay > defaultRetryDelay {
		return defaultRetryDelay
	}
	return delay
}
