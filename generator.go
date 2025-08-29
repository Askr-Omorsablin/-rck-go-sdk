package rck

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Generator provides access to the RCK image generation functionalities.
type Generator struct {
	client *HttpClient
}

// NewGenerator creates a new Generator instance.
func NewGenerator(client *HttpClient) *Generator {
	return &Generator{client: client}
}

// Generate creates images based on the provided parameters.
func (g *Generator) Generate(ctx context.Context, params GenerateParams) (*ImageResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	payload := &UnifiedAPIRequest{
		Config: &APIConfig{Engine: EngineImage},
		Program: APIProgram{
			Input: APIInput{Input: params.Input},
			Pipeline: APIPipeline{
				FrameComposition: params.FrameComposition,
				Lighting:         params.Lighting,
				Style:            params.Style,
			},
		},
	}

	rawResponse, err := g.client.Post(ctx, unifiedEndpoint, payload)
	if err != nil {
		return nil, err
	}

	var output []string
	if err := json.Unmarshal(rawResponse.Output, &output); err != nil {
		return nil, &APIError{
			StatusCode:   200, // Assuming 200 OK but bad format
			ResponseData: rawResponse,
		}
	}
	return NewImageResponse(output, *rawResponse), nil
}

// SaveImages saves all images from an ImageResponse to a specified directory.
func (g *Generator) SaveImages(imageResponse *ImageResponse, outputDir, baseFilename string) (savedFiles []string, errs []error) {
	if !imageResponse.Success() {
		errs = append(errs, fmt.Errorf("image response is not successful, cannot save"))
		return
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		errs = append(errs, fmt.Errorf("failed to create output directory: %w", err))
		return
	}

	for _, img := range imageResponse.Images {
		var filename string
		ext := img.GetFileExtension()
		if imageResponse.Count == 1 {
			filename = fmt.Sprintf("%s.%s", baseFilename, ext)
		} else {
			filename = fmt.Sprintf("%s_%d.%s", baseFilename, img.Index+1, ext)
		}

		filepath := filepath.Join(outputDir, filename)
		if err := os.WriteFile(filepath, img.ImageData, 0644); err != nil {
			errs = append(errs, fmt.Errorf("failed to save image %s: %w", filepath, err))
		} else {
			savedFiles = append(savedFiles, filepath)
		}
	}
	return
}
