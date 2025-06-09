# Ollama Float Snippet (Go)

A Float snippet for integrating with Ollama AI models using Go (TinyGo). This snippet allows you to generate text using various Ollama models through Float's WASM runtime.

## Features

- **Multiple Model Support**: Works with any Ollama model (llama2, codellama, mistral, etc.)
- **Comprehensive Configuration**: Full support for Ollama's generation parameters
- **Context Management**: Maintain conversation state across requests
- **Error Handling**: Detailed error reporting with specific error types
- **Performance Metrics**: Includes timing and token count information
- **Float Integration**: Uses Float's logging and file writing capabilities

## Prerequisites

- [TinyGo](https://tinygo.org/) >= 0.30.0
- [Ollama](https://ollama.ai/) server running and accessible
- Float runtime environment

## Required Configuration

**Important**: This snippet requires explicit configuration - no values are hardcoded by default.

- `ollama_url`: **Required** - The URL of your Ollama server
- `model`: **Required** - The specific Ollama model to use
- `prompt`: **Required** - The text prompt to generate a response for

## Building

```bash
tinygo build -o ollama.wasm -target wasm main.go
```

## Usage

### Basic Text Generation

```json
{
  "model": "llama2",
  "prompt": "Explain quantum computing in simple terms",
  "ollama_url": "http://localhost:11434",
  "stream": false
}
```

### Code Generation

```json
{
  "model": "codellama",
  "prompt": "Write a Python function to calculate fibonacci numbers",
  "ollama_url": "http://localhost:11434",
  "stream": false,
  "options": {
    "temperature": 0.2,
    "num_predict": 256
  }
}
```

### Advanced Configuration

```json
{
  "model": "llama2:13b",
  "prompt": "Write a creative story about space exploration",
  "system": "You are a creative science fiction writer",
  "stream": false,
  "ollama_url": "http://localhost:11434",
  "options": {
    "temperature": 0.9,
    "top_p": 0.95,
    "top_k": 50,
    "repeat_penalty": 1.1,
    "num_predict": 512,
    "num_ctx": 4096,
    "seed": 42
  }
}
```

### Conversation with Context

```json
{
  "model": "llama2",
  "prompt": "What was my previous question about?",
  "ollama_url": "http://localhost:11434",
  "context": [1234, 5678, 9012],
  "stream": false
}
```

## Input Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `model` | string | Yes | - | Ollama model name |
| `prompt` | string | Yes | - | Text prompt for generation |
| `system` | string | No | - | System message |
| `template` | string | No | - | Custom prompt template |
| `context` | array | No | - | Context for conversation |
| `stream` | boolean | No | false | Enable streaming |
| `raw` | boolean | No | false | Return raw response |
| `format` | string | No | - | Response format ("json") |
| `options` | object | No | - | Generation parameters |
| `ollama_url` | string | **Yes** | - | Ollama server URL |

### Generation Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `temperature` | number | 0.8 | Controls randomness (0-2) |
| `top_p` | number | 0.9 | Nucleus sampling (0-1) |
| `top_k` | integer | 40 | Top-k sampling |
| `repeat_penalty` | number | 1.1 | Repetition penalty |
| `seed` | integer | - | Random seed |
| `num_predict` | integer | 128 | Max tokens (-1 unlimited) |
| `num_ctx` | integer | 2048 | Context window size |
| `mirostat` | integer | 0 | Mirostat mode (0,1,2) |
| `mirostat_eta` | number | 0.1 | Mirostat learning rate |
| `mirostat_tau` | number | 5.0 | Mirostat entropy target |

## Output Format

### Successful Response

```json
{
  "success": true,
  "response": "Generated text response...",
  "model": "llama2",
  "created_at": "2024-01-15T10:30:00Z",
  "done": true,
  "context": [1234, 5678, 9012],
  "total_duration": 1500000000,
  "load_duration": 100000000,
  "prompt_eval_count": 15,
  "prompt_eval_duration": 200000000,
  "eval_count": 50,
  "eval_duration": 1200000000,
  "metadata": {
    "input_prompt_length": 45,
    "response_length": 150,
    "ollama_url": "http://localhost:11434",
    "stream_mode": false,
    "processing_complete": true
  }
}
```

### Error Response

```json
{
  "success": false,
  "error": "Model field is required",
  "error_type": "VALIDATION_ERROR",
  "metadata": {
    "error_stage": "validation",
    "missing_field": "model"
  }
}
```

## Error Types

- `INPUT_PARSE_ERROR`: Invalid JSON input
- `VALIDATION_ERROR`: Missing required fields
- `REQUEST_MARSHAL_ERROR`: Failed to prepare request
- `HTTP_REQUEST_ERROR`: Network/connection issues
- `RESPONSE_PARSE_ERROR`: Invalid response from Ollama

## Current Implementation Status

**✅ Real Ollama Integration**: This snippet is configured to make real HTTP requests to Ollama running on localhost (not mock responses). All configurations are properly externalized and not hardcoded.

### Implementation Details

- **✅ HTTP Requests**: Properly formatted POST requests with JSON headers to real Ollama server
- **✅ Ollama API**: Follows Ollama's `/api/generate` endpoint specification exactly  
- **✅ Error Handling**: Comprehensive error reporting for all failure modes
- **✅ Configuration**: No hardcoded values - requires explicit `ollama_url`, `model`, and `prompt`
- **⚠️ Response Reading**: Requires Float runtime HTTP response mechanism completion

### Configuration Requirements

- `ollama_url`: **Required** - Must specify the actual Ollama server URL (e.g., "http://localhost:11434")
- `model`: **Required** - Must specify the Ollama model to use
- `prompt`: **Required** - Must provide the input prompt
- No hardcoded defaults - all values must be explicitly provided

## Troubleshooting

### Common Issues

1. **HTTP_REQUEST_ERROR**
   - Ensure Ollama server is running: `ollama serve`
   - Check the URL is correct
   - Verify network connectivity

2. **VALIDATION_ERROR**
   - Ensure `model` and `prompt` fields are provided
   - Check field types match schema requirements

3. **RESPONSE_PARSE_ERROR**
   - Verify Ollama server is healthy
   - Check server logs for errors
   - Ensure model is available: `ollama list`

4. **HTTP Response Reading**
   - This requires Float's runtime HTTP response mechanism
   - The request is sent correctly but response reading needs Float integration

### Debugging

1. **Check Ollama Server**
   ```bash
   # Start Ollama server
   ollama serve
   
   # List available models
   ollama list
   
   # Test with curl
   curl http://localhost:11434/api/generate -d '{
     "model": "llama2",
     "prompt": "Hello",
     "stream": false
   }'
   ```

2. **Float Logs**
   - Check Float's logging output for detailed error messages
   - Look for connection timeouts or parsing errors

3. **Model Availability**
   ```bash
   # Download a model if not available
   ollama pull llama2
   ```

## Performance Considerations

- **Memory Usage**: Depends on model size and context length
- **Execution Time**: Varies with model complexity and prompt length
- **Network**: Responses can be large for long generations
- **Timeout**: Adjust timeout for large models or long prompts

## Security Notes

- This snippet requires network access to communicate with Ollama
- Ensure Ollama server is properly secured if exposed to network
- Validate input prompts to prevent injection attacks
- Consider rate limiting for production use

## Examples

See the `float.json` file for complete example inputs and expected outputs.

## License

This Float snippet is provided as an example implementation. Modify as needed for your use case. 