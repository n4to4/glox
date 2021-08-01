package scanner

import (
	"github.com/n4to4/glox/error"
	"github.com/n4to4/glox/tokens"
)

type Scanner struct {
	source string
	tokens []tokens.Token

	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []tokens.Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() []tokens.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, tokens.NewToken(tokens.EOF, "", "", s.line))
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	case "(":
		s.addToken(tokens.LEFT_PAREN, "(")
	case ")":
		s.addToken(tokens.RIGHT_PAREN, ")")
	case "{":
		s.addToken(tokens.LEFT_BRACE, "{")
	case "}":
		s.addToken(tokens.RIGHT_BRACE, "}")
	case ",":
		s.addToken(tokens.COMMA, ",")
	case ".":
		s.addToken(tokens.DOT, ".")
	case "-":
		s.addToken(tokens.MINUS, "-")
	case "+":
		s.addToken(tokens.PLUS, "+")
	case ";":
		s.addToken(tokens.SEMICOLON, ";")
	case "*":
		s.addToken(tokens.STAR, "*")
	case "!":
		if s.match("=") {
			s.addToken(tokens.BANG_EQUAL, "!=")
		} else {
			s.addToken(tokens.BANG, "!")
		}
	case "=":
		if s.match("=") {
			s.addToken(tokens.EQUAL_EQUAL, "==")
		} else {
			s.addToken(tokens.EQUAL, "=")
		}
	case "<":
		if s.match("=") {
			s.addToken(tokens.LESS_EQUAL, "<=")
		} else {
			s.addToken(tokens.LESS, "<")
		}
	case ">":
		if s.match("=") {
			s.addToken(tokens.GREATER_EQUAL, ">=")
		} else {
			s.addToken(tokens.GREATER, ">")
		}
	case "/":
		if s.match("/") {
			for s.peek() != "\n" && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(tokens.SLASH, "/")
		}
	case " ", "\r", "\t":
	case "\n":
		s.line++
	case `"`:
		s.scanString()
	default:
		error.ErrorReport(s.line, "Unexpected character.")
	}
}

func (s *Scanner) advance() string {
	c := s.source[s.current : s.current+1]
	s.current++
	return c
}

func (s *Scanner) addToken(ttype, literal string) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, tokens.NewToken(ttype, text, literal, s.line))
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
		error.ErrorReport(s.line, "Unterminated string.")
		return
	}

	// the closing ".
	s.advance()

	// Trim the surrounding quotes.
	value := s.source[s.start+1 : s.current-1]
	s.addToken(tokens.STRING, value)
}
