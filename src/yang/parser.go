//line parser.y:2
package yang

import __yyfmt__ "fmt"

//line parser.y:2
import (
	"fmt"
	"strings"
)

type yangError struct {
	s string
}

func (err *yangError) Error() string {
	return err.s
}

func tokenString(s string) string {
	return s[1 : len(s)-1]
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

//line parser.y:61
type yySymType struct {
	yys      int
	ident    string
	token    string
	dataType *DataType
	stack    *yangMetaStack
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

//line parser.y:489

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

const yyNprod = 110
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 351

var yyAct = [...]int{

	193, 138, 140, 131, 120, 115, 106, 121, 82, 101,
	94, 89, 80, 133, 24, 160, 119, 122, 3, 9,
	7, 194, 7, 24, 9, 9, 113, 109, 43, 111,
	9, 21, 55, 48, 47, 86, 87, 45, 46, 49,
	9, 210, 50, 108, 51, 209, 85, 113, 109, 43,
	111, 133, 57, 202, 48, 47, 191, 194, 45, 46,
	49, 201, 92, 50, 108, 51, 200, 104, 112, 118,
	95, 199, 78, 103, 6, 9, 81, 90, 8, 155,
	96, 154, 102, 107, 116, 126, 141, 54, 129, 53,
	177, 92, 136, 198, 9, 197, 99, 81, 124, 196,
	146, 95, 123, 104, 150, 195, 90, 189, 112, 103,
	158, 96, 161, 188, 187, 125, 179, 118, 102, 173,
	168, 169, 149, 107, 172, 167, 9, 162, 99, 9,
	97, 172, 116, 159, 180, 151, 147, 117, 43, 137,
	134, 127, 185, 48, 47, 52, 157, 45, 46, 49,
	9, 15, 50, 9, 51, 99, 9, 124, 99, 43,
	97, 123, 135, 208, 48, 47, 9, 145, 45, 46,
	49, 9, 85, 50, 125, 51, 204, 86, 87, 91,
	43, 186, 56, 203, 184, 48, 47, 183, 85, 45,
	46, 49, 178, 206, 50, 43, 51, 181, 144, 9,
	48, 47, 35, 36, 45, 46, 49, 117, 43, 50,
	143, 51, 9, 48, 47, 70, 68, 45, 46, 49,
	67, 43, 50, 66, 51, 65, 48, 47, 205, 64,
	45, 46, 49, 9, 85, 50, 62, 51, 59, 58,
	163, 91, 43, 19, 176, 175, 174, 48, 47, 170,
	207, 45, 46, 49, 43, 166, 50, 165, 51, 48,
	47, 35, 36, 45, 46, 49, 43, 164, 50, 152,
	51, 48, 47, 148, 142, 45, 46, 49, 18, 171,
	50, 17, 51, 9, 16, 99, 182, 124, 6, 9,
	14, 123, 8, 5, 156, 77, 76, 75, 12, 74,
	73, 72, 71, 69, 125, 61, 60, 10, 42, 41,
	110, 105, 39, 100, 63, 38, 88, 26, 84, 83,
	79, 25, 114, 40, 192, 190, 153, 98, 93, 37,
	132, 130, 128, 44, 34, 33, 32, 31, 30, 29,
	28, 27, 139, 23, 22, 13, 20, 11, 4, 2,
	1,
}
var yyPact = [...]int{

	-7, -1000, 63, 303, 277, 142, 279, -1000, 276, 273,
	236, 233, 136, 80, 22, -1000, -1000, -1000, -1000, -1000,
	174, -1000, -1000, -1000, -1000, 232, 231, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 302, 301, 229, 222, 218,
	216, 213, 209, 299, 208, 298, 297, 296, 295, 293,
	292, 291, -1000, -1000, 18, -1000, -1000, -1000, 12, 221,
	-1000, -1000, 144, -1000, 200, 28, 187, 141, 141, 132,
	13, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 131, 154,
	-1000, 130, -1000, 245, 245, 269, 203, 191, 159, -1000,
	127, 268, -1000, 114, -1000, -1000, 126, 264, 72, 290,
	138, -1000, 124, -1000, -1000, 7, -1000, 118, 234, 262,
	-1000, 252, -1000, 250, 117, -1000, 112, 244, -1000, 271,
	-1000, -1000, 110, 241, 240, 239, 82, -1000, 184, 107,
	-25, -1000, 190, 282, 179, -1000, -1000, -1000, 176, 245,
	-1000, 173, 105, -1000, -1000, -1000, -1000, -1000, 104, -1000,
	-1000, -1000, 98, -1000, -1000, 39, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 96, 90, 86, 84, -1000, -1000, -1000,
	62, -1000, -1000, -1000, 57, 52, 44, -1000, -1000, -1000,
	-1000, 245, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	168, 223, 3, -1000, 246, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 155, -1000, 36, -1000, 32, -1000, -1000,
	-1000,
}
var yyPgo = [...]int{

	0, 350, 349, 348, 347, 346, 345, 17, 293, 31,
	344, 343, 2, 1, 342, 341, 340, 339, 338, 337,
	336, 335, 334, 333, 332, 331, 3, 330, 329, 328,
	10, 7, 327, 326, 325, 324, 323, 322, 5, 321,
	320, 12, 8, 319, 318, 317, 316, 11, 315, 314,
	313, 9, 312, 311, 6, 310, 309, 16, 4, 308,
	0,
}
var yyR1 = [...]int{

	0, 1, 2, 6, 4, 4, 7, 3, 3, 8,
	8, 8, 9, 9, 9, 5, 5, 13, 13, 12,
	12, 12, 12, 12, 12, 12, 12, 14, 14, 22,
	24, 24, 24, 23, 25, 25, 26, 27, 15, 28,
	29, 29, 30, 30, 30, 31, 32, 33, 33, 34,
	34, 18, 36, 37, 37, 38, 38, 38, 21, 10,
	39, 40, 40, 41, 41, 41, 41, 43, 44, 11,
	45, 46, 46, 47, 47, 47, 16, 49, 48, 50,
	50, 51, 51, 51, 17, 52, 53, 53, 54, 54,
	54, 54, 54, 54, 55, 19, 56, 57, 57, 58,
	58, 58, 58, 58, 20, 59, 35, 35, 60, 42,
}
var yyR2 = [...]int{

	0, 5, 3, 2, 2, 5, 2, 2, 3, 2,
	1, 2, 1, 1, 1, 1, 2, 0, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 2, 4,
	0, 2, 1, 2, 1, 2, 4, 2, 4, 2,
	1, 2, 1, 2, 3, 2, 2, 1, 3, 3,
	1, 4, 2, 1, 2, 2, 3, 1, 3, 4,
	2, 1, 2, 2, 1, 3, 3, 2, 2, 4,
	2, 1, 2, 2, 3, 1, 2, 3, 2, 1,
	2, 2, 1, 1, 4, 2, 1, 2, 2, 3,
	3, 1, 3, 1, 3, 4, 2, 1, 2, 1,
	2, 3, 3, 3, 4, 2, 1, 2, 3, 3,
}
var yyChk = [...]int{

	-1000, -1, -2, 25, -3, -8, 11, -7, 15, 12,
	4, -4, -8, -6, 13, 9, 5, 5, 5, 7,
	-5, -9, -10, -11, -12, -39, -45, -15, -16, -17,
	-18, -19, -20, -21, -22, 28, 29, -28, -48, -52,
	-36, -56, -59, 21, -23, 30, 31, 27, 26, 32,
	35, 37, 9, 9, 7, 10, 8, -9, 7, 7,
	4, 4, 7, -49, 7, 7, 7, 7, 7, 4,
	7, 4, 4, 4, 4, 4, 4, 4, -7, -40,
	-41, -7, -42, -43, -44, 34, 23, 24, -46, -47,
	-7, 20, -12, -29, -30, -31, -7, 16, -32, 14,
	-50, -51, -7, -42, -12, -53, -54, -7, 36, 20,
	-55, 22, -12, 19, -37, -38, -7, 20, -12, -57,
	-58, -31, -7, 20, 16, 33, -57, 9, -24, -7,
	-25, -26, -27, 38, 9, 8, -41, 9, -13, -14,
	-12, -13, 5, 7, 7, 8, -47, 9, 5, 8,
	-30, 9, 5, -33, 9, 7, 4, 8, -51, 9,
	8, -54, 9, 6, 5, 5, 5, 8, -38, 9,
	5, 8, -58, 9, 5, 5, 5, 8, 8, 9,
	-26, 7, 4, 8, 8, -12, 8, 9, 9, 9,
	-34, 17, -35, -60, 18, 9, 9, 9, 9, 9,
	9, 9, 9, -13, 8, 5, -60, 4, 8, 9,
	9,
}
var yyDef = [...]int{

	0, -2, 0, 0, 0, 0, 0, 10, 0, 0,
	0, 0, 0, 0, 0, 7, 9, 11, 6, 2,
	0, 15, 12, 13, 14, 0, 0, 19, 20, 21,
	22, 23, 24, 25, 26, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 8, 4, 0, 3, 1, 16, 0, 0,
	60, 70, 0, 76, 0, 0, 0, 0, 0, 0,
	30, 39, 78, 85, 52, 96, 105, 33, 0, 0,
	61, 0, 64, 17, 17, 0, 0, 0, 0, 71,
	0, 0, 75, 0, 40, 42, 0, 0, 0, 0,
	0, 79, 0, 82, 83, 0, 86, 0, 0, 0,
	91, 0, 93, 0, 0, 53, 0, 0, 57, 0,
	97, 99, 0, 0, 0, 0, 0, 58, 0, 0,
	32, 34, 0, 0, 0, 59, 62, 63, 0, 18,
	27, 0, 0, 67, 68, 69, 72, 73, 0, 38,
	41, 43, 0, 45, 47, 0, 46, 77, 80, 81,
	84, 87, 88, 0, 0, 0, 0, 51, 54, 55,
	0, 95, 98, 100, 0, 0, 0, 104, 29, 31,
	35, 17, 37, 5, 65, 28, 66, 109, 74, 44,
	0, 0, 50, 106, 0, 89, 90, 92, 94, 56,
	101, 102, 103, 0, 48, 0, 107, 0, 36, 49,
	108,
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
		//line parser.y:116
		{
			m := &Module{Ident: yyDollar[2].token}
			yylval.stack.Push(m)
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:122
		{
			d := yylval.stack.Peek()
			r := &Revision{Ident: yyDollar[2].token}
			d.(*Module).Revision = r
			yylval.stack.Push(r)
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:130
		{
			yylval.stack.Pop()
		}
	case 5:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:133
		{
			yylval.stack.Pop()
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:137
		{
			yylval.stack.Peek().(Describable).SetDescription(tokenString(yyDollar[2].token))
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:147
		{
			d := yylval.stack.Peek()
			d.(*Module).Namespace = tokenString(yyDollar[2].token)
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:152
		{
			m := yylval.stack.Peek().(*Module)
			m.Prefix = tokenString(yyDollar[2].token)
		}
	case 29:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:187
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 33:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:199
		{
			yylval.stack.Push(&Choice{Ident: yyDollar[2].token})
		}
	case 36:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:210
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:217
		{
			yylval.stack.Push(&ChoiceCase{Ident: yyDollar[2].token})
		}
	case 38:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:225
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:232
		{
			yylval.stack.Push(&Typedef{Ident: yyDollar[2].token})
		}
	case 45:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:247
		{
			y := yylval.stack.Peek().(HasDataType)
			y.SetDataType(yylval.dataType)
		}
	case 46:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:252
		{
			yylval.dataType = NewDataType(yyDollar[2].token)
		}
	case 49:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:261
		{
			var err error
			if err = yylval.dataType.DecodeLength(tokenString(yyDollar[2].token)); err != nil {
				yylex.Error(err.Error())
				goto ret1
			}
		}
	case 51:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:274
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 52:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:281
		{
			yylval.stack.Push(&Container{Ident: yyDollar[2].token})
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:295
		{
			yylval.stack.Push(&Uses{Ident: yyDollar[2].token})
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 59:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:307
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 60:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:314
		{
			yylval.stack.Push(&Rpc{Ident: yyDollar[2].token})
		}
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:325
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:330
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 67:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:337
		{
			yylval.stack.Push(&RpcInput{})
		}
	case 68:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:342
		{
			yylval.stack.Push(&RpcOutput{})
		}
	case 69:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:350
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 70:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:357
		{
			yylval.stack.Push(&Notification{Ident: yyDollar[2].token})
		}
	case 76:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:373
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 78:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:385
		{
			yylval.stack.Push(&Grouping{Ident: yyDollar[2].token})
		}
	case 84:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:401
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 85:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:408
		{
			yylval.stack.Push(&List{Ident: yyDollar[2].token})
		}
	case 94:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:424
		{
			if list, valid := yylval.stack.Peek().(*List); valid {
				list.Keys = strings.Split(tokenString(yyDollar[2].token), " ")
			} else {
				yylex.Error("expected a list for key statement")
				goto ret1
			}
		}
	case 95:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:438
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 96:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:445
		{
			yylval.stack.Push(&Leaf{Ident: yyDollar[2].token})
		}
	case 104:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:466
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 105:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:473
		{
			yylval.stack.Push(&LeafList{Ident: yyDollar[2].token})
		}
	case 108:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:482
		{
			yylval.dataType.Enumeration = append(yylval.dataType.Enumeration, yyDollar[2].token)
		}
	}
	goto yystack /* stack new state and value */
}
