//line parser.y:2
package yang

import __yyfmt__ "fmt"

//line parser.y:2
import (
	"fmt"
	"schema"
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
	lval.importer = l.importer
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
	if def, ok := i.(schema.Meta); ok {
		parent := yylval.stack.Peek()
		if parentList, ok := parent.(schema.MetaList); ok {
			return parentList.AddMeta(def)
		} else {
			return &yangError{fmt.Sprintf("Cannot add \"%s\" to \"%s\"; not collection type.", i.GetIdent(), parent.GetIdent())}
		}
	} else {
		return &yangError{fmt.Sprintf("\"%s\" cannot be stored in a collection type.", i.GetIdent())}
	}
}

//line parser.y:63
type yySymType struct {
	yys      int
	ident    string
	token    string
	dataType *schema.DataType
	stack    *yangMetaStack
	importer ImportModule
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
const kywd_import = 57381
const kywd_include = 57382

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
	"kywd_import",
	"kywd_include",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line parser.y:509

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 112
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 355

var yyAct = [...]int{

	196, 141, 143, 134, 123, 118, 109, 124, 85, 104,
	97, 92, 83, 10, 136, 3, 27, 125, 122, 163,
	7, 24, 7, 10, 194, 197, 27, 197, 10, 58,
	116, 112, 46, 114, 158, 213, 157, 51, 50, 136,
	212, 48, 49, 52, 10, 60, 53, 111, 54, 205,
	204, 116, 112, 46, 114, 203, 6, 10, 51, 50,
	9, 202, 48, 49, 52, 95, 10, 53, 111, 54,
	107, 115, 121, 98, 201, 81, 106, 89, 90, 84,
	93, 200, 199, 99, 11, 105, 110, 119, 88, 144,
	129, 132, 198, 192, 95, 139, 191, 6, 10, 16,
	84, 9, 190, 149, 98, 182, 107, 153, 57, 93,
	56, 115, 106, 161, 99, 164, 176, 172, 165, 162,
	121, 105, 154, 171, 150, 11, 110, 175, 170, 10,
	140, 102, 10, 100, 175, 119, 137, 183, 130, 55,
	120, 46, 17, 211, 207, 188, 51, 50, 189, 160,
	48, 49, 52, 10, 187, 53, 10, 54, 102, 186,
	127, 181, 46, 184, 126, 138, 147, 51, 50, 10,
	148, 48, 49, 52, 10, 88, 53, 128, 54, 146,
	89, 90, 94, 46, 166, 59, 206, 73, 51, 50,
	71, 88, 48, 49, 52, 70, 209, 53, 46, 54,
	69, 68, 10, 51, 50, 38, 39, 48, 49, 52,
	120, 46, 53, 67, 54, 10, 51, 50, 65, 62,
	48, 49, 52, 61, 46, 53, 22, 54, 208, 51,
	50, 179, 178, 48, 49, 52, 10, 88, 53, 177,
	54, 45, 173, 169, 94, 46, 168, 167, 155, 151,
	51, 50, 145, 210, 48, 49, 52, 46, 20, 53,
	19, 54, 51, 50, 38, 39, 48, 49, 52, 46,
	18, 53, 5, 54, 51, 50, 185, 14, 48, 49,
	52, 159, 180, 53, 80, 54, 10, 174, 102, 79,
	127, 10, 152, 102, 126, 127, 10, 78, 102, 126,
	100, 77, 76, 75, 74, 72, 64, 128, 63, 21,
	12, 44, 128, 113, 108, 42, 103, 66, 41, 91,
	29, 87, 86, 82, 28, 117, 43, 195, 193, 156,
	101, 96, 40, 135, 133, 131, 47, 37, 36, 35,
	34, 33, 32, 31, 30, 142, 26, 25, 8, 15,
	23, 13, 4, 2, 1,
}
var yyPact = [...]int{

	-10, -1000, 45, 306, 86, 133, 265, -1000, -1000, 255,
	253, 305, 219, 236, 130, 101, 19, -1000, -1000, -1000,
	-1000, -1000, -1000, 177, -1000, -1000, -1000, -1000, 216, 212,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 304, 302,
	211, 206, 194, 193, 188, 183, 301, 180, 300, 299,
	298, 297, 293, 285, 280, -1000, -1000, 16, -1000, -1000,
	-1000, 54, 224, -1000, -1000, 117, -1000, 203, 32, 190,
	144, 144, 129, 1, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 127, 157, -1000, 121, -1000, 248, 248, 247, 172,
	159, 162, -1000, 115, 244, -1000, 284, -1000, -1000, 113,
	243, 27, 277, 141, -1000, 110, -1000, -1000, 11, -1000,
	109, 178, 242, -1000, 241, -1000, 238, 120, -1000, 108,
	237, -1000, 279, -1000, -1000, 107, 234, 227, 226, 274,
	-1000, 153, 96, -24, -1000, 156, 272, 151, -1000, -1000,
	-1000, 146, 248, -1000, 140, 93, -1000, -1000, -1000, -1000,
	-1000, 87, -1000, -1000, -1000, 84, -1000, -1000, 7, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 83, 73, 72, 65,
	-1000, -1000, -1000, 52, -1000, -1000, -1000, 46, 41, 40,
	-1000, -1000, -1000, -1000, 248, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 136, 223, 9, -1000, 249, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 135, -1000, 31, -1000,
	26, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 354, 353, 352, 351, 350, 349, 17, 272, 348,
	21, 347, 346, 2, 1, 345, 344, 343, 342, 341,
	340, 339, 338, 337, 336, 335, 334, 3, 333, 332,
	331, 10, 7, 330, 329, 328, 327, 326, 325, 5,
	324, 323, 12, 8, 322, 321, 320, 319, 11, 318,
	317, 316, 9, 315, 314, 6, 313, 311, 18, 4,
	241, 0,
}
var yyR1 = [...]int{

	0, 1, 2, 6, 4, 4, 7, 3, 3, 8,
	8, 8, 8, 9, 10, 10, 10, 5, 5, 14,
	14, 13, 13, 13, 13, 13, 13, 13, 13, 15,
	15, 23, 25, 25, 25, 24, 26, 26, 27, 28,
	16, 29, 30, 30, 31, 31, 31, 32, 33, 34,
	34, 35, 35, 19, 37, 38, 38, 39, 39, 39,
	22, 11, 40, 41, 41, 42, 42, 42, 42, 44,
	45, 12, 46, 47, 47, 48, 48, 48, 17, 50,
	49, 51, 51, 52, 52, 52, 18, 53, 54, 54,
	55, 55, 55, 55, 55, 55, 56, 20, 57, 58,
	58, 59, 59, 59, 59, 59, 21, 60, 36, 36,
	61, 43,
}
var yyR2 = [...]int{

	0, 5, 3, 2, 2, 5, 2, 2, 3, 2,
	1, 1, 2, 2, 1, 1, 1, 1, 2, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	2, 4, 0, 2, 1, 2, 1, 2, 4, 2,
	4, 2, 1, 2, 1, 2, 3, 2, 2, 1,
	3, 3, 1, 4, 2, 1, 2, 2, 3, 1,
	3, 4, 2, 1, 2, 2, 1, 3, 3, 2,
	2, 4, 2, 1, 2, 2, 3, 1, 2, 3,
	2, 1, 2, 2, 1, 1, 4, 2, 1, 2,
	2, 3, 3, 1, 3, 1, 3, 4, 2, 1,
	2, 1, 2, 3, 3, 3, 4, 2, 1, 2,
	3, 3,
}
var yyChk = [...]int{

	-1000, -1, -2, 25, -3, -8, 11, -7, -9, 15,
	12, 39, 4, -4, -8, -6, 13, 9, 5, 5,
	5, 4, 7, -5, -10, -11, -12, -13, -40, -46,
	-16, -17, -18, -19, -20, -21, -22, -23, 28, 29,
	-29, -49, -53, -37, -57, -60, 21, -24, 30, 31,
	27, 26, 32, 35, 37, 9, 9, 7, 10, 8,
	-10, 7, 7, 4, 4, 7, -50, 7, 7, 7,
	7, 7, 4, 7, 4, 4, 4, 4, 4, 4,
	4, -7, -41, -42, -7, -43, -44, -45, 34, 23,
	24, -47, -48, -7, 20, -13, -30, -31, -32, -7,
	16, -33, 14, -51, -52, -7, -43, -13, -54, -55,
	-7, 36, 20, -56, 22, -13, 19, -38, -39, -7,
	20, -13, -58, -59, -32, -7, 20, 16, 33, -58,
	9, -25, -7, -26, -27, -28, 38, 9, 8, -42,
	9, -14, -15, -13, -14, 5, 7, 7, 8, -48,
	9, 5, 8, -31, 9, 5, -34, 9, 7, 4,
	8, -52, 9, 8, -55, 9, 6, 5, 5, 5,
	8, -39, 9, 5, 8, -59, 9, 5, 5, 5,
	8, 8, 9, -27, 7, 4, 8, 8, -13, 8,
	9, 9, 9, -35, 17, -36, -61, 18, 9, 9,
	9, 9, 9, 9, 9, 9, -14, 8, 5, -61,
	4, 8, 9, 9,
}
var yyDef = [...]int{

	0, -2, 0, 0, 0, 0, 0, 10, 11, 0,
	0, 0, 0, 0, 0, 0, 0, 7, 9, 12,
	6, 13, 2, 0, 17, 14, 15, 16, 0, 0,
	21, 22, 23, 24, 25, 26, 27, 28, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 8, 4, 0, 3, 1,
	18, 0, 0, 62, 72, 0, 78, 0, 0, 0,
	0, 0, 0, 32, 41, 80, 87, 54, 98, 107,
	35, 0, 0, 63, 0, 66, 19, 19, 0, 0,
	0, 0, 73, 0, 0, 77, 0, 42, 44, 0,
	0, 0, 0, 0, 81, 0, 84, 85, 0, 88,
	0, 0, 0, 93, 0, 95, 0, 0, 55, 0,
	0, 59, 0, 99, 101, 0, 0, 0, 0, 0,
	60, 0, 0, 34, 36, 0, 0, 0, 61, 64,
	65, 0, 20, 29, 0, 0, 69, 70, 71, 74,
	75, 0, 40, 43, 45, 0, 47, 49, 0, 48,
	79, 82, 83, 86, 89, 90, 0, 0, 0, 0,
	53, 56, 57, 0, 97, 100, 102, 0, 0, 0,
	106, 31, 33, 37, 19, 39, 5, 67, 30, 68,
	111, 76, 46, 0, 0, 52, 108, 0, 91, 92,
	94, 96, 58, 103, 104, 105, 0, 50, 0, 109,
	0, 38, 51, 110,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40,
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

	case 2:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:121
		{
			m := &schema.Module{Ident: yyDollar[2].token}
			yylval.stack.Push(m)
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:127
		{
			d := yylval.stack.Peek()
			r := &schema.Revision{Ident: yyDollar[2].token}
			d.(*schema.Module).Revision = r
			yylval.stack.Push(r)
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:135
		{
			yylval.stack.Pop()
		}
	case 5:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:138
		{
			yylval.stack.Pop()
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:142
		{
			yylval.stack.Peek().(schema.Describable).SetDescription(tokenString(yyDollar[2].token))
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:152
		{
			d := yylval.stack.Peek()
			d.(*schema.Module).Namespace = tokenString(yyDollar[2].token)
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:158
		{
			m := yylval.stack.Peek().(*schema.Module)
			m.Prefix = tokenString(yyDollar[2].token)
		}
	case 13:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:163
		{
			var err error
			if yylval.importer == nil {
				yylex.Error("No importer defined")
				goto ret1
			} else {
				m := yylval.stack.Peek().(*schema.Module)
				if err = yylval.importer(m, yyDollar[2].token); err != nil {
					yylex.Error(err.Error())
					goto ret1
				}
			}
		}
	case 31:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:207
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 35:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:219
		{
			yylval.stack.Push(&schema.Choice{Ident: yyDollar[2].token})
		}
	case 38:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:230
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:237
		{
			yylval.stack.Push(&schema.ChoiceCase{Ident: yyDollar[2].token})
		}
	case 40:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:245
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 41:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:252
		{
			yylval.stack.Push(&schema.Typedef{Ident: yyDollar[2].token})
		}
	case 47:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:267
		{
			y := yylval.stack.Peek().(schema.HasDataType)
			y.SetDataType(yylval.dataType)
		}
	case 48:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:272
		{
			yylval.dataType = schema.NewDataType(yyDollar[2].token)
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:281
		{
			var err error
			if err = yylval.dataType.DecodeLength(tokenString(yyDollar[2].token)); err != nil {
				yylex.Error(err.Error())
				goto ret1
			}
		}
	case 53:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:294
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 54:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:301
		{
			yylval.stack.Push(&schema.Container{Ident: yyDollar[2].token})
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:315
		{
			yylval.stack.Push(&schema.Uses{Ident: yyDollar[2].token})
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 61:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:327
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 62:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:334
		{
			yylval.stack.Push(&schema.Rpc{Ident: yyDollar[2].token})
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:345
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:350
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 69:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:357
		{
			yylval.stack.Push(&schema.RpcInput{})
		}
	case 70:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:362
		{
			yylval.stack.Push(&schema.RpcOutput{})
		}
	case 71:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:370
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 72:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:377
		{
			yylval.stack.Push(&schema.Notification{Ident: yyDollar[2].token})
		}
	case 78:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:393
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 80:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:405
		{
			yylval.stack.Push(&schema.Grouping{Ident: yyDollar[2].token})
		}
	case 86:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:421
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 87:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:428
		{
			yylval.stack.Push(&schema.List{Ident: yyDollar[2].token})
		}
	case 96:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:444
		{
			if list, valid := yylval.stack.Peek().(*schema.List); valid {
				list.Keys = strings.Split(tokenString(yyDollar[2].token), " ")
			} else {
				yylex.Error("expected a list for key statement")
				goto ret1
			}
		}
	case 97:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:458
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 98:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:465
		{
			yylval.stack.Push(&schema.Leaf{Ident: yyDollar[2].token})
		}
	case 106:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:486
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 107:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:493
		{
			yylval.stack.Push(&schema.LeafList{Ident: yyDollar[2].token})
		}
	case 110:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:502
		{
			yylval.dataType.Enumeration = append(yylval.dataType.Enumeration, yyDollar[2].token)
		}
	}
	goto yystack /* stack new state and value */
}
