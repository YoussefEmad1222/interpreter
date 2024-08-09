package main

type ErrorType string
type operationType string
type TokenType string

const (
	SYNTAX_ERROR  ErrorType = "Syntax Error"
	RUNTIME_ERROR ErrorType = "Runtime Error"
	NONE          ErrorType = "None"
)

const (
	NUMBER_OP  operationType = "NUMBER"
	STRING_OP  operationType = "STRING"
	LOGICAL_OP operationType = "LOGICAL"
	NO_OP      operationType = "NO_OP"
)

var operand = map[TokenType]bool{
	NUMBER:     true,
	STRING:     true,
	TRUE:       true,
	FALSE:      true,
	NIL:        true,
	IDENTIFIER: true,
}

var operators = map[TokenType]bool{
	PLUS:          true,
	MINUS:         true,
	STAR:          true,
	SLASH:         true,
	GREATER:       true,
	GREATER_EQUAL: true,
	LESS:          true,
	LESS_EQUAL:    true,
	EQUAL_EQUAL:   true,
	BANG_EQUAL:    true,
	LEFT_PAREN:    true,
}

const (
	LEFT_PAREN    TokenType = "LEFT_PAREN"
	RIGHT_PAREN   TokenType = "RIGHT_PAREN"
	LEFT_BRACE    TokenType = "LEFT_BRACE"
	RIGHT_BRACE   TokenType = "RIGHT_BRACE"
	COMMA         TokenType = "COMMA"
	DOT           TokenType = "DOT"
	MINUS         TokenType = "MINUS"
	PLUS          TokenType = "PLUS"
	SEMICOLON     TokenType = "SEMICOLON"
	SLASH         TokenType = "SLASH"
	STAR          TokenType = "STAR"
	BANG          TokenType = "BANG"
	BANG_EQUAL    TokenType = "BANG_EQUAL"
	EQUAL         TokenType = "EQUAL"
	EQUAL_EQUAL   TokenType = "EQUAL_EQUAL"
	GREATER       TokenType = "GREATER"
	GREATER_EQUAL TokenType = "GREATER_EQUAL"
	LESS          TokenType = "LESS"
	LESS_EQUAL    TokenType = "LESS_EQUAL"
	IDENTIFIER    TokenType = "IDENTIFIER"
	STRING        TokenType = "STRING"
	NUMBER        TokenType = "NUMBER"
	AND           TokenType = "AND"
	CLASS         TokenType = "CLASS"
	ELSE          TokenType = "ELSE"
	FALSE         TokenType = "FALSE"
	FUN           TokenType = "FUN"
	FOR           TokenType = "FOR"
	IF            TokenType = "IF"
	NIL           TokenType = "NIL"
	OR            TokenType = "OR"
	PRINT         TokenType = "PRINT"
	RETURN        TokenType = "RETURN"
	SUPER         TokenType = "SUPER"
	THIS          TokenType = "THIS"
	TRUE          TokenType = "TRUE"
	VAR           TokenType = "VAR"
	WHILE         TokenType = "WHILE"
	EOF           TokenType = "EOF"
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}
