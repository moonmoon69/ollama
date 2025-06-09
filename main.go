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
func floatHttpRequest(urlPtr, urlLen, methodPtr, methodLen, bodyPtr, bodyLen uint32) uint32

//go:wasmimport env float_write_file
func floatWriteFile(pathPtr, pathLen, dataPtr, dataLen uint32) uint32

// Input structure for Ollama request
type OllamaInput struct {
	Model     string                 `json:"model"`
	Prompt    string                 `json:"prompt"`
	System    string                 `json:"system,omitempty"`
	Template  string                 `json:"template,omitempty"`
	Context   []int                  `json:"context,omitempty"`
	Stream    bool                   `json:"stream"`
	Raw       bool                   `json:"raw,omitempty"`
	Format    string                 `json:"format,omitempty"`
	Options   map[string]interface{} `json:"options,omitempty"`
	OllamaURL string                 `json:"ollama_url"`
	Timeout   int                    `json:"timeout,omitempty"`
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
	Model    string                 `json:"model"`
	Prompt   string                 `json:"prompt"`
	System   string                 `json:"system,omitempty"`
	Template string                 `json:"template,omitempty"`
	Context  []int                  `json:"context,omitempty"`
	Stream   bool                   `json:"stream"`
	Raw      bool                   `json:"raw,omitempty"`
	Format   string                 `json:"format,omitempty"`
	Options  map[string]interface{} `json:"options,omitempty"`
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
	urlBytes := []byte(url)
	methodBytes := []byte(method)
	bodyBytes := []byte(body)

	urlPtr := uintptr(unsafe.Pointer(&urlBytes[0]))
	methodPtr := uintptr(unsafe.Pointer(&methodBytes[0]))
	bodyPtr := uintptr(unsafe.Pointer(&bodyBytes[0]))

	// Note: This is a simplified HTTP request implementation
	// In a real Float environment, this would return response data
	result := floatHttpRequest(
		uint32(urlPtr), uint32(len(urlBytes)),
		uint32(methodPtr), uint32(len(methodBytes)),
		uint32(bodyPtr), uint32(len(bodyBytes)),
	)

	if result != 0 {
		return "", fmt.Errorf("HTTP request failed with code: %d", result)
	}

	// In a real implementation, Float would provide the response data
	// For this example, we'll simulate a successful response
	return `{"response": "Generated response from Ollama", "done": true, "model": "llama2"}`, nil
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
		logToFloat("Using default Ollama URL")
		input.OllamaURL = "http://localhost:11434"
	}

	// Ensure URL doesn't end with slash
	input.OllamaURL = strings.TrimRight(input.OllamaURL, "/")

	logToFloat(fmt.Sprintf("Generating with model: %s", input.Model))
	logToFloat(fmt.Sprintf("Prompt length: %d characters", len(input.Prompt)))

	// Prepare Ollama API request
	apiRequest := OllamaAPIRequest{
		Model:    input.Model,
		Prompt:   input.Prompt,
		System:   input.System,
		Template: input.Template,
		Context:  input.Context,
		Stream:   input.Stream,
		Raw:      input.Raw,
		Format:   input.Format,
		Options:  input.Options,
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
