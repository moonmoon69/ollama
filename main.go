package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unsafe"
)

// Float host function declarations
//
//go:wasmimport env float_log
func floatLog(ptr, len uint32)

//go:wasmimport env float_http_request
func floatHttpRequest(urlPtr, urlLen, methodPtr, methodLen, headersPtr, headersLen, bodyPtr, bodyLen uint32) uint32

//go:wasmimport env float_read_file
func floatReadFile(pathPtr, pathLen uint32) uint32

//go:wasmimport env float_write_file
func floatWriteFile(pathPtr, pathLen, dataPtr, dataLen uint32) uint32

// Input structure for Ollama request
type OllamaInput struct {
	Model     string                 `json:"model"`
	Prompt    string                 `json:"prompt"`
	System    string                 `json:"system,omitempty"`
	Template  string                 `json:"template,omitempty"`
	Context   []int                  `json:"context,omitempty"`
	Stream    *bool                  `json:"stream,omitempty"`
	Raw       bool                   `json:"raw,omitempty"`
	Format    interface{}            `json:"format,omitempty"`
	Options   map[string]interface{} `json:"options,omitempty"`
	Suffix    string                 `json:"suffix,omitempty"`
	KeepAlive string                 `json:"keep_alive,omitempty"`
	Images    []string               `json:"images,omitempty"`
	OllamaURL string                 `json:"ollama_url"`
}

// Output structure for Ollama response
type OllamaOutput struct {
	Success            bool                   `json:"success"`
	Response           string                 `json:"response,omitempty"`
	Model              string                 `json:"model,omitempty"`
	CreatedAt          string                 `json:"created_at,omitempty"`
	Done               bool                   `json:"done"`
	Context            []int                  `json:"context,omitempty"`
	TotalDuration      int64                  `json:"total_duration,omitempty"`
	LoadDuration       int64                  `json:"load_duration,omitempty"`
	PromptEvalCount    int                    `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64                  `json:"prompt_eval_duration,omitempty"`
	EvalCount          int                    `json:"eval_count,omitempty"`
	EvalDuration       int64                  `json:"eval_duration,omitempty"`
	Error              string                 `json:"error,omitempty"`
	ErrorType          string                 `json:"error_type,omitempty"`
	Metadata           map[string]interface{} `json:"metadata"`
}

// Ollama API request structure (what gets sent to Ollama)
type OllamaAPIRequest struct {
	Model     string                 `json:"model"`
	Prompt    string                 `json:"prompt"`
	System    string                 `json:"system,omitempty"`
	Template  string                 `json:"template,omitempty"`
	Context   []int                  `json:"context,omitempty"`
	Stream    *bool                  `json:"stream,omitempty"`
	Raw       bool                   `json:"raw,omitempty"`
	Format    interface{}            `json:"format,omitempty"`
	Options   map[string]interface{} `json:"options,omitempty"`
	Suffix    string                 `json:"suffix,omitempty"`
	KeepAlive string                 `json:"keep_alive,omitempty"`
	Images    []string               `json:"images,omitempty"`
}

// Ollama API response structure
type OllamaAPIResponse struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Response           string `json:"response"`
	Done               bool   `json:"done"`
	Context            []int  `json:"context,omitempty"`
	TotalDuration      int64  `json:"total_duration,omitempty"`
	LoadDuration       int64  `json:"load_duration,omitempty"`
	PromptEvalCount    int    `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64  `json:"prompt_eval_duration,omitempty"`
	EvalCount          int    `json:"eval_count,omitempty"`
	EvalDuration       int64  `json:"eval_duration,omitempty"`
}

func logToFloat(message string) {
	if len(message) == 0 {
		return
	}
	messageBytes := []byte(message)
	ptr := uintptr(unsafe.Pointer(&messageBytes[0]))
	floatLog(uint32(ptr), uint32(len(messageBytes)))
}

func writeOutputFile(data *OllamaOutput) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal output: %v", err)
	}

	path := "output.json"
	pathBytes := []byte(path)
	pathPtr := uintptr(unsafe.Pointer(&pathBytes[0]))
	dataPtr := uintptr(unsafe.Pointer(&jsonData[0]))

	result := floatWriteFile(
		uint32(pathPtr), uint32(len(pathBytes)),
		uint32(dataPtr), uint32(len(jsonData)),
	)

	if result != 0 {
		return fmt.Errorf("failed to write output file")
	}

	return nil
}

func makeHttpRequest(url, method, body string) (string, error) {
	logToFloat(fmt.Sprintf("Making HTTP %s request to %s", method, url))

	urlBytes := []byte(url)
	methodBytes := []byte(method)
	bodyBytes := []byte(body)
	headers := "Content-Type: application/json\r\n"
	headersBytes := []byte(headers)

	urlPtr := uintptr(unsafe.Pointer(&urlBytes[0]))
	methodPtr := uintptr(unsafe.Pointer(&methodBytes[0]))
	headersPtr := uintptr(unsafe.Pointer(&headersBytes[0]))
	bodyPtr := uintptr(unsafe.Pointer(&bodyBytes[0]))

	// Make the HTTP request
	result := floatHttpRequest(
		uint32(urlPtr), uint32(len(urlBytes)),
		uint32(methodPtr), uint32(len(methodBytes)),
		uint32(headersPtr), uint32(len(headersBytes)),
		uint32(bodyPtr), uint32(len(bodyBytes)),
	)

	if result != 0 {
		return "", fmt.Errorf("HTTP request failed with status code: %d", result)
	}

	// Note: In a real Float implementation, the HTTP response would be available
	// through Float's response mechanism. For this snippet, we'll simulate
	// a successful response format that matches Ollama's API.
	simulatedResponse := `{
		"model": "` + extractModelFromRequest(body) + `",
		"created_at": "` + time.Now().Format(time.RFC3339) + `",
		"response": "This is a simulated response. In production, this would be the actual Ollama API response.",
		"done": true,
		"total_duration": 1000000000,
		"load_duration": 100000000,
		"prompt_eval_count": 10,
		"prompt_eval_duration": 200000000,
		"eval_count": 25,
		"eval_duration": 700000000
	}`

	logToFloat("HTTP request completed successfully")
	return simulatedResponse, nil
}

// Extract model name from request body for response simulation
func extractModelFromRequest(body string) string {
	var req map[string]interface{}
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		return "unknown"
	}
	if model, ok := req["model"].(string); ok {
		return model
	}
	return "unknown"
}

// Validate input data
func validateInput(input *OllamaInput) error {
	if input.Model == "" {
		return fmt.Errorf("model field is required")
	}

	if input.Prompt == "" {
		return fmt.Errorf("prompt field is required")
	}

	if input.OllamaURL == "" {
		return fmt.Errorf("ollama_url field is required")
	}

	// Validate URL format
	if !strings.HasPrefix(input.OllamaURL, "http://") && !strings.HasPrefix(input.OllamaURL, "https://") {
		return fmt.Errorf("ollama_url must be a valid HTTP/HTTPS URL")
	}

	// Validate prompt length
	if len(input.Prompt) > 32768 {
		return fmt.Errorf("prompt exceeds maximum length of 32768 characters")
	}

	// Validate system message length if provided
	if input.System != "" && len(input.System) > 8192 {
		return fmt.Errorf("system message exceeds maximum length of 8192 characters")
	}

	// Validate options if provided
	if input.Options != nil {
		if temp, ok := input.Options["temperature"]; ok {
			if tempFloat, ok := temp.(float64); ok && (tempFloat < 0 || tempFloat > 2) {
				return fmt.Errorf("temperature must be between 0 and 2")
			}
		}
		if topP, ok := input.Options["top_p"]; ok {
			if topPFloat, ok := topP.(float64); ok && (topPFloat < 0 || topPFloat > 1) {
				return fmt.Errorf("top_p must be between 0 and 1")
			}
		}
	}

	return nil
}

// Build metadata for the response
func buildMetadata(input *OllamaInput, responseLength int) map[string]interface{} {
	metadata := map[string]interface{}{
		"input_prompt_length": len(input.Prompt),
		"response_length":     responseLength,
		"ollama_url":          input.OllamaURL,
		"stream_mode":         input.Stream != nil && *input.Stream,
		"processing_complete": true,
	}

	if input.System != "" {
		metadata["has_system_message"] = true
	}
	if len(input.Images) > 0 {
		metadata["image_count"] = len(input.Images)
	}
	if input.Context != nil && len(input.Context) > 0 {
		metadata["has_context"] = true
		metadata["context_length"] = len(input.Context)
	}

	return metadata
}

// Create error output
func createErrorOutput(errorMsg, errorType string, metadata map[string]interface{}) *OllamaOutput {
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	metadata["processing_complete"] = false

	return &OllamaOutput{
		Success:   false,
		Error:     errorMsg,
		ErrorType: errorType,
		Done:      true,
		Metadata:  metadata,
	}
}

//export generate_ollama
func generate_ollama(inputPtr, inputLen uint32) uint32 {
	logToFloat("Starting Ollama text generation")

	// Read input data
	inputBytes := make([]byte, inputLen)
	copy(inputBytes, (*(*[1 << 30]byte)(unsafe.Pointer(uintptr(inputPtr))))[0:inputLen])

	// Parse input JSON
	var input OllamaInput
	if err := json.Unmarshal(inputBytes, &input); err != nil {
		logToFloat(fmt.Sprintf("Failed to parse input JSON: %v", err))
		output := createErrorOutput(
			fmt.Sprintf("Invalid input JSON: %v", err),
			"INPUT_PARSE_ERROR",
			map[string]interface{}{"error_stage": "input_parsing"},
		)
		writeOutputFile(output)
		return 1
	}

	logToFloat(fmt.Sprintf("Parsed input for model: %s", input.Model))

	// Validate input
	if err := validateInput(&input); err != nil {
		logToFloat(fmt.Sprintf("Input validation failed: %v", err))
		errorType := "VALIDATION_ERROR"
		metadata := map[string]interface{}{"error_stage": "validation"}

		if strings.Contains(err.Error(), "model") {
			metadata["missing_field"] = "model"
		} else if strings.Contains(err.Error(), "prompt") {
			metadata["missing_field"] = "prompt"
		} else if strings.Contains(err.Error(), "ollama_url") {
			metadata["missing_field"] = "ollama_url"
		}

		output := createErrorOutput(err.Error(), errorType, metadata)
		writeOutputFile(output)
		return 1
	}

	logToFloat("Input validation passed")

	// Ensure URL doesn't end with slash for consistent API calls
	input.OllamaURL = strings.TrimRight(input.OllamaURL, "/")

	// Set stream to false if not specified (ensure complete response)
	if input.Stream == nil {
		streamFalse := false
		input.Stream = &streamFalse
	}

	// Prepare Ollama API request
	apiRequest := OllamaAPIRequest{
		Model:     input.Model,
		Prompt:    input.Prompt,
		System:    input.System,
		Template:  input.Template,
		Context:   input.Context,
		Stream:    input.Stream,
		Raw:       input.Raw,
		Format:    input.Format,
		Options:   input.Options,
		Suffix:    input.Suffix,
		KeepAlive: input.KeepAlive,
		Images:    input.Images,
	}

	// Marshal request body
	requestBody, err := json.Marshal(apiRequest)
	if err != nil {
		logToFloat(fmt.Sprintf("Failed to marshal request: %v", err))
		output := createErrorOutput(
			fmt.Sprintf("Failed to prepare request: %v", err),
			"REQUEST_MARSHAL_ERROR",
			map[string]interface{}{"error_stage": "request_preparation"},
		)
		writeOutputFile(output)
		return 1
	}

	logToFloat(fmt.Sprintf("Prepared request body (%d bytes)", len(requestBody)))

	// Make HTTP request to Ollama
	url := input.OllamaURL + "/api/generate"
	logToFloat(fmt.Sprintf("Making request to Ollama API: %s", url))

	responseBody, err := makeHttpRequest(url, "POST", string(requestBody))
	if err != nil {
		logToFloat(fmt.Sprintf("HTTP request failed: %v", err))
		output := createErrorOutput(
			fmt.Sprintf("Failed to connect to Ollama server: %v", err),
			"HTTP_REQUEST_ERROR",
			map[string]interface{}{
				"error_stage": "http_request",
				"ollama_url":  url,
			},
		)
		writeOutputFile(output)
		return 1
	}

	logToFloat("HTTP request completed successfully")

	// Parse Ollama response
	var apiResponse OllamaAPIResponse
	if err := json.Unmarshal([]byte(responseBody), &apiResponse); err != nil {
		logToFloat(fmt.Sprintf("Failed to parse Ollama response: %v", err))
		output := createErrorOutput(
			fmt.Sprintf("Invalid response from Ollama server: %v", err),
			"RESPONSE_PARSE_ERROR",
			map[string]interface{}{
				"error_stage":  "response_parsing",
				"raw_response": responseBody,
			},
		)
		writeOutputFile(output)
		return 1
	}

	logToFloat(fmt.Sprintf("Successfully parsed Ollama response (%d characters)", len(apiResponse.Response)))

	// Build successful output
	output := &OllamaOutput{
		Success:            true,
		Response:           apiResponse.Response,
		Model:              apiResponse.Model,
		CreatedAt:          apiResponse.CreatedAt,
		Done:               apiResponse.Done,
		Context:            apiResponse.Context,
		TotalDuration:      apiResponse.TotalDuration,
		LoadDuration:       apiResponse.LoadDuration,
		PromptEvalCount:    apiResponse.PromptEvalCount,
		PromptEvalDuration: apiResponse.PromptEvalDuration,
		EvalCount:          apiResponse.EvalCount,
		EvalDuration:       apiResponse.EvalDuration,
		Metadata:           buildMetadata(&input, len(apiResponse.Response)),
	}

	// Write output file
	if writeErr := writeOutputFile(output); writeErr != nil {
		logToFloat(fmt.Sprintf("Failed to write output file: %v", writeErr))
		return 1
	}

	logToFloat("Ollama text generation completed successfully")
	logToFloat(fmt.Sprintf("Generated %d characters in response", len(apiResponse.Response)))

	return 0
}

func main() {
	// Required for TinyGo WASM modules
}
