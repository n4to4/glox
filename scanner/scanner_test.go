package scanner

import (
	"fmt"
	"testing"

	"github.com/n4to4/glox/tokens"
)

func TestScanner(t *testing.T) {
	cases := []struct {
		source   string
		expected []string
	}{
		{"(", []string{tokens.LEFT_PAREN, tokens.EOF}},
		{",.-", []string{tokens.COMMA, tokens.DOT, tokens.MINUS, tokens.EOF}},
		{"== != >= <=", []string{
			tokens.EQUAL_EQUAL,
			tokens.BANG_EQUAL,
			tokens.GREATER_EQUAL,
			tokens.LESS_EQUAL,
			tokens.EOF,
		}},
		{`"string"!=`, []string{tokens.STRING, tokens.BANG_EQUAL, tokens.EOF}},
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
	if tok.Ttype != tokens.STRING {
		t.Errorf("token type want %q got %q", tokens.STRING, tok.Ttype)
	}

	wantLiteral := "string"
	if tok.Literal != wantLiteral {
		t.Errorf("literal want %q got %q", wantLiteral, tok.Literal)
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
			if tok.Ttype != tokens.NUMBER {
				t.Errorf("token type want %q got %q", tokens.NUMBER, tok.Ttype)
			}

			if num, ok := tok.Literal.(float64); !ok || num != cc.expected {
				t.Errorf("want %f got %f", cc.expected, num)
			}
		})
	}
}

func TestIdentifier(t *testing.T) {
	cases := []struct {
		source string
		lexeme string
		ttype  string
	}{
		{"id", "id", tokens.IDENTIFIER},
		{"name", "name", tokens.IDENTIFIER},
		{"or", "or", tokens.OR},
		{" if ", "if", tokens.IF},
	}

	for _, cc := range cases {
		t.Run(cc.source, func(t *testing.T) {
			toks := NewScanner(cc.source).ScanTokens()

			if len(toks) != 2 {
				t.Fatalf("len want %d got %d", 2, len(toks))
			}

			tok := toks[0]
			if tok.Ttype != cc.ttype {
				t.Errorf("token type want %q got %q", cc.ttype, tok.Ttype)
			}

			if tok.Lexeme != cc.lexeme {
				t.Errorf("lexeme want %q got %q", cc.lexeme, tok.Lexeme)
			}
		})
	}
}

func assertTokenTypes(t *testing.T, toks []tokens.Token, ttypes ...string) {
	t.Helper()

	if len(toks) != len(ttypes) {
		t.Fatalf("two lengths don't match: want %d got %d", len(ttypes), len(toks))
	}

	for i, tok := range toks {
		if tok.Ttype != ttypes[i] {
			t.Errorf("want token type %q, got %q", ttypes[i], tok.Ttype)
		}
	}
}
