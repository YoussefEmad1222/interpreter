package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

var tokenError ErrorType

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func createToken(TypeOfToken TokenType, lexeme string, literal interface{}, line int) Token {
	return Token{TypeOfToken, lexeme, literal, line}
}

type Lexer struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func NewLexer(source string) *Lexer {
	return &Lexer{
		source:  source,
		tokens:  make([]Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}
func isAlpha(c byte) bool {
	return unicode.IsLetter(rune(c)) || c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

func isDigit(c byte) bool {
	return unicode.IsDigit(rune(c))
}

func (l *Lexer) isTheEnd() bool {
	return l.current >= len(l.source)
}
func (l *Lexer) advance() byte {
	l.current++
	return l.source[l.current-1]
}
func (l *Lexer) addToken(tokenType TokenType, literal interface{}) {
	text := l.source[l.start:l.current]
	l.tokens = append(l.tokens, createToken(tokenType, text, literal, l.line))
}
func (l *Lexer) match(expected byte) bool {
	if l.isTheEnd() {
		return false
	}
	if l.source[l.current] != expected {
		return false
	}
	l.current++
	return true
}
func (l *Lexer) peek() byte {
	if l.isTheEnd() {
		return 0
	}
	return l.source[l.current]
}
func peekNext(l *Lexer) byte {
	if l.current+1 >= len(l.source) {
		return 0
	}
	return l.source[l.current+1]
}

func (l *Lexer) skipComment() {
	for l.peek() != '\n' && !l.isTheEnd() {
		l.advance()
	}
	if l.isTheEnd() {
		return
	}
}

func (l *Lexer) skipWhitespace() {
	for l.peek() == ' ' || l.peek() == '\t' || l.peek() == '\r' && !l.isTheEnd() {
		if l.peek() == '\n' {
			l.line++
		}
		l.advance()
	}
}
func (l *Lexer) stringToken() {
	for l.peek() != '"' && !l.isTheEnd() {
		if l.peek() == '\n' {
			l.line++
		}
		l.advance()
	}
	if l.isTheEnd() {
		fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", l.line)
		tokenError = SYNTAX_ERROR
		return
	}
	l.advance()
	text := l.source[l.start+1 : l.current-1]
	l.addToken(STRING, text)
}
func (l *Lexer) numberToken() {
	for isDigit(l.peek()) {
		l.advance()
	}
	if l.peek() == '.' && isDigit(peekNext(l)) {
		l.advance()
		for isDigit(l.peek()) {
			l.advance()
		}
	}
	text := l.source[l.start:l.current]
	value, _ := strconv.ParseFloat(text, 64)
	l.addToken(NUMBER, value)
}
func (l *Lexer) identifierOrKeywordToken() {
	for isAlphaNumeric(l.peek()) && !l.isTheEnd() {
		l.advance()
	}
	text := l.source[l.start:l.current]
	if tokenType, ok := keywords[text]; ok {
		l.addToken(tokenType, nil)
	} else {
		l.addToken(IDENTIFIER, nil)
	}
}
func (l *Lexer) nextToken() {
	char := l.advance()
	switch char {
	case '(':
		l.addToken(LEFT_PAREN, nil)
	case ')':
		l.addToken(RIGHT_PAREN, nil)
	case '{':
		l.addToken(LEFT_BRACE, nil)
	case '}':
		l.addToken(RIGHT_BRACE, nil)
	case ',':
		l.addToken(COMMA, nil)
	case '.':
		l.addToken(DOT, nil)
	case '-':
		l.addToken(MINUS, nil)
	case '+':
		l.addToken(PLUS, nil)
	case ';':
		l.addToken(SEMICOLON, nil)
	case '*':
		l.addToken(STAR, nil)
	case '!':
		if l.match('=') {
			l.addToken(BANG_EQUAL, nil)
		} else {
			l.addToken(BANG, nil)
		}
	case '=':
		if l.match('=') {
			l.addToken(EQUAL_EQUAL, nil)
		} else {
			l.addToken(EQUAL, nil)
		}
	case '<':
		if l.match('=') {
			l.addToken(LESS_EQUAL, nil)
		} else {
			l.addToken(LESS, nil)
		}
	case '>':
		if l.match('=') {
			l.addToken(GREATER_EQUAL, nil)
		} else {
			l.addToken(GREATER, nil)
		}
	case '/':
		if l.match('/') {
			l.skipComment()
		} else {
			l.addToken(SLASH, nil)
		}
	case ' ', '\r', '\t':
		l.skipWhitespace()
	case '\n':
		l.line++
	case '"':
		l.stringToken()
	default:
		if isDigit(char) {
			l.numberToken()
		} else if isAlpha(char) {
			l.identifierOrKeywordToken()
		} else {
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", l.line, char)
			tokenError = SYNTAX_ERROR
		}
	}
}

func (l *Lexer) tokenize() {
	for !l.isTheEnd() {
		l.start = l.current
		l.nextToken()
	}
	l.tokens = append(l.tokens, createToken(EOF, "", nil, l.line))
}
func (l *Lexer) printTokens() {
	i := 0
	for i < len(l.tokens) {
		fmt.Printf("%s %s %s\n", l.tokens[i].Type, l.tokens[i].Lexeme, printInterface(l.tokens[i].Literal))
		i++
	}
}

func (l *Lexer) resetCurrent() {
	l.current = 0
	l.start = 0
	l.line = 1
}
func printInterface(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return "null"
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		if v == float64(int(v)) {
			return fmt.Sprintf("%.1f", v)
		}
		return fmt.Sprintf("%g", v)
	case string:
		return v
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return fmt.Sprintf("Unknown type: %T", v)
	}
}
