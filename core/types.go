package core

import "encoding/json"

// Engine defines the type for the compute engine.
type Engine string

// Available engines
const (
	EngineStandard  Engine = "standard"
	EngineAttractor Engine = "attractor"
	EngineImage     Engine = "image"
	EnginePure      Engine = "pure"
	EngineAuto      Engine = "auto"
)

// Speed defines the optimization strategy.
type Speed string

// Available speeds
const (
	SpeedFast     Speed = "fast"
	SpeedBalanced Speed = "balanced"
	SpeedQuality  Speed = "quality"
)

// Scale defines the resource allocation size.
type Scale string

// Available scales
const (
	ScaleLow    Scale = "low"
	ScaleMedium Scale = "medium"
	ScaleHigh   Scale = "high"
)

// ClientOptions holds configuration for the client.
type ClientOptions struct {
	Timeout int // Request timeout in milliseconds
	BaseURL string
}

// ComputeConfig holds execution configuration for a compute request.
// It corresponds to the 'config' object in the API, excluding the 'engine'.
type ComputeConfig struct {
	Speed       Speed    `json:"speed,omitempty"`
	Scale       Scale    `json:"scale,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
}

// APIConfig is the complete configuration sent to the API.
type APIConfig struct {
	ComputeConfig
	Engine Engine `json:"engine"`
}

// APIInput represents the input data for a program.
type APIInput struct {
	Input    string              `json:"input"`
	Resource []map[string]string `json:"resource,omitempty"`
}

// APIExample is a single input-output pair for the attractor engine.
// The Output field must be a JSON-encoded string.
type APIExample struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

// APIPipeline defines the processing logic for a program.
type APIPipeline struct {
	FunctionName     string            `json:"FunctionName,omitempty"`
	OutputDataClass  string            `json:"OutputDataClass,omitempty"`
	FunctionLogic    string            `json:"FunctionLogic,omitempty"`
	CustomLogic      map[string]string `json:"CustomLogic,omitempty"`
	Examples         []APIExample      `json:"Examples,omitempty"`
	FrameComposition string            `json:"frame_composition,omitempty"`
	Lighting         string            `json:"lighting,omitempty"`
	Style            string            `json:"style,omitempty"`
}

// APIProgram is the core computation task definition.
type APIProgram struct {
	Input    APIInput    `json:"input"`
	Pipeline APIPipeline `json:"Pipeline"`
}

// UnifiedAPIRequest is the top-level structure for an API request.
type UnifiedAPIRequest struct {
	Config  *APIConfig `json:"config,omitempty"`
	Program APIProgram `json:"program"`
}

// UnifiedAPIResponse is the top-level structure for an API response.
type UnifiedAPIResponse struct {
	Output  json.RawMessage `json:"output"`
	Error   string          `json:"error,omitempty"`
	Details string          `json:"details,omitempty"`
}
