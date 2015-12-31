package yapl
import "testing"


func TestYaplIndent(t *testing.T) {
	l := &lexer{input:""}
	if l.acceptSpaceIndent() {
		t.Error("expected no indent")
	}
	l = &lexer{input:"  "}
	if ! l.acceptSpaceIndent() {
		t.Error("expected one ident")
	}
	if l.acceptSpaceIndent() {
		t.Error("expected only one ident")
	}
}

func TestYaplLexer(t *testing.T) {
	l := lex(
`foo
  a = b
  let q = foo(bar,6)
  if x
    z = "y"
  select c into d
    goto bleep

goo
`)
	expectations := []Token {
		Token{token_ident, "foo"}, Token{token_eol, "\n"},
		Token{token_space_indent, "  "}, Token{token_ident, "a"}, Token{token_equal, "="}, Token{token_ident, "b"}, Token{token_eol, "\n"},
		Token{token_space_indent, "  "}, Token{kywd_let, "let"}, Token{token_ident, "q"}, Token{token_equal, "="}, Token{token_function, "foo"},
			Token{token_open_paren, "("}, Token{token_ident, "bar"},
		    Token{token_comma, ","}, Token{token_number, "6"},
		    Token{token_close_paren, ")"}, Token{token_eol, "\n"},
		Token{token_space_indent, "  "}, Token{kywd_if, "if"}, Token{token_ident, "x"}, Token{token_eol, "\n"},
		Token{token_space_indent, "  "}, Token{token_space_indent, "  "},
			Token{token_ident, "z"}, Token{token_equal, "="}, Token{token_string, `"y"`}, Token{token_eol, "\n"},
		Token{token_space_indent, "  "}, Token{kywd_select, "select"},
			Token{token_ident, "c"}, Token{kywd_into, "into"}, Token{token_ident, "d"}, Token{token_eol, "\n"},
		Token{token_space_indent, "  "}, Token{token_space_indent, "  "}, Token{kywd_goto, "goto"},
			Token{token_ident, "bleep"}, Token{token_eol, "\n"},
		Token{token_ident, "goo"}, Token{token_eol, "\n"},
		Token{1, "EOF"},
	}
	for i, expected := range expectations {
		tok, err := l.nextToken()
		if err != nil {
			t.Fatal(err)
		}
		if tok.typ != expected.typ ||  tok.val != expected.val {
			t.Errorf("Token #%d\nExpected:%s(%d)\n  Actual:%s(%d)", i, expected.String(), expected.typ,
				tok.String(), tok.typ)
		}
	}
}
