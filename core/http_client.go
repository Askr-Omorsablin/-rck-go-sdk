package core

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const sdkVersion = "1.0.0"

// HttpClient is responsible for making authenticated HTTP requests to the RCK API.
type HttpClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewHttpClient creates a new instance of the HttpClient.
func NewHttpClient(apiKey, baseURL string, timeout time.Duration) *HttpClient {
	return &HttpClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// Post sends a POST request to the specified endpoint.
func (c *HttpClient) Post(ctx context.Context, endpoint string, payload *UnifiedAPIRequest) (*UnifiedAPIResponse, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, &NetworkError{Message: "failed to marshal request payload", OriginalError: err}
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, &NetworkError{Message: "failed to create HTTP request", OriginalError: err}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.apiKey)
	req.Header.Set("User-Agent", "RCK-GO-SDK/"+sdkVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, &NetworkError{Message: "request timeout"}
		}
		return nil, &NetworkError{Message: "network request failed", OriginalError: err}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &NetworkError{Message: "failed to read response body", OriginalError: err}
	}

	var apiResponse UnifiedAPIResponse
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		if resp.StatusCode >= 400 {
			return nil, &APIError{
				StatusCode: resp.StatusCode,
				ResponseData: &UnifiedAPIResponse{
					Error:   "Invalid response format",
					Details: string(bodyBytes),
				},
			}
		}
		return nil, &NetworkError{Message: "failed to unmarshal response JSON", OriginalError: err}
	}

	if resp.StatusCode >= 400 {
		return nil, c.handleErrorResponse(resp.StatusCode, &apiResponse)
	}

	return &apiResponse, nil
}

func (c *HttpClient) handleErrorResponse(statusCode int, responseData *UnifiedAPIResponse) error {
	if statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden {
		return ErrAuthentication
	}

	return &APIError{
		StatusCode:   statusCode,
		ResponseData: responseData,
	}
}
