package main

import (
	"fmt"
	"os"
)

func isNumber(value interface{}) bool {
	_, ok := value.(float64)
	return ok
}

func evaluate(lexer *Lexer) {
	lexer.resetCurrent()
	postFixLexer := NewLexer("")
	answer := NewLexer("")
	infixToPostfix(lexer, postFixLexer)
	evaluatePostfix(postFixLexer, answer)
	if answer.tokens[0].Type == STRING || answer.tokens[0].Type == NUMBER {
		fmt.Println(answer.tokens[0].Literal)
	} else {
		fmt.Println(answer.tokens[0].Lexeme)
	}
}

func evaluatePostfix(lexer *Lexer, answer *Lexer) {
	var stack []Token
	for i := 0; i < len(lexer.tokens); i++ {
		token := lexer.tokens[i]
		switch token.Type {
		case NUMBER, STRING, IDENTIFIER, TRUE, FALSE, NIL:
			stack = append(stack, token)
		default:
			operand2 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			operand1 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			result, opType := applyOperator(token.Type, operand1.Literal, operand2.Literal)
			tokenType := NUMBER
			if opType == NO_OP {
				fmt.Fprintf(os.Stderr, "Invalid operator\n")
				os.Exit(65)
			}
			if opType == STRING_OP {
				tokenType = STRING
			} else if opType == LOGICAL_OP {
				tokenType = TRUE
				if result == false {
					tokenType = FALSE
				}
			}
			lexeme := ""
			if tokenType == STRING {
				lexeme = fmt.Sprintf("\"%v\"", result)
			} else if tokenType == TRUE || tokenType == FALSE {
				lexeme = fmt.Sprintf("%v", result)
				result = nil
			} else {
				lexeme = fmt.Sprintf("%g", result)
			}
			stack = append(stack, Token{Type: tokenType, Lexeme: lexeme, Literal: result})
		}
	}
	answer.tokens = stack
}
func applyOperator(op TokenType, operand1, operand2 interface{}) (interface{}, operationType) {
	valid := checkOperands(op, operand1, operand2)
	if !valid {
		fmt.Fprintf(os.Stderr, "Operand must be a number.\n")
		tokenError = RUNTIME_ERROR
		os.Exit(70)
	}
	switch op {
	case PLUS:
		str1, isStr1 := operand1.(string)
		str2, isStr2 := operand2.(string)
		if isStr1 && isStr2 {
			return str1 + str2, STRING_OP
		}
		return operand1.(float64) + operand2.(float64), NUMBER_OP
	case MINUS:
		return operand1.(float64) - operand2.(float64), NUMBER_OP
	case STAR:
		return operand1.(float64) * operand2.(float64), NUMBER_OP
	case SLASH:
		if operand2.(float64) == 0 {
			_, _ = fmt.Fprintln(os.Stderr, "Division by zero")
			os.Exit(65)
		}
		return operand1.(float64) / operand2.(float64), NUMBER_OP
	case GREATER:
		if isNumber(operand1) && isNumber(operand2) {
			return operand1.(float64) > operand2.(float64), LOGICAL_OP
		}
	case GREATER_EQUAL:
		if isNumber(operand1) && isNumber(operand2) {
			return operand1.(float64) >= operand2.(float64), LOGICAL_OP
		}
	case LESS:
		if isNumber(operand1) && isNumber(operand2) {
			return operand1.(float64) < operand2.(float64), LOGICAL_OP
		}
	case LESS_EQUAL:
		if isNumber(operand1) && isNumber(operand2) {
			return operand1.(float64) <= operand2.(float64), LOGICAL_OP
		}
	case EQUAL_EQUAL:
		if isNumber(operand1) && isNumber(operand2) {
			return operand1.(float64) == operand2.(float64), LOGICAL_OP
		}
		// For strings
		str1, isStr1 := operand1.(string)
		str2, isStr2 := operand2.(string)
		if isStr1 && isStr2 {
			return str1 == str2, LOGICAL_OP
		}
		return false, LOGICAL_OP
	case BANG_EQUAL:
		if isNumber(operand1) && isNumber(operand2) {
			return operand1.(float64) != operand2.(float64), LOGICAL_OP
		}
		// For strings
		str1, isStr1 := operand1.(string)
		str2, isStr2 := operand2.(string)
		if isStr1 && isStr2 {
			return str1 != str2, LOGICAL_OP
		}
		return true, LOGICAL_OP
	default:
		return nil, NO_OP
	}
	return nil, NO_OP
}

func checkOperands(op TokenType, operand1 interface{}, operand2 interface{}) bool {

	_, isNum1 := operand1.(float64)
	_, isNum2 := operand2.(float64)
	_, isStr1 := operand1.(string)
	_, isStr2 := operand2.(string)
	if op == EQUAL_EQUAL || op == BANG_EQUAL {
		return true
	} else if op == PLUS {
		return isStr1 && isStr2 || isNum1 && isNum2
	} else {
		return isNum1 && isNum2
	}
}

func getNotOperands(lexer *Lexer, i *int) []Token {
	var notAnswer []Token
	for *i < len(lexer.tokens)-1 {
		token := lexer.tokens[*i]
		switch token.Type {
		case BANG:
			*i++
			notAnswer = getNotOperands(lexer, i)
			notAnswer[0] = negateToken(notAnswer[0])
			return notAnswer
		case LEFT_PAREN:
			*i++
			notAnswer = processParentheses(lexer, i)
			return notAnswer
		default:
			notAnswer = append(notAnswer, token)
			*i++
			return notAnswer
		}
	}
	return notAnswer
}

func negateToken(token Token) Token {
	switch token.Type {
	case TRUE:
		token.Type = FALSE
		token.Lexeme = "false"
	case FALSE:
		token.Type = TRUE
		token.Lexeme = "true"
	case NIL:
		token.Type = TRUE
		token.Lexeme = "true"
	default:
		if isNumber(token.Literal) {
			if token.Literal.(float64) != 0.0 {
				token.Type = FALSE
				token.Lexeme = "false"
			} else {
				token.Type = TRUE
				token.Lexeme = "true"
			}
		} else {
			fmt.Fprintf(os.Stderr, "Invalid operand for NOT operator\n")
			os.Exit(65)
		}
	}
	token.Literal = nil
	return token
}

func processParentheses(lexer *Lexer, i *int) []Token {
	var stack []Token
	for *i < len(lexer.tokens)-1 && lexer.tokens[*i].Type != RIGHT_PAREN {
		stack = append(stack, lexer.tokens[*i])
		*i++
	}
	*i++ // Skip the RIGHT_PAREN
	answer := NewLexer("")
	tokens := append(stack, createToken(EOF, "", nil, lexer.line))
	infixToPostfix(&Lexer{tokens: tokens}, answer)
	evaluatePostfix(answer, answer)
	return []Token{answer.tokens[0]}
}

func infixToPostfix(lexer *Lexer, postFixLexer *Lexer) {
	var stack []Token
	var postfix []Token
	var lastTokenType TokenType
	var token Token
	var i = 0
	for i < len(lexer.tokens)-1 {
		token = lexer.tokens[i]
		switch token.Type {
		case NUMBER, STRING, IDENTIFIER, TRUE, FALSE, NIL:
			postfix = append(postfix, token)
		case LEFT_PAREN:
			stack = append(stack, token)
		case RIGHT_PAREN:
			for len(stack) > 0 && stack[len(stack)-1].Type != LEFT_PAREN {
				postfix = append(postfix, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 && stack[len(stack)-1].Type == LEFT_PAREN {
				stack = stack[:len(stack)-1]
			}
		default:
			prevIsOperator := lastTokenType == "" || operators[lastTokenType]

			if token.Type == MINUS && (prevIsOperator) {
				postfix = append(postfix, Token{Type: NUMBER, Lexeme: "-1", Literal: float64(-1)})
				stack = append(stack, Token{Type: STAR, Lexeme: "*", Literal: nil})
			} else if token.Type == BANG {
				notAnswer := getNotOperands(lexer, &i)
				stack = append(stack, notAnswer[0])
				i--
			} else {
				for len(stack) > 0 && precedence(stack[len(stack)-1].Type) >= precedence(token.Type) {
					postfix = append(postfix, stack[len(stack)-1])
					stack = stack[:len(stack)-1]
				}
				stack = append(stack, token)
			}

		}
		lastTokenType = token.Type
		i++
	}
	for len(stack) > 0 {
		postfix = append(postfix, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	postFixLexer.tokens = postfix
}
func precedence(op TokenType) int {
	switch op {
	case PLUS, MINUS:
		return 2
	case STAR, SLASH:
		return 3
	case GREATER, GREATER_EQUAL, LESS, LESS_EQUAL, EQUAL_EQUAL, BANG_EQUAL:
		return 1
	default:
		return 0
	}
}
