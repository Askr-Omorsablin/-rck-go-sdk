package rck

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
		return NewValidationError("Input", "is required")
	}
	if p.FrameComposition == "" {
		return NewValidationError("FrameComposition", "is required")
	}
	if p.Lighting == "" {
		return NewValidationError("Lighting", "is required")
	}
	if p.Style == "" {
		return NewValidationError("Style", "is required")
	}
	return nil
}
