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

//line parser.y:517

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 113
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 361

var yyAct = [...]int{

	119, 194, 118, 122, 110, 135, 124, 85, 98, 92,
	105, 120, 142, 83, 7, 121, 7, 27, 10, 164,
	123, 137, 125, 10, 6, 10, 16, 27, 9, 3,
	117, 96, 46, 115, 195, 139, 161, 51, 50, 10,
	10, 48, 49, 52, 137, 10, 53, 112, 54, 46,
	89, 90, 11, 10, 51, 50, 89, 90, 48, 49,
	52, 88, 88, 53, 58, 54, 95, 88, 159, 81,
	158, 108, 116, 84, 93, 107, 208, 100, 94, 106,
	111, 24, 126, 126, 113, 133, 127, 127, 99, 202,
	144, 144, 130, 207, 84, 95, 140, 192, 195, 200,
	145, 150, 57, 93, 56, 60, 154, 94, 108, 100,
	199, 198, 107, 116, 165, 162, 106, 46, 197, 171,
	99, 111, 51, 50, 196, 113, 48, 49, 52, 187,
	174, 53, 10, 54, 103, 126, 101, 174, 185, 127,
	181, 10, 126, 190, 189, 188, 127, 186, 117, 96,
	46, 115, 180, 175, 172, 51, 50, 166, 163, 48,
	49, 52, 155, 151, 53, 112, 54, 141, 138, 206,
	153, 6, 10, 10, 10, 9, 103, 131, 101, 55,
	17, 96, 46, 184, 179, 201, 182, 51, 50, 148,
	147, 48, 49, 52, 73, 204, 53, 170, 54, 11,
	178, 10, 171, 71, 10, 70, 103, 69, 128, 96,
	46, 68, 96, 67, 65, 51, 50, 167, 149, 48,
	49, 52, 10, 62, 53, 129, 54, 61, 22, 203,
	96, 46, 177, 59, 176, 169, 51, 50, 168, 156,
	48, 49, 52, 45, 152, 53, 46, 54, 146, 20,
	10, 51, 50, 38, 39, 48, 49, 52, 96, 46,
	53, 19, 54, 10, 51, 50, 18, 44, 48, 49,
	52, 5, 46, 53, 205, 54, 14, 51, 50, 183,
	160, 48, 49, 52, 46, 88, 53, 80, 54, 51,
	50, 38, 39, 48, 49, 52, 79, 173, 53, 78,
	54, 10, 10, 103, 103, 128, 128, 77, 76, 96,
	96, 75, 74, 72, 64, 63, 21, 12, 114, 109,
	42, 104, 129, 129, 66, 41, 91, 29, 87, 86,
	82, 28, 43, 193, 191, 157, 102, 97, 40, 136,
	134, 132, 47, 37, 36, 35, 34, 33, 32, 31,
	30, 143, 26, 25, 8, 15, 23, 13, 4, 2,
	1,
}
var yyPact = [...]int{

	4, -1000, 160, 313, 13, 171, 261, -1000, -1000, 256,
	244, 312, 221, 263, 170, 95, 54, -1000, -1000, -1000,
	-1000, -1000, -1000, 225, -1000, -1000, -1000, -1000, 220, 216,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 311, 310,
	207, 206, 204, 200, 198, 196, 309, 187, 308, 307,
	304, 303, 295, 292, 283, -1000, -1000, 41, -1000, -1000,
	-1000, 33, 238, -1000, -1000, 120, -1000, 251, 129, 238,
	290, 290, 168, 6, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 159, 27, -1000, 158, -1000, 96, 96, 243, 183,
	182, 210, -1000, 154, -1000, -1000, 239, 162, -1000, -1000,
	153, 234, 61, 276, 28, -1000, 149, -1000, -1000, 11,
	-1000, 148, 211, -1000, -1000, 233, -1000, 230, 189, -1000,
	145, -1000, -1000, 289, -1000, -1000, 144, -1000, 229, 227,
	192, -1000, 176, 143, -17, -1000, 179, 275, 175, -1000,
	-1000, -1000, 130, 96, -1000, 121, 136, -1000, -1000, -1000,
	-1000, -1000, 135, -1000, -1000, -1000, 134, -1000, -1000, 80,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 115, 109, 102,
	-1000, -1000, -1000, -1000, -1000, -1000, 101, 90, -1000, -1000,
	-1000, -1000, 238, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 81, 224, 16, -1000, 270, -1000, -1000, -1000, -1000,
	-1000, 161, -1000, 84, -1000, 67, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 360, 359, 358, 357, 356, 355, 11, 271, 354,
	81, 353, 352, 3, 12, 351, 350, 349, 348, 347,
	346, 345, 344, 343, 342, 341, 340, 5, 339, 2,
	338, 337, 8, 22, 336, 335, 334, 333, 332, 0,
	15, 331, 330, 13, 7, 329, 328, 327, 326, 9,
	325, 324, 321, 10, 320, 319, 4, 318, 267, 20,
	6, 243, 1,
}
var yyR1 = [...]int{

	0, 1, 2, 6, 4, 4, 7, 3, 3, 8,
	8, 8, 8, 9, 10, 10, 10, 5, 5, 14,
	14, 13, 13, 13, 13, 13, 13, 13, 13, 15,
	15, 23, 25, 25, 25, 24, 26, 26, 27, 28,
	16, 30, 31, 31, 32, 32, 32, 33, 34, 35,
	35, 36, 36, 19, 38, 29, 29, 39, 39, 39,
	22, 11, 41, 42, 42, 43, 43, 43, 43, 45,
	46, 12, 47, 48, 48, 49, 49, 49, 17, 51,
	50, 52, 52, 53, 53, 53, 18, 54, 55, 55,
	56, 56, 56, 56, 56, 56, 57, 20, 58, 59,
	59, 60, 60, 60, 60, 60, 40, 21, 61, 37,
	37, 62, 44,
}
var yyR2 = [...]int{

	0, 5, 3, 2, 2, 5, 2, 2, 3, 2,
	1, 1, 2, 2, 1, 1, 1, 1, 2, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	2, 4, 0, 2, 1, 2, 1, 2, 4, 2,
	4, 2, 1, 2, 1, 2, 3, 2, 2, 1,
	3, 3, 1, 4, 2, 1, 2, 2, 1, 1,
	3, 4, 2, 1, 2, 2, 1, 3, 3, 2,
	2, 4, 2, 1, 2, 2, 1, 1, 2, 3,
	2, 1, 2, 2, 1, 1, 4, 2, 1, 2,
	2, 3, 1, 1, 3, 1, 3, 4, 2, 1,
	2, 1, 2, 1, 3, 3, 3, 4, 2, 1,
	2, 3, 3,
}
var yyChk = [...]int{

	-1000, -1, -2, 25, -3, -8, 11, -7, -9, 15,
	12, 39, 4, -4, -8, -6, 13, 9, 5, 5,
	5, 4, 7, -5, -10, -11, -12, -13, -41, -47,
	-16, -17, -18, -19, -20, -21, -22, -23, 28, 29,
	-30, -50, -54, -38, -58, -61, 21, -24, 30, 31,
	27, 26, 32, 35, 37, 9, 9, 7, 10, 8,
	-10, 7, 7, 4, 4, 7, -51, 7, 7, 7,
	7, 7, 4, 7, 4, 4, 4, 4, 4, 4,
	4, -7, -42, -43, -7, -44, -45, -46, 34, 23,
	24, -48, -49, -7, -40, -13, 20, -31, -32, -33,
	-7, 16, -34, 14, -52, -53, -7, -44, -13, -55,
	-56, -7, 36, -40, -57, 22, -13, 19, -29, -39,
	-7, -40, -13, -59, -60, -33, -7, -40, 16, 33,
	-59, 9, -25, -7, -26, -27, -28, 38, 9, 8,
	-43, 9, -14, -15, -13, -14, 5, 7, 7, 8,
	-49, 9, 5, 8, -32, 9, 5, -35, 9, 7,
	4, 8, -53, 9, 8, -56, 9, 6, 5, 5,
	8, -39, 9, 8, -60, 9, 5, 5, 8, 8,
	9, -27, 7, 4, 8, 8, -13, 8, 9, 9,
	9, -36, 17, -37, -62, 18, 9, 9, 9, 9,
	9, -29, 8, 5, -62, 4, 8, 9, 9,
}
var yyDef = [...]int{

	0, -2, 0, 0, 0, 0, 0, 10, 11, 0,
	0, 0, 0, 0, 0, 0, 0, 7, 9, 12,
	6, 13, 2, 0, 17, 14, 15, 16, 0, 0,
	21, 22, 23, 24, 25, 26, 27, 28, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 8, 4, 0, 3, 1,
	18, 0, 0, 62, 72, 0, 78, 0, 0, 0,
	0, 0, 0, 32, 41, 80, 87, 54, 98, 108,
	35, 0, 0, 63, 0, 66, 19, 19, 0, 0,
	0, 0, 73, 0, 76, 77, 0, 0, 42, 44,
	0, 0, 0, 0, 0, 81, 0, 84, 85, 0,
	88, 0, 0, 92, 93, 0, 95, 0, 0, 55,
	0, 58, 59, 0, 99, 101, 0, 103, 0, 0,
	0, 60, 0, 0, 34, 36, 0, 0, 0, 61,
	64, 65, 0, 20, 29, 0, 0, 69, 70, 71,
	74, 75, 0, 40, 43, 45, 0, 47, 49, 0,
	48, 79, 82, 83, 86, 89, 90, 0, 0, 0,
	53, 56, 57, 97, 100, 102, 0, 0, 107, 31,
	33, 37, 0, 39, 5, 67, 30, 68, 112, 106,
	46, 0, 0, 52, 109, 0, 91, 94, 96, 104,
	105, 0, 50, 0, 110, 0, 38, 51, 111,
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
		//line parser.y:229
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:236
		{
			yylval.stack.Push(&schema.ChoiceCase{Ident: yyDollar[2].token})
		}
	case 40:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:244
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 41:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:251
		{
			yylval.stack.Push(&schema.Typedef{Ident: yyDollar[2].token})
		}
	case 47:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:266
		{
			y := yylval.stack.Peek().(schema.HasDataType)
			y.SetDataType(yylval.dataType)
		}
	case 48:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:271
		{
			yylval.dataType = schema.NewDataType(yyDollar[2].token)
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:280
		{
			var err error
			if err = yylval.dataType.DecodeLength(tokenString(yyDollar[2].token)); err != nil {
				yylex.Error(err.Error())
				goto ret1
			}
		}
	case 53:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:293
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 54:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:300
		{
			yylval.stack.Push(&schema.Container{Ident: yyDollar[2].token})
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:314
		{
			yylval.stack.Push(&schema.Uses{Ident: yyDollar[2].token})
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 61:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:326
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 62:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:333
		{
			yylval.stack.Push(&schema.Rpc{Ident: yyDollar[2].token})
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:344
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:349
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 69:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:356
		{
			yylval.stack.Push(&schema.RpcInput{})
		}
	case 70:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:361
		{
			yylval.stack.Push(&schema.RpcOutput{})
		}
	case 71:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:369
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 72:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:376
		{
			yylval.stack.Push(&schema.Notification{Ident: yyDollar[2].token})
		}
	case 78:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:392
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 80:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:404
		{
			yylval.stack.Push(&schema.Grouping{Ident: yyDollar[2].token})
		}
	case 86:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:420
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 87:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:427
		{
			yylval.stack.Push(&schema.List{Ident: yyDollar[2].token})
		}
	case 96:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:443
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
		//line parser.y:457
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 98:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:464
		{
			yylval.stack.Push(&schema.Leaf{Ident: yyDollar[2].token})
		}
	case 106:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:480
		{
			if hasDetails, valid := yylval.stack.Peek().(schema.HasDetails); valid {
				hasDetails.Details().SetConfig("true" == yyDollar[2].token)
			} else {
				yylex.Error("expected config statement on schema supporting details")
				goto ret1
			}
		}
	case 107:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:494
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 108:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:501
		{
			yylval.stack.Push(&schema.LeafList{Ident: yyDollar[2].token})
		}
	case 111:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:510
		{
			yylval.dataType.Enumeration = append(yylval.dataType.Enumeration, yyDollar[2].token)
		}
	}
	goto yystack /* stack new state and value */
}
