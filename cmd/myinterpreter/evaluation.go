package main

import (
	"fmt"
	"os"
)

type Environment struct {
	values map[string]Token
}

func newEnvironment() *Environment {
	return &Environment{values: make(map[string]Token)}
}
func (e *Environment) define(name string, value Token) {
	e.values[name] = value
}
func (e *Environment) get(name string) Token {
	if value, ok := e.values[name]; ok {
		return value
	}
	handleError("Undefined variable '" + name + "'")
	return Token{}
}
func isNumber(value interface{}) bool {
	_, ok := value.(float64)
	return ok
}

func isString(value interface{}) bool {
	_, ok := value.(string)
	return ok
}

func evaluate(lexer *Lexer) {
	// Create a Parser instance with the tokens from the lexer
	parser := NewParser(lexer.tokens)

	// Parse the tokens to build the AST
	parser.parse()

	// Check for syntax errors
	if tokenError == SYNTAX_ERROR {
		fmt.Fprintf(os.Stderr, "Syntax Error!\n")
		return
	}
	environment := newEnvironment()
	// Evaluate the AST and print the result
	result := evaluateStmts(parser.stmts, environment)
	for _, token := range result {
		printResult(token)
	}
}

func evaluateStmts(stmts []Stmt, environment *Environment) []Token {
	var results []Token
	for _, stmt := range stmts {
		ans := evaluateStmt(stmt, environment)
		if ans != (Token{}) {
			results = append(results, ans)
		}
	}
	return results
}

func evaluateStmt(stmt Stmt, environment *Environment) Token {
	switch s := stmt.(type) {
	case PrintStmt:
		return evaluateAST(s.expr, environment)
	case ExprStmt:
		return evaluateAST(s.expr, environment)
	case VarStmt:
		value := evaluateAST(s.value, environment)
		environment.define(s.name.Lexeme, value)
		return value
	case BlockStmt:
		return Token{}
	default:
		return Token{}
	}
}

func evaluateAST(expr Expr, environment *Environment) Token {
	switch e := expr.(type) {
	case Binary:
		left := evaluateAST(e.left, environment)
		right := evaluateAST(e.right, environment)
		return applyBinaryOperator(e.operator, left, right)
	case Unary:
		right := evaluateAST(e.right, environment)
		return applyUnaryOperator(e.operator, right)
	case Variable:
		return environment.get(e.name.Lexeme)
	case Literal:
		return e.token
	case Assign:
		environment.get(e.name.Lexeme) //check if variable exists in the environment
		value := evaluateAST(e.value, environment)
		environment.define(e.name.Lexeme, value)
		return value
	case Grouping:
		return evaluateAST(e.expr, environment)
	default:
		return Token{}
	}
}

func applyUnaryOperator(operator Token, right Token) Token {
	switch operator.Type {
	case BANG:
		if right.Type == NIL {
			return Token{Type: TRUE, Literal: nil, Lexeme: "true"}
		}
		if isBoolean(right) {
			return Token{Type: negateBoolean(right), Lexeme: negateLexeme(right)}
		}
		if isNumber(right.Literal) {
			return Token{Type: FALSE, Literal: nil, Lexeme: "false"}
		}
		handleError("Invalid operand for NOT operator")
	case MINUS:
		if isNumber(right.Literal) {
			return Token{Type: NUMBER, Literal: -right.Literal.(float64), Lexeme: "-" + right.Lexeme}
		}
		handleError("Operand must be a number")
	default:
		handleError("Invalid operator")
	}
	return Token{}
}

func applyBinaryOperator(operator Token, left Token, right Token) Token {
	switch operator.Type {
	case PLUS:
		if validStringOrNumberOperands(left, right) {
			if isString(left.Literal) || isString(right.Literal) {
				return Token{Type: STRING, Literal: concatenateStrings(left, right), Lexeme: concatenateLexemes(left, right)}
			}
			return Token{Type: NUMBER, Literal: left.Literal.(float64) + right.Literal.(float64), Lexeme: fmt.Sprintf("%v", left.Literal.(float64)+right.Literal.(float64))}
		}
		handleError("Operand must be a number or a string")
	case MINUS, STAR, SLASH:
		if isNumberOperands(left, right) {
			return applyNumericOperator(operator.Type, left, right)
		}
		handleError("Operand must be a number")
	case GREATER, GREATER_EQUAL, LESS, LESS_EQUAL:
		if isNumberOperands(left, right) {
			return compareNumbers(operator.Type, left, right)
		}
		handleError("Operand must be a number")
	case EQUAL_EQUAL, BANG_EQUAL:
		return compareEquality(operator.Type, left, right)
	default:
		return Token{}
	}
	return Token{}
}

func isBoolean(token Token) bool {
	return token.Type == TRUE || token.Type == FALSE
}

func negateBoolean(token Token) TokenType {
	if token.Type == TRUE {
		return FALSE
	}
	return TRUE
}

func negateLexeme(token Token) string {
	if token.Type == TRUE {
		return "false"
	}
	return "true"
}

func validStringOrNumberOperands(left Token, right Token) bool {
	return isNumberOperands(left, right) || isStringOperands(left, right)
}

func isNumberOperands(left Token, right Token) bool {
	return isNumber(left.Literal) && isNumber(right.Literal)
}

func isStringOperands(left Token, right Token) bool {
	return isString(left.Literal) && isString(right.Literal)
}

func applyNumericOperator(op TokenType, left Token, right Token) Token {
	switch op {
	case MINUS:
		return Token{Type: NUMBER, Literal: left.Literal.(float64) - right.Literal.(float64), Lexeme: fmt.Sprintf("%v", left.Literal.(float64)-right.Literal.(float64))}
	case STAR:
		return Token{Type: NUMBER, Literal: left.Literal.(float64) * right.Literal.(float64), Lexeme: fmt.Sprintf("%v", left.Literal.(float64)*right.Literal.(float64))}
	case SLASH:
		if right.Literal.(float64) == 0 {
			handleError("Division by zero")
		}
		return Token{Type: NUMBER, Literal: left.Literal.(float64) / right.Literal.(float64), Lexeme: fmt.Sprintf("%v", left.Literal.(float64)/right.Literal.(float64))}
	}
	return Token{}
}

func compareNumbers(op TokenType, left Token, right Token) Token {
	switch op {
	case GREATER:
		return booleanToken(left.Literal.(float64) > right.Literal.(float64))
	case GREATER_EQUAL:
		return booleanToken(left.Literal.(float64) >= right.Literal.(float64))
	case LESS:
		return booleanToken(left.Literal.(float64) < right.Literal.(float64))
	case LESS_EQUAL:
		return booleanToken(left.Literal.(float64) <= right.Literal.(float64))
	}
	return Token{}
}

func printResult(token Token) {
	switch token.Type {
	case STRING, NUMBER:
		fmt.Println(token.Literal)
	case TRUE, FALSE, NIL:
		fmt.Println(token.Lexeme)
	default:
		fmt.Println("nil")
	}
}

func compareEquality(op TokenType, left Token, right Token) Token {
	switch op {
	case EQUAL_EQUAL:
		if (left.Type == TRUE || left.Type == FALSE) && (right.Type == TRUE || right.Type == FALSE) {
			return booleanToken(left.Type == right.Type)
		}
		return booleanToken(left.Literal == right.Literal)
	case BANG_EQUAL:
		if (left.Type == TRUE || left.Type == FALSE) && (right.Type == TRUE || right.Type == FALSE) {
			return booleanToken(left.Type != right.Type)
		}
		return booleanToken(left.Literal != right.Literal)
	}
	return Token{}
}

func booleanToken(condition bool) Token {
	if condition {
		return Token{Type: TRUE, Literal: nil, Lexeme: "true"}
	}
	return Token{Type: FALSE, Literal: nil, Lexeme: "false"}
}

func handleError(message string) {
	fmt.Fprintf(os.Stderr, message+"\n")
	tokenError = RUNTIME_ERROR
	os.Exit(70)
}

func concatenateStrings(left Token, right Token) string {
	return left.Literal.(string) + right.Literal.(string)
}

func concatenateLexemes(left Token, right Token) string {
	return "\"" + left.Literal.(string) + right.Literal.(string) + "\""
}
