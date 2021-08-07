package main

import (
	"strconv"
)

type Scanner struct {
	source string
	tokens []Token

	start   int
	current int
	line    int

	keywords map[string]TokenType
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:   source,
		tokens:   []Token{},
		start:    0,
		current:  0,
		line:     1,
		keywords: newKeywords(),
	}
}

func newKeywords() map[string]TokenType {
	var ks = make(map[string]TokenType)

	ks["and"] = AND
	ks["class"] = CLASS
	ks["else"] = ELSE
	ks["false"] = FALSE
	ks["for"] = FOR
	ks["fun"] = FUN
	ks["if"] = IF
	ks["nil"] = NIL
	ks["or"] = OR
	ks["print"] = PRINT
	ks["return"] = RETURN
	ks["super"] = SUPER
	ks["this"] = THIS
	ks["true"] = TRUE
	ks["var"] = VAR
	ks["while"] = WHILE

	return ks
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, NewToken(EOF, "", "", s.line))
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	case "(":
		s.addToken(LEFT_PAREN, "(")
	case ")":
		s.addToken(RIGHT_PAREN, ")")
	case "{":
		s.addToken(LEFT_BRACE, "{")
	case "}":
		s.addToken(RIGHT_BRACE, "}")
	case ",":
		s.addToken(COMMA, ",")
	case ".":
		s.addToken(DOT, ".")
	case "-":
		s.addToken(MINUS, "-")
	case "+":
		s.addToken(PLUS, "+")
	case ";":
		s.addToken(SEMICOLON, ";")
	case "*":
		s.addToken(STAR, "*")
	case "!":
		if s.match("=") {
			s.addToken(BANG_EQUAL, "!=")
		} else {
			s.addToken(BANG, "!")
		}
	case "=":
		if s.match("=") {
			s.addToken(EQUAL_EQUAL, "==")
		} else {
			s.addToken(EQUAL, "=")
		}
	case "<":
		if s.match("=") {
			s.addToken(LESS_EQUAL, "<=")
		} else {
			s.addToken(LESS, "<")
		}
	case ">":
		if s.match("=") {
			s.addToken(GREATER_EQUAL, ">=")
		} else {
			s.addToken(GREATER, ">")
		}
	case "/":
		if s.match("/") {
			for s.peek() != "\n" && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH, "/")
		}
	case " ", "\r", "\t":
	case "\n":
		s.line++
	case `"`:
		s.scanString()
	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			ErrorReport(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) advance() string {
	c := s.source[s.current : s.current+1]
	s.current++
	return c
}

func (s *Scanner) addToken(ttype TokenType, literal interface{}) {
	lexeme := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(ttype, lexeme, literal, s.line))
}

func (s *Scanner) match(expected string) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current:s.current+1] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() string {
	if s.isAtEnd() {
		return ""
	} else {
		return s.source[s.current : s.current+1]
	}
}

func (s *Scanner) peekNext() string {
	if s.current+1 >= len(s.source) {
		return ""
	}
	return s.source[s.current+1 : s.current+2]
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanString() {
	for s.peek() != `"` && !s.isAtEnd() {
		if s.peek() == "\n" {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		ErrorReport(s.line, "Unterminated string.")
		return
	}

	// the closing ".
	s.advance()

	// Trim the surrounding quotes.
	value := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, value)
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == "." && isDigit(s.peekNext()) {
		s.advance()
	}

	for isDigit(s.peek()) {
		s.advance()
	}

	lexeme := s.source[s.start:s.current]
	f, _ := strconv.ParseFloat(lexeme, 64)
	s.addToken(NUMBER, f)
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	ttype, ok := s.keywords[text]
	if !ok {
		ttype = IDENTIFIER
	}
	s.addToken(ttype, nil)
}

func isDigit(str string) bool {
	if len(str) != 1 {
		return false
	}

	c := str[0]
	return '0' <= c && c <= '9'
}

func isAlpha(str string) bool {
	if len(str) != 1 {
		return false
	}

	c := str[0]
	return c == '_' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

func isAlphaNumeric(str string) bool {
	return isDigit(str) || isAlpha(str)
}
