package image

import "github.com/Askr-Omorsablin/rck-go-sdk/core"

// GenerateParams defines the parameters for generating an image.
type GenerateParams struct {
	Input            string
	FrameComposition string
	Lighting         string
	Style            string
}

// Validate checks if the parameters are valid.
func (p *GenerateParams) Validate() error {
	if p.Input == "" {
		return core.NewValidationError("Input", "is required")
	}
	if p.FrameComposition == "" {
		return core.NewValidationError("FrameComposition", "is required")
	}
	if p.Lighting == "" {
		return core.NewValidationError("Lighting", "is required")
	}
	if p.Style == "" {
		return core.NewValidationError("Style", "is required")
	}
	return nil
}
