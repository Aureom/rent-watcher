package geolocation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"rent-watcher/internal/models"
)

const (
	baseURL      = "https://maps.googleapis.com/maps/api/distancematrix/json"
	timeout      = 10 * time.Second
	maxRetries   = 3
	retryBackoff = 1 * time.Second
)

type DistanceMatrixResponse struct {
	Rows []struct {
		Elements []struct {
			Distance struct {
				Value int `json:"value"`
			} `json:"distance"`
		} `json:"elements"`
	} `json:"rows"`
}

type GoogleMapsClient struct {
	APIKey     string
	HTTPClient *http.Client
}

func NewGoogleMapsClient(apiKey string) *GoogleMapsClient {
	return &GoogleMapsClient{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *GoogleMapsClient) CalculateDistance(ctx context.Context, property *models.Property, destLat, destLng float64) (int, error) {
	origin := fmt.Sprintf("%s,%s,%s", property.Logradouro, property.Bairro, property.Cidade)
	destination := fmt.Sprintf("%.7f,%.7f", destLat, destLng)

	u, err := url.Parse(baseURL)
	if err != nil {
		return 0, fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := u.Query()
	q.Set("origins", origin)
	q.Set("destinations", destination)
	q.Set("key", c.APIKey)
	u.RawQuery = q.Encode()

	var distance int
	err = c.doWithRetry(ctx, u.String(), maxRetries, func(resp *http.Response) error {
		var result DistanceMatrixResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}

		if len(result.Rows) == 0 || len(result.Rows[0].Elements) == 0 {
			return fmt.Errorf("no distance data in response: %w", err)
		}

		distance = result.Rows[0].Elements[0].Distance.Value
		return nil
	})

	if err != nil {
		return 0, err
	}

	return distance, nil
}

func (c *GoogleMapsClient) doWithRetry(ctx context.Context, url string, maxRetries int, process func(*http.Response) error) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(retryBackoff * time.Duration(i+1))
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			time.Sleep(retryBackoff * time.Duration(i+1))
			continue
		}

		if err := process(resp); err != nil {
			lastErr = err
			time.Sleep(retryBackoff * time.Duration(i+1))
			continue
		}

		return nil
	}

	return fmt.Errorf("max retries reached: %w", lastErr)
}
