package pluggedinkit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultBaseURL = "https://plugged.in"
	DefaultTimeout = 60 * time.Second
	Version        = "1.0.1"
)

// Client is the main Plugged.in API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	userAgent  string

	// Services
	Documents *DocumentsService
	RAG       *RAGService
	Uploads   *UploadsService
}

// NewClient creates a new Plugged.in API client
func NewClient(apiKey string) *Client {
	return NewClientWithOptions(apiKey, DefaultBaseURL, nil)
}

// NewClientWithOptions creates a new client with custom options
func NewClientWithOptions(apiKey, baseURL string, httpClient *http.Client) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: DefaultTimeout,
		}
	}

	c := &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: httpClient,
		userAgent:  fmt.Sprintf("pluggedinkit-go/%s", Version),
	}

	// Initialize services
	c.Documents = &DocumentsService{client: c}
	c.RAG = &RAGService{client: c}
	c.Uploads = &UploadsService{client: c}

	return c
}

// SetAPIKey updates the API key
func (c *Client) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
}

// SetBaseURL updates the base URL
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// request performs an HTTP request
func (c *Client) request(ctx context.Context, method, path string, body interface{}, v interface{}) error {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return err
	}
	u.Path = path

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return &APIError{
				StatusCode: resp.StatusCode,
				Message:    fmt.Sprintf("API error: %s", resp.Status),
			}
		}
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    errResp.Error,
			Details:    errResp.Details,
		}
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}

	return nil
}

// get performs a GET request
func (c *Client) get(ctx context.Context, path string, v interface{}) error {
	return c.request(ctx, http.MethodGet, path, nil, v)
}

// post performs a POST request
func (c *Client) post(ctx context.Context, path string, body, v interface{}) error {
	return c.request(ctx, http.MethodPost, path, body, v)
}

// put performs a PUT request
func (c *Client) put(ctx context.Context, path string, body, v interface{}) error {
	return c.request(ctx, http.MethodPut, path, body, v)
}

// patch performs a PATCH request
func (c *Client) patch(ctx context.Context, path string, body, v interface{}) error {
	return c.request(ctx, http.MethodPatch, path, body, v)
}

// delete performs a DELETE request
func (c *Client) delete(ctx context.Context, path string, v interface{}) error {
	return c.request(ctx, http.MethodDelete, path, nil, v)
}

// download performs a GET request and returns raw binary data
func (c *Client) download(ctx context.Context, path string) ([]byte, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}
	u.Path = path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, &APIError{
				StatusCode: resp.StatusCode,
				Message:    fmt.Sprintf("API error: %s", resp.Status),
			}
		}
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    errResp.Error,
			Details:    errResp.Details,
		}
	}

	// Read all the binary data
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
