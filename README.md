# Go Interpreter

This repository contains the implementation of a Go-based interpreter, inspired by the concepts and examples from the book. The project is structured to follow the design and implementation outlined in the book's PDF, with the core functionality written in Go.

## Features

- **Lexer**: Tokenizes the input source code into meaningful tokens.
- **Parser**: Parses the tokens into an Abstract Syntax Tree (AST).
- **Evaluator**: Evaluates the AST to produce results.
- **REPL**: A Read-Eval-Print Loop for interactive usage.

## Structure

The project is organized into the following directories:

- **`internal/lexer`**: Contains the lexer implementation for tokenizing input.
- **`internal/parser`**: Contains the parser implementation for building the AST.
- **`internal/evaluator`**: Contains the evaluator for interpreting the AST.
- **`internal/repl`**: Contains the REPL for interactive usage.

## How to Run

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/go-interpreter.git
   cd go-interpreter
2. Build the project:
    - go build ./...
2. Run the REPL:
    - go run cmd/main.go

## Example Usage
`>>` let x = 5;
`>>` x + 10;
15
`>>` quit