package scanner

import (
	"testing"

	"github.com/n4to4/glox/tokens"
)

func TestScanner(t *testing.T) {
	cases := []struct {
		source   string
		expected []string
	}{
		{"(", []string{tokens.LEFT_PAREN, tokens.EOF}},
	}

	for _, cc := range cases {
		s := NewScanner(cc.source)
		toks := s.ScanTokens()

		assertTokenTypes(t, toks, cc.expected...)
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
