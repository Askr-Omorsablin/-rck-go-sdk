package rck

import (
	"context"
	"time"
)

const (
	defaultBaseURL = "https://rck-aehhddpisa.us-west-1.fcapp.run"
	defaultTimeout = 60 * time.Second
)

// Client is the main entry point for the RCK SDK.
type Client struct {
	Compute *Kernel
	Image   *Generator
	client  *HttpClient
}

// NewClient creates a new RCK client.
// An API key is required. Options can be nil.
func NewClient(apiKey string, options *ClientOptions) (*Client, error) {
	if apiKey == "" {
		return nil, ErrAPIKeyRequired
	}

	baseURL := defaultBaseURL
	timeout := defaultTimeout

	if options != nil {
		if options.BaseURL != "" {
			baseURL = options.BaseURL
		}
		if options.Timeout > 0 {
			timeout = time.Duration(options.Timeout) * time.Millisecond
		}
	}

	httpClient := NewHttpClient(apiKey, baseURL, timeout)

	return &Client{
		Compute: NewKernel(httpClient),
		Image:   NewGenerator(httpClient),
		client:  httpClient,
	}, nil
}

// TestConnection sends a simple request to the API to verify connectivity and authentication.
func (c *Client) TestConnection(ctx context.Context) error {
	params := StructuredTransformParams{
		Input:         "test",
		FunctionLogic: "simple analysis",
		OutputDataClass: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"result": map[string]string{"type": "string"},
			},
		},
	}
	_, err := c.Compute.StructuredTransform(ctx, params)
	return err
}
