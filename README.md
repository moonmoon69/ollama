# Ollama AI Text Generator - Float Snippet

A robust, production-ready Float snippet for integrating with Ollama AI models using Go (TinyGo). This snippet provides comprehensive text generation capabilities with full error handling, validation, and metadata tracking.

## Overview

This snippet allows you to generate text using any Ollama model through a simple JSON interface. It includes comprehensive validation, detailed error reporting, and extensive metadata for monitoring and debugging.

## Features

- **Multiple Model Support**: Works with any Ollama model (llama2, codellama, mistral, phi3, gemma, etc.)
- **Comprehensive Configuration**: Full support for Ollama's generation parameters
- **Robust Error Handling**: Detailed error reporting with specific error types and stages
- **Input Validation**: Strict validation of all input parameters with helpful error messages
- **Conversation Context**: Support for maintaining conversation state across requests
- **Multimodal Support**: Base64 image input for vision-capable models
- **Performance Monitoring**: Detailed timing and token count metadata
- **Production Ready**: Built following Float's Snippet Development Guide best practices

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
  "system": "You are an expert programmer. Write clean, well-documented code.",
  "ollama_url": "http://localhost:11434",
  "stream": false,
  "options": {
    "temperature": 0.2,
    "top_p": 0.9,
    "top_k": 40,
    "num_predict": 256,
    "seed": 42
  }
}
```

### Conversation with Context

```json
{
  "model": "llama2",
  "prompt": "What was my previous question about?",
  "system": "You are a helpful AI assistant with excellent memory.",
  "context": [1, 2, 3, 4, 5],
  "ollama_url": "http://localhost:11434",
  "options": {
    "temperature": 0.7
  }
}
```

## Input Schema

### Required Fields

- **model** (string): The Ollama model to use (e.g., "llama2", "codellama", "mistral")
- **prompt** (string): The text prompt to generate a response for (1-32768 characters)
- **ollama_url** (string): Ollama server URL (must start with http:// or https://)

### Optional Fields

- **system** (string): System message to set model behavior (max 8192 characters)
- **template** (string): Custom prompt template
- **context** (array): Context from previous response for conversation continuity
- **stream** (boolean): Whether to stream response (default: false)
- **raw** (boolean): Return raw response without formatting (default: false)
- **format** (string|object): Response format ("json" or JSON schema object)
- **suffix** (string): Text after model response (for code completion)
- **keep_alive** (string): How long to keep model loaded (e.g., "5m", "10s", "1h")
- **images** (array): Base64-encoded images for multimodal models
- **options** (object): Model-specific generation options

### Generation Options

- **temperature** (number): Controls randomness (0-2, default: 0.8)
- **top_p** (number): Nucleus sampling parameter (0-1, default: 0.9)
- **top_k** (integer): Top-k sampling parameter (default: 40)
- **repeat_penalty** (number): Penalty for repeating tokens (default: 1.1)
- **seed** (integer): Random seed for reproducible generation
- **num_predict** (integer): Maximum tokens to generate (-1 for unlimited, default: 128)
- **num_ctx** (integer): Context window size (default: 2048)
- **mirostat** (integer): Mirostat sampling mode (0, 1, or 2)
- **mirostat_eta** (number): Mirostat learning rate (default: 0.1)
- **mirostat_tau** (number): Mirostat target entropy (default: 5.0)

## Output Schema

### Success Response

```json
{
  "success": true,
  "response": "Generated text response from Ollama",
  "model": "llama2",
  "created_at": "2024-01-15T10:30:00Z",
  "done": true,
  "context": [1, 2, 3, 4, 5, 6, 7, 8],
  "total_duration": 1000000000,
  "load_duration": 100000000,
  "prompt_eval_count": 10,
  "prompt_eval_duration": 200000000,
  "eval_count": 25,
  "eval_duration": 700000000,
  "metadata": {
    "input_prompt_length": 45,
    "response_length": 128,
    "ollama_url": "http://localhost:11434",
    "stream_mode": false,
    "processing_complete": true,
    "has_system_message": true,
    "has_context": true,
    "context_length": 5
  }
}
```

### Error Response

```json
{
  "success": false,
  "error": "Model field is required",
  "error_type": "VALIDATION_ERROR",
  "done": true,
  "metadata": {
    "error_stage": "validation",
    "missing_field": "model",
    "processing_complete": false
  }
}
```

## Error Types

- **INPUT_PARSE_ERROR**: Invalid input JSON format
- **VALIDATION_ERROR**: Input validation failed (missing required fields, invalid values)
- **REQUEST_MARSHAL_ERROR**: Failed to prepare request for Ollama
- **HTTP_REQUEST_ERROR**: Failed to connect to Ollama server
- **RESPONSE_PARSE_ERROR**: Unable to parse Ollama server response

## Build Instructions

### Prerequisites

- TinyGo >= 0.30.0
- Ollama server running and accessible

### Building the WASM Module

```bash
tinygo build -o ollama.wasm -target wasm main.go
```

### Local Testing

```bash
# Start Ollama server (if not already running)
ollama serve

# Test with the provided test input
cat test_input.json | your-float-runner ollama.wasm
```

## Performance Characteristics

- **Memory Usage**: Low to Medium (depends on model size and context)
- **Execution Time**: Variable (depends on model size, prompt complexity, and server performance)
- **Network Bandwidth**: Medium (responses can be large for long generations)

## Troubleshooting

### Common Issues

#### HTTP_REQUEST_ERROR
- Ensure Ollama server is running and accessible at the specified URL
- Check network connectivity and firewall settings
- Verify the ollama_url format is correct (must start with http:// or https://)

#### VALIDATION_ERROR
- Check that required fields (model, prompt, ollama_url) are provided
- Ensure prompt length doesn't exceed 32768 characters
- Verify system message doesn't exceed 8192 characters
- Check that temperature is between 0 and 2, top_p is between 0 and 1

#### RESPONSE_PARSE_ERROR
- Verify Ollama server is responding with valid JSON
- Check Ollama server logs for errors
- Ensure the specified model is available and loaded

### Debugging

1. **Enable Logging**: Check Float logs for detailed error messages and request/response information
2. **Verify Connectivity**: Test Ollama server connectivity separately using curl:
   ```bash
   curl -X POST http://localhost:11434/api/generate \
     -H "Content-Type: application/json" \
     -d '{"model":"llama2","prompt":"Hello"}'
   ```
3. **Check Model**: Ensure the specified model is available:
   ```bash
   ollama list
   ```
4. **Validate Input**: Use the schema to validate input before making requests

## Security Considerations

- **Network Access**: This snippet requires network access to communicate with Ollama servers
- **File Write**: Writes output.json with generation results
- **Input Validation**: All inputs are validated before processing
- **URL Validation**: Ollama URLs must be valid HTTP/HTTPS endpoints

## Version History

- **0.1.0**: Initial production-ready release
  - Comprehensive error handling and validation
  - Full Ollama API support
  - Detailed metadata and performance tracking
  - Robust logging and debugging capabilities

## License

MIT License - see LICENSE file for details.

## Contributing

This snippet follows the Float Snippet Development Guide standards. When contributing:

1. Maintain comprehensive error handling
2. Follow TinyGo patterns for WASM modules
3. Update schemas when adding new features
4. Include appropriate tests and documentation
5. Ensure all changes are production-ready 