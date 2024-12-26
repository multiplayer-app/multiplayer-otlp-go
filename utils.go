package multiplayer

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// GenerateSchema recursively generates a JSON schema for a given interface{}.
func generateSchema(data interface{}) map[string]interface{} {
	switch v := data.(type) {
	case map[string]interface{}: // Object
		properties := make(map[string]interface{})
		for key, value := range v {
			properties[key] = generateSchema(value)
		}
		return map[string]interface{}{
			"type":       "object",
			"properties": properties,
		}
	case []interface{}: // Array
		if len(v) > 0 {
			return map[string]interface{}{
				"type":  "array",
				"items": generateSchema(v[0]),
			}
		}
		// If the array is empty, assume it's an array of `null` values
		return map[string]interface{}{
			"type": "array",
		}
	case string: // String
		return map[string]interface{}{"type": "string"}
	case float64: // Number
		return map[string]interface{}{"type": "number"}
	case bool: // Boolean
		return map[string]interface{}{"type": "boolean"}
	case nil: // Null
		return map[string]interface{}{"type": "null"}
	default:
		// For unsupported types, return a generic object type
		return map[string]interface{}{"type": "object"}
	}
}

// GenerateJSONSchema takes a byte slice, attempts to unmarshal it into JSON, and returns the JSON schema.
func GenerateJSONSchema(input []byte) []byte {
	var parsed interface{}

	if err := json.Unmarshal(input, &parsed); err != nil {
		return nil
	}
	schema := generateSchema(parsed)
	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return nil
	}

	return schemaBytes
}

// maskJSONValues recursively masks values, respecting MAX_MASK_DEPTH.
func maskJSONValues(data interface{}, depth, maxDepth int) interface{} {
	if depth > maxDepth {
		return MASK_PLACEHOLDER
	}

	switch v := data.(type) {
	case map[string]interface{}: // Object
		masked := make(map[string]interface{})
		for key, value := range v {
			masked[key] = maskJSONValues(value, depth+1, maxDepth)
		}
		return masked
	case []interface{}: // Array
		for i, item := range v {
			v[i] = maskJSONValues(item, depth+1, maxDepth)
		}
		return v
	case string, float64, bool, nil: // Primitive types
		return MASK_PLACEHOLDER
	default:
		return v
	}
}

func MaskJSONIfValid(input []byte, maxDepth int) []byte {
	var parsed interface{}

	decoder := json.NewDecoder(bytes.NewReader(input))
	if err := decoder.Decode(&parsed); err != nil {
		return []byte(MASK_PLACEHOLDER)
	}
	masked := maskJSONValues(parsed, 0, maxDepth)
	maskedBytes, err := json.Marshal(masked)
	if err != nil {
		return []byte(MASK_PLACEHOLDER)
	}
	return maskedBytes
}

func maskHeaders(headers http.Header, customHeaderNamesToMask []string) string {
	defaultHeaderNamesToMask := []string{"set-cookie", "cookie", "authorization", "proxy-authorization"}
	headersToMask := append(defaultHeaderNamesToMask, customHeaderNamesToMask...)
	headersToMaskMap := make(map[string]bool, len(headersToMask))
	for _, value := range headersToMask {
		headersToMaskMap[value] = true
	}
	normalizedHeaders := make(map[string]any, len(headers))

	for name, values := range headers {
		lowerCaseHeaderName := strings.ToLower(name)
		if _, exists := headersToMaskMap[lowerCaseHeaderName]; exists {
			normalizedHeaders[name] = MASK_PLACEHOLDER
		} else {
			normalizedHeaders[name] = values
		}
	}
	if reqHeadersBytes, err := json.Marshal(normalizedHeaders); err != nil {
		log.Println("Could not Marshal Req Headers")
		return ""
	} else {
		return string(reqHeadersBytes)
	}
}

func truncateIfNeeded(data string, maxPayloadSize int) string {
	if len(data) > maxPayloadSize {
		return data[:maxPayloadSize] + "...[TRUNCATED]"
	}
	return data
}

func IsDebugTrace(traceId string) bool {
	return strings.HasPrefix(traceId, MULTIPLAYER_TRACE_DEBUG_PREFIX)
}

func IsDocTrace(traceId string) bool {
	return strings.HasPrefix(traceId, MULTIPLAYER_TRACE_DOC_PREFIX)
}

func IsMultiplayerTrace(traceId string) bool {
	return IsDebugTrace(traceId) || IsDocTrace(traceId)
}
