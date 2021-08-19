package main

import (
	"fmt"
	"testing"
)

func TestScanner(t *testing.T) {
	cases := []struct {
		source   string
		expected []TokenType
	}{
		{"(", []TokenType{LEFT_PAREN, EOF}},
		{",.-", []TokenType{COMMA, DOT, MINUS, EOF}},
		{"== != >= <=", []TokenType{
			EQUAL_EQUAL,
			BANG_EQUAL,
			GREATER_EQUAL,
			LESS_EQUAL,
			EOF,
		}},
		{`"string"!=`, []TokenType{STRING, BANG_EQUAL, EOF}},
	}

	for _, cc := range cases {
		toks := NewScanner(cc.source).ScanTokens()
		t.Run(fmt.Sprintf("with source %q", cc.source), func(t *testing.T) {
			assertTokenTypes(t, toks, cc.expected...)
		})
	}
}

func TestScanString(t *testing.T) {
	source := `"string"`
	toks := NewScanner(source).ScanTokens()

	if len(toks) != 2 {
		t.Fatalf("len want %d got %d", 2, len(toks))
	}

	tok := toks[0]
	if tok.ttype != STRING {
		t.Errorf("token type want %q got %q", STRING, tok.ttype)
	}

	wantLiteral := "string"
	if tok.literal != wantLiteral {
		t.Errorf("literal want %q got %q", wantLiteral, tok.literal)
	}
}

func TestScanNumber(t *testing.T) {
	cases := []struct {
		source   string
		expected float64
	}{
		{"123", 123},
		{"3.14", 3.14},
	}

	for _, cc := range cases {
		t.Run(cc.source, func(t *testing.T) {
			toks := NewScanner(cc.source).ScanTokens()

			if len(toks) != 2 {
				t.Fatalf("len want %d got %d", 2, len(toks))
			}

			tok := toks[0]
			if tok.ttype != NUMBER {
				t.Errorf("token type want %q got %q", NUMBER, tok.ttype)
			}

			if num, ok := tok.literal.(float64); !ok || num != cc.expected {
				t.Errorf("want %f got %f", cc.expected, num)
			}
		})
	}
}

func TestIdentifier(t *testing.T) {
	cases := []struct {
		source string
		lexeme string
		ttype  TokenType
	}{
		{"id", "id", IDENTIFIER},
		{"name", "name", IDENTIFIER},
		{"or", "or", OR},
		{" if ", "if", IF},
	}

	for _, cc := range cases {
		t.Run(cc.source, func(t *testing.T) {
			toks := NewScanner(cc.source).ScanTokens()

			if len(toks) != 2 {
				t.Fatalf("len want %d got %d", 2, len(toks))
			}

			tok := toks[0]
			if tok.ttype != cc.ttype {
				t.Errorf("token type want %q got %q", cc.ttype, tok.ttype)
			}

			if tok.lexeme != cc.lexeme {
				t.Errorf("lexeme want %q got %q", cc.lexeme, tok.lexeme)
			}
		})
	}
}

func assertTokenTypes(t *testing.T, toks []Token, ttypes ...TokenType) {
	t.Helper()

	if len(toks) != len(ttypes) {
		t.Fatalf("two lengths don't match: want %d got %d", len(ttypes), len(toks))
	}

	for i, tok := range toks {
		if tok.ttype != ttypes[i] {
			t.Errorf("want token type %q, got %q", ttypes[i], tok.ttype)
		}
	}
}
