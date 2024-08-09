package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

// Parser grammar for the parser
// expression     → equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )*  ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )*;
// unary          → ( "!" | "-" ) unary | primary ;
// primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
//
// The parser is responsible for taking the tokens from the lexer and turning them into an abstract syntax tree (AST).
type Parser struct {
	tokens  []Token
	current int
	expr    Expr
}
type Expr interface{}
type Binary struct {
	left     Expr
	operator Token
	right    Expr
}
type Unary struct {
	operator Token
	right    Expr
}
type Literal struct {
	value interface{}
}
type Grouping struct {
	expr Expr
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}

func (p *Parser) advance() {
	if !p.isAtEnd() {
		p.current++
	}
}

func (p *Parser) peek() Token {
	if p.isAtEnd() {
		return Token{}
	}
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	if p.current == 0 {
		return Token{}
	}
	return p.tokens[p.current-1]
}

func (p *Parser) parsePrimary() Expr {
	if operand[p.peek().Type] {
		p.advance()
		if p.previous().Type == NUMBER || p.previous().Type == STRING {
			return Literal{p.previous().Literal}
		}
		return Literal{p.previous().Lexeme}
	}
	if p.peek().Type == LEFT_PAREN {
		p.advance()
		expr := p.parseExpression()
		if p.peek().Type != RIGHT_PAREN {
			fmt.Fprintf(os.Stderr, "Expect ')' after expression")
			tokenError = SYNTAX_ERROR
			return nil
		}
		p.advance()
		return Grouping{expr}
	}
	tokenError = SYNTAX_ERROR
	return nil
}
func (p *Parser) parseUnary() Expr {
	if p.peek().Type == BANG || p.peek().Type == MINUS {
		operator := p.peek()
		p.advance()
		right := p.parseUnary()
		return Unary{operator, right}
	}
	return p.parsePrimary()
}

func (p *Parser) parseFactor() Expr {
	expr := p.parseUnary()
	for p.peek().Type == STAR || p.peek().Type == SLASH {
		operator := p.peek()
		p.advance()
		right := p.parseUnary()
		expr = Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) parseTerm() Expr {
	expr := p.parseFactor()
	for p.peek().Type == PLUS || p.peek().Type == MINUS {
		operator := p.peek()
		p.advance()
		right := p.parseFactor()
		expr = Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) parseComparison() Expr {
	expr := p.parseTerm()
	for p.peek().Type == GREATER || p.peek().Type == GREATER_EQUAL || p.peek().Type == LESS || p.peek().Type == LESS_EQUAL {
		operator := p.peek()
		p.advance()
		right := p.parseTerm()
		expr = Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) parseEquality() Expr {
	expr := p.parseComparison()
	for p.peek().Type == BANG_EQUAL || p.peek().Type == EQUAL_EQUAL {
		operator := p.peek()
		p.advance()
		right := p.parseComparison()
		expr = Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) parseExpression() Expr {
	return p.parseEquality()
}
func (p *Parser) parse() {
	p.expr = p.parseExpression()
}

func (p *Parser) printExpr() {
	var strBuilder strings.Builder
	printAST(p.expr, &strBuilder)
	if tokenError == SYNTAX_ERROR {
		return
	}
	fmt.Println(strBuilder.String())
}

func printAST(expr Expr, s *strings.Builder) {
	if tokenError == SYNTAX_ERROR {
		return
	}

	switch e := expr.(type) {
	case Binary:
		s.WriteString("(")
		s.WriteString(e.operator.Lexeme + " ")
		printAST(e.left, s)
		s.WriteString(" ")
		printAST(e.right, s)
		s.WriteString(")")
	case Unary:
		s.WriteString("(")
		s.WriteString(e.operator.Lexeme + " ")
		printAST(e.right, s)
		s.WriteString(")")
	case Literal:
		if isNumber(e.value) {
			if math.Floor(e.value.(float64)) == e.value.(float64) {
				s.WriteString(fmt.Sprintf("%.1f", e.value.(float64)))
			} else {
				s.WriteString(fmt.Sprintf("%v", e.value))
			}
		} else {
			s.WriteString(e.value.(string))
		}
	case Grouping:
		s.WriteString("(")
		s.WriteString("group ")
		printAST(e.expr, s)
		s.WriteString(")")
	}

}
