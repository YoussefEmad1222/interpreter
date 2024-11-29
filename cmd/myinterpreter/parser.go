package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

// statement: printStmt | expressionStmt | varStmt | blockStmt ;
// varStmt: "var" IDENTIFIER ( "=" expression )? ";"
// printStmt: "print" expression ";"
// expressionStmt: expression ";"
// expression: equality
// equality: comparison ( ( "!=" | "==" ) comparison )*
// comparison: term ( ( ">" | ">=" | "<" | "<=" ) term )*
// term: factor ( ( "-" | "+" ) factor )*
// factor: unary ( ( "/" | "*" ) unary )*
// unary: ( "!" | "-" ) unary | primary
// primary: NUMBER | STRING | "false" | "true" | "nil" | "(" expression ")" | IDENTIFIER;

type Parser struct {
	tokens  []Token
	current int
	stmts   []Stmt
}
type Stmt interface{}
type PrintStmt struct {
	expr Expr
}
type VarStmt struct {
	name  Token
	value Expr
}
type ExprStmt struct {
	expr Expr
}
type BlockStmt struct {
	stmts []Stmt
}

func (b *BlockStmt) getStmts() []Stmt {
	return b.stmts
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
	token Token
}
type Grouping struct {
	expr Expr
}

type Assign struct {
	name  Token
	value Expr
}
type Variable struct {
	name Token
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens) || p.tokens[p.current].Type == EOF || p.tokens[p.current].Type == SEMICOLON
}

func (p *Parser) advance() {
	if !p.isAtEnd() {
		p.current++
	}
}
func (p *Parser) consume(tokenType TokenType, errorMsg string) Token {
	peekType := p.tokens[p.current].Type
	if isRunning {
		if peekType == tokenType {
			p.current++
			return p.tokens[p.current-1]
		}
	} else {
		if peekType == tokenType || peekType == EOF {
			p.current++
			return p.tokens[p.current-1]
		}
	}
	fmt.Fprintf(os.Stderr, errorMsg)
	tokenError = SYNTAX_ERROR
	return Token{}
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
	if p.peek().Type == SEMICOLON {
		return nil
	}
	if p.peek().Type == IDENTIFIER {
		p.advance()
		return Variable{p.previous()}
	}
	if operand[p.peek().Type] {
		p.advance()
		return Literal{p.previous()}
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
func (p *Parser) parseAssignment() Expr {
	expr := p.parseEquality()
	for p.peek().Type == EQUAL {
		//operator := p.peek()
		p.advance()
		right := p.parseAssignment()
		if _, ok := expr.(Variable); ok {
			return Assign{expr.(Variable).name, right}
		}
		fmt.Fprintf(os.Stderr, "Invalid assignment target")
		tokenError = SYNTAX_ERROR
		return nil
	}
	return expr
}

func (p *Parser) parseExpression() Expr {
	return p.parseAssignment()
}

func (p *Parser) parsePrintStmt() Stmt {
	p.advance()
	expr := p.parseExpression()
	p.consume(SEMICOLON, "Expect ';' after value")
	return PrintStmt{expr}
}
func (p *Parser) parseExpressionStmt() Stmt {
	expr := p.parseExpression()
	p.consume(SEMICOLON, "Expect ';' after expression")
	return ExprStmt{expr}
}
func (p *Parser) parseVarStmt() Stmt {
	p.advance()
	name := p.consume(IDENTIFIER, "Expect variable name")
	var value Expr
	if p.peek().Type == EQUAL {
		p.advance()
		value = p.parseExpression()
	}
	p.consume(SEMICOLON, "Expect ';' after value")
	return VarStmt{name, value}
}
func (p *Parser) parseBlockStmt() Stmt {
	var statements []Stmt
	p.advance()
	for (p.peek().Type != RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.parseDeclaration())
	}
	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return BlockStmt{stmts: statements}
}
func (p *Parser) parseStmt() Stmt {
	if p.peek().Type == PRINT {
		return p.parsePrintStmt()
	}
	if p.peek().Type == LEFT_BRACE {
		return p.parseBlockStmt()
	}
	return p.parseExpressionStmt()
}
func (p *Parser) parseDeclaration() Stmt {
	if p.peek().Type == VAR {
		return p.parseVarStmt()
	}
	return p.parseStmt()
}

func (p *Parser) parse() {
	var statements []Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.parseDeclaration())
	}
	p.stmts = statements
}

func (p *Parser) printExpr() {
	var strBuilder strings.Builder
	for _, stmt := range p.stmts {
		//convert stmt to exprStmt or printStmt
		strBuilder.Reset()
		switch s := stmt.(type) {
		case PrintStmt:
			strBuilder.WriteString("print ")
			printAST(s.expr, &strBuilder)
		case ExprStmt:
			printAST(s.expr, &strBuilder)
		case VarStmt:
			strBuilder.WriteString("var ")
			strBuilder.WriteString(s.name.Lexeme)
			if s.value != nil {
				strBuilder.WriteString(" = ")
				printAST(s.value, &strBuilder)
			}
		}
		if tokenError == SYNTAX_ERROR {
			return
		}
		fmt.Println(strBuilder.String())
	}
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
		if e.token.Type == NUMBER {
			num := e.token.Literal.(float64)
			if math.Floor(num) == num {
				s.WriteString(fmt.Sprintf("%.1f", num))
			} else {
				s.WriteString(fmt.Sprintf("%v", num))
			}
		} else if e.token.Type == STRING {
			s.WriteString(e.token.Literal.(string))
		} else {
			s.WriteString(e.token.Lexeme)
		}
	case Grouping:
		s.WriteString("(")
		s.WriteString("group ")
		printAST(e.expr, s)
		s.WriteString(")")
	}
}
