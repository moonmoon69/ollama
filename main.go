package main

import (
	"encoding/json"
	"fmt"
	"strings"
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
	Suffix    string                 `json:"suffix,omitempty"`
	System    string                 `json:"system,omitempty"`
	Template  string                 `json:"template,omitempty"`
	Context   []int                  `json:"context,omitempty"`
	Stream    *bool                  `json:"stream,omitempty"` // Pointer to allow nil (default)
	Raw       bool                   `json:"raw,omitempty"`
	Format    interface{}            `json:"format,omitempty"`     // Can be string "json" or JSON schema object
	KeepAlive string                 `json:"keep_alive,omitempty"` // Duration string like "5m"
	Images    []string               `json:"images,omitempty"`     // Base64-encoded images for multimodal models
	Options   map[string]interface{} `json:"options,omitempty"`
	OllamaURL string                 `json:"ollama_url,omitempty"` // Configuration field, not sent to API
}

// Output structure for Ollama response
type OllamaOutput struct {
	Success            bool                   `json:"success"`
	Response           string                 `json:"response"`
	Model              string                 `json:"model"`
	CreatedAt          string                 `json:"created_at"`
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

// Ollama API request structure
type OllamaAPIRequest struct {
	Model     string                 `json:"model"`
	Prompt    string                 `json:"prompt"`
	Suffix    string                 `json:"suffix,omitempty"`
	System    string                 `json:"system,omitempty"`
	Template  string                 `json:"template,omitempty"`
	Context   []int                  `json:"context,omitempty"`
	Stream    *bool                  `json:"stream,omitempty"`
	Raw       bool                   `json:"raw,omitempty"`
	Format    interface{}            `json:"format,omitempty"`
	KeepAlive string                 `json:"keep_alive,omitempty"`
	Images    []string               `json:"images,omitempty"`
	Options   map[string]interface{} `json:"options,omitempty"`
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
	logToFloat(fmt.Sprintf("Request body: %s", body))

	urlBytes := []byte(url)
	methodBytes := []byte(method)
	bodyBytes := []byte(body)

	// Set headers for JSON content type
	headers := "Content-Type: application/json\r\n"
	headersBytes := []byte(headers)

	urlPtr := uintptr(unsafe.Pointer(&urlBytes[0]))
	methodPtr := uintptr(unsafe.Pointer(&methodBytes[0]))
	headersPtr := uintptr(unsafe.Pointer(&headersBytes[0]))
	bodyPtr := uintptr(unsafe.Pointer(&bodyBytes[0]))

	// Make the HTTP request to Float host
	result := floatHttpRequest(
		uint32(urlPtr), uint32(len(urlBytes)),
		uint32(methodPtr), uint32(len(methodBytes)),
		uint32(headersPtr), uint32(len(headersBytes)),
		uint32(bodyPtr), uint32(len(bodyBytes)),
	)

	// Check the result status
	if result != 0 {
		logToFloat(fmt.Sprintf("HTTP request failed with status: %d", result))
		return "", fmt.Errorf("HTTP request failed with status code: %d", result)
	}

	// Read the response from Float using file reading
	// Float typically writes HTTP responses to a temporary file
	responseBody, err := readHttpResponseFromFloat()
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	logToFloat(fmt.Sprintf("Received response (%d bytes)", len(responseBody)))

	return responseBody, nil
}

// Read HTTP response from Float using file reading mechanism
func readHttpResponseFromFloat() (string, error) {
	// Float typically writes HTTP responses to a temporary file that can be read
	responsePath := "/tmp/float_http_response"
	responsePathBytes := []byte(responsePath)
	responsePathPtr := uintptr(unsafe.Pointer(&responsePathBytes[0]))

	logToFloat("Reading HTTP response from Float temporary file")

	// Try to read the response file
	result := floatReadFile(uint32(responsePathPtr), uint32(len(responsePathBytes)))
	if result != 0 {
		return "", fmt.Errorf("failed to read HTTP response file: error code %d", result)
	}

	// Since floatReadFile doesn't return the content directly, we need to read it
	// from a known location or use another mechanism
	// For now, we'll implement a direct approach that works with Float's actual API

	logToFloat("Successfully read HTTP response from Float")

	// In a real Float implementation, this would read the actual response
	// For demonstration purposes, let's assume the response is made available
	// This needs to be replaced with actual Float response reading
	return readActualOllamaResponse()
}

func readActualOllamaResponse() (string, error) {
	// This represents the actual Ollama API response
	// In practice, Float would handle the HTTP request and provide the real response

	logToFloat("Processing real Ollama API response")

	// This is where the real Ollama response would be processed
	// For now, we'll return an error indicating this needs real implementation
	return "", fmt.Errorf("real Ollama API integration requires Float runtime HTTP response mechanism")
}

//export generate_ollama
func generate_ollama(inputPtr, inputLen uint32) uint32 {
	logToFloat("Starting Ollama generation request")

	// Read input data
	inputBytes := make([]byte, inputLen)
	copy(inputBytes, (*(*[1 << 30]byte)(unsafe.Pointer(uintptr(inputPtr))))[0:inputLen])

	var input OllamaInput
	if err := json.Unmarshal(inputBytes, &input); err != nil {
		logToFloat(fmt.Sprintf("Failed to parse input: %v", err))
		output := &OllamaOutput{
			Success:   false,
			Error:     fmt.Sprintf("Invalid input JSON: %v", err),
			ErrorType: "INPUT_PARSE_ERROR",
			Metadata: map[string]interface{}{
				"error_stage": "input_parsing",
			},
		}
		writeOutputFile(output)
		return 1
	}

	// Validate required fields
	if input.Model == "" {
		logToFloat("Model is required")
		output := &OllamaOutput{
			Success:   false,
			Error:     "Model field is required",
			ErrorType: "VALIDATION_ERROR",
			Metadata: map[string]interface{}{
				"error_stage":   "validation",
				"missing_field": "model",
			},
		}
		writeOutputFile(output)
		return 1
	}

	if input.Prompt == "" {
		logToFloat("Prompt is required")
		output := &OllamaOutput{
			Success:   false,
			Error:     "Prompt field is required",
			ErrorType: "VALIDATION_ERROR",
			Metadata: map[string]interface{}{
				"error_stage":   "validation",
				"missing_field": "prompt",
			},
		}
		writeOutputFile(output)
		return 1
	}

	if input.OllamaURL == "" {
		logToFloat("ERROR: ollama_url is required")
		output := &OllamaOutput{
			Success:   false,
			Error:     "ollama_url field is required - please specify the Ollama server URL (e.g., 'http://localhost:11434')",
			ErrorType: "VALIDATION_ERROR",
			Metadata: map[string]interface{}{
				"error_stage":   "validation",
				"missing_field": "ollama_url",
			},
		}
		writeOutputFile(output)
		return 1
	}

	// Ensure URL doesn't end with slash
	input.OllamaURL = strings.TrimRight(input.OllamaURL, "/")

	logToFloat(fmt.Sprintf("Generating with model: %s", input.Model))
	logToFloat(fmt.Sprintf("Prompt length: %d characters", len(input.Prompt)))

	// Prepare Ollama API request
	apiRequest := OllamaAPIRequest{
		Model:     input.Model,
		Prompt:    input.Prompt,
		Suffix:    input.Suffix,
		System:    input.System,
		Template:  input.Template,
		Context:   input.Context,
		Stream:    input.Stream,
		Raw:       input.Raw,
		Format:    input.Format,
		KeepAlive: input.KeepAlive,
		Images:    input.Images,
		Options:   input.Options,
	}

	// Marshal request body
	requestBody, err := json.Marshal(apiRequest)
	if err != nil {
		logToFloat(fmt.Sprintf("Failed to marshal request: %v", err))
		output := &OllamaOutput{
			Success:   false,
			Error:     fmt.Sprintf("Failed to prepare request: %v", err),
			ErrorType: "REQUEST_MARSHAL_ERROR",
			Metadata: map[string]interface{}{
				"error_stage": "request_preparation",
			},
		}
		writeOutputFile(output)
		return 1
	}

	// Make HTTP request to Ollama
	url := input.OllamaURL + "/api/generate"
	logToFloat(fmt.Sprintf("Making request to: %s", url))

	responseBody, err := makeHttpRequest(url, "POST", string(requestBody))
	if err != nil {
		logToFloat(fmt.Sprintf("HTTP request failed: %v", err))
		output := &OllamaOutput{
			Success:   false,
			Error:     fmt.Sprintf("Failed to connect to Ollama: %v", err),
			ErrorType: "HTTP_REQUEST_ERROR",
			Metadata: map[string]interface{}{
				"error_stage": "http_request",
				"ollama_url":  url,
			},
		}
		writeOutputFile(output)
		return 1
	}

	// Parse Ollama response
	var apiResponse OllamaAPIResponse
	if err := json.Unmarshal([]byte(responseBody), &apiResponse); err != nil {
		logToFloat(fmt.Sprintf("Failed to parse response: %v", err))
		output := &OllamaOutput{
			Success:   false,
			Error:     fmt.Sprintf("Invalid response from Ollama: %v", err),
			ErrorType: "RESPONSE_PARSE_ERROR",
			Metadata: map[string]interface{}{
				"error_stage":  "response_parsing",
				"raw_response": responseBody,
			},
		}
		writeOutputFile(output)
		return 1
	}

	// Build output
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
		Metadata: map[string]interface{}{
			"input_prompt_length": len(input.Prompt),
			"response_length":     len(apiResponse.Response),
			"ollama_url":          input.OllamaURL,
			"stream_mode":         input.Stream,
			"has_suffix":          input.Suffix != "",
			"has_system":          input.System != "",
			"has_images":          len(input.Images) > 0,
			"processing_complete": true,
		},
	}

	// Write output file
	if writeErr := writeOutputFile(output); writeErr != nil {
		logToFloat(fmt.Sprintf("Failed to write output: %v", writeErr))
		return 1
	}

	logToFloat("Ollama generation completed successfully")
	logToFloat(fmt.Sprintf("Generated response length: %d characters", len(apiResponse.Response)))

	return 0
}

func main() {
	// Required for TinyGo WASM modules
}
