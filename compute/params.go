package compute

import "github.com/Askr-Omorsablin/rck-go-sdk/core"

// Example is a pair of input and structured output for learning.
type Example struct {
	Input  string
	Output map[string]interface{}
}

// StructuredTransformParams are the parameters for a structured data transformation task.
type StructuredTransformParams struct {
	Input           string
	FunctionLogic   string
	OutputDataClass interface{} // Can be a map[string]interface{} (JSON Schema) or a string
	CustomLogic     map[string]string
	Resource        []map[string]string
}

// Validate checks if the parameters are valid.
func (p *StructuredTransformParams) Validate() error {
	if p.Input == "" {
		return core.NewValidationError("Input", "is required")
	}
	if p.FunctionLogic == "" {
		return core.NewValidationError("FunctionLogic", "is required")
	}
	if p.OutputDataClass == nil {
		return core.NewValidationError("OutputDataClass", "is required")
	}
	return nil
}

// AnalyzeParams are the parameters for an analysis task.
type AnalyzeParams struct {
	Input         string
	FunctionLogic string
	OutputFormat  string // Predefined schema name
	CustomLogic   map[string]string
}

// Validate checks if the parameters are valid.
func (p *AnalyzeParams) Validate() error {
	if p.Input == "" {
		return core.NewValidationError("Input", "is required")
	}
	if p.FunctionLogic == "" {
		return core.NewValidationError("FunctionLogic", "is required")
	}
	if p.OutputFormat == "" {
		return core.NewValidationError("OutputFormat", "is required")
	}
	if !HasSchema(p.OutputFormat) {
		return core.NewValidationError("OutputFormat", "unknown schema name")
	}
	return nil
}

// TranslateParams are the parameters for a translation task.
type TranslateParams struct {
	Input                string
	TargetLanguage       string
	IncludeCulturalNotes bool
}

// Validate checks if the parameters are valid.
func (p *TranslateParams) Validate() error {
	if p.Input == "" {
		return core.NewValidationError("Input", "is required")
	}
	if p.TargetLanguage == "" {
		return core.NewValidationError("TargetLanguage", "is required")
	}
	return nil
}

// LearnFromExamplesParams are the parameters for learning from examples.
type LearnFromExamplesParams struct {
	Input       string
	Examples    []Example
	CustomLogic map[string]string
	Resource    []map[string]string
}

// Validate checks if the parameters are valid.
func (p *LearnFromExamplesParams) Validate() error {
	if p.Input == "" {
		return core.NewValidationError("Input", "is required")
	}
	if len(p.Examples) < 1 {
		return core.NewValidationError("Examples", "requires at least one example")
	}
	return nil
}

// GenerateTextParams are the parameters for a text generation task.
type GenerateTextParams struct {
	Input         string
	FunctionLogic string
	CustomLogic   map[string]string
	Resource      []map[string]string
}

// Validate checks if the parameters are valid.
func (p *GenerateTextParams) Validate() error {
	if p.Input == "" {
		return core.NewValidationError("Input", "is required")
	}
	if p.FunctionLogic == "" {
		return core.NewValidationError("FunctionLogic", "is required")
	}
	return nil
}

// AutoParams are the parameters for the automatic engine selection task.
type AutoParams struct {
	Input            string
	Resource         []map[string]string
	FunctionLogic    string
	OutputDataClass  interface{} // map[string]interface{} or string
	CustomLogic      map[string]string
	Examples         []Example
	FrameComposition string
	Lighting         string
	Style            string
}

// Validate checks if the parameters are valid.
func (p *AutoParams) Validate() error {
	if p.Input == "" {
		return core.NewValidationError("Input", "is required")
	}
	if p.FunctionLogic == "" && len(p.Examples) == 0 && p.FrameComposition == "" && p.Lighting == "" && p.Style == "" {
		return core.NewValidationError("Logic", "At least one logic-defining property is required (FunctionLogic, Examples, or image parameters)")
	}
	return nil
}
