package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		_, _ = fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]
	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	tokenError = NONE
	lexer := NewLexer(string(fileContents))
	lexer.tokenize()
	if command == "tokenize" {
		lexer.printTokens()
	} else if command == "evaluate" {
		evaluate(lexer)
	} else if command == "parse" {
		parser := NewParser(lexer.tokens)
		parser.parse()
		parser.printExpr()
	} else {
		_, _ = fmt.Fprintln(os.Stderr, "Invalid command")
	}
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	if tokenError == SYNTAX_ERROR {
		os.Exit(65)
	} else if tokenError == RUNTIME_ERROR {
		os.Exit(70)
	}
	os.Exit(0)
	//lexer := NewLexer("78")
	//lexer.tokenize()
	//parser := NewParser(lexer.tokens)
	//parser.parse()
	//parser.printExpr()
}
