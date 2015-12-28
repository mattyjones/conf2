package yapl
import (
	"errors"
	"unicode"
	"unicode/utf8"
	"strings"
	"process"
	"fmt"
)

// This uses the go feature call go tools in the build process. To ensure this gets
//  called before compilation, make this call before building
//
//    go generate process/yapl
//
//go:generate go tool yacc -o parser.go parser.y

type Token struct {
	typ int
	val string
}

func (t Token) String() string {
	if t.typ == ParseErr {
		return fmt.Sprintf("ERROR: %q", t.val)
	}
	if len(t.val) > 10 {
		return fmt.Sprintf("%.10q...", t.val)
	}
	return fmt.Sprintf("%q", t.val)
}

type stateFunc func(*lexer) stateFunc

const eof rune = 0

const (
	ParseEof = iota + 1
	ParseErr
)

const (
	char_doublequote  = '"'
	char_eol = '\n'
	char_backslash    = '\\'
	str_comment_start = '#'
)

// needs to be in-sync w/ %token list in parser.y
var keywords = [...]string{
	"[ident]",
	"[string]",
	"[number]",
	"  ",
	"\\n",
	"(",
	")",
	",",
	"=",

	// KEEP LIST IN SYNC WITH parser.y
	"select",
	"into",
	"let",
	"goto",
	"if",
}


type operatorStack struct {
	scripts   []*process.Script
	defs  [64]process.Operation
	count int
}

func (s *operatorStack) PushScript(script *process.Script) {
	if s.scripts == nil {
		s.scripts = make([]*process.Script, 0, 1)
	}
	s.scripts = append(s.scripts, script)
}

func (s *operatorStack) Push(def process.Operation) {
	s.defs[s.count] = def
	s.count++
}

func (s *operatorStack) Pop() process.Operation {
	s.count--
	return s.defs[s.count]
}

func (s *operatorStack) Peek() process.Operation {
	return s.defs[s.count-1]
}

type lexer struct {
	pos       int
	start     int
	width     int
	state     stateFunc
	input     string
	tokens    []Token
	stack     operatorStack
	head      int
	tail      int
	lastError error
}

func (l *lexer) acceptSpaceIndent() bool {
	return l.acceptToken(token_space_indent)
}

func (l *lexer) acceptWS() {
	if l.next() != ' ' {
		l.backup()
		return
	}
	l.ignore()
}

func (l *lexer) emit(t int) {
	l.pushToken(Token{t, l.input[l.start:l.pos]})
	l.start = l.pos
	switch t {
	case token_space_indent, token_eol:
		break
	default:
		l.acceptWS()
	}
}

func (l *lexer) popToken() Token {
	token := l.tokens[l.tail]
	l.tail = (l.tail + 1) % len(l.tokens)
	return token
}

func (l *lexer) pushToken(t Token) {
	l.tokens[l.head] = t
	l.head = (l.head + 1) % len(l.tokens)
}

func (l *lexer) nextToken() (Token, error) {
	for {
		if l.head != l.tail {
			token := l.popToken()
			if token.typ == ParseEof {
				return token, errors.New(token.val)
			}
			return token, nil
		} else {
			if l.state == nil {
				return Token{ParseEof, "EOF"}, nil
			}
			l.state = l.state(l)
		}
	}
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) Position() (line, col int) {
	for p := 0; p < l.pos; p++ {
		if l.input[p] == '\n' {
			line += 1
			col = 0
		} else {
			col += 1
		}
	}
	return
}

func (l *lexer) isEof() bool {
	return l.pos >= len(l.input)
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) acceptAlphaNumeric(ttype int) bool {
	accepted := false
	for {
		r := l.next()
		// TODO: review spec on legal chars
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) && !(r == '-') && !(r == '_') {
			l.backup()
			if accepted {
				l.emit(ttype)
			}
			return accepted
		}
		accepted = true
	}
}

func (l *lexer) acceptString(ttype int) bool {
	r := l.next()
	if r != char_doublequote {
		l.backup()
		return false
	}
	for {
		r = l.next()
		if r == char_backslash {
			l.next()
		} else if r == char_doublequote {
			l.emit(ttype)
			return true
		} else if r == eof {
			// bad format?
			return false
		}
	}
}

func (l *lexer) acceptNumber(ttype int) bool {
	accepted := false
	for {
		r := l.next()
		if !unicode.IsDigit(r) {
			l.backup()
			if accepted {
				l.emit(ttype)
			}
			return accepted
		}
		accepted = true
	}
}

func (l *lexer) keyword(ttype int) string {
	if ttype < token_ident {
		panic("Not a keyword")
	}
	return keywords[ttype-token_ident]
}

func (l *lexer) acceptToken(ttype int) bool {
	if l.isEof() {
		return false
	}
	var keyword string
	switch ttype {
	case token_ident:
		return l.acceptAlphaNumeric(token_ident)
	case token_string:
		return l.acceptString(token_string)
	case token_number:
		return l.acceptNumber(token_number)
	case token_space_indent:
		keyword = "  "
	case token_equal:
		keyword = "="
	case token_comma:
		keyword = ","
	case token_open_paren:
		keyword = "("
	case token_close_paren:
		keyword = ")"
	case token_eol:
		keyword = "\n"
	default:
		keyword = l.keyword(ttype)
	}
	if !strings.HasPrefix(l.input[l.pos:], keyword) {
		return false
	}
	l.pos += len(keyword)
	l.emit(ttype)
	return true
}

func lexExpression(l *lexer) stateFunc {
	expectedTokens := []int {
		token_open_paren,
		token_close_paren,
		token_comma,
		token_eol,
		token_number,
		token_string,
		// least specific should be least first
		token_ident,
	}
	for _, t := range expectedTokens {
		if l.acceptToken(t) {
			if t == token_eol {
				return lexBegin(l)
			}
			return lexExpression(l)
		}
	}

	return nil
}

func (l *lexer) Error(e string) {
	line, col := l.Position()
	msg := fmt.Sprintf("%s - line %d, col %d", e, line, col)
	l.lastError = errors.New(msg)
}

func (l *lexer) error(msg string) stateFunc {
	l.tokens = append(l.tokens, Token{
		ParseErr,
		msg,
	})
	fmt.Println("Setting err ", msg)
	l.Error(msg)
	return nil
}

func lexBegin(l *lexer) stateFunc {
	for l.acceptToken(token_eol) {
	}

	for l.acceptToken(token_space_indent) {
	}

	if l.acceptToken(kywd_goto) {
		if ! l.acceptToken(token_ident) {
			return l.error("Expected identifier after 'goto'")
		}
		if ! l.acceptToken(token_eol) {
			return l.error("Expected end of statement")
		}

		return lexBegin(l)
	}

	if l.acceptToken(kywd_if) {
		return lexExpression(l)
	}

	if l.acceptToken(kywd_select) {
		if ! l.acceptToken(token_ident) {
			return l.error("Expected identifier after 'select'")
		}
		if l.acceptToken(kywd_into) {
			if ! l.acceptToken(token_ident) {
				return l.error("Expected identifier after 'into'")
			}
		}
		if ! l.acceptToken(token_eol) {
			return l.error("Expected end of statement")
		}
		return lexBegin(l)
	}

	l.acceptToken(kywd_let)

	if l.acceptToken(token_ident) {
		if l.acceptToken(token_eol) {
			return lexBegin(l)
		}
		if l.acceptToken(token_equal) {
			return lexExpression(l)
		}
	}

	return nil
}

func lex(input string) *lexer {
	l := &lexer{
		input: input,
		tokens: make([]Token, 64),
		state : lexBegin,
	}
	l.state = l.state(l)
	return l
}