//line parser.y:2
package yang

import __yyfmt__ "fmt"

//line parser.y:2
import (
	"fmt"
)

func (l *lexer) Lex(lval *yySymType) int {
	t, _ := l.nextToken()
	if t.typ == ParseEof {
		return 0
	}
	lval.token = t.val
	lval.stack = l.stack
	return int(t.typ)
}

func (l *lexer) Error(e string) {
	line, col := l.Position()
	msg := fmt.Sprintf("%s - line %d, col %d", e, line, col)
	l.lastError = &yangError{msg}
}

func HasError(l yyLexer, e error) bool {
	if e == nil {
		return false
	}
	l.Error(e.Error())
	return true
}

func popAndAddDef(yylval *yySymType) error {
	i := yylval.stack.Pop()
	if def, ok := i.(Def); ok {
		parent := yylval.stack.Peek()
		if parentList, ok := parent.(DefList); ok {
			return parentList.AddDef(def)
		} else {
			return &yangError{fmt.Sprintf("Cannot add \"%s\" to \"%s\"; not collection type.", i.GetIdent(), parent.GetIdent())}
		}
	} else {
		return &yangError{fmt.Sprintf("\"%s\" cannot be stored in a collection type.", i.GetIdent())}
	}
}

//line parser.y:48
type yySymType struct {
	yys   int
	ident string
	token string
	stack *yangDefStack
}

const token_ident = 57346
const token_string = 57347
const token_int = 57348
const token_curly_open = 57349
const token_curly_close = 57350
const token_semi = 57351
const token_rev_ident = 57352
const kywd_namespace = 57353
const kywd_description = 57354
const kywd_revision = 57355
const kywd_type = 57356
const kywd_prefix = 57357
const kywd_default = 57358
const kywd_length = 57359
const kywd_enum = 57360
const kywd_key = 57361
const kywd_config = 57362
const kywd_uses = 57363
const kywd_unique = 57364
const kywd_input = 57365
const kywd_output = 57366
const kywd_module = 57367
const kywd_container = 57368
const kywd_list = 57369
const kywd_rpc = 57370
const kywd_notification = 57371
const kywd_typedef = 57372
const kywd_grouping = 57373
const kywd_leaf = 57374
const kywd_mandatory = 57375
const kywd_reference = 57376
const kywd_leaf_list = 57377
const kywd_max_elements = 57378
const kywd_choice = 57379
const kywd_case = 57380

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"token_ident",
	"token_string",
	"token_int",
	"token_curly_open",
	"token_curly_close",
	"token_semi",
	"token_rev_ident",
	"kywd_namespace",
	"kywd_description",
	"kywd_revision",
	"kywd_type",
	"kywd_prefix",
	"kywd_default",
	"kywd_length",
	"kywd_enum",
	"kywd_key",
	"kywd_config",
	"kywd_uses",
	"kywd_unique",
	"kywd_input",
	"kywd_output",
	"kywd_module",
	"kywd_container",
	"kywd_list",
	"kywd_rpc",
	"kywd_notification",
	"kywd_typedef",
	"kywd_grouping",
	"kywd_leaf",
	"kywd_mandatory",
	"kywd_reference",
	"kywd_leaf_list",
	"kywd_max_elements",
	"kywd_choice",
	"kywd_case",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line parser.y:416

func parse(yang string) int {
	l := lex(yang)
	err := yyParse(l)
	return err
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 102
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 303

var yyAct = [...]int{

	177, 156, 70, 155, 105, 91, 159, 81, 139, 97,
	67, 158, 61, 124, 7, 69, 7, 9, 157, 3,
	178, 9, 9, 43, 86, 85, 89, 87, 9, 37,
	36, 37, 36, 102, 103, 78, 194, 78, 79, 193,
	79, 83, 37, 36, 71, 175, 178, 21, 78, 192,
	71, 79, 88, 95, 59, 109, 84, 93, 63, 107,
	68, 82, 92, 98, 106, 62, 9, 99, 45, 152,
	73, 153, 63, 112, 108, 89, 191, 117, 68, 62,
	37, 36, 25, 88, 190, 181, 78, 84, 125, 79,
	173, 25, 82, 95, 6, 9, 133, 93, 8, 42,
	169, 41, 92, 141, 141, 168, 137, 109, 98, 146,
	142, 107, 99, 9, 167, 166, 106, 165, 164, 72,
	86, 85, 89, 87, 154, 163, 145, 37, 36, 150,
	9, 24, 147, 78, 138, 134, 79, 83, 108, 89,
	24, 126, 132, 171, 37, 36, 9, 136, 144, 118,
	78, 9, 116, 79, 94, 89, 9, 180, 113, 110,
	37, 36, 102, 103, 40, 180, 78, 15, 143, 79,
	37, 36, 9, 71, 186, 172, 78, 188, 71, 79,
	94, 89, 170, 149, 127, 185, 37, 36, 121, 9,
	120, 65, 78, 161, 52, 79, 179, 160, 89, 51,
	9, 50, 65, 9, 161, 65, 49, 161, 160, 89,
	162, 160, 89, 44, 48, 37, 36, 38, 39, 34,
	35, 162, 111, 9, 162, 65, 9, 64, 65, 46,
	64, 37, 36, 38, 39, 34, 35, 6, 9, 14,
	19, 8, 187, 184, 183, 182, 148, 135, 130, 129,
	128, 119, 114, 18, 17, 16, 5, 189, 131, 123,
	122, 12, 115, 58, 57, 56, 55, 54, 53, 10,
	77, 76, 80, 30, 66, 47, 29, 104, 33, 101,
	100, 96, 32, 90, 31, 176, 174, 151, 60, 28,
	75, 74, 140, 27, 26, 23, 22, 13, 20, 11,
	4, 2, 1,
}
var yyPact = [...]int{

	-6, -1000, 83, 265, 226, 158, 250, -1000, 249, 248,
	233, 189, 155, 92, 13, -1000, -1000, -1000, -1000, -1000,
	205, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 222, 207,
	199, 194, 192, 187, 264, 263, 262, 261, 260, 259,
	-1000, -1000, 9, -1000, -1000, -1000, 211, -1000, 16, 101,
	160, 10, 54, -1000, -1000, -1000, -1000, -1000, -1000, 150,
	214, -1000, -1000, 149, 247, 258, 144, -1000, 140, -1000,
	-1000, 246, -1000, -1000, -1000, -1000, 183, 181, 256, 255,
	5, -1000, 132, 178, -1000, 245, 244, 243, -1000, 254,
	134, -1000, 126, -1000, 242, -1000, 139, -1000, 125, -1000,
	3, 3, 161, 141, 118, -1000, 123, -1000, 241, -1000,
	175, -1000, -1000, -1000, 120, 62, -1000, -1000, -1000, 115,
	191, 191, -1000, -1000, -1000, -1000, -1000, 109, 108, 106,
	105, 96, -1000, -1000, -1000, 91, -1000, -1000, -1000, 174,
	3, -1000, 167, -1000, -1000, -1000, -1000, -1000, 81, -1000,
	-1000, -1000, 28, -1000, -1000, 188, -1000, -1000, 76, -1000,
	240, 239, 238, 177, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, 166, 237, 2, -1000, 253, -1000,
	-1000, -1000, 75, 67, 40, -1000, -1000, 30, -1000, 27,
	-1000, -1000, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 302, 301, 300, 299, 298, 297, 11, 256, 47,
	296, 295, 119, 70, 294, 293, 8, 292, 2, 291,
	290, 289, 288, 12, 18, 287, 286, 285, 284, 283,
	5, 6, 282, 281, 9, 15, 280, 279, 278, 277,
	4, 276, 275, 274, 10, 273, 272, 7, 271, 3,
	1, 270, 0,
}
var yyR1 = [...]int{

	0, 1, 2, 6, 4, 4, 7, 3, 3, 8,
	8, 8, 9, 9, 9, 9, 9, 9, 5, 5,
	16, 16, 18, 18, 18, 18, 17, 17, 10, 21,
	22, 22, 23, 23, 23, 24, 25, 25, 26, 26,
	13, 28, 29, 29, 30, 30, 30, 30, 31, 14,
	32, 33, 33, 34, 34, 34, 34, 36, 37, 15,
	38, 39, 39, 40, 40, 40, 40, 11, 42, 41,
	43, 43, 44, 44, 44, 12, 45, 46, 46, 47,
	47, 47, 47, 47, 47, 47, 19, 48, 49, 49,
	50, 50, 50, 50, 50, 50, 20, 51, 27, 27,
	52, 35,
}
var yyR2 = [...]int{

	0, 5, 3, 2, 2, 5, 2, 2, 3, 2,
	1, 2, 1, 1, 1, 1, 1, 1, 1, 2,
	0, 1, 1, 1, 1, 1, 1, 2, 4, 2,
	1, 2, 1, 2, 3, 3, 3, 1, 3, 1,
	4, 2, 1, 2, 2, 1, 3, 1, 3, 4,
	2, 1, 2, 2, 1, 3, 3, 2, 2, 4,
	2, 1, 2, 2, 1, 3, 1, 2, 3, 2,
	1, 2, 2, 1, 1, 4, 2, 1, 2, 2,
	3, 1, 3, 3, 3, 1, 4, 2, 1, 2,
	1, 2, 1, 3, 3, 3, 4, 2, 1, 2,
	3, 3,
}
var yyChk = [...]int{

	-1000, -1, -2, 25, -3, -8, 11, -7, 15, 12,
	4, -4, -8, -6, 13, 9, 5, 5, 5, 7,
	-5, -9, -10, -11, -12, -13, -14, -15, -21, -41,
	-45, -28, -32, -38, 30, 31, 27, 26, 28, 29,
	9, 9, 7, 10, 8, -9, 7, -42, 7, 7,
	7, 7, 7, 4, 4, 4, 4, 4, 4, -7,
	-22, -23, -24, -7, 16, 14, -43, -44, -7, -35,
	-18, 34, -12, -13, -19, -20, -48, -51, 32, 35,
	-46, -47, -7, 36, -31, 20, 19, 22, -18, 21,
	-29, -30, -7, -31, 20, -18, -33, -34, -7, -35,
	-36, -37, 23, 24, -39, -40, -7, -31, 20, -18,
	9, 8, -23, 9, 5, 4, 8, -44, 9, 5,
	7, 7, 4, 4, 8, -47, 9, 6, 5, 5,
	5, 4, 8, -30, 9, 5, 8, -34, 9, -16,
	-17, -18, -16, 7, 7, 8, -40, 9, 5, 8,
	9, -25, 7, 9, 9, -49, -50, -24, -7, -31,
	20, 16, 33, -49, 9, 9, 9, 9, 9, 9,
	8, -18, 8, 9, -26, 17, -27, -52, 18, 8,
	-50, 9, 5, 5, 5, 8, 8, 5, -52, 4,
	9, 9, 9, 9, 9,
}
var yyDef = [...]int{

	0, -2, 0, 0, 0, 0, 0, 10, 0, 0,
	0, 0, 0, 0, 0, 7, 9, 11, 6, 2,
	0, 18, 12, 13, 14, 15, 16, 17, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	8, 4, 0, 3, 1, 19, 0, 67, 0, 0,
	0, 0, 0, 29, 69, 76, 41, 50, 60, 0,
	0, 30, 32, 0, 0, 0, 0, 70, 0, 73,
	74, 0, 22, 23, 24, 25, 0, 0, 0, 0,
	0, 77, 0, 0, 81, 0, 0, 0, 85, 0,
	0, 42, 0, 45, 0, 47, 0, 51, 0, 54,
	20, 20, 0, 0, 0, 61, 0, 64, 0, 66,
	0, 28, 31, 33, 0, 0, 68, 71, 72, 0,
	0, 0, 87, 97, 75, 78, 79, 0, 0, 0,
	0, 0, 40, 43, 44, 0, 49, 52, 53, 0,
	21, 26, 0, 57, 58, 59, 62, 63, 0, 5,
	34, 35, 0, 37, 101, 0, 88, 90, 0, 92,
	0, 0, 0, 0, 80, 82, 83, 84, 48, 46,
	55, 27, 56, 65, 0, 0, 39, 98, 0, 86,
	89, 91, 0, 0, 0, 96, 36, 0, 99, 0,
	93, 94, 95, 38, 100,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38,
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

	case 2:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:101
		{
			m := &Module{Ident: yyDollar[2].token}
			yylval.stack.Push(m)
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:107
		{
			d := yylval.stack.Peek()
			r := &Revision{Ident: yyDollar[2].token}
			d.(*Module).Revision = r
			yylval.stack.Push(r)
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:115
		{
			yylval.stack.Pop()
		}
	case 5:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:118
		{
			yylval.stack.Pop()
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:122
		{
			yylval.stack.Peek().(Describable).SetDescription(yyDollar[2].token)
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:132
		{
			d := yylval.stack.Peek()
			d.(*Module).Namespace = yyDollar[2].token
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:137
		{
			m := yylval.stack.Peek().(*Module)
			m.Prefix = yyDollar[2].token
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:151
		{
			if HasError(yylex, popAndAddDef(&yylval)) {
				goto ret1
			}
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:156
		{
			if HasError(yylex, popAndAddDef(&yylval)) {
				goto ret1
			}
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:173
		{
			if HasError(yylex, popAndAddDef(&yylval)) {
				goto ret1
			}
		}
	case 27:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:178
		{
			if HasError(yylex, popAndAddDef(&yylval)) {
				goto ret1
			}
		}
	case 29:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:191
		{
			yylval.stack.Push(&Typedef{Ident: yyDollar[2].token})
		}
	case 41:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:225
		{
			yylval.stack.Push(&Container{Ident: yyDollar[2].token})
		}
	case 47:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:237
		{
			if HasError(yylex, popAndAddDef(&yylval)) {
				goto ret1
			}
		}
	case 48:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:244
		{
			yylval.stack.Push(&Uses{Ident: yyDollar[2].token})
			if HasError(yylex, popAndAddDef(&yylval)) {
				goto ret1
			}
		}
	case 50:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:259
		{
			yylval.stack.Push(&Rpc{Ident: yyDollar[2].token})
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:270
		{
			input := yylval.stack.Pop().(*RpcInput)
			rpc := yylval.stack.Peek().(*Rpc)
			rpc.Input = input
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:275
		{
			output := yylval.stack.Pop().(*RpcOutput)
			rpc := yylval.stack.Peek().(*Rpc)
			rpc.Output = output
		}
	case 57:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:282
		{
			yylval.stack.Push(&RpcInput{})
		}
	case 58:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:287
		{
			yylval.stack.Push(&RpcOutput{})
		}
	case 60:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:298
		{
			yylval.stack.Push(&Notification{Ident: yyDollar[2].token})
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:311
		{
			if HasError(yylex, popAndAddDef(&yylval)) {
				goto ret1
			}
		}
	case 69:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:327
		{
			yylval.stack.Push(&Grouping{Ident: yyDollar[2].token})
		}
	case 74:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:338
		{
			if HasError(yylex, popAndAddDef(&yylval)) {
				goto ret1
			}
		}
	case 76:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:350
		{
			yylval.stack.Push(&List{Ident: yyDollar[2].token})
		}
	case 85:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:365
		{
			if HasError(yylex, popAndAddDef(&yylval)) {
				goto ret1
			}
		}
	case 87:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:378
		{
			yylval.stack.Push(&Leaf{Ident: yyDollar[2].token})
		}
	case 97:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:403
		{
			yylval.stack.Push(&LeafList{Ident: yyDollar[2].token})
		}
	}
	goto yystack /* stack new state and value */
}
