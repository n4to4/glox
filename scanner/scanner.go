package scanner

import "github.com/n4to4/glox/tokens"

type Scanner struct {
	source string
	tokens []tokens.Token

	start   int
	current int
	line    int
}

func NewScanner(source string) Scanner {
	return Scanner{
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
	return nil
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case "(":
		s.addToken(tokens.LEFT_PAREN, "(")
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

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}
