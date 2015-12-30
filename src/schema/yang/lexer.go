package yang

import "schema"

// This uses the go feature call go tools in the build process. To ensure this gets
//  called before compilation, make this call before building
//
//    go generate schema/yang
//
//go:generate go tool yacc -o parser.go parser.y

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Token struct {
	typ int
	val string
}

type stateFunc func(*lexer) stateFunc
type tokenSink func(*Token)

const (
	ParseEof = iota
	ParseErr
)

const (
	char_doublequote  = '"'
	char_backslash    = '\\'
	str_comment_start = "/*"
	str_comment_end   = "*/"
)

// needs to be in-sync w/ %token list in parser.y
var keywords = [...]string{
	"[ident]",
	"[string]",
	"[int]",
	"{",
	"}",
	";",
	"[revision]",

	// KEEP LIST IN SYNC WITH parser.y
	"namespace",
	"description",
	"revision",
	"type",
	"prefix",
	"default",
	"length",
	"enum",
	"key",
	"config",
	"uses",
	"unique",
	"input",
	"output",
	"module",
	"container",
	"list",
	"rpc",
	"notification",
	"typedef",
	"grouping",
	"leaf",
	"mandatory",
	"reference",
	"leaf-list",
	"max-elements",
	"choice",
	"case",
	"import",
	"include",
	"action",
	"anyxml",
}

const eof rune = 0

func (l *lexer) keyword(ttype int) string {
	if ttype < token_ident {
		panic("Not a keyword")
	}
	return keywords[ttype-token_ident]
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

func (l *lexer) error(msg string) stateFunc {
	l.tokens = append(l.tokens, Token{
		ParseErr,
		msg,
	})
	fmt.Println("Setting err ", msg)
	l.Error(msg)
	return nil
}

func (l *lexer) importModule(into *schema.Module, moduleName string) error {
	fmt.Printf("lexer.go - Import module here %s\n", moduleName)
	return nil
}

type yangMetaStack struct {
	defs  []schema.Identifiable
	count int
}

func (s *yangMetaStack) Push(def schema.Identifiable) {
	s.defs[s.count] = def
	s.count++
}

func (s *yangMetaStack) Pop() schema.Identifiable {
	s.count--
	return s.defs[s.count]
}

func (s *yangMetaStack) Peek() schema.Identifiable {
	return s.defs[s.count-1]
}

func newDefStack(size int) *yangMetaStack {
	return &yangMetaStack{
		defs:  make([]schema.Identifiable, size),
		count: 0,
	}
}

type lexer struct {
	pos       int
	start     int
	width     int
	state     stateFunc
	input     string
	tokens    []Token
	head      int
	tail      int
	stack     *yangMetaStack
	importer  ImportModule
	lastError error
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

func (l *lexer) acceptWS() {
	for {
		for unicode.IsSpace(l.next()) {
		}
		l.backup()

		if strings.HasPrefix(l.input[l.pos:], str_comment_start) {
			for {
				l.next()
				if strings.HasPrefix(l.input[l.pos:], str_comment_end) {
					l.pos += len(str_comment_end)
					break
				}
			}
		} else {
			break
		}
	}
	l.ignore()
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

func (l *lexer) acceptToken(ttype int) bool {
	var keyword string
	switch ttype {
	case token_ident:
		return l.acceptAlphaNumeric(token_ident)
	case token_string:
		return l.acceptString(token_ident)
	case token_int:
		return l.acceptInteger(token_int)
	case token_curly_open:
		keyword = "{"
		break
	case token_curly_close:
		keyword = "}"
		break
	case token_semi:
		keyword = ";"
		break
	case token_rev_ident:
		return l.acceptAlphaNumeric(token_rev_ident)
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

func (l *lexer) acceptRun(ttype int, valid string) bool {
	found := false
	for strings.IndexRune(valid, l.next()) >= 0 {
		found = true
	}
	l.backup()
	if found {
		l.emit(ttype)
	}
	return found
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

func (l *lexer) acceptInteger(ttype int) bool {
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

func lexBegin(l *lexer) stateFunc {
	if l.isEof() {
		return nil
	}

	// FORMAT: xxx zzz { ...
	// order from longest keyword to shorted to ensure "foobar" doesn't get picked
	// up as "foo"
	defTypes := [...]int{
		kywd_notification,
		kywd_container,
		kywd_leaf_list,
		kywd_grouping,
		kywd_typedef,
		kywd_module,
		kywd_choice,
		kywd_leaf,
		kywd_list,
		kywd_case,
		kywd_rpc,
		kywd_action,
	}
	for _, ttype := range defTypes {
		if l.acceptToken(ttype) {
			if !l.acceptToken(token_ident) {
				return l.error("expected ident")
			}
			if !l.acceptToken(token_curly_open) {
				return l.error("expected {")
			}
			return lexBegin
		}
	}

	// FORMAT : aaa { ...
	defSpecialReference := [...]int{
		kywd_input,
		kywd_output,
	}
	for _, ttype := range defSpecialReference {
		if l.acceptToken(ttype) {
			if !l.acceptToken(token_curly_open) {
				return l.error("expected {")
			}
			return lexBegin
		}
	}

	if l.acceptToken(kywd_revision) {
		if !l.acceptRun(token_rev_ident, "0123456789-") {
			return l.error("expected identifier")
		}
		if l.acceptToken(token_curly_open) {
			return lexBegin
		}
		if !l.acceptToken(token_semi) {
			return l.error("expected { or ;")
		}
		return lexBegin
	}

	// FORMAT: Either
	//  xxx zzz;
	// or
	//  xxx zzz { ...
	defOrReference := [...]int{
		kywd_type,
		kywd_import,
		kywd_include,
		kywd_anyxml,
	}
	for _, ttype := range defOrReference {
		if l.acceptToken(ttype) {
			if !l.acceptToken(token_ident) {
				return l.error("expecting string")
			}
			if l.acceptToken(token_curly_open) {
				return lexBegin
			}
			if !l.acceptToken(token_semi) {
				return l.error("expecting { or ;")
			}
			return lexBegin
		}
	}

	// FORMAT:
	//  xxx zzz;
	tokenIdentPair := [...]int{
		kywd_uses,
		kywd_enum,
	}
	for _, ttype := range tokenIdentPair {
		if l.acceptToken(ttype) {
			if !l.acceptToken(token_ident) {
				return l.error("expecting string")
			}
			if !l.acceptToken(token_semi) {
				return l.error("expecting ;")
			}
			return lexBegin
		}
	}

	// FORMAT: xxx "zzz";
	tokenStringPair := [...]int{
		kywd_prefix,
		kywd_namespace,
		kywd_description,
		kywd_config,
		kywd_type,
		kywd_length,
		kywd_key,
		kywd_mandatory,
		kywd_default,
		kywd_unique,
	}
	for _, ttype := range tokenStringPair {
		if l.acceptToken(ttype) {
			if !l.acceptString(token_string) {
				return l.error("expecting string")
			}
			if !l.acceptToken(token_semi) {
				return l.error("expecting semicolon")
			}
			return lexBegin
		}
	}

	tokenIntPair := [...]int{
		kywd_max_elements,
	}
	for _, ttype := range tokenIntPair {
		if l.acceptToken(ttype) {
			if !l.acceptInteger(token_int) {
				return l.error("expecting integer")
			}
			if !l.acceptToken(token_semi) {
				return l.error("expecting semicolon")
			}
			return lexBegin
		}
	}

	if l.acceptToken(token_curly_close) {
		return lexBegin
	}
	return l.error("unknown statement")
}

func (l *lexer) emit(t int) {
	l.pushToken(Token{t, l.input[l.start:l.pos]})
	l.start = l.pos
	l.acceptWS()
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

const (
	lexRingBufferSize = 64
	nestedYangDefMax  = 256
)

func lex(input string, importer ImportModule) *lexer {
	l := &lexer{
		input:    input,
		tokens:   make([]Token, lexRingBufferSize),
		head:     0,
		tail:     0,
		state:    lexBegin,
		stack:    newDefStack(256),
		importer: importer,
	}
	l.acceptWS()
	return l
}
