//line parser.y:2
package yapl

import __yyfmt__ "fmt"

//line parser.y:2
import (
	"process"
)

func (l *lexer) Lex(lval *yySymType) int {
	t, _ := l.nextToken()
	if t.typ == ParseEof {
		return 0
	}
	lval.token = t.val
	lval.stack = &l.stack
	return int(t.typ)
}

//line parser.y:20
type yySymType struct {
	yys   int
	ident string
	token string
	stack *operatorStack
	sel   *process.Select
	args  []process.Expression
	expr  process.Expression
}

const token_ident = 57346
const token_string = 57347
const token_number = 57348
const token_space_indent = 57349
const token_eol = 57350
const token_open_paren = 57351
const token_close_paren = 57352
const token_comma = 57353
const token_equal = 57354
const token_function = 57355
const kywd_select = 57356
const kywd_into = 57357
const kywd_let = 57358
const kywd_goto = 57359
const kywd_if = 57360

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"token_ident",
	"token_string",
	"token_number",
	"token_space_indent",
	"token_eol",
	"token_open_paren",
	"token_close_paren",
	"token_comma",
	"token_equal",
	"token_function",
	"kywd_select",
	"kywd_into",
	"kywd_let",
	"kywd_goto",
	"kywd_if",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line parser.y:158

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 33
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 55

var yyAct = [...]int{

	31, 21, 40, 12, 10, 11, 27, 32, 33, 34,
	30, 26, 52, 22, 25, 24, 36, 43, 9, 12,
	10, 41, 38, 7, 37, 28, 4, 15, 39, 29,
	13, 2, 20, 5, 19, 18, 17, 42, 16, 14,
	8, 46, 6, 44, 50, 45, 3, 47, 1, 35,
	53, 48, 51, 49, 23,
}
var yyPact = [...]int{

	22, 22, -1000, 13, 11, -1000, 13, -1000, -3, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -6, 21, -5, 3, 20, 18, 3, -10, -1000,
	17, 11, -1000, -1000, -1000, -1000, 8, 11, -1000, 11,
	3, 11, -1000, 3, -1000, -1000, 11, -1000, 2, 3,
	-1000, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 54, 53, 51, 0, 49, 48, 31, 46, 42,
	23, 40, 39, 38, 36, 35, 34, 32, 5, 18,
}
var yyR1 = [...]int{

	0, 6, 6, 7, 9, 9, 10, 12, 12, 12,
	12, 12, 18, 13, 14, 16, 17, 1, 15, 15,
	4, 4, 4, 4, 5, 3, 3, 2, 2, 19,
	11, 11, 8,
}
var yyR2 = [...]int{

	0, 1, 2, 2, 1, 2, 2, 1, 1, 1,
	1, 1, 1, 4, 5, 3, 3, 2, 2, 4,
	1, 1, 1, 1, 4, 0, 1, 1, 2, 1,
	1, 2, 2,
}
var yyChk = [...]int{

	-1000, -6, -7, -8, 4, -7, -9, -10, -11, -19,
	7, -18, 8, -10, -12, -19, -13, -14, -15, -16,
	-17, 4, 16, -1, 18, 17, 14, 12, 4, -18,
	15, -4, 4, 5, 6, -5, 13, 4, 4, -4,
	12, 4, -18, 9, -18, -18, -4, -18, -3, -2,
	-4, -18, 10, -4,
}
var yyDef = [...]int{

	0, -2, 1, 0, 0, 2, 3, 4, 0, 30,
	29, 32, 12, 5, 6, 31, 7, 8, 9, 10,
	11, 0, 0, 0, 0, 0, 0, 0, 0, 18,
	0, 0, 20, 21, 22, 23, 0, 0, 17, 0,
	0, 0, 15, 25, 16, 13, 0, 19, 0, 26,
	27, 14, 24, 28,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lookahead func() int
}

func (p *yyParserImpl) Lookahead() int {
	return p.lookahead()
}

func yyNewParser() yyParser {
	p := &yyParserImpl{
		lookahead: func() int { return -1 },
	}
	return p
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yytoken := -1 // yychar translated into internal numbering
	yyrcvr.lookahead = func() int { return yychar }
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yychar = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar, yytoken = yylex1(yylex, &yylval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yychar = -1
		yytoken = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar, yytoken = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yychar = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:74
		{
			yylval.stack.nextDepth = 0
		}
	case 13:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:79
		{
			op := &process.Set{Name: yyDollar[1].token}
			yylval.stack.Push(op)
		}
	case 14:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:85
		{
			op := &process.Let{Name: yyDollar[2].token}
			yylval.stack.Push(op)
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:91
		{
			op := &process.If{}
			yylval.stack.Push(op)
		}
	case 16:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:97
		{
			op := &process.Goto{Script: yyDollar[2].token}
			yylval.stack.Push(op)
		}
	case 17:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:103
		{
			op := &process.Select{On: yyDollar[2].token}
			yylval.stack.Push(op)
			yyVAL.sel = op
		}
	case 19:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:111
		{
			yyDollar[1].sel.Into = yyDollar[3].token
		}
	case 20:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:116
		{
			yyVAL.expr = &process.Primative{Var: yyDollar[1].token}
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:119
		{
			yyVAL.expr = &process.Primative{Str: yyDollar[1].token}
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:122
		{
			yyVAL.expr = &process.Primative{Num: 0}
		}
	case 24:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:128
		{
			yyVAL.expr = &process.Function{Name: yyDollar[1].token, Arguments: yyDollar[3].args}
		}
	case 25:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:133
		{
			yyVAL.args = []process.Expression{}
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:139
		{
			yyVAL.args = []process.Expression{yyDollar[1].expr}
		}
	case 28:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:142
		{
			yyVAL.args = append(yyDollar[1].args, yyDollar[2].expr)
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:146
		{
			yylval.stack.depth = yylval.stack.nextDepth
			yylval.stack.nextDepth++
		}
	case 32:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:154
		{
			s := &process.Script{Name: yyDollar[1].token}
			yylval.stack.PushScript(s)
		}
	}
	goto yystack /* stack new state and value */
}
