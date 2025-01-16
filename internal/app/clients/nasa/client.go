package nasa

import (
	"biggest-mars-pictures/internal/app/clients/nasa/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sheepla/go-urlbuilder"
	"io"
	"net/http"
	"net/url"
	"time"
)

var ApiEndpoint = urlbuilder.MustParse("https://api.nasa.gov/mars-photos/api/v1/rovers/curiosity/photos")

var HTTPClient = http.Client{
	Timeout: time.Second * 10,
}

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
	}
}

func (c *Client) FindNasaPhotos(ctx context.Context, sol string) (*models.NasaPhotos, error) {
	var currentUrl = buildUrl(c.apiKey, sol)

	req, err := http.NewRequestWithContext(ctx, "GET", currentUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("could not close response body: %v\n", err)
		}
	}()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, response body: %s\n", resp.StatusCode, responseBytes)
	}
	var nasaPhotos models.NasaPhotos
	if err := json.Unmarshal(responseBytes, &nasaPhotos); err != nil {
		return nil, fmt.Errorf("could not parse response body: %w", err)
	}
	return &nasaPhotos, nil
}

func (c *Client) FindPhotoSize(ctx context.Context, imgUrl string) (int64, error) {
	req, err := http.NewRequestWithContext(ctx, "HEAD", imgUrl, nil)
	if err != nil {
		return 0, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := HTTPClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("could not send request: %w", err)
	}
	return resp.ContentLength, nil
}

func buildUrl(apiKey, sol string) string {
	ApiEndpoint.
		EditQuery(func(q url.Values) url.Values {
			q.Set("sol", sol)
			q.Set("api_key", apiKey)
			return q
		})
	return ApiEndpoint.MustString()
}

func (c *Client) ConvertToBytes(ctx context.Context, imgUrl string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", imgUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := HTTPClient.Do(req)
	respBytes := new(bytes.Buffer)
	_, err = resp.Body.Read(respBytes.Bytes())
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("could not close response body: %v\n", err)
		}
	}(resp.Body)
	return respBytes.Bytes(), nil
}
