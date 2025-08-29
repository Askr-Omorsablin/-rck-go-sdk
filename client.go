package rck

import (
	"context"
	"time"

	"github.com/Askr-Omorsablin/rck-go-sdk/compute"
	"github.com/Askr-Omorsablin/rck-go-sdk/core"
	"github.com/Askr-Omorsablin/rck-go-sdk/image"
)

const (
	defaultBaseURL = "https://rck-aehhddpisa.us-west-1.fcapp.run"
	defaultTimeout = 60 * time.Second
)

// Client is the main entry point for the RCK SDK.
type Client struct {
	Compute *compute.Kernel
	Image   *image.Generator
	client  *core.HttpClient
}

// NewClient creates a new RCK client.
// An API key is required. Options can be nil.
func NewClient(apiKey string, options *core.ClientOptions) (*Client, error) {
	if apiKey == "" {
		return nil, core.ErrAPIKeyRequired
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

	httpClient := core.NewHttpClient(apiKey, baseURL, timeout)

	return &Client{
		Compute: compute.NewKernel(httpClient),
		Image:   image.NewGenerator(httpClient),
		client:  httpClient,
	}, nil
}

// TestConnection sends a simple request to the API to verify connectivity and authentication.
func (c *Client) TestConnection(ctx context.Context) error {
	params := compute.StructuredTransformParams{
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
