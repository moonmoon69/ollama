# Ollama AI Text Generator - Float Reagent

A Float WASM reagent for integrating with Ollama AI models to generate text responses with comprehensive configuration options.

## Overview

This reagent provides a seamless interface to Ollama AI models, allowing you to generate text responses with fine-grained control over model parameters, system prompts, and generation options. Built using Go (TinyGo) for optimal WASM performance.

## Features

- **Multiple Model Support**: Compatible with all Ollama models (llama2, codellama, mistral, phi3, gemma, etc.)
- **Advanced Configuration**: Full control over temperature, top-p, top-k, and other generation parameters  
- **System Prompts**: Set model behavior with custom system messages
- **Context Preservation**: Maintain conversation state across requests
- **Multimodal Support**: Process images with compatible models (e.g., llava)
- **Error Handling**: Comprehensive error reporting with detailed metadata
- **Validation**: Input validation according to Ollama API specifications

## Usage

### Basic Text Generation

```json
{
  "model": "llama2",
  "prompt": "Explain quantum computing in simple terms",
  "ollama_url": "http://localhost:11434"
}
```

### Advanced Configuration

```json
{
  "model": "codellama",
  "prompt": "Write a Python function to calculate fibonacci numbers",
  "system": "You are a senior Python developer. Write clean, well-documented code with examples.",
  "ollama_url": "http://localhost:11434",
  "options": {
    "temperature": 0.7,
    "top_p": 0.9,
    "top_k": 40,
    "num_predict": 200,
    "seed": 42
  },
  "keep_alive": "10m"
}
```

### Multimodal Request (with images)

```json
{
  "model": "llava",
  "prompt": "Describe what you see in this image",
  "images": ["base64_encoded_image_data"],
  "ollama_url": "http://localhost:11434"
}
```

## Configuration Options

### Required Fields

- `model` (string): The Ollama model to use for generation
- `prompt` (string): The text prompt to generate a response for
- `ollama_url` (string): URL of the Ollama server (e.g., "http://localhost:11434")

### Optional Fields

- `system` (string): System message to set model behavior (max 8192 chars)
- `template` (string): Custom prompt template 
- `context` (array): Context from previous response for conversation continuity
- `stream` (boolean): Enable streaming response (default: false)
- `raw` (boolean): Return raw response without formatting (default: false)
- `format` (string|object): Response format specification ("json" or JSON schema)
- `suffix` (string): Text after the model response (for code completion)
- `keep_alive` (string): How long to keep model loaded ("5m", "10s", "1h")
- `images` (array): Base64-encoded images for multimodal models

### Generation Options

The `options` object supports all Ollama parameters:

- `temperature` (0-2): Controls randomness (default: 0.8)
- `top_p` (0-1): Nucleus sampling parameter (default: 0.9) 
- `top_k` (integer): Top-k sampling parameter (default: 40)
- `repeat_penalty` (number): Penalty for repeating tokens (default: 1.1)
- `seed` (integer): Random seed for reproducible generation
- `num_predict` (integer): Max tokens to generate, -1 for unlimited (default: 128)
- `num_ctx` (integer): Context window size (default: 2048)
- `mirostat` (0|1|2): Mirostat sampling mode (default: 0)
- `mirostat_eta` (number): Mirostat learning rate (default: 0.1)
- `mirostat_tau` (number): Mirostat target entropy (default: 5.0)

## Response Format

### Successful Response

```json
{
  "success": true,
  "response": "Generated text response from the AI model...",
  "model": "llama2",
  "created_at": "2024-12-16T10:30:00Z",
  "done": true,
  "context": [1, 2, 3, ...],
  "total_duration": 1000000000,
  "load_duration": 100000000,
  "prompt_eval_count": 15,
  "prompt_eval_duration": 200000000,
  "eval_count": 25,
  "eval_duration": 700000000,
  "metadata": {
    "input_prompt_length": 45,
    "response_length": 120,
    "processing_complete": true,
    "go_version": "tinygo",
    "timestamp": "2024-12-16T10:30:00Z"
  }
}
```

### Error Response

```json
{
  "success": false,
  "done": true,
  "error": "model field is required",
  "error_type": "VALIDATION_ERROR",
  "metadata": {
    "processing_complete": false,
    "error_stage": "validation",
    "error_timestamp": "2024-12-16T10:30:00Z"
  }
}
```

## Error Types

- `INPUT_PARSE_ERROR`: Invalid JSON input format
- `VALIDATION_ERROR`: Missing required fields or invalid values
- `REQUEST_MARSHAL_ERROR`: Failed to prepare request for Ollama
- `HTTP_REQUEST_ERROR`: Failed to connect to Ollama server
- `RESPONSE_PARSE_ERROR`: Invalid response from Ollama server

## Building

This reagent is built using TinyGo for WASM:

```bash
tinygo build -o ollama.wasm -target wasm main.go
```

## Requirements

- **TinyGo**: >= 0.30.0
- **Ollama Server**: Running and accessible at the specified URL
- **Float Runtime**: For WASM execution with host function support

## Float Host Functions Used

This reagent utilizes Float's host functions for:

- **Logging**: `float.log()` for debug and status messages
- **HTTP Requests**: `float.http_request()` for Ollama API calls  
- **File Operations**: `float.write_file()` for output file generation

## Development

### Testing

Use the provided `test_input.json` for testing:

```bash
# Test input example
{
  "model": "llama2",
  "prompt": "Explain machine learning in simple terms",
  "system": "You are a helpful AI assistant",
  "ollama_url": "http://localhost:11434",
  "options": {
    "temperature": 0.7,
    "num_predict": 200
  }
}
```

### Adding New Features

1. Update the input/output structs in `main.go`
2. Modify the `float.json` schema accordingly
3. Implement validation in `validateInput()`
4. Test with various input configurations

## License

This Float reagent is part of the Float platform ecosystem. See Float documentation for license details.

## Contributing

Please follow the [Float Reagent Development Guide](Reagent_Development_Guide.md) when contributing to this reagent. 