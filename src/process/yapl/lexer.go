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
	char_space = ' '
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
	"[func]",

	// KEEP LIST IN SYNC WITH parser.y
	"select",
	"into",
	"let",
	"goto",
	"if",
}

const MaxCodeDepth = 64

type operatorStack struct {
	scripts   map[string]*process.Script
	code 	  [MaxCodeDepth]process.CodeBlock
	depth	  int
	nextDepth int
}

func (s *operatorStack) PushScript(script *process.Script) {
	if s.scripts == nil {
		s.scripts = make(map[string]*process.Script)
	}
	s.scripts[script.Name] = script
	s.code[0] = script
}

func (s *operatorStack) Push(def process.Operation) {
	s.code[s.depth].AddOperation(def)
	if next, isCodeBlock := def.(process.CodeBlock); isCodeBlock {
		s.code[s.depth + 1] = next
	}
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

func (l *lexer) acceptEmptyLine() bool {
	if l.isEof() {
		return false
	}
	var i int
	var c rune
	for i, c = range l.input[l.start:] {
		switch c {
		case char_space:
			continue
		case char_eol:
			goto empty_line
		default:
			return false
		}
	}

empty_line:
	l.pos += (i + 1)
	l.ignore()
	return true
}

func (l *lexer) acceptWS() {
	for {
		if l.next() != ' ' {
			l.backup()
			return
		}
		l.ignore()
		if l.isEof() {
			return
		}
	}
}

func (l *lexer) emit(t int) {
	s := l.input[l.start:l.pos]
	l.pushToken(Token{t, s})
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

func (l *lexer) isIdent() bool {
	for i := 0; true; i++ {
		r := l.next()
		if (unicode.IsDigit(r) && i > 0) ||
			unicode.IsLetter(r) ||
			(r == '-' && i > 0) ||
			(r == '.' && i > 0) ||
			r == '_' {
			continue
		}
		if i == 0 {
			break
		}
		l.backup()
		return true
	}
	l.pos = l.start
	return false
}

func (l *lexer) acceptFunction() bool {
	if l.isIdent() {
		for !l.isEof() {
			r := l.next()
			if unicode.IsSpace(r) {
				continue
			}
			if r == '(' {
				l.backup()
				l.emit(token_function)
				return true
			} else {
				break
			}
		}
	}
	l.pos = l.start
	return false
}

func (l *lexer) acceptAlphaNumeric(ttype int) bool {
	if l.isIdent() {
		l.emit(ttype)
		return true
	}
	l.pos = l.start
	return false
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
	case token_function:
		return l.acceptFunction()
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
		token_function,
		// least specific should be least first
		token_ident,
	}
	for _, t := range expectedTokens {
		if l.acceptToken(t) {
			if t == token_eol {
				return lexBegin
			}
			return lexExpression
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
	for l.acceptEmptyLine() {
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

		return lexBegin
	}

	if l.acceptToken(kywd_if) {
		return lexExpression
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
		return lexBegin
	}

	l.acceptToken(kywd_let)

	if l.acceptToken(token_ident) {
		if l.acceptToken(token_eol) {
			return lexBegin
		}
		if l.acceptToken(token_equal) {
			return lexExpression
		}
	}

	return nil
}

func lex(input string) *lexer {
	l := &lexer{
		input: input,
		tokens: make([]Token, 128),
		state : lexBegin,
	}
	l.state = l.state(l)
	return l
}