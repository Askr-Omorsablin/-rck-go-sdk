package rck

import (
	"encoding/json"
)

// ComputeResponse wraps the structured data returned by the compute API.
type ComputeResponse struct {
	rawData UnifiedAPIResponse
}

// NewComputeResponse creates a response from the raw API output.
func NewComputeResponse(raw UnifiedAPIResponse) *ComputeResponse {
	return &ComputeResponse{rawData: raw}
}

// Decode unmarshals the response data into the provided struct.
func (r *ComputeResponse) Decode(v interface{}) error {
	return json.Unmarshal(r.rawData.Output, v)
}

// AsMap returns the response data as a map.
func (r *ComputeResponse) AsMap() (map[string]interface{}, error) {
	var data map[string]interface{}
	err := r.Decode(&data)
	return data, err
}

// Raw returns the raw JSON bytes of the output field.
func (r *ComputeResponse) Raw() json.RawMessage {
	return r.rawData.Output
}
