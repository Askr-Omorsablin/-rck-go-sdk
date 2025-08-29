package rck

import (
	"encoding/base64"
	"regexp"
	"strings"
)

var dataURLRegex = regexp.MustCompile(`^data:(.+?);base64,(.*)$`)

// ImageInfo holds metadata and data for a single generated image.
type ImageInfo struct {
	ImageData []byte // Raw image bytes
	Index     int
	MimeType  string
	DataURL   string // The full data URL
}

// ImageResponse wraps the results from an image generation request.
type ImageResponse struct {
	Images  []ImageInfo
	Count   int
	RawData UnifiedAPIResponse
}

// NewImageResponse creates an ImageResponse from the API output.
func NewImageResponse(dataUrls []string, rawData UnifiedAPIResponse) *ImageResponse {
	images := make([]ImageInfo, 0, len(dataUrls))
	for i, url := range dataUrls {
		matches := dataURLRegex.FindStringSubmatch(url)
		if len(matches) != 3 {
			continue // Skip invalid Data URLs
		}
		mimeType := matches[1]
		b64Data := matches[2]

		// The standard library's base64 decoder handles padding correctly
		imgData, err := base64.StdEncoding.DecodeString(b64Data)
		if err != nil {
			continue // Skip on decoding error
		}

		images = append(images, ImageInfo{
			ImageData: imgData,
			Index:     i,
			MimeType:  mimeType,
			DataURL:   url,
		})
	}
	return &ImageResponse{
		Images:  images,
		Count:   len(images),
		RawData: rawData,
	}
}

// Success checks if the request resulted in at least one image.
func (r *ImageResponse) Success() bool {
	return r.Count > 0
}

// GetFirstImage returns the first image, or nil if none exist.
func (r *ImageResponse) GetFirstImage() *ImageInfo {
	if r.Count > 0 {
		return &r.Images[0]
	}
	return nil
}

// GetFileExtension returns the common file extension for the image's MIME type.
func (i *ImageInfo) GetFileExtension() string {
	parts := strings.Split(i.MimeType, "/")
	if len(parts) == 2 {
		ext := parts[1]
		if ext == "jpeg" {
			return "jpg"
		}
		return ext
	}
	return "png" // Default
}
