package compute

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/Askr-Omorsablin/rck-go-sdk/core"
	"github.com/Askr-Omorsablin/rck-go-sdk/image"
)

const unifiedEndpoint = "/calculs"

// Kernel provides access to the RCK compute functionalities.
type Kernel struct {
	client *core.HttpClient
}

// NewKernel creates a new Kernel instance.
func NewKernel(client *core.HttpClient) *Kernel {
	return &Kernel{client: client}
}

func (k *Kernel) execute(ctx context.Context, program core.APIProgram, config *core.APIConfig) (*core.UnifiedAPIResponse, error) {
	payload := &core.UnifiedAPIRequest{
		Program: program,
		Config:  config,
	}
	return k.client.Post(ctx, unifiedEndpoint, payload)
}

func outputToInterface(output json.RawMessage) (interface{}, error) {
	var strVal string
	if err := json.Unmarshal(output, &strVal); err == nil {
		return strVal, nil
	}
	var strArrVal []string
	if err := json.Unmarshal(output, &strArrVal); err == nil {
		return strArrVal, nil
	}
	var mapVal map[string]interface{}
	if err := json.Unmarshal(output, &mapVal); err == nil {
		return mapVal, nil
	}
	return nil, errors.New("failed to determine output type")
}

// Auto determines the engine automatically based on parameters.
func (k *Kernel) Auto(ctx context.Context, params AutoParams) (interface{}, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	var pipeline core.APIPipeline
	pipeline.FunctionLogic = params.FunctionLogic
	pipeline.CustomLogic = params.CustomLogic
	pipeline.FrameComposition = params.FrameComposition
	pipeline.Lighting = params.Lighting
	pipeline.Style = params.Style

	if params.OutputDataClass != nil {
		outputClassStr, err := stringifyOutputClass(params.OutputDataClass)
		if err != nil {
			return nil, err
		}
		pipeline.OutputDataClass = outputClassStr
	}

	if len(params.Examples) > 0 {
		apiExamples := make([]core.APIExample, len(params.Examples))
		for i, ex := range params.Examples {
			outputBytes, err := json.Marshal(ex.Output)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal example output at index %d: %w", i, err)
			}
			apiExamples[i] = core.APIExample{Input: ex.Input, Output: string(outputBytes)}
		}
		pipeline.Examples = apiExamples
	}

	program := core.APIProgram{
		Input: core.APIInput{
			Input:    params.Input,
			Resource: params.Resource,
		},
		Pipeline: pipeline,
	}

	response, err := k.execute(ctx, program, nil) // No config, let server decide
	if err != nil {
		return nil, err
	}
	// Handle cases where the API returns a success status but an empty output.
	if response.Output == nil || string(response.Output) == "null" {
		return nil, &core.APIError{
			StatusCode:   200,
			ResponseData: response,
		}
	}

	val, err := outputToInterface(response.Output)
	if err != nil {
		return nil, err
	}

	switch v := val.(type) {
	case string:
		return v, nil
	case []string:
		return image.NewImageResponse(v, *response), nil
	case map[string]interface{}:
		return NewComputeResponse(*response), nil
	default:
		return response.Output, nil
	}
}

// StructuredTransform performs a data transformation based on a schema and logic.
func (k *Kernel) StructuredTransform(ctx context.Context, params StructuredTransformParams, config ...core.ComputeConfig) (*ComputeResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	outputClassStr, err := stringifyOutputClass(params.OutputDataClass)
	if err != nil {
		return nil, err
	}

	program := core.APIProgram{
		Input: core.APIInput{
			Input:    params.Input,
			Resource: params.Resource,
		},
		Pipeline: core.APIPipeline{
			OutputDataClass: outputClassStr,
			FunctionLogic:   params.FunctionLogic,
			CustomLogic:     params.CustomLogic,
		},
	}
	apiConfig := &core.APIConfig{Engine: core.EngineStandard}
	if len(config) > 0 {
		apiConfig.ComputeConfig = config[0]
	}

	response, err := k.execute(ctx, program, apiConfig)
	if err != nil {
		return nil, err
	}
	return NewComputeResponse(*response), nil
}

// LearnFromExamples learns a transformation from input-output examples.
func (k *Kernel) LearnFromExamples(ctx context.Context, params LearnFromExamplesParams, config ...core.ComputeConfig) (*ComputeResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}
	apiExamples := make([]core.APIExample, len(params.Examples))
	for i, ex := range params.Examples {
		outputBytes, err := json.Marshal(ex.Output)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal example output at index %d: %w", i, err)
		}
		apiExamples[i] = core.APIExample{Input: ex.Input, Output: string(outputBytes)}
	}

	program := core.APIProgram{
		Input: core.APIInput{
			Input:    params.Input,
			Resource: params.Resource,
		},
		Pipeline: core.APIPipeline{
			Examples:    apiExamples,
			CustomLogic: params.CustomLogic,
		},
	}
	apiConfig := &core.APIConfig{Engine: core.EngineAttractor}
	if len(config) > 0 {
		apiConfig.ComputeConfig = config[0]
	}

	response, err := k.execute(ctx, program, apiConfig)
	if err != nil {
		return nil, err
	}
	return NewComputeResponse(*response), nil
}

// GenerateText generates free-form text based on a prompt and logic.
func (k *Kernel) GenerateText(ctx context.Context, params GenerateTextParams, config ...core.ComputeConfig) (string, error) {
	if err := params.Validate(); err != nil {
		return "", err
	}
	program := core.APIProgram{
		Input: core.APIInput{
			Input:    params.Input,
			Resource: params.Resource,
		},
		Pipeline: core.APIPipeline{
			FunctionLogic: params.FunctionLogic,
			CustomLogic:   params.CustomLogic,
		},
	}
	apiConfig := &core.APIConfig{Engine: core.EnginePure}
	if len(config) > 0 {
		apiConfig.ComputeConfig = config[0]
	}

	response, err := k.execute(ctx, program, apiConfig)
	if err != nil {
		return "", err
	}
	var output string
	if err := json.Unmarshal(response.Output, &output); err != nil {
		return string(response.Output), nil
	}
	return output, nil
}

// Analyze performs structured analysis using a predefined output format.
func (k *Kernel) Analyze(ctx context.Context, params AnalyzeParams, config ...core.ComputeConfig) (*ComputeResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}
	schema, ok := GetPredefinedSchema(params.OutputFormat)
	if !ok {
		return nil, core.NewValidationError("OutputFormat", "unknown schema name")
	}
	transformParams := StructuredTransformParams{
		Input:           params.Input,
		FunctionLogic:   params.FunctionLogic,
		OutputDataClass: schema,
		CustomLogic:     params.CustomLogic,
	}
	return k.StructuredTransform(ctx, transformParams, config...)
}

// Translate translates text to a target language.
func (k *Kernel) Translate(ctx context.Context, params TranslateParams, config ...core.ComputeConfig) (*ComputeResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}
	functionLogic := fmt.Sprintf("Translate text to %s", params.TargetLanguage)
	if params.IncludeCulturalNotes {
		functionLogic += " and provide cultural background notes"
	}
	schema, _ := GetPredefinedSchema("translation")
	customLogic := map[string]string{
		"target_language":        params.TargetLanguage,
		"include_cultural_notes": strconv.FormatBool(params.IncludeCulturalNotes),
	}
	transformParams := StructuredTransformParams{
		Input:           params.Input,
		FunctionLogic:   functionLogic,
		OutputDataClass: schema,
		CustomLogic:     customLogic,
	}
	return k.StructuredTransform(ctx, transformParams, config...)
}

func stringifyOutputClass(outputClass interface{}) (string, error) {
	switch v := outputClass.(type) {
	case string:
		// Assume it's already a JSON string or a description
		return v, nil
	case map[string]interface{}:
		bytes, err := json.Marshal(v)
		if err != nil {
			return "", fmt.Errorf("failed to marshal OutputDataClass map to JSON: %w", err)
		}
		return string(bytes), nil
	default:
		return "", core.NewValidationError("OutputDataClass", "must be a string or a map[string]interface{}")
	}
}
