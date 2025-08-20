// Package httpclient provides an abstraction for the HTTP client with automatic token refreshing.
package httpclient

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

type Client struct {
	baseURL        string
	httpClient     *http.Client
	tokenGetter    TokenGetter
	tokenRefresher TokenRefresher
}

func NewClient(baseURL string, tokenGetter TokenGetter, tokenRefresher TokenRefresher) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		tokenGetter:    tokenGetter,
		tokenRefresher: tokenRefresher,
	}
}

func (c *Client) GetBaseURL() string {
	return c.baseURL
}

// Do performs a HTTP Request and automatically refreshes the access token if it is expired
func (c *Client) Do(ctx context.Context, req Request) (*Response, error) {
	if req.Protected {
		if c.tokenGetter.IsAccessTokenExpired() || time.Until(c.tokenGetter.GetAccessTokenExpiresAt()) < 5*time.Minute {
			if err := c.tokenRefresher.RefreshToken(ctx); err != nil {
				return nil, fmt.Errorf("failed to refresh token: %w", err)
			}
		}
		resp, err := c.doRequest(ctx, req, true)
		if err != nil {
			return nil, fmt.Errorf("failed to perform request: %w", err)
		}

		if resp.StatusCode == http.StatusUnauthorized {
			if err := c.tokenRefresher.RefreshToken(ctx); err != nil {
				return nil, fmt.Errorf("failed to refresh token: %w", err)
			}
			resp, err = c.doRequest(ctx, req, true)
			if err != nil {
				return nil, fmt.Errorf("failed to perform request: %w", err)
			}
		}
		return resp, nil
	} else {
		resp, err := c.doRequest(ctx, req, false)
		if err != nil {
			return nil, fmt.Errorf("failed to perform request: %w", err)
		}
		return resp, nil
	}
}

// doRequest performs a HTTP Request
func (c *Client) doRequest(ctx context.Context, req Request, useRefreshToken bool) (*Response, error) {
	var body io.Reader
	if req.Body != nil {
		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		body = bytes.NewBuffer(jsonBody)
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, c.baseURL+req.URL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	//Default Headers
	httpReq.Header.Set("Content-Type", "application/json")

	if len(req.QueryParams) > 0 {
		queryParams := url.Values{}
		for key, value := range req.QueryParams {
			queryParams.Add(key, value)
		}
		httpReq.URL.RawQuery = queryParams.Encode()
	}

	if req.Protected {
		if c.tokenGetter.GetAccessToken() != "" {
			httpReq.Header.Set("Authorization", "Bearer "+c.tokenGetter.GetAccessToken())
		}
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
	}, nil
}

// DoJSON performs a request and unmarshals the response body into the provided struct
func (c *Client) DoJSON(ctx context.Context, req Request, respBody interface{}) error {
	resp, err := c.Do(ctx, req)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	if respBody != nil {
		if err := json.Unmarshal(resp.Body, respBody); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}

	return nil
}

func (c *Client) Get(ctx context.Context, url string, headers map[string]string, queryParams map[string]string, protected bool) (*Response, error) {
	return c.Do(ctx, Request{
		Method:      http.MethodGet,
		Headers:     headers,
		URL:         url,
		QueryParams: queryParams,
		Protected:   protected,
	})
}

func (c *Client) Post(ctx context.Context, url string, body interface{}, headers map[string]string, protected bool, queryParams map[string]string) (*Response, error) {
	return c.Do(ctx, Request{
		Method:      http.MethodPost,
		Headers:     headers,
		URL:         url,
		Body:        body,
		QueryParams: queryParams,
		Protected:   protected,
	})
}

func (c *Client) Put(ctx context.Context, url string, body interface{}, headers map[string]string, protected bool, queryParams map[string]string) (*Response, error) {
	return c.Do(ctx, Request{
		Method:      http.MethodPut,
		Headers:     headers,
		URL:         url,
		Body:        body,
		QueryParams: queryParams,
		Protected:   protected,
	})
}

func (c *Client) Delete(ctx context.Context, url string, body interface{}, headers map[string]string, protected bool, queryParams map[string]string) (*Response, error) {
	return c.Do(ctx, Request{
		Method:      http.MethodDelete,
		Headers:     headers,
		URL:         url,
		Body:        body,
		QueryParams: queryParams,
		Protected:   protected,
	})
}

func (c *Client) Patch(ctx context.Context, url string, body interface{}, headers map[string]string, protected bool, queryParams map[string]string) (*Response, error) {
	return c.Do(ctx, Request{
		Method:      http.MethodPatch,
		Headers:     headers,
		URL:         url,
		Body:        body,
		QueryParams: queryParams,
		Protected:   protected,
	})
}
