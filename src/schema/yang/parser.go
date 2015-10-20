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
const kywd_action = 57383

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
	"kywd_action",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line parser.y:562

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 124
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 403

var yyAct = [...]int{

	124, 214, 123, 127, 144, 155, 140, 129, 103, 110,
	115, 125, 130, 90, 7, 97, 7, 27, 88, 128,
	177, 142, 3, 215, 10, 10, 10, 27, 24, 61,
	230, 122, 101, 47, 120, 126, 229, 174, 53, 52,
	220, 10, 50, 51, 54, 219, 10, 55, 117, 56,
	47, 142, 63, 57, 218, 53, 52, 149, 150, 50,
	51, 54, 217, 93, 55, 172, 56, 171, 93, 100,
	57, 60, 86, 59, 113, 121, 89, 98, 212, 215,
	105, 104, 111, 116, 112, 131, 131, 216, 138, 145,
	210, 146, 209, 208, 135, 157, 157, 199, 158, 89,
	100, 99, 193, 188, 185, 179, 153, 118, 98, 132,
	132, 167, 163, 113, 105, 104, 6, 10, 121, 175,
	9, 111, 176, 112, 184, 178, 116, 10, 168, 164,
	154, 224, 99, 6, 10, 16, 187, 9, 94, 95,
	131, 151, 136, 187, 11, 58, 194, 131, 198, 93,
	118, 157, 157, 200, 201, 145, 17, 146, 10, 197,
	206, 11, 223, 10, 132, 122, 101, 47, 120, 222,
	207, 132, 53, 52, 149, 150, 50, 51, 54, 205,
	204, 55, 117, 56, 228, 93, 166, 57, 10, 10,
	10, 108, 108, 106, 106, 192, 101, 47, 221, 203,
	202, 195, 53, 52, 161, 160, 50, 51, 54, 77,
	76, 55, 183, 56, 180, 226, 10, 57, 74, 10,
	73, 108, 184, 133, 101, 47, 72, 101, 152, 71,
	53, 52, 10, 70, 50, 51, 54, 225, 162, 55,
	134, 56, 10, 94, 95, 57, 68, 65, 64, 22,
	101, 47, 46, 190, 93, 62, 53, 52, 189, 182,
	50, 51, 54, 181, 169, 55, 165, 56, 47, 159,
	20, 57, 10, 53, 52, 39, 40, 50, 51, 54,
	101, 47, 55, 19, 56, 10, 53, 52, 57, 18,
	50, 51, 54, 5, 47, 55, 227, 56, 14, 53,
	52, 57, 196, 50, 51, 54, 45, 93, 55, 173,
	56, 47, 85, 84, 57, 83, 53, 52, 39, 40,
	50, 51, 54, 47, 82, 55, 81, 56, 53, 52,
	80, 57, 50, 51, 54, 79, 191, 55, 78, 56,
	10, 186, 108, 57, 133, 10, 75, 108, 101, 133,
	67, 66, 21, 101, 12, 119, 114, 43, 109, 69,
	42, 134, 96, 29, 148, 147, 134, 143, 49, 92,
	91, 87, 28, 44, 213, 211, 170, 107, 102, 41,
	141, 139, 137, 48, 38, 37, 36, 35, 34, 33,
	32, 31, 30, 156, 26, 25, 8, 15, 23, 13,
	4, 2, 1,
}
var yyPact = [...]int{

	-3, -1000, 105, 350, 122, 147, 284, -1000, -1000, 278,
	265, 348, 242, 290, 136, 64, 19, -1000, -1000, -1000,
	-1000, -1000, -1000, 247, -1000, -1000, -1000, -1000, 241, 240,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 347,
	346, 239, 226, 222, 219, 213, 211, 342, 203, 202,
	334, 331, 326, 322, 320, 311, 309, 308, -1000, -1000,
	14, -1000, -1000, -1000, 115, 260, -1000, -1000, 177, -1000,
	273, 146, 260, 207, 207, 133, 13, 34, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 132, 220, -1000, 121,
	-1000, 302, 302, 264, 198, 197, 230, -1000, 120, -1000,
	-1000, 261, 178, -1000, -1000, 119, 259, 58, 305, 29,
	-1000, 113, -1000, -1000, 12, -1000, 96, 208, -1000, -1000,
	258, -1000, 254, 204, -1000, 95, -1000, -1000, 333, -1000,
	-1000, 94, -1000, 253, 248, 328, -1000, 187, 93, -17,
	-1000, 194, 298, 151, -1000, 88, -1000, 302, 302, 193,
	192, 172, -1000, -1000, -1000, 171, 302, -1000, 162, 84,
	-1000, -1000, -1000, -1000, -1000, 83, -1000, -1000, -1000, 81,
	-1000, -1000, 61, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	78, 53, 45, -1000, -1000, -1000, -1000, -1000, -1000, 36,
	31, -1000, -1000, -1000, -1000, 260, -1000, -1000, -1000, -1000,
	161, 154, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 123, 232, 5, -1000, 292, -1000, -1000, -1000, -1000,
	-1000, 176, -1000, -1000, -1000, 27, -1000, 21, -1000, -1000,
	-1000,
}
var yyPgo = [...]int{

	0, 402, 401, 400, 399, 398, 397, 11, 293, 396,
	28, 395, 394, 3, 5, 393, 392, 391, 390, 389,
	388, 387, 386, 385, 384, 383, 382, 381, 6, 380,
	2, 379, 378, 8, 12, 377, 376, 375, 374, 373,
	0, 35, 372, 371, 18, 13, 370, 369, 368, 367,
	4, 365, 364, 363, 362, 15, 360, 359, 358, 9,
	357, 356, 10, 355, 306, 19, 7, 252, 1,
}
var yyR1 = [...]int{

	0, 1, 2, 6, 4, 4, 7, 3, 3, 8,
	8, 8, 8, 9, 10, 10, 10, 5, 5, 14,
	14, 13, 13, 13, 13, 13, 13, 13, 13, 13,
	15, 15, 23, 26, 26, 26, 25, 27, 27, 28,
	29, 16, 31, 32, 32, 33, 33, 33, 34, 35,
	36, 36, 37, 37, 19, 39, 30, 30, 40, 40,
	40, 22, 11, 42, 43, 43, 44, 44, 44, 44,
	46, 47, 24, 48, 49, 49, 50, 50, 50, 50,
	51, 52, 12, 53, 54, 54, 55, 55, 55, 17,
	57, 56, 58, 58, 59, 59, 59, 18, 60, 61,
	61, 62, 62, 62, 62, 62, 62, 63, 20, 64,
	65, 65, 66, 66, 66, 66, 66, 41, 21, 67,
	38, 38, 68, 45,
}
var yyR2 = [...]int{

	0, 5, 3, 2, 2, 5, 2, 2, 3, 2,
	1, 1, 2, 2, 1, 1, 1, 1, 2, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 2, 4, 0, 2, 1, 2, 1, 2, 4,
	2, 4, 2, 1, 2, 1, 2, 3, 2, 2,
	1, 3, 3, 1, 4, 2, 1, 2, 2, 1,
	1, 3, 4, 2, 1, 2, 2, 1, 3, 3,
	2, 2, 4, 2, 1, 2, 2, 1, 3, 3,
	2, 2, 4, 2, 1, 2, 2, 1, 1, 2,
	3, 2, 1, 2, 2, 1, 1, 4, 2, 1,
	2, 2, 3, 1, 1, 3, 1, 3, 4, 2,
	1, 2, 1, 2, 1, 3, 3, 3, 4, 2,
	1, 2, 3, 3,
}
var yyChk = [...]int{

	-1000, -1, -2, 25, -3, -8, 11, -7, -9, 15,
	12, 39, 4, -4, -8, -6, 13, 9, 5, 5,
	5, 4, 7, -5, -10, -11, -12, -13, -42, -53,
	-16, -17, -18, -19, -20, -21, -22, -23, -24, 28,
	29, -31, -56, -60, -39, -64, -67, 21, -25, -48,
	30, 31, 27, 26, 32, 35, 37, 41, 9, 9,
	7, 10, 8, -10, 7, 7, 4, 4, 7, -57,
	7, 7, 7, 7, 7, 4, 7, 7, 4, 4,
	4, 4, 4, 4, 4, 4, -7, -43, -44, -7,
	-45, -46, -47, 34, 23, 24, -54, -55, -7, -41,
	-13, 20, -32, -33, -34, -7, 16, -35, 14, -58,
	-59, -7, -45, -13, -61, -62, -7, 36, -41, -63,
	22, -13, 19, -30, -40, -7, -41, -13, -65, -66,
	-34, -7, -41, 16, 33, -65, 9, -26, -7, -27,
	-28, -29, 38, -49, -50, -7, -45, -51, -52, 23,
	24, 9, 8, -44, 9, -14, -15, -13, -14, 5,
	7, 7, 8, -55, 9, 5, 8, -33, 9, 5,
	-36, 9, 7, 4, 8, -59, 9, 8, -62, 9,
	6, 5, 5, 8, -40, 9, 8, -66, 9, 5,
	5, 8, 8, 9, -28, 7, 4, 8, -50, 9,
	-14, -14, 7, 7, 8, 8, -13, 8, 9, 9,
	9, -37, 17, -38, -68, 18, 9, 9, 9, 9,
	9, -30, 8, 8, 8, 5, -68, 4, 8, 9,
	9,
}
var yyDef = [...]int{

	0, -2, 0, 0, 0, 0, 0, 10, 11, 0,
	0, 0, 0, 0, 0, 0, 0, 7, 9, 12,
	6, 13, 2, 0, 17, 14, 15, 16, 0, 0,
	21, 22, 23, 24, 25, 26, 27, 28, 29, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 8, 4,
	0, 3, 1, 18, 0, 0, 63, 83, 0, 89,
	0, 0, 0, 0, 0, 0, 33, 0, 42, 91,
	98, 55, 109, 119, 36, 73, 0, 0, 64, 0,
	67, 19, 19, 0, 0, 0, 0, 84, 0, 87,
	88, 0, 0, 43, 45, 0, 0, 0, 0, 0,
	92, 0, 95, 96, 0, 99, 0, 0, 103, 104,
	0, 106, 0, 0, 56, 0, 59, 60, 0, 110,
	112, 0, 114, 0, 0, 0, 61, 0, 0, 35,
	37, 0, 0, 0, 74, 0, 77, 19, 19, 0,
	0, 0, 62, 65, 66, 0, 20, 30, 0, 0,
	70, 71, 82, 85, 86, 0, 41, 44, 46, 0,
	48, 50, 0, 49, 90, 93, 94, 97, 100, 101,
	0, 0, 0, 54, 57, 58, 108, 111, 113, 0,
	0, 118, 32, 34, 38, 0, 40, 72, 75, 76,
	0, 0, 80, 81, 5, 68, 31, 69, 123, 117,
	47, 0, 0, 53, 120, 0, 102, 105, 107, 115,
	116, 0, 78, 79, 51, 0, 121, 0, 39, 52,
	122,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
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
		//line parser.y:122
		{
			m := &schema.Module{Ident: yyDollar[2].token}
			yylval.stack.Push(m)
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:128
		{
			d := yylval.stack.Peek()
			r := &schema.Revision{Ident: yyDollar[2].token}
			d.(*schema.Module).Revision = r
			yylval.stack.Push(r)
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:136
		{
			yylval.stack.Pop()
		}
	case 5:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:139
		{
			yylval.stack.Pop()
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:143
		{
			yylval.stack.Peek().(schema.Describable).SetDescription(tokenString(yyDollar[2].token))
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:153
		{
			d := yylval.stack.Peek()
			d.(*schema.Module).Namespace = tokenString(yyDollar[2].token)
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:159
		{
			m := yylval.stack.Peek().(*schema.Module)
			m.Prefix = tokenString(yyDollar[2].token)
		}
	case 13:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:164
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
	case 32:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:209
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:221
		{
			yylval.stack.Push(&schema.Choice{Ident: yyDollar[2].token})
		}
	case 39:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:231
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 40:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:238
		{
			yylval.stack.Push(&schema.ChoiceCase{Ident: yyDollar[2].token})
		}
	case 41:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:246
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 42:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:253
		{
			yylval.stack.Push(&schema.Typedef{Ident: yyDollar[2].token})
		}
	case 48:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:268
		{
			y := yylval.stack.Peek().(schema.HasDataType)
			y.SetDataType(yylval.dataType)
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:273
		{
			yylval.dataType = schema.NewDataType(yyDollar[2].token)
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:282
		{
			var err error
			if err = yylval.dataType.DecodeLength(tokenString(yyDollar[2].token)); err != nil {
				yylex.Error(err.Error())
				goto ret1
			}
		}
	case 54:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:295
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 55:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:302
		{
			yylval.stack.Push(&schema.Container{Ident: yyDollar[2].token})
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:316
		{
			yylval.stack.Push(&schema.Uses{Ident: yyDollar[2].token})
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 62:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:328
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 63:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:335
		{
			yylval.stack.Push(&schema.Rpc{Ident: yyDollar[2].token})
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:346
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:351
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 70:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:358
		{
			yylval.stack.Push(&schema.RpcInput{})
		}
	case 71:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:363
		{
			yylval.stack.Push(&schema.RpcOutput{})
		}
	case 72:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:371
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 73:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:378
		{
			yylval.stack.Push(&schema.Rpc{Ident: yyDollar[2].token})
		}
	case 78:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:389
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:394
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 80:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:401
		{
			yylval.stack.Push(&schema.RpcInput{})
		}
	case 81:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:406
		{
			yylval.stack.Push(&schema.RpcOutput{})
		}
	case 82:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:414
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 83:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:421
		{
			yylval.stack.Push(&schema.Notification{Ident: yyDollar[2].token})
		}
	case 89:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:437
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 91:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:449
		{
			yylval.stack.Push(&schema.Grouping{Ident: yyDollar[2].token})
		}
	case 97:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:465
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 98:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:472
		{
			yylval.stack.Push(&schema.List{Ident: yyDollar[2].token})
		}
	case 107:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:488
		{
			if list, valid := yylval.stack.Peek().(*schema.List); valid {
				list.Keys = strings.Split(tokenString(yyDollar[2].token), " ")
			} else {
				yylex.Error("expected a list for key statement")
				goto ret1
			}
		}
	case 108:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:502
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 109:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:509
		{
			yylval.stack.Push(&schema.Leaf{Ident: yyDollar[2].token})
		}
	case 117:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:525
		{
			if hasDetails, valid := yylval.stack.Peek().(schema.HasDetails); valid {
				hasDetails.Details().SetConfig("true" == yyDollar[2].token)
			} else {
				yylex.Error("expected config statement on schema supporting details")
				goto ret1
			}
		}
	case 118:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:539
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 119:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:546
		{
			yylval.stack.Push(&schema.LeafList{Ident: yyDollar[2].token})
		}
	case 122:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:555
		{
			yylval.dataType.Enumeration = append(yylval.dataType.Enumeration, yyDollar[2].token)
		}
	}
	goto yystack /* stack new state and value */
}
