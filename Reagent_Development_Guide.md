# Float Reagent Development Guide

**Version**: 1.0.11  
**Date**: December 16, 2024  
**Audience**: Developers, Alchemist AI, Component Developers

## Table of Contents

1. [Overview](#overview)
2. [Float.json Configuration Guide](#floatjson-configuration-guide)
3. [Schema Definition Guidelines](#schema-definition-guidelines)
4. [Best Practices](#best-practices)
5. [Testing Your Reagents](#testing-your-reagents)
6. [Language-Specific Examples](#language-specific-examples)
   - [Rust Reagent Examples](#rust-reagent-examples)
   - [Go (TinyGo) Reagent Examples](#go-tinygo-reagent-examples)
   - [Python (Pyodide) Reagent Examples](#python-pyodide-reagent-examples)
   - [JavaScript Reagent Examples](#javascript-reagent-examples)
   - [TypeScript Reagent Examples](#typescript-reagent-examples)

## Overview

This guide provides comprehensive guidance for developing WASM Reagents in all languages supported by Float. The guide covers:

- Configuration and schema design principles
- Best practices for development and testing
- Complete implementation examples for each supported language
- Integration patterns with Float's host functions

### Float Host Functions Available to All Languages

All WASM snippets have access to these Float host functions:

- `float.log(message)` - Logging from WASM modules
- `float.http_request(params)` - Controlled HTTP client access
- `float.read_file(path)` - File reading with path validation  
- `float.write_file(path, data)` - Atomic file writing operations

> **Note**: For complete technical specifications of these host functions, see [Data Passing Implementation Guide](Data_Passing_Implementation_Guide.md).

## Float.json Configuration Guide

The `float.json` file is the cornerstone of every Float snippet, serving as both a manifest and API contract. This section provides a comprehensive guide to constructing and understanding this crucial configuration file.

### Overview

The `float.json` file defines:
- **Snippet metadata** (ID, version, description)
- **Entry point descriptions** (what operations users can call and what each does)
- **Input/output schemas** (API contracts for each entry point)
- **Runtime configuration** (language, dependencies, build commands)

> **Important**: Entry point values should describe **what the operation does** when called, not internal function names. Think of them as user-facing operation descriptions.

### Core Structure

```json
{
  "snippet_id": "my-snippet",           // Unique identifier
  "version": "1.0.0",                   // Semantic version
  "name": "My Snippet",                 // Human-readable name
  "description": "What this snippet does",
  "language": "rust",                   // Programming language
  "entry_points": {                     // Callable operations with descriptions
    "main": "Primary data processing operation",
    "validate": "Validates input data before processing"
  },
  "schemas": {                          // API contracts
    "main": { "input": {...}, "output": {...} },
    "validate": { "input": {...}, "output": {...} }
  }
}
```

### Entry Points Explained

Entry points define **what operations can be called** from your snippet and **what each operation does**:

#### Simple Entry Point (Single Operation)
```json
{
  "entry_points": {
    "main": "Performs arithmetic calculations on input numbers"
  },
  "schemas": {
    "main": {
      "input": { "type": "object", "properties": {...} },
      "output": { "type": "object", "properties": {...} }
    }
  }
}
```

#### Multiple Entry Points (Multi-Operation)
```json
{
  "entry_points": {
    "main": "Processes and analyzes the input data",
    "validate": "Validates input data format and constraints", 
    "status": "Returns current processing status and statistics",
    "cleanup": "Cleans up temporary files and resources"
  },
  "schemas": {
    "main": {
      "input": { /* structured data schema */ },
      "output": { /* processed results schema */ }
    },
    "validate": {
      "input": { /* same as main, but for validation */ },
      "output": { /* validation results schema */ }
    },
    "status": {
      "input": {},                      // No input required
      "output": { /* status information schema */ }
    },
    "cleanup": {
      "input": { /* optional cleanup parameters */ },
      "output": { /* cleanup confirmation schema */ }
    }
  }
}
```

### Language-Agnostic Entry Point Patterns

The `float.json` entry points remain consistent across all languages - they describe **what operations are available** to users, regardless of the underlying implementation language.

#### For Calculator Snippet
```json
{
  "entry_points": {
    "main": "Performs arithmetic calculations and returns formatted results",
    "validate": "Validates mathematical expressions before calculation"
  }
}
```

#### For Data Analysis Snippet  
```json
{
  "entry_points": {
    "main": "Runs complete statistical analysis on the dataset",
    "analyze": "Performs comprehensive data analysis with multiple operations", 
    "stats": "Calculates basic statistics (mean, median, mode) quickly"
  }
}
```

#### For Data Transformation Snippet
```json
{
  "entry_points": {
    "main": "Transforms JSON data according to specified operations",
    "transform": "Applies data transformations like key conversion and filtering",
    "validate": "Validates JSON structure and data types before processing"
  }
}
```

> **Key Principle**: The `float.json` describes the **user interface** of your snippet - what operations are available and what they do. The internal implementation (function names, file structure) is abstracted away and handled by Float's runtime based on the programming language.

### Schema Design Patterns

#### Pattern 1: Structured Input/Output
For complex data processing:
```json
{
  "schemas": {
    "main": {
      "input": {
        "type": "object",
        "properties": {
          "data": { "type": "array", "items": { "type": "number" } },
          "operations": { "type": "array", "items": { "type": "string" } },
          "options": { "type": "object" }
        },
        "required": ["data"]
      },
      "output": {
        "type": "object",
        "properties": {
          "success": { "type": "boolean" },
          "results": { "type": "object" },
          "metadata": { "type": "object" }
        },
        "required": ["success"]
      }
    }
  }
}
```

#### Pattern 2: String-Based Input/Output
For simple or raw data processing:
```json
{
  "schemas": {
    "process": {
      "input": {
        "type": "string",
        "description": "Raw JSON string containing processing instructions"
      },
      "output": {
        "type": "string",
        "description": "JSON string containing processing results"
      }
    }
  }
}
```

#### Pattern 3: Mixed Entry Points
Different entry points with different schemas:
```json
{
  "schemas": {
    "main": {
      "input": {
        "type": "object",
        "properties": {
          "numbers": { "type": "array", "items": { "type": "number" } },
          "operation": { "type": "string", "enum": ["add", "multiply"] }
        }
      },
      "output": {
        "type": "object",
        "properties": {
          "result": { "type": "number" },
          "operation": { "type": "string" }
        }
      }
    },
    "validate": {
      "input": {
        "type": "object", 
        "properties": {
          "numbers": { "type": "array" },
          "operation": { "type": "string" }
        }
      },
      "output": {
        "type": "object",
        "properties": {
          "valid": { "type": "boolean" },
          "errors": { "type": "array", "items": { "type": "string" } }
        }
      }
    },
    "info": {
      "input": {},                      // No input
      "output": {
        "type": "object",
        "properties": {
          "name": { "type": "string" },
          "version": { "type": "string" },
          "capabilities": { "type": "array", "items": { "type": "string" } }
        }
      }
    }
  }
}
```

### Validation and Best Practices

#### Required Fields Checklist
- ✅ `snippet_id`: Unique, descriptive identifier
- ✅ `version`: Valid semantic version
- ✅ `name`: Clear, human-readable name  
- ✅ `language`: Supported language (rust|go|python|javascript|typescript)
- ✅ `entry_points.main`: Default entry point defined
- ✅ `schemas.main`: Schema for main entry point

#### Optional but Recommended
- ✅ `description`: Clear purpose description
- ✅ `author`: Creator identification
- ✅ Additional entry points for complex snippets
- ✅ `dependencies`: Required libraries/packages
- ✅ `build_command`: Custom build instructions

#### Schema Design Guidelines
1. **Be Specific**: Use precise types and constraints
2. **Be Descriptive**: Include helpful descriptions
3. **Be Consistent**: Use consistent naming patterns
4. **Be Flexible**: Allow optional parameters where appropriate
5. **Be Validatable**: Ensure schemas can be automatically validated

This configuration system enables Float to:
- **Route calls** to the correct functions
- **Validate inputs** before execution
- **Validate outputs** after execution  
- **Generate documentation** automatically
- **Provide IDE support** through schema definitions

## Schema Definition Guidelines

### Schema File Structure

Every Float snippet should include a `float.json` file with the following structure:

```json
{
  "snippet_id": "unique-identifier",              // REQUIRED: Universal unique identifier
  "version": "semver-version",                    // REQUIRED: Semantic version
  "name": "Human-readable name",                  // REQUIRED: Display name
  "description": "Human-readable description",   // OPTIONAL: Purpose description
  "author": "Author Name",                        // OPTIONAL: Creator information
  "language": "rust|go|python|javascript|typescript", // REQUIRED: Programming language
  "runtime": "native|pyodide|quickjs",          // OPTIONAL: Runtime environment (defaults based on language)
  "entry_points": {                              // REQUIRED: Entry point definitions
    "main": "Description of what main operation does",          // REQUIRED: Default entry point
    "method_name": "Description of what this operation does"    // OPTIONAL: Additional callable methods
  },
  "build_command": "optional build command",     // OPTIONAL: Custom build command
  "dependencies": {                              // OPTIONAL: Language-specific dependencies
    "language_specific": "dependencies"
  },
  "schemas": {                                   // REQUIRED: Entry point schemas
    "main": {                                  // REQUIRED: Schema for main entry point
      "input": { /* JSON Schema for input */ },
      "output": { /* JSON Schema for output */ }
    },
    "method_name": {                           // OPTIONAL: Schema for additional entry points
      "input": { /* JSON Schema for input */ },
      "output": { /* JSON Schema for output */ }
    }
  }
}
```

**Required vs Optional Fields:**
- **REQUIRED**: `snippet_id`, `version`, `name`, `language`, `entry_points.main`, `schemas.main`
- **OPTIONAL**: `description`, `author`, `runtime`, `build_command`, `dependencies`, additional entry points, additional schemas

**Per-Entry-Point Schemas:**
Each entry point can have its own input/output schema, allowing for different function signatures:
- `schemas.main` is required and defines the schema for the default entry point
- `schemas.{method_name}` defines schemas for additional entry points
- This enables snippets with multiple methods that have different input/output requirements

### Best Practices for Schema Design

1. **Use Descriptive Field Names**
   ```json
   {
     "user_profile": {
       "type": "object",
       "properties": {
         "full_name": { "type": "string" },
         "email_address": { "type": "string" },
         "age_in_years": { "type": "integer" }
       }
     }
   }
   ```

2. **Include Validation Constraints**
   ```json
   {
     "email": {
       "type": "string",
       "format": "email",
       "pattern": "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"
     },
     "age": {
       "type": "integer",
       "minimum": 0,
       "maximum": 150
     }
   }
   ```

3. **Provide Clear Descriptions**
   ```json
   {
     "processing_options": {
       "type": "object",
       "description": "Configuration options for data processing",
       "properties": {
         "batch_size": {
           "type": "integer",
           "description": "Number of records to process in each batch",
           "minimum": 1,
           "maximum": 1000,
           "default": 100
         }
       }
     }
   }
   ```

## Best Practices

### Error Handling

1. **Always Use Try-Catch Blocks**
   ```rust
   // Rust
   match parse_input(&input) {
       Ok(data) => process_data(data),
       Err(e) => return Err(format!("Input parsing failed: {}", e))
   }
   ```

   ```python
   # Python
   try:
       result = process_data(input_data)
   except ValueError as e:
       return {"error": f"Invalid input: {str(e)}"}
   ```

2. **Provide Meaningful Error Messages**
   ```javascript
   // JavaScript
   if (!Array.isArray(data)) {
       throw new Error(`Expected array for data processing, got ${typeof data}`);
   }
   ```

### Logging

1. **Use Float's Logging Function**
   ```rust
   log_to_float("Starting data processing");
   log_to_float(&format!("Processing {} items", items.len()));
   ```

2. **Log Key Operations and Progress**
   ```python
   float_log(f"Loaded {len(dataset)} records")
   float_log(f"Applying {len(operations)} transformations")
   float_log("Processing completed successfully")
   ```

### Performance

1. **Minimize Memory Allocations**
   ```rust
   // Reuse buffers when possible
   let mut buffer = Vec::with_capacity(1000);
   for item in items {
       buffer.clear();
       process_item(item, &mut buffer);
   }
   ```

2. **Stream Large Data Sets**
   ```python
   # Process data in chunks
   def process_large_dataset(data, chunk_size=1000):
       for i in range(0, len(data), chunk_size):
           chunk = data[i:i + chunk_size]
           yield process_chunk(chunk)
   ```

### Security

1. **Validate All Inputs**
   ```typescript
   function validateInput(data: any): ValidationResult {
       if (typeof data !== 'object') {
           throw new Error('Input must be an object');
       }
       // Additional validation logic
   }
   ```

2. **Sanitize String Inputs**
   ```python
   import re
   
   def sanitize_string(input_str: str) -> str:
       # Remove potentially dangerous characters
       return re.sub(r'[<>"\'']', '', input_str.strip())
   ```

## Testing Your Snippets

### Unit Testing

Create test files alongside your snippets:

**File: `test_snippet.py` (for Python snippets)**
```python
import json
import unittest
from main import process_data_analysis

class TestDataAnalysis(unittest.TestCase):
    
    def test_basic_stats(self):
        input_data = {
            "dataset": [1, 2, 3, 4, 5],
            "operations": ["basic_stats"]
        }
        
        result = json.loads(process_data_analysis(json.dumps(input_data)))
        
        self.assertTrue(result['success'])
        self.assertEqual(result['results']['basic_stats']['mean'], 3.0)
        self.assertEqual(result['results']['basic_stats']['count'], 5)
    
    def test_empty_dataset(self):
        input_data = {
            "dataset": [],
            "operations": ["basic_stats"]
        }
        
        result = json.loads(process_data_analysis(json.dumps(input_data)))
        
        self.assertFalse(result['success'])
        self.assertIn('error', result)

if __name__ == '__main__':
    unittest.main()
```

### Integration Testing

**File: `integration_test.json`**
```json
{
  "test_cases": [
    {
      "name": "valid_input_test",
      "input": {
        "data": {"name": "John", "age": 30},
        "rules": {
          "name": {"required": true, "min_length": 2},
          "age": {"required": true}
        }
      },
      "expected_output": {
        "valid": true,
        "errors": []
      }
    },
    {
      "name": "validation_error_test", 
      "input": {
        "data": {"name": "J"},
        "rules": {
          "name": {"required": true, "min_length": 2}
        }
      },
      "expected_output": {
        "valid": false,
        "errors": [
          {
            "code": "MIN_LENGTH_VIOLATION"
          }
        ]
      }
    }
  ]
}
```

### Performance Testing

**File: `benchmark.py`**
```python
import time
import json
from main import process_data_analysis

def benchmark_processing():
    """Benchmark data analysis performance"""
    
    # Generate test dataset
    large_dataset = list(range(10000))
    
    input_data = {
        "dataset": large_dataset,
        "operations": ["basic_stats", "descriptive_stats"]
    }
    
    start_time = time.time()
    result = process_data_analysis(json.dumps(input_data))
    end_time = time.time()
    
    processing_time = end_time - start_time
    
    print(f"Processed {len(large_dataset)} items in {processing_time:.2f} seconds")
    print(f"Processing rate: {len(large_dataset) / processing_time:.0f} items/second")
    
    return processing_time < 5.0  # Should complete within 5 seconds

if __name__ == "__main__":
    success = benchmark_processing()
    print(f"Performance test: {'PASSED' if success else 'FAILED'}")
```

## Language-Specific Examples

This section provides complete implementation examples for each supported language. All examples include proper error handling, Float host function integration, and corresponding `float.json` configurations.

> **Note**: The `float.json` configurations remain consistent across languages, focusing on **what operations are available** rather than implementation details.
## Rust Snippet Examples

### Example 1: Data Validator Snippet

**File: `src/lib.rs`**
```rust
use wasm_bindgen::prelude::*;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

// Import Float host functions
extern "C" {
    fn float_log(ptr: *const u8, len: usize);
    fn float_write_file(path_ptr: *const u8, path_len: usize, data_ptr: *const u8, data_len: usize) -> i32;
}

#[derive(Deserialize)]
struct ValidationInput {
    data: serde_json::Value,
    rules: HashMap<String, ValidationRule>,
}

#[derive(Deserialize)]
struct ValidationRule {
    required: Option<bool>,
    min_length: Option<usize>,
    max_length: Option<usize>,
    pattern: Option<String>,
}

#[derive(Serialize)]
struct ValidationOutput {
    valid: bool,
    errors: Vec<ValidationError>,
    validated_data: serde_json::Value,
}

#[derive(Serialize)]
struct ValidationError {
    field: String,
    message: String,
    code: String,
}

// Helper function to call Float's log function
fn log_to_float(message: &str) {
    let bytes = message.as_bytes();
    unsafe {
        float_log(bytes.as_ptr(), bytes.len());
    }
}

// Helper function to write output file
fn write_output_file(data: &ValidationOutput) -> Result<(), String> {
    let json_output = serde_json::to_string_pretty(data)
        .map_err(|e| format!("Failed to serialize output: {}", e))?;
    
    let path = "output.json";
    let result = unsafe {
        float_write_file(
            path.as_ptr(),
            path.len(),
            json_output.as_ptr(),
            json_output.len()
        )
    };
    
    if result == 0 {
        Ok(())
    } else {
        Err("Failed to write output file".to_string())
    }
}

#[wasm_bindgen]
pub fn validate_data(input_json: &str) -> Result<(), JsValue> {
    log_to_float("Starting data validation");
    
    // Parse input
    let input: ValidationInput = serde_json::from_str(input_json)
        .map_err(|e| JsValue::from_str(&format!("Invalid input JSON: {}", e)))?;
    
    let mut errors = Vec::new();
    let mut validated_data = input.data.clone();
    
    // Apply validation rules
    if let serde_json::Value::Object(ref data_map) = input.data {
        for (field_name, rule) in &input.rules {
            if let Some(field_value) = data_map.get(field_name) {
                validate_field(field_name, field_value, rule, &mut errors);
            } else if rule.required.unwrap_or(false) {
                errors.push(ValidationError {
                    field: field_name.clone(),
                    message: format!("Field '{}' is required", field_name),
                    code: "REQUIRED_FIELD_MISSING".to_string(),
                });
            }
        }
    }
    
    let output = ValidationOutput {
        valid: errors.is_empty(),
        errors,
        validated_data,
    };
    
    // Write output
    write_output_file(&output)
        .map_err(|e| JsValue::from_str(&e))?;
    
    log_to_float(&format!("Validation complete. Valid: {}", output.valid));
    Ok(())
}

fn validate_field(
    field_name: &str,
    field_value: &serde_json::Value,
    rule: &ValidationRule,
    errors: &mut Vec<ValidationError>,
) {
    if let serde_json::Value::String(ref s) = field_value {
        // Check minimum length
        if let Some(min_len) = rule.min_length {
            if s.len() < min_len {
                errors.push(ValidationError {
                    field: field_name.to_string(),
                    message: format!("Field '{}' must be at least {} characters", field_name, min_len),
                    code: "MIN_LENGTH_VIOLATION".to_string(),
                });
            }
        }
        
        // Check maximum length
        if let Some(max_len) = rule.max_length {
            if s.len() > max_len {
                errors.push(ValidationError {
                    field: field_name.to_string(),
                    message: format!("Field '{}' must be at most {} characters", field_name, max_len),
                    code: "MAX_LENGTH_VIOLATION".to_string(),
                });
            }
        }
        
        // Check pattern (basic regex)
        if let Some(ref pattern) = rule.pattern {
            if pattern == "email" && !s.contains('@') {
                errors.push(ValidationError {
                    field: field_name.to_string(),
                    message: format!("Field '{}' must be a valid email", field_name),
                    code: "INVALID_EMAIL".to_string(),
                });
            }
        }
    }
}
```

**File: `Cargo.toml`**
```toml
[package]
name = "data-validator"
version = "0.1.0"
edition = "2021"

[lib]
crate-type = ["cdylib"]

[dependencies]
wasm-bindgen = "0.2"
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"

[dependencies.web-sys]
version = "0.3"
features = [
  "console",
]
```

**File: `float.json` (Schema Definition)**
```json
{
  "snippet_id": "data-validator",
  "version": "0.1.0",
  "name": "Data Validator",
  "description": "Validates JSON data against configurable rules",
  "language": "rust",
  "entry_points": {
    "main": "Validates JSON data against configurable validation rules"
  },
  "schemas": {
    "main": {
      "input": {
        "type": "object",
        "properties": {
          "data": {
            "type": "object",
            "description": "The data to validate"
          },
          "rules": {
            "type": "object",
            "description": "Validation rules for each field",
            "additionalProperties": {
              "type": "object",
              "properties": {
                "required": { "type": "boolean" },
                "min_length": { "type": "integer", "minimum": 0 },
                "max_length": { "type": "integer", "minimum": 0 },
                "pattern": { "type": "string" }
              }
            }
          }
        },
        "required": ["data", "rules"]
      },
      "output": {
        "type": "object",
        "properties": {
          "valid": { "type": "boolean" },
          "errors": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "field": { "type": "string" },
                "message": { "type": "string" },
                "code": { "type": "string" }
              },
              "required": ["field", "message", "code"]
            }
          },
          "validated_data": { "type": "object" }
        },
        "required": ["valid", "errors", "validated_data"]
      }
    }
  }
}
```

**Build Command:**
```bash
wasm-pack build --target web --out-dir pkg
```

### Example 2: Text Processing Snippet

**File: `src/lib.rs`**
```rust
use wasm_bindgen::prelude::*;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

extern "C" {
    fn float_log(ptr: *const u8, len: usize);
    fn float_http_request(
        url_ptr: *const u8, url_len: usize,
        method_ptr: *const u8, method_len: usize,
        body_ptr: *const u8, body_len: usize
    ) -> i32;
}

#[derive(Deserialize)]
struct TextProcessingInput {
    text: String,
    operations: Vec<String>,
    options: HashMap<String, serde_json::Value>,
}

#[derive(Serialize)]
struct TextProcessingOutput {
    original_text: String,
    processed_text: String,
    operations_applied: Vec<String>,
    statistics: TextStatistics,
}

#[derive(Serialize)]
struct TextStatistics {
    character_count: usize,
    word_count: usize,
    line_count: usize,
    operations_count: usize,
}

fn log_to_float(message: &str) {
    let bytes = message.as_bytes();
    unsafe {
        float_log(bytes.as_ptr(), bytes.len());
    }
}

#[wasm_bindgen]
pub fn process_text(input_json: &str) -> String {
    log_to_float("Starting text processing");
    
    let input: TextProcessingInput = match serde_json::from_str(input_json) {
        Ok(input) => input,
        Err(e) => {
            let error_msg = format!("Invalid input JSON: {}", e);
            log_to_float(&error_msg);
            return serde_json::to_string(&serde_json::json!({
                "error": error_msg
            })).unwrap_or_default();
        }
    };
    
    let original_text = input.text.clone();
    let mut processed_text = input.text;
    let mut operations_applied = Vec::new();
    
    // Apply operations
    for operation in &input.operations {
        match operation.as_str() {
            "uppercase" => {
                processed_text = processed_text.to_uppercase();
                operations_applied.push("uppercase".to_string());
                log_to_float("Applied uppercase transformation");
            },
            "lowercase" => {
                processed_text = processed_text.to_lowercase();
                operations_applied.push("lowercase".to_string());
                log_to_float("Applied lowercase transformation");
            },
            "trim" => {
                processed_text = processed_text.trim().to_string();
                operations_applied.push("trim".to_string());
                log_to_float("Applied trim operation");
            },
            "reverse" => {
                processed_text = processed_text.chars().rev().collect();
                operations_applied.push("reverse".to_string());
                log_to_float("Applied reverse operation");
            },
            "word_count" => {
                // This operation doesn't modify text, just logs stats
                let word_count = processed_text.split_whitespace().count();
                log_to_float(&format!("Word count: {}", word_count));
                operations_applied.push("word_count".to_string());
            },
            _ => {
                log_to_float(&format!("Unknown operation: {}", operation));
            }
        }
    }
    
    let statistics = TextStatistics {
        character_count: processed_text.len(),
        word_count: processed_text.split_whitespace().count(),
        line_count: processed_text.lines().count(),
        operations_count: operations_applied.len(),
    };
    
    let output = TextProcessingOutput {
        original_text,
        processed_text,
        operations_applied,
        statistics,
    };
    
    log_to_float("Text processing completed");
    serde_json::to_string_pretty(&output).unwrap_or_default()
}
```

## Go (TinyGo) Snippet Examples

### Example 1: Simple Calculator

**File: `main.go`**
```go
package main

import (
    "encoding/json"
    "fmt"
    "strconv"
    "unsafe"
)

// Float host function declarations
//go:wasmimport env float_log
func floatLog(ptr, len uint32)

//go:wasmimport env float_write_file
func floatWriteFile(pathPtr, pathLen, dataPtr, dataLen uint32) uint32

type CalculatorInput struct {
    Operation string  `json:"operation"`
    A         float64 `json:"a"`
    B         float64 `json:"b"`
}

type CalculatorOutput struct {
    Result    float64 `json:"result"`
    Operation string  `json:"operation"`
    Success   bool    `json:"success"`
    Error     string  `json:"error,omitempty"`
}

func logToFloat(message string) {
    ptr := uintptr(unsafe.Pointer(&[]byte(message)[0]))
    floatLog(uint32(ptr), uint32(len(message)))
}

func writeOutputFile(data *CalculatorOutput) error {
    jsonData, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal output: %v", err)
    }
    
    path := "output.json"
    pathPtr := uintptr(unsafe.Pointer(&[]byte(path)[0]))
    dataPtr := uintptr(unsafe.Pointer(&jsonData[0]))
    
    result := floatWriteFile(
        uint32(pathPtr), uint32(len(path)),
        uint32(dataPtr), uint32(len(jsonData)),
    )
    
    if result != 0 {
        return fmt.Errorf("failed to write output file")
    }
    
    return nil
}

//export calculate
func calculate(inputPtr, inputLen uint32) uint32 {
    logToFloat("Starting calculation")
    
    // Read input data
    inputBytes := make([]byte, inputLen)
    copy(inputBytes, (*(*[1 << 30]byte)(unsafe.Pointer(uintptr(inputPtr))))[0:inputLen])
    
    var input CalculatorInput
    if err := json.Unmarshal(inputBytes, &input); err != nil {
        logToFloat(fmt.Sprintf("Failed to parse input: %v", err))
        output := &CalculatorOutput{
            Success: false,
            Error:   fmt.Sprintf("Invalid input JSON: %v", err),
        }
        writeOutputFile(output)
        return 1
    }
    
    var result float64
    var err error
    
    switch input.Operation {
    case "add":
        result = input.A + input.B
        logToFloat(fmt.Sprintf("Adding %f + %f = %f", input.A, input.B, result))
    case "subtract":
        result = input.A - input.B
        logToFloat(fmt.Sprintf("Subtracting %f - %f = %f", input.A, input.B, result))
    case "multiply":
        result = input.A * input.B
        logToFloat(fmt.Sprintf("Multiplying %f * %f = %f", input.A, input.B, result))
    case "divide":
        if input.B == 0 {
            err = fmt.Errorf("division by zero")
        } else {
            result = input.A / input.B
            logToFloat(fmt.Sprintf("Dividing %f / %f = %f", input.A, input.B, result))
        }
    default:
        err = fmt.Errorf("unknown operation: %s", input.Operation)
    }
    
    output := &CalculatorOutput{
        Result:    result,
        Operation: input.Operation,
        Success:   err == nil,
    }
    
    if err != nil {
        output.Error = err.Error()
        logToFloat(fmt.Sprintf("Calculation failed: %v", err))
    } else {
        logToFloat("Calculation completed successfully")
    }
    
    if writeErr := writeOutputFile(output); writeErr != nil {
        logToFloat(fmt.Sprintf("Failed to write output: %v", writeErr))
        return 1
    }
    
    if err != nil {
        return 1
    }
    return 0
}

func main() {
    // Required for TinyGo WASM modules
}
```

**File: `float.json`**
```json
{
  "snippet_id": "simple-calculator",
  "version": "0.1.0",
  "description": "Basic arithmetic calculator",
  "language": "go",
  "entry_points": {
    "main": "Performs basic arithmetic operations (add, subtract, multiply, divide)",
    "validate": "Validates arithmetic expression syntax and operands"
  },
  "schemas": {
    "main": {
      "input": {
        "type": "object",
        "properties": {
          "operation": {
            "type": "string",
            "enum": ["add", "subtract", "multiply", "divide"]
          },
          "a": { "type": "number" },
          "b": { "type": "number" }
        },
        "required": ["operation", "a", "b"]
      },
      "output": {
        "type": "object",
        "properties": {
          "result": { "type": "number" },
          "operation": { "type": "string" },
          "success": { "type": "boolean" },
          "error": { "type": "string" }
        },
        "required": ["result", "operation", "success"]
      }
    },
    "validate": {
      "input": {
        "type": "object",
        "properties": {
          "operation": { "type": "string" },
          "a": { "type": "number" },
          "b": { "type": "number" }
        },
        "required": ["operation", "a", "b"]
      },
      "output": {
        "type": "object",
        "properties": {
          "valid": { "type": "boolean" },
          "errors": {
            "type": "array",
            "items": { "type": "string" }
          }
        },
        "required": ["valid", "errors"]
      }
    }
  }
}
```

**Build Command:**
```bash
tinygo build -o calculator.wasm -target wasm main.go
```

## Python (Pyodide) Snippet Examples

### Example 1: Data Analysis Snippet

**File: `main.py`**
```python
import json
import sys
import math
import statistics
from typing import Dict, List, Any, Optional

def float_log(message: str) -> None:
    """Log message to Float's logging system"""
    # In Pyodide, console.log is available and captured by Float
    print(f"[FLOAT_LOG] {message}")

def process_data_analysis(input_data: str) -> str:
    """
    Main entry point for data analysis snippet
    """
    float_log("Starting data analysis")
    
    try:
        # Parse input
        data = json.loads(input_data)
        float_log(f"Parsed input data with {len(data)} items")
        
        # Validate input structure
        if not isinstance(data, dict) or 'dataset' not in data:
            raise ValueError("Input must contain 'dataset' field")
        
        dataset = data['dataset']
        operations = data.get('operations', ['basic_stats'])
        
        if not isinstance(dataset, list) or not dataset:
            raise ValueError("Dataset must be a non-empty list")
        
        # Perform analysis
        results = perform_analysis(dataset, operations)
        
        output = {
            'success': True,
            'dataset_size': len(dataset),
            'operations_performed': operations,
            'results': results,
            'metadata': {
                'python_version': sys.version,
                'processing_complete': True
            }
        }
        
        float_log("Data analysis completed successfully")
        return json.dumps(output, indent=2)
        
    except Exception as e:
        float_log(f"Error during data analysis: {str(e)}")
        error_output = {
            'success': False,
            'error': str(e),
            'error_type': type(e).__name__
        }
        return json.dumps(error_output, indent=2)

def perform_analysis(dataset: List[float], operations: List[str]) -> Dict[str, Any]:
    """
    Perform statistical analysis on the dataset
    """
    results = {}
    
    # Convert dataset to numbers, filtering out non-numeric values
    numeric_data = []
    for item in dataset:
        try:
            if isinstance(item, (int, float)):
                numeric_data.append(float(item))
            elif isinstance(item, str):
                numeric_data.append(float(item))
        except (ValueError, TypeError):
            float_log(f"Skipping non-numeric value: {item}")
    
    if not numeric_data:
        raise ValueError("No numeric data found in dataset")
    
    float_log(f"Processing {len(numeric_data)} numeric values")
    
    for operation in operations:
        try:
            if operation == 'basic_stats':
                results['basic_stats'] = calculate_basic_stats(numeric_data)
            elif operation == 'descriptive_stats':
                results['descriptive_stats'] = calculate_descriptive_stats(numeric_data)
            elif operation == 'distribution_analysis':
                results['distribution_analysis'] = analyze_distribution(numeric_data)
            elif operation == 'outlier_detection':
                results['outlier_detection'] = detect_outliers(numeric_data)
            else:
                float_log(f"Unknown operation: {operation}")
                
        except Exception as e:
            float_log(f"Error in operation {operation}: {str(e)}")
            results[operation] = {'error': str(e)}
    
    return results

def calculate_basic_stats(data: List[float]) -> Dict[str, float]:
    """Calculate basic statistical measures"""
    return {
        'count': len(data),
        'sum': sum(data),
        'mean': statistics.mean(data),
        'median': statistics.median(data),
        'min': min(data),
        'max': max(data),
        'range': max(data) - min(data)
    }

def calculate_descriptive_stats(data: List[float]) -> Dict[str, float]:
    """Calculate descriptive statistics"""
    mean_val = statistics.mean(data)
    
    stats = {
        'variance': statistics.variance(data),
        'standard_deviation': statistics.stdev(data),
        'coefficient_of_variation': statistics.stdev(data) / mean_val if mean_val != 0 else 0
    }
    
    # Add percentiles
    sorted_data = sorted(data)
    n = len(sorted_data)
    
    stats.update({
        'q1': sorted_data[n // 4],
        'q3': sorted_data[3 * n // 4],
        'iqr': sorted_data[3 * n // 4] - sorted_data[n // 4]
    })
    
    return stats

def analyze_distribution(data: List[float]) -> Dict[str, Any]:
    """Analyze data distribution characteristics"""
    mean_val = statistics.mean(data)
    std_dev = statistics.stdev(data)
    
    # Simple skewness calculation
    n = len(data)
    skewness = sum(((x - mean_val) / std_dev) ** 3 for x in data) / n
    
    # Simple kurtosis calculation
    kurtosis = sum(((x - mean_val) / std_dev) ** 4 for x in data) / n - 3
    
    return {
        'skewness': skewness,
        'kurtosis': kurtosis,
        'is_normal_like': abs(skewness) < 0.5 and abs(kurtosis) < 0.5
    }

def detect_outliers(data: List[float]) -> Dict[str, Any]:
    """Detect outliers using IQR method"""
    sorted_data = sorted(data)
    n = len(sorted_data)
    
    q1 = sorted_data[n // 4]
    q3 = sorted_data[3 * n // 4]
    iqr = q3 - q1
    
    lower_bound = q1 - 1.5 * iqr
    upper_bound = q3 + 1.5 * iqr
    
    outliers = [x for x in data if x < lower_bound or x > upper_bound]
    
    return {
        'outlier_count': len(outliers),
        'outlier_percentage': (len(outliers) / len(data)) * 100,
        'outliers': outliers[:10],  # Limit to first 10 outliers
        'lower_bound': lower_bound,
        'upper_bound': upper_bound
    }

# Entry point for Float WASM execution
if __name__ == "__main__":
    # Read input from stdin or command line argument
    if len(sys.argv) > 1:
        input_data = sys.argv[1]
    else:
        input_data = sys.stdin.read()
    
    result = process_data_analysis(input_data)
    print(result)
```

**File: `float.json`**
```json
{
  "snippet_id": "data-analysis-python",
  "version": "0.1.0",
  "description": "Statistical data analysis using Python",
  "language": "python",
  "runtime": "pyodide",
  "entry_points": {
    "main": "Performs comprehensive statistical analysis on numeric datasets",
    "analyze": "Analyzes data with customizable operations and detailed reporting"
  },
  "schemas": {
    "main": {
      "input": {
        "type": "object",
        "properties": {
          "dataset": {
            "type": "array",
            "items": { "type": "number" },
            "minItems": 1,
            "description": "Array of numeric values to analyze"
          },
          "operations": {
            "type": "array",
            "items": {
              "type": "string",
              "enum": ["basic_stats", "descriptive_stats", "distribution_analysis", "outlier_detection"]
            },
            "default": ["basic_stats"],
            "description": "List of analysis operations to perform"
          }
        },
        "required": ["dataset"]
      },
      "output": {
        "type": "object",
        "properties": {
          "success": { "type": "boolean" },
          "dataset_size": { "type": "integer" },
          "operations_performed": {
            "type": "array",
            "items": { "type": "string" }
          },
          "results": { "type": "object" },
          "metadata": { "type": "object" },
          "error": { "type": "string" },
          "error_type": { "type": "string" }
        },
        "required": ["success"]
      }
    },
    "analyze": {
      "input": {
        "type": "string",
        "description": "Raw JSON string containing dataset and operations"
      },
      "output": {
        "type": "string",
        "description": "JSON string containing analysis results"
      }
    }
  }
}
```

### Example 2: Text Processing with NLP

**File: `nlp_processor.py`**
```python
import json
import re
from typing import Dict, List, Any
from collections import Counter

def float_log(message: str) -> None:
    """Log message to Float's logging system"""
    print(f"[FLOAT_LOG] {message}")

def process_text_nlp(input_data: str) -> str:
    """
    Natural Language Processing text analysis
    """
    float_log("Starting NLP text processing")
    
    try:
        data = json.loads(input_data)
        text = data.get('text', '')
        operations = data.get('operations', ['tokenize', 'word_frequency'])
        
        if not text:
            raise ValueError("Text field is required and cannot be empty")
        
        float_log(f"Processing text of length {len(text)} characters")
        
        results = {}
        
        for operation in operations:
            float_log(f"Performing operation: {operation}")
            
            if operation == 'tokenize':
                results['tokenization'] = tokenize_text(text)
            elif operation == 'word_frequency':
                results['word_frequency'] = analyze_word_frequency(text)
            elif operation == 'sentiment_basic':
                results['sentiment'] = basic_sentiment_analysis(text)
            elif operation == 'readability':
                results['readability'] = calculate_readability(text)
            elif operation == 'extract_entities':
                results['entities'] = extract_basic_entities(text)
            else:
                float_log(f"Unknown operation: {operation}")
        
        output = {
            'success': True,
            'text_length': len(text),
            'operations_performed': operations,
            'results': results
        }
        
        float_log("NLP processing completed successfully")
        return json.dumps(output, indent=2)
        
    except Exception as e:
        float_log(f"Error during NLP processing: {str(e)}")
        error_output = {
            'success': False,
            'error': str(e),
            'error_type': type(e).__name__
        }
        return json.dumps(error_output, indent=2)

def tokenize_text(text: str) -> Dict[str, Any]:
    """Basic text tokenization"""
    # Simple tokenization
    sentences = re.split(r'[.!?]+', text)
    sentences = [s.strip() for s in sentences if s.strip()]
    
    # Word tokenization
    words = re.findall(r'\b\w+\b', text.lower())
    
    # Paragraph tokenization
    paragraphs = [p.strip() for p in text.split('\n\n') if p.strip()]
    
    return {
        'sentence_count': len(sentences),
        'word_count': len(words),
        'paragraph_count': len(paragraphs),
        'unique_words': len(set(words)),
        'sentences': sentences[:5],  # First 5 sentences
        'average_sentence_length': sum(len(s.split()) for s in sentences) / len(sentences) if sentences else 0
    }

def analyze_word_frequency(text: str) -> Dict[str, Any]:
    """Analyze word frequency"""
    words = re.findall(r'\b\w+\b', text.lower())
    
    # Remove common stop words
    stop_words = {'the', 'a', 'an', 'and', 'or', 'but', 'in', 'on', 'at', 'to', 'for', 'of', 'with', 'by', 'is', 'are', 'was', 'were', 'be', 'been', 'being', 'have', 'has', 'had', 'do', 'does', 'did', 'will', 'would', 'could', 'should', 'may', 'might', 'must', 'shall', 'can', 'this', 'that', 'these', 'those', 'i', 'you', 'he', 'she', 'it', 'we', 'they', 'me', 'him', 'her', 'us', 'them'}
    
    filtered_words = [word for word in words if word not in stop_words]
    word_freq = Counter(filtered_words)
    
    return {
        'total_words': len(words),
        'unique_words': len(set(words)),
        'filtered_words': len(filtered_words),
        'top_10_words': word_freq.most_common(10),
        'vocabulary_richness': len(set(words)) / len(words) if words else 0
    }

def basic_sentiment_analysis(text: str) -> Dict[str, Any]:
    """Basic rule-based sentiment analysis"""
    positive_words = {'good', 'great', 'excellent', 'amazing', 'wonderful', 'fantastic', 'awesome', 'love', 'like', 'happy', 'joy', 'pleased', 'satisfied', 'perfect', 'best'}
    negative_words = {'bad', 'terrible', 'awful', 'horrible', 'hate', 'dislike', 'angry', 'sad', 'disappointed', 'worst', 'pathetic', 'disgusting', 'annoying', 'frustrated'}
    
    words = re.findall(r'\b\w+\b', text.lower())
    
    positive_count = sum(1 for word in words if word in positive_words)
    negative_count = sum(1 for word in words if word in negative_words)
    
    total_sentiment_words = positive_count + negative_count
    
    if total_sentiment_words == 0:
        sentiment = 'neutral'
        score = 0.0
    else:
        score = (positive_count - negative_count) / total_sentiment_words
        if score > 0.1:
            sentiment = 'positive'
        elif score < -0.1:
            sentiment = 'negative'
        else:
            sentiment = 'neutral'
    
    return {
        'sentiment': sentiment,
        'score': score,
        'positive_words_count': positive_count,
        'negative_words_count': negative_count,
        'confidence': abs(score)
    }

def calculate_readability(text: str) -> Dict[str, Any]:
    """Calculate basic readability metrics"""
    sentences = re.split(r'[.!?]+', text)
    sentences = [s.strip() for s in sentences if s.strip()]
    
    words = re.findall(r'\b\w+\b', text)
    
    if not sentences or not words:
        return {'error': 'Insufficient text for readability analysis'}
    
    avg_sentence_length = len(words) / len(sentences)
    
    # Count syllables (simple approximation)
    def count_syllables(word):
        word = word.lower()
        syllables = 0
        vowels = 'aeiouy'
        prev_was_vowel = False
        for char in word:
            if char in vowels:
                if not prev_was_vowel:
                    syllables += 1
                prev_was_vowel = True
            else:
                prev_was_vowel = False
        if word.endswith('e'):
            syllables -= 1
        if syllables == 0:
            syllables = 1
        return syllables
    
    total_syllables = sum(count_syllables(word) for word in words)
    avg_syllables_per_word = total_syllables / len(words)
    
    # Simple Flesch Reading Ease approximation
    flesch_score = 206.835 - (1.015 * avg_sentence_length) - (84.6 * avg_syllables_per_word)
    
    if flesch_score >= 90:
        reading_level = 'Very Easy'
    elif flesch_score >= 80:
        reading_level = 'Easy'
    elif flesch_score >= 70:
        reading_level = 'Fairly Easy'
    elif flesch_score >= 60:
        reading_level = 'Standard'
    elif flesch_score >= 50:
        reading_level = 'Fairly Difficult'
    elif flesch_score >= 30:
        reading_level = 'Difficult'
    else:
        reading_level = 'Very Difficult'
    
    return {
        'flesch_score': round(flesch_score, 2),
        'reading_level': reading_level,
        'avg_sentence_length': round(avg_sentence_length, 2),
        'avg_syllables_per_word': round(avg_syllables_per_word, 2),
        'total_sentences': len(sentences),
        'total_words': len(words)
    }

def extract_basic_entities(text: str) -> Dict[str, Any]:
    """Extract basic entities using pattern matching"""
    # Email addresses
    emails = re.findall(r'\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b', text)
    
    # Phone numbers (simple patterns)
    phones = re.findall(r'\b(?:\+?1[-.\s]?)?\(?([0-9]{3})\)?[-.\s]?([0-9]{3})[-.\s]?([0-9]{4})\b', text)
    
    # URLs
    urls = re.findall(r'https?://[^\s]+', text)
    
    # Dates (simple patterns)
    dates = re.findall(r'\b\d{1,2}[/-]\d{1,2}[/-]\d{2,4}\b', text)
    
    # Numbers
    numbers = re.findall(r'\b\d+(?:\.\d+)?\b', text)
    
    # Capitalized words (potential proper nouns)
    proper_nouns = re.findall(r'\b[A-Z][a-z]+\b', text)
    
    return {
        'emails': emails,
        'phone_numbers': ['-'.join(phone) for phone in phones],
        'urls': urls,
        'dates': dates,
        'numbers': numbers[:20],  # Limit to first 20 numbers
        'potential_proper_nouns': list(set(proper_nouns))[:20]  # Unique proper nouns, limited to 20
    }

# Entry point
if __name__ == "__main__":
    import sys
    
    if len(sys.argv) > 1:
        input_data = sys.argv[1]
    else:
        input_data = sys.stdin.read()
    
    result = process_text_nlp(input_data)
    print(result)
```

## JavaScript Snippet Examples

### Example 1: JSON Data Transformer

**File: `transformer.js`**
```javascript
// Float host function declarations
const Float = {
    log: function(message) {
        if (typeof console !== 'undefined' && console.log) {
            console.log(`[FLOAT_LOG] ${message}`);
        }
    },
    
    writeFile: function(path, data) {
        // This would be implemented by Float's JavaScript runtime
        if (typeof floatWriteFile !== 'undefined') {
            return floatWriteFile(path, data);
        }
        return false;
    }
};

class DataTransformer {
    constructor() {
        this.transformations = {
            'uppercase_keys': this.uppercaseKeys.bind(this),
            'lowercase_keys': this.lowercaseKeys.bind(this),
            'snake_to_camel': this.snakeToCamel.bind(this),
            'camel_to_snake': this.camelToSnake.bind(this),
            'flatten_object': this.flattenObject.bind(this),
            'filter_null_values': this.filterNullValues.bind(this),
            'sort_arrays': this.sortArrays.bind(this),
            'normalize_strings': this.normalizeStrings.bind(this)
        };
    }
    
    transform(data, operations) {
        Float.log(`Starting transformation with ${operations.length} operations`);
        
        let result = JSON.parse(JSON.stringify(data)); // Deep clone
        const appliedOperations = [];
        
        for (const operation of operations) {
            try {
                if (this.transformations[operation]) {
                    Float.log(`Applying transformation: ${operation}`);
                    result = this.transformations[operation](result);
                    appliedOperations.push(operation);
                } else {
                    Float.log(`Unknown transformation: ${operation}`);
                }
            } catch (error) {
                Float.log(`Error in transformation ${operation}: ${error.message}`);
                throw new Error(`Transformation failed: ${operation} - ${error.message}`);
            }
        }
        
        Float.log(`Transformation completed. Applied ${appliedOperations.length} operations`);
        
        return {
            success: true,
            original_data: data,
            transformed_data: result,
            applied_operations: appliedOperations,
            transformation_count: appliedOperations.length
        };
    }
    
    uppercaseKeys(obj) {
        if (Array.isArray(obj)) {
            return obj.map(item => this.uppercaseKeys(item));
        } else if (obj !== null && typeof obj === 'object') {
            const newObj = {};
            for (const [key, value] of Object.entries(obj)) {
                newObj[key.toUpperCase()] = this.uppercaseKeys(value);
            }
            return newObj;
        }
        return obj;
    }
    
    lowercaseKeys(obj) {
        if (Array.isArray(obj)) {
            return obj.map(item => this.lowercaseKeys(item));
        } else if (obj !== null && typeof obj === 'object') {
            const newObj = {};
            for (const [key, value] of Object.entries(obj)) {
                newObj[key.toLowerCase()] = this.lowercaseKeys(value);
            }
            return newObj;
        }
        return obj;
    }
    
    snakeToCamel(obj) {
        if (Array.isArray(obj)) {
            return obj.map(item => this.snakeToCamel(item));
        } else if (obj !== null && typeof obj === 'object') {
            const newObj = {};
            for (const [key, value] of Object.entries(obj)) {
                const camelKey = key.replace(/_([a-z])/g, (match, letter) => letter.toUpperCase());
                newObj[camelKey] = this.snakeToCamel(value);
            }
            return newObj;
        }
        return obj;
    }
    
    camelToSnake(obj) {
        if (Array.isArray(obj)) {
            return obj.map(item => this.camelToSnake(item));
        } else if (obj !== null && typeof obj === 'object') {
            const newObj = {};
            for (const [key, value] of Object.entries(obj)) {
                const snakeKey = key.replace(/([A-Z])/g, '_$1').toLowerCase();
                newObj[snakeKey] = this.camelToSnake(value);
            }
            return newObj;
        }
        return obj;
    }
    
    flattenObject(obj, prefix = '', result = {}) {
        for (const [key, value] of Object.entries(obj)) {
            const newKey = prefix ? `${prefix}.${key}` : key;
            
            if (value !== null && typeof value === 'object' && !Array.isArray(value)) {
                this.flattenObject(value, newKey, result);
            } else {
                result[newKey] = value;
            }
        }
        return result;
    }
    
    filterNullValues(obj) {
        if (Array.isArray(obj)) {
            return obj
                .filter(item => item !== null && item !== undefined)
                .map(item => this.filterNullValues(item));
        } else if (obj !== null && typeof obj === 'object') {
            const newObj = {};
            for (const [key, value] of Object.entries(obj)) {
                if (value !== null && value !== undefined) {
                    newObj[key] = this.filterNullValues(value);
                }
            }
            return newObj;
        }
        return obj;
    }
    
    sortArrays(obj) {
        if (Array.isArray(obj)) {
            const sorted = [...obj].sort((a, b) => {
                if (typeof a === 'string' && typeof b === 'string') {
                    return a.localeCompare(b);
                } else if (typeof a === 'number' && typeof b === 'number') {
                    return a - b;
                } else {
                    return String(a).localeCompare(String(b));
                }
            });
            return sorted.map(item => this.sortArrays(item));
        } else if (obj !== null && typeof obj === 'object') {
            const newObj = {};
            for (const [key, value] of Object.entries(obj)) {
                newObj[key] = this.sortArrays(value);
            }
            return newObj;
        }
        return obj;
    }
    
    normalizeStrings(obj) {
        if (Array.isArray(obj)) {
            return obj.map(item => this.normalizeStrings(item));
        } else if (obj !== null && typeof obj === 'object') {
            const newObj = {};
            for (const [key, value] of Object.entries(obj)) {
                newObj[key] = this.normalizeStrings(value);
            }
            return newObj;
        } else if (typeof obj === 'string') {
            return obj.trim().replace(/\s+/g, ' ');
        }
        return obj;
    }
}

function processTransformation(inputJson) {
    try {
        Float.log("Starting JSON transformation process");
        
        const input = JSON.parse(inputJson);
        
        // Validate input
        if (!input.data) {
            throw new Error("Input must contain 'data' field");
        }
        
        if (!input.operations || !Array.isArray(input.operations)) {
            throw new Error("Input must contain 'operations' array");
        }
        
        const transformer = new DataTransformer();
        const result = transformer.transform(input.data, input.operations);
        
        // Write output
        const outputJson = JSON.stringify(result, null, 2);
        Float.writeFile('output.json', outputJson);
        
        Float.log("JSON transformation completed successfully");
        return outputJson;
        
    } catch (error) {
        Float.log(`Error during transformation: ${error.message}`);
        
        const errorOutput = {
            success: false,
            error: error.message,
            error_type: error.constructor.name,
            stack: error.stack
        };
        
        return JSON.stringify(errorOutput, null, 2);
    }
}

// Export for WASM module
if (typeof module !== 'undefined' && module.exports) {
    module.exports = { processTransformation, DataTransformer };
}

// Global function for WASM execution
this.processTransformation = processTransformation;
```

**File: `float.json`**
```json
{
  "snippet_id": "json-transformer-js",
  "version": "0.1.0",
  "description": "JSON data transformation utilities",
  "language": "javascript",
  "runtime": "quickjs",
  "entry_points": {
    "main": "Transforms JSON data using configurable transformation operations",
    "transform": "Applies data transformations like key conversion and object flattening"
  },
  "schemas": {
    "main": {
      "input": {
        "type": "object",
        "properties": {
          "data": {
            "type": "object",
            "description": "The JSON data to transform"
          },
          "operations": {
            "type": "array",
            "items": {
              "type": "string",
              "enum": [
                "uppercase_keys",
                "lowercase_keys", 
                "snake_to_camel",
                "camel_to_snake",
                "flatten_object",
                "filter_null_values",
                "sort_arrays",
                "normalize_strings"
              ]
            },
            "description": "List of transformation operations to apply"
          }
        },
        "required": ["data", "operations"]
      },
      "output": {
        "type": "object",
        "properties": {
          "success": { "type": "boolean" },
          "original_data": { "type": "object" },
          "transformed_data": { "type": "object" },
          "applied_operations": {
            "type": "array",
            "items": { "type": "string" }
          },
          "transformation_count": { "type": "integer" },
          "error": { "type": "string" },
          "error_type": { "type": "string" }
        },
        "required": ["success"]
      }
    },
    "transform": {
      "input": {
        "type": "string",
        "description": "JSON string containing data and operations to transform"
      },
      "output": {
        "type": "string",
        "description": "JSON string containing transformation results"
      }
    }
  }
}
```

## TypeScript Snippet Examples

### Example 1: Type-Safe Data Validator

**File: `validator.ts`**
```typescript
// Float host function declarations
declare function floatLog(message: string): void;
declare function floatWriteFile(path: string, data: string): boolean;

// Type definitions
interface ValidationRule {
    type: 'string' | 'number' | 'boolean' | 'array' | 'object';
    required?: boolean;
    minLength?: number;
    maxLength?: number;
    min?: number;
    max?: number;
    pattern?: string;
    enum?: any[];
    properties?: { [key: string]: ValidationRule };
    items?: ValidationRule;
}

interface ValidationSchema {
    [key: string]: ValidationRule;
}

interface ValidationInput {
    data: any;
    schema: ValidationSchema;
    strict?: boolean;
}

interface ValidationError {
    path: string;
    message: string;
    code: string;
    value?: any;
}

interface ValidationResult {
    valid: boolean;
    errors: ValidationError[];
    validatedData: any;
    summary: {
        totalFields: number;
        validFields: number;
        errorCount: number;
    };
}

class TypeSafeValidator {
    private errors: ValidationError[] = [];
    private strict: boolean = false;
    
    constructor(strict: boolean = false) {
        this.strict = strict;
    }
    
    validate(data: any, schema: ValidationSchema): ValidationResult {
        this.errors = [];
        
        floatLog("Starting type-safe validation");
        floatLog(`Validating ${Object.keys(schema).length} schema fields`);
        
        const validatedData = this.validateObject(data, schema, '');
        
        const result: ValidationResult = {
            valid: this.errors.length === 0,
            errors: this.errors,
            validatedData: validatedData,
            summary: {
                totalFields: Object.keys(schema).length,
                validFields: Object.keys(schema).length - this.errors.length,
                errorCount: this.errors.length
            }
        };
        
        floatLog(`Validation completed: ${result.valid ? 'PASSED' : 'FAILED'} (${this.errors.length} errors)`);
        
        return result;
    }
    
    private validateObject(data: any, schema: ValidationSchema, basePath: string): any {
        if (typeof data !== 'object' || data === null) {
            this.addError(basePath, 'Expected object', 'TYPE_MISMATCH', data);
            return {};
        }
        
        const validatedObj: any = {};
        
        // Validate each property in schema
        for (const [key, rule] of Object.entries(schema)) {
            const path = basePath ? `${basePath}.${key}` : key;
            const value = data[key];
            
            if (value === undefined || value === null) {
                if (rule.required) {
                    this.addError(path, `Required field '${key}' is missing`, 'REQUIRED_FIELD_MISSING');
                }
                continue;
            }
            
            const validatedValue = this.validateValue(value, rule, path);
            if (validatedValue !== undefined) {
                validatedObj[key] = validatedValue;
            }
        }
        
        // Check for extra properties in strict mode
        if (this.strict) {
            for (const key of Object.keys(data)) {
                if (!(key in schema)) {
                    const path = basePath ? `${basePath}.${key}` : key;
                    this.addError(path, `Unexpected property '${key}' in strict mode`, 'UNEXPECTED_PROPERTY', data[key]);
                }
            }
        }
        
        return validatedObj;
    }
    
    private validateValue(value: any, rule: ValidationRule, path: string): any {
        // Type validation
        if (!this.validateType(value, rule.type, path)) {
            return undefined;
        }
        
        // Type-specific validations
        switch (rule.type) {
            case 'string':
                return this.validateString(value as string, rule, path);
            case 'number':
                return this.validateNumber(value as number, rule, path);
            case 'array':
                return this.validateArray(value as any[], rule, path);
            case 'object':
                return this.validateNestedObject(value, rule, path);
            case 'boolean':
                return value as boolean;
            default:
                return value;
        }
    }
    
    private validateType(value: any, expectedType: string, path: string): boolean {
        let actualType = typeof value;
        
        // Handle array type
        if (expectedType === 'array') {
            if (!Array.isArray(value)) {
                this.addError(path, `Expected array, got ${actualType}`, 'TYPE_MISMATCH', value);
                return false;
            }
            return true;
        }
        
        // Handle object type
        if (expectedType === 'object') {
            if (actualType !== 'object' || value === null || Array.isArray(value)) {
                this.addError(path, `Expected object, got ${actualType}`, 'TYPE_MISMATCH', value);
                return false;
            }
            return true;
        }
        
        // Handle primitive types
        if (actualType !== expectedType) {
            this.addError(path, `Expected ${expectedType}, got ${actualType}`, 'TYPE_MISMATCH', value);
            return false;
        }
        
        return true;
    }
    
    private validateString(value: string, rule: ValidationRule, path: string): string {
        // Length validation
        if (rule.minLength !== undefined && value.length < rule.minLength) {
            this.addError(path, `String must be at least ${rule.minLength} characters`, 'MIN_LENGTH_VIOLATION', value);
        }
        
        if (rule.maxLength !== undefined && value.length > rule.maxLength) {
            this.addError(path, `String must be at most ${rule.maxLength} characters`, 'MAX_LENGTH_VIOLATION', value);
        }
        
        // Pattern validation
        if (rule.pattern) {
            try {
                const regex = new RegExp(rule.pattern);
                if (!regex.test(value)) {
                    this.addError(path, `String does not match pattern: ${rule.pattern}`, 'PATTERN_MISMATCH', value);
                }
            } catch (e) {
                floatLog(`Invalid regex pattern: ${rule.pattern}`);
            }
        }
        
        // Enum validation
        if (rule.enum && !rule.enum.includes(value)) {
            this.addError(path, `Value must be one of: ${rule.enum.join(', ')}`, 'ENUM_VIOLATION', value);
        }
        
        return value;
    }
    
    private validateNumber(value: number, rule: ValidationRule, path: string): number {
        // Range validation
        if (rule.min !== undefined && value < rule.min) {
            this.addError(path, `Number must be at least ${rule.min}`, 'MIN_VALUE_VIOLATION', value);
        }
        
        if (rule.max !== undefined && value > rule.max) {
            this.addError(path, `Number must be at most ${rule.max}`, 'MAX_VALUE_VIOLATION', value);
        }
        
        // Enum validation
        if (rule.enum && !rule.enum.includes(value)) {
            this.addError(path, `Value must be one of: ${rule.enum.join(', ')}`, 'ENUM_VIOLATION', value);
        }
        
        return value;
    }
    
    private validateArray(value: any[], rule: ValidationRule, path: string): any[] {
        const validatedArray: any[] = [];
        
        // Length validation
        if (rule.minLength !== undefined && value.length < rule.minLength) {
            this.addError(path, `Array must have at least ${rule.minLength} items`, 'MIN_LENGTH_VIOLATION', value);
        }
        
        if (rule.maxLength !== undefined && value.length > rule.maxLength) {
            this.addError(path, `Array must have at most ${rule.maxLength} items`, 'MAX_LENGTH_VIOLATION', value);
        }
        
        // Item validation
        if (rule.items) {
            value.forEach((item, index) => {
                const itemPath = `${path}[${index}]`;
                const validatedItem = this.validateValue(item, rule.items!, itemPath);
                if (validatedItem !== undefined) {
                    validatedArray.push(validatedItem);
                }
            });
        } else {
            validatedArray.push(...value);
        }
        
        return validatedArray;
    }
    
    private validateNestedObject(value: any, rule: ValidationRule, path: string): any {
        if (rule.properties) {
            return this.validateObject(value, rule.properties, path);
        }
        return value;
    }
    
    private addError(path: string, message: string, code: string, value?: any): void {
        this.errors.push({
            path,
            message,
            code,
            value
        });
    }
}

// Main processing function
function processValidation(inputJson: string): string {
    try {
        floatLog("Starting TypeScript validation process");
        
        const input: ValidationInput = JSON.parse(inputJson);
        
        // Validate input structure
        if (!input.data) {
            throw new Error("Input must contain 'data' field");
        }
        
        if (!input.schema) {
            throw new Error("Input must contain 'schema' field");
        }
        
        const validator = new TypeSafeValidator(input.strict || false);
        const result = validator.validate(input.data, input.schema);
        
        // Add metadata
        const output = {
            ...result,
            metadata: {
                typescript_version: "5.0.0",
                validation_timestamp: new Date().toISOString(),
                strict_mode: input.strict || false
            }
        };
        
        // Write output file
        const outputJson = JSON.stringify(output, null, 2);
        floatWriteFile('output.json', outputJson);
        
        floatLog("TypeScript validation completed successfully");
        return outputJson;
        
    } catch (error: any) {
        floatLog(`Error during validation: ${error.message}`);
        
        const errorOutput = {
            valid: false,
            error: error.message,
            error_type: error.constructor.name,
            stack: error.stack
        };
        
        return JSON.stringify(errorOutput, null, 2);
    }
}

// Export for WASM module
declare global {
    var processValidation: (input: string) => string;
}

globalThis.processValidation = processValidation;
```

**File: `tsconfig.json`**
```json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "ESNext",
    "lib": ["ES2020"],
    "strict": true,
    "noImplicitAny": true,
    "strictNullChecks": true,
    "strictFunctionTypes": true,
    "noImplicitThis": true,
    "noImplicitReturns": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "declaration": false,
    "outDir": "./dist",
    "rootDir": "./",
    "removeComments": true,
    "skipLibCheck": true
  },
  "include": [
    "validator.ts"
  ],
  "exclude": [
    "node_modules",
    "dist"
  ]
}
```

**File: `float.json`**
```json
{
  "snippet_id": "type-safe-validator-ts",
  "version": "0.1.0",
  "description": "Type-safe data validation with comprehensive error reporting",
  "language": "typescript",
  "runtime": "quickjs",
  "entry_points": {
    "main": "Performs type-safe validation with comprehensive error reporting",
    "validate": "Validates data against TypeScript-style schemas with strict type checking"
  },
  "build_command": "tsc && node dist/validator.js",
  "schemas": {
    "main": {
      "input": {
        "type": "object",
        "properties": {
          "data": {
            "type": "object",
            "description": "The data to validate"
          },
          "schema": {
            "type": "object",
            "description": "Validation schema with type definitions",
            "additionalProperties": {
              "type": "object",
              "properties": {
                "type": {
                  "type": "string",
                  "enum": ["string", "number", "boolean", "array", "object"]
                },
                "required": { "type": "boolean" },
                "minLength": { "type": "integer" },
                "maxLength": { "type": "integer" },
                "min": { "type": "number" },
                "max": { "type": "number" },
                "pattern": { "type": "string" },
                "enum": { "type": "array" }
              },
              "required": ["type"]
            }
          },
          "strict": {
            "type": "boolean",
            "default": false,
            "description": "Enable strict mode validation"
          }
        },
        "required": ["data", "schema"]
      },
      "output": {
        "type": "object",
        "properties": {
          "valid": { "type": "boolean" },
          "errors": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "path": { "type": "string" },
                "message": { "type": "string" },
                "code": { "type": "string" },
                "value": {}
              },
              "required": ["path", "message", "code"]
            }
          },
          "validatedData": { "type": "object" },
          "summary": {
            "type": "object",
            "properties": {
              "totalFields": { "type": "integer" },
              "validFields": { "type": "integer" },
              "errorCount": { "type": "integer" }
            },
            "required": ["totalFields", "validFields", "errorCount"]
          },
          "metadata": { "type": "object" }
        },
        "required": ["valid", "errors", "validatedData", "summary"]
      }
    },
    "validate": {
      "input": {
        "type": "string",
        "description": "JSON string containing data, schema, and options to validate"
      },
      "output": {
        "type": "string",
        "description": "JSON string containing validation results with errors and metadata"
      }
    }
  }
}
```

## Float.json Configuration Guide

The `float.json` file is the cornerstone of every Float snippet, serving as both a manifest and API contract. This section provides a comprehensive guide to constructing and understanding this crucial configuration file.

### Overview

The `float.json` file defines:
- **Snippet metadata** (ID, version, description)
- **Entry point descriptions** (what operations users can call and what each does)
- **Input/output schemas** (API contracts for each entry point)
- **Runtime configuration** (language, dependencies, build commands)

> **Important**: Entry point values should describe **what the operation does** when called, not internal function names. Think of them as user-facing operation descriptions.

### Core Structure

```json
{
  "snippet_id": "my-snippet",           // Unique identifier
  "version": "1.0.0",                   // Semantic version
  "name": "My Snippet",                 // Human-readable name
  "description": "What this snippet does",
  "language": "rust",                   // Programming language
  "entry_points": {                     // Callable operations with descriptions
    "main": "Primary data processing operation",
    "validate": "Validates input data before processing"
  },
  "schemas": {                          // API contracts
    "main": { "input": {...}, "output": {...} },
    "validate": { "input": {...}, "output": {...} }
  }
}
```

### Entry Points Explained

Entry points define **what operations can be called** from your snippet and **what each operation does**:

#### Simple Entry Point (Single Operation)
```json
