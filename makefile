# @echo "Tests for the Lexer ::"
# @go test ./internal/lexer
# @echo "Tests for the AST::"
# @go test ./internal/ast
run:
	@go run cmd/main.go
test:
	@echo "Tests for the Parser::"
	@go test ./internal/parser
	