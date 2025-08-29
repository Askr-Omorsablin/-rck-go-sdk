package compute

import (
	"encoding/json"

	"github.com/Askr-Omorsablin/rck-go-sdk/core"
)

var predefinedSchemas = map[string]string{
	"basic_analysis": `{
    "type": "object",
    "properties": {
      "emotion": { "type": "string", "description": "Emotion analysis result" },
      "theme": { "type": "string", "description": "Theme analysis" },
      "analysis": { "type": "string", "description": "Detailed analysis" }
    },
    "required": ["emotion", "theme", "analysis"]
  }`,
	"poem_creation": `{
    "type": "object",
    "properties": {
      "poem": { "type": "string", "description": "Created poem" },
      "creative_process": { "type": "string", "description": "Creative process" },
      "style_notes": { "type": "string", "description": "Style notes" }
    },
    "required": ["poem"]
  }`,
	"scene_description": `{
    "type": "object",
    "properties": {
      "scene_description": {
        "type": "object",
        "properties": {
          "main_subjects": { "type": "string", "description": "Main objects and spatial relationships" },
          "lighting": { "type": "string", "description": "Lighting conditions and atmosphere" },
          "composition": { "type": "string", "description": "Picture composition" },
          "style": { "type": "string", "description": "Artistic style" }
        },
        "required": ["main_subjects", "lighting", "composition", "style"]
      }
    },
    "required": ["scene_description"]
  }`,
	"translation": `{
    "type": "object",
    "properties": {
      "translation": { "type": "string", "description": "Translation result" },
      "original_language": { "type": "string", "description": "Source language" },
      "target_language": { "type": "string", "description": "Target language" },
      "cultural_notes": { "type": "string", "description": "Cultural background notes" }
    },
    "required": ["translation"]
  }`,
}

// GetPredefinedSchema returns the JSON string for a predefined schema.
func GetPredefinedSchema(schemaName string) (string, bool) {
	schema, ok := predefinedSchemas[schemaName]
	return schema, ok
}

// GetPredefinedSchemaAsMap returns the schema as a map.
func GetPredefinedSchemaAsMap(schemaName string) (map[string]interface{}, error) {
	schemaStr, ok := GetPredefinedSchema(schemaName)
	if !ok {
		return nil, core.NewValidationError("schemaName", "unknown schema name")
	}

	var schemaMap map[string]interface{}
	err := json.Unmarshal([]byte(schemaStr), &schemaMap)
	if err != nil {
		return nil, err
	}
	return schemaMap, nil
}

// GetAvailableSchemas returns a list of all available schema names.
func GetAvailableSchemas() []string {
	keys := make([]string, 0, len(predefinedSchemas))
	for k := range predefinedSchemas {
		keys = append(keys, k)
	}
	return keys
}

// HasSchema checks if a predefined schema exists.
func HasSchema(schemaName string) bool {
	_, ok := predefinedSchemas[schemaName]
	return ok
}
