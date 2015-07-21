//line parser.y:2
package yang

import __yyfmt__ "fmt"

//line parser.y:2
import (
	"fmt"
)

type yangError struct {
	s string
}

func (err *yangError) Error() string {
	return err.s
}

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

func popAndAddMeta(yylval *yySymType) error {
	i := yylval.stack.Pop()
	if def, ok := i.(Meta); ok {
		parent := yylval.stack.Peek()
		if parentList, ok := parent.(MetaList); ok {
			return parentList.AddMeta(def)
		} else {
			return &yangError{fmt.Sprintf("Cannot add \"%s\" to \"%s\"; not collection type.", i.GetIdent(), parent.GetIdent())}
		}
	} else {
		return &yangError{fmt.Sprintf("\"%s\" cannot be stored in a collection type.", i.GetIdent())}
	}
}

//line parser.y:56
type yySymType struct {
	yys   int
	ident string
	token string
	stack *yangMetaStack
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

//line parser.y:462

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

const yyNprod = 108
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 331

var yyAct = [...]int{

	200, 145, 147, 129, 118, 120, 113, 105, 7, 100,
	7, 92, 80, 86, 26, 117, 119, 155, 131, 88,
	3, 9, 201, 26, 198, 201, 9, 55, 109, 108,
	45, 110, 181, 21, 182, 48, 47, 54, 208, 53,
	203, 49, 45, 9, 50, 107, 51, 48, 47, 39,
	40, 37, 38, 49, 57, 9, 50, 84, 51, 83,
	78, 207, 195, 89, 82, 103, 87, 93, 101, 131,
	111, 116, 9, 106, 114, 81, 6, 9, 127, 194,
	8, 94, 193, 97, 98, 82, 124, 192, 89, 191,
	190, 87, 134, 189, 90, 188, 81, 93, 148, 139,
	133, 202, 103, 143, 9, 101, 84, 111, 83, 152,
	106, 94, 156, 176, 187, 116, 9, 183, 114, 163,
	179, 174, 167, 109, 108, 45, 110, 168, 164, 167,
	48, 47, 175, 56, 150, 157, 49, 153, 144, 50,
	107, 51, 140, 142, 149, 135, 45, 9, 132, 185,
	125, 48, 47, 39, 40, 37, 38, 49, 97, 98,
	50, 162, 51, 52, 15, 9, 186, 184, 178, 90,
	45, 173, 72, 115, 45, 48, 47, 70, 196, 48,
	47, 49, 158, 151, 50, 49, 51, 9, 50, 69,
	51, 68, 67, 62, 61, 102, 45, 60, 58, 19,
	205, 48, 47, 204, 138, 171, 170, 49, 9, 169,
	50, 9, 51, 84, 165, 122, 161, 45, 160, 121,
	159, 154, 48, 47, 141, 136, 9, 18, 49, 17,
	90, 50, 123, 51, 115, 45, 44, 16, 206, 9,
	48, 47, 177, 9, 137, 5, 49, 102, 45, 50,
	12, 51, 45, 48, 47, 77, 76, 48, 47, 49,
	75, 74, 50, 49, 51, 90, 50, 172, 51, 73,
	166, 9, 71, 84, 9, 122, 84, 66, 122, 121,
	65, 64, 121, 6, 9, 14, 63, 8, 10, 43,
	104, 41, 123, 85, 59, 123, 28, 99, 30, 96,
	95, 91, 29, 112, 42, 199, 197, 180, 79, 27,
	130, 128, 126, 46, 36, 35, 34, 33, 32, 31,
	146, 25, 24, 23, 22, 13, 20, 11, 4, 2,
	1,
}
var yyPact = [...]int{

	-5, -1000, 65, 284, 272, 155, 232, -1000, 224, 222,
	192, 21, 154, 30, 17, -1000, -1000, -1000, -1000, -1000,
	125, -1000, -1000, -1000, -1000, -1000, -1000, 191, 190, 187,
	186, -1000, -1000, -1000, -1000, -1000, -1000, 282, 277, 276,
	273, 185, 184, 182, 170, 268, 165, 265, 257, 256,
	252, 251, -1000, -1000, 14, -1000, -1000, -1000, 43, -1000,
	231, 60, 227, -1000, -1000, -1000, -1000, 104, 214, 199,
	199, 141, 31, -1000, -1000, -1000, -1000, -1000, 139, 92,
	-1000, -1000, 136, 220, 240, 196, -1000, 133, -1000, -1000,
	219, 135, -1000, 129, -1000, 149, 149, 137, 127, 175,
	-1000, 128, 216, -1000, 9, -1000, 126, 176, 215, 213,
	211, -1000, 153, -1000, 119, 209, -1000, 262, -1000, -1000,
	118, 204, 201, 200, 259, -1000, 163, 112, -20, -1000,
	106, 238, 160, -1000, -1000, -1000, 111, 25, -1000, -1000,
	-1000, 108, -1000, -1000, -1000, 159, 149, -1000, 158, -1000,
	-1000, -1000, -1000, -1000, 105, -1000, -1000, -1000, 86, 84,
	81, 80, -1000, -1000, -1000, 78, -1000, -1000, -1000, 73,
	70, 53, -1000, -1000, -1000, -1000, 149, -1000, -1000, -1000,
	-1000, 7, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 93, 32, 198, 4,
	-1000, 234, -1000, -1000, 52, -1000, 29, -1000, -1000,
}
var yyPgo = [...]int{

	0, 330, 329, 328, 327, 326, 325, 5, 245, 33,
	324, 323, 322, 321, 2, 1, 320, 319, 318, 317,
	316, 315, 314, 313, 312, 311, 3, 310, 309, 308,
	12, 16, 307, 306, 305, 304, 303, 6, 302, 301,
	11, 19, 300, 299, 298, 297, 9, 296, 294, 293,
	13, 291, 290, 7, 289, 15, 4, 236, 0,
}
var yyR1 = [...]int{

	0, 1, 2, 6, 4, 4, 7, 3, 3, 8,
	8, 8, 9, 9, 9, 9, 9, 5, 5, 15,
	15, 14, 14, 14, 14, 14, 14, 16, 16, 22,
	24, 24, 24, 23, 25, 25, 26, 27, 10, 28,
	29, 29, 30, 30, 30, 31, 32, 32, 33, 33,
	18, 35, 36, 36, 37, 37, 37, 21, 12, 38,
	39, 39, 40, 40, 40, 40, 42, 43, 13, 44,
	45, 45, 46, 46, 46, 11, 48, 47, 49, 49,
	50, 50, 50, 17, 51, 52, 52, 53, 53, 53,
	53, 53, 53, 19, 54, 55, 55, 56, 56, 56,
	56, 56, 20, 57, 34, 34, 58, 41,
}
var yyR2 = [...]int{

	0, 5, 3, 2, 2, 5, 2, 2, 3, 2,
	1, 2, 1, 1, 1, 1, 1, 1, 2, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 2, 4,
	0, 2, 1, 2, 1, 2, 4, 2, 4, 2,
	1, 2, 1, 2, 3, 3, 3, 1, 3, 1,
	4, 2, 1, 2, 2, 3, 1, 3, 4, 2,
	1, 2, 2, 1, 3, 3, 2, 2, 4, 2,
	1, 2, 2, 3, 1, 2, 3, 2, 1, 2,
	2, 1, 1, 4, 2, 1, 2, 2, 3, 3,
	3, 3, 1, 4, 2, 1, 2, 1, 2, 3,
	3, 3, 4, 2, 1, 2, 3, 3,
}
var yyChk = [...]int{

	-1000, -1, -2, 25, -3, -8, 11, -7, 15, 12,
	4, -4, -8, -6, 13, 9, 5, 5, 5, 7,
	-5, -9, -10, -11, -12, -13, -14, -28, -47, -38,
	-44, -17, -18, -19, -20, -21, -22, 30, 31, 28,
	29, -51, -35, -54, -57, 21, -23, 27, 26, 32,
	35, 37, 9, 9, 7, 10, 8, -9, 7, -48,
	7, 7, 7, 4, 4, 4, 4, 7, 7, 7,
	7, 4, 7, 4, 4, 4, 4, 4, -7, -29,
	-30, -31, -7, 16, 14, -49, -50, -7, -41, -14,
	34, -39, -40, -7, -41, -42, -43, 23, 24, -45,
	-46, -7, 20, -14, -52, -53, -7, 36, 20, 19,
	22, -14, -36, -37, -7, 20, -14, -55, -56, -31,
	-7, 20, 16, 33, -55, 9, -24, -7, -25, -26,
	-27, 38, 9, 8, -30, 9, 5, 4, 8, -50,
	9, 5, 8, -40, 9, -15, -16, -14, -15, 7,
	7, 8, -46, 9, 5, 8, -53, 9, 6, 5,
	5, 5, 8, -37, 9, 5, 8, -56, 9, 5,
	5, 5, 8, 8, 9, -26, 7, 4, 8, 9,
	-32, 7, 9, 9, 8, -14, 8, 9, 9, 9,
	9, 9, 9, 9, 9, 9, -15, -33, 17, -34,
	-58, 18, 8, 8, 5, -58, 4, 9, 9,
}
var yyDef = [...]int{

	0, -2, 0, 0, 0, 0, 0, 10, 0, 0,
	0, 0, 0, 0, 0, 7, 9, 11, 6, 2,
	0, 17, 12, 13, 14, 15, 16, 0, 0, 0,
	0, 21, 22, 23, 24, 25, 26, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 8, 4, 0, 3, 1, 18, 0, 75,
	0, 0, 0, 39, 77, 59, 69, 0, 0, 0,
	0, 0, 30, 84, 51, 94, 103, 33, 0, 0,
	40, 42, 0, 0, 0, 0, 78, 0, 81, 82,
	0, 0, 60, 0, 63, 19, 19, 0, 0, 0,
	70, 0, 0, 74, 0, 85, 0, 0, 0, 0,
	0, 92, 0, 52, 0, 0, 56, 0, 95, 97,
	0, 0, 0, 0, 0, 57, 0, 0, 32, 34,
	0, 0, 0, 38, 41, 43, 0, 0, 76, 79,
	80, 0, 58, 61, 62, 0, 20, 27, 0, 66,
	67, 68, 71, 72, 0, 83, 86, 87, 0, 0,
	0, 0, 50, 53, 54, 0, 93, 96, 98, 0,
	0, 0, 102, 29, 31, 35, 19, 37, 5, 44,
	45, 0, 47, 107, 64, 28, 65, 73, 88, 89,
	90, 91, 55, 99, 100, 101, 0, 0, 0, 49,
	104, 0, 36, 46, 0, 105, 0, 48, 106,
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
		//line parser.y:110
		{
			m := &Module{Ident: yyDollar[2].token}
			yylval.stack.Push(m)
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:116
		{
			d := yylval.stack.Peek()
			r := &Revision{Ident: yyDollar[2].token}
			d.(*Module).Revision = r
			yylval.stack.Push(r)
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:124
		{
			yylval.stack.Pop()
		}
	case 5:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:127
		{
			yylval.stack.Pop()
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:131
		{
			yylval.stack.Peek().(Describable).SetDescription(yyDollar[2].token)
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:141
		{
			d := yylval.stack.Peek()
			d.(*Module).Namespace = yyDollar[2].token
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:146
		{
			m := yylval.stack.Peek().(*Module)
			m.Prefix = yyDollar[2].token
		}
	case 29:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:181
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 33:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:193
		{
			yylval.stack.Push(&Choice{Ident: yyDollar[2].token})
		}
	case 36:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:204
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:211
		{
			yylval.stack.Push(&ChoiceCase{Ident: yyDollar[2].token})
		}
	case 38:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:219
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:226
		{
			yylval.stack.Push(&Typedef{Ident: yyDollar[2].token})
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:241
		{
			y := yylval.stack.Peek().(HasType)
			y.SetType(yyDollar[2].token)
		}
	case 50:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:260
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 51:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:267
		{
			yylval.stack.Push(&Container{Ident: yyDollar[2].token})
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:281
		{
			yylval.stack.Push(&Uses{Ident: yyDollar[2].token})
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 58:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:293
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 59:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:300
		{
			yylval.stack.Push(&Rpc{Ident: yyDollar[2].token})
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:311
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:316
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 66:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:323
		{
			yylval.stack.Push(&RpcInput{})
		}
	case 67:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:328
		{
			yylval.stack.Push(&RpcOutput{})
		}
	case 68:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:336
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 69:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:343
		{
			yylval.stack.Push(&Notification{Ident: yyDollar[2].token})
		}
	case 75:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:359
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 77:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:371
		{
			yylval.stack.Push(&Grouping{Ident: yyDollar[2].token})
		}
	case 83:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:387
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 84:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:394
		{
			yylval.stack.Push(&List{Ident: yyDollar[2].token})
		}
	case 93:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:414
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 94:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:421
		{
			yylval.stack.Push(&Leaf{Ident: yyDollar[2].token})
		}
	case 102:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:442
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 103:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:449
		{
			yylval.stack.Push(&LeafList{Ident: yyDollar[2].token})
		}
	}
	goto yystack /* stack new state and value */
}
