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
const kywd_anyxml = 57384
const kywd_path = 57385

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
	"kywd_anyxml",
	"kywd_path",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line parser.y:594

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 133
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 416

var yyAct = [...]int{

	131, 221, 130, 134, 152, 163, 147, 136, 122, 117,
	140, 103, 96, 137, 109, 132, 94, 27, 7, 149,
	7, 185, 135, 3, 222, 10, 10, 27, 10, 24,
	133, 64, 129, 107, 49, 127, 218, 222, 239, 55,
	54, 238, 237, 52, 53, 56, 227, 226, 57, 124,
	59, 10, 149, 66, 60, 58, 225, 224, 129, 107,
	49, 127, 220, 6, 10, 55, 54, 9, 223, 52,
	53, 56, 106, 216, 57, 124, 59, 120, 128, 91,
	60, 58, 112, 95, 104, 110, 119, 111, 215, 118,
	123, 11, 138, 138, 154, 206, 145, 153, 200, 105,
	142, 165, 165, 196, 166, 125, 106, 139, 139, 95,
	161, 6, 10, 16, 171, 9, 10, 193, 104, 112,
	120, 187, 110, 175, 111, 128, 183, 157, 158, 119,
	186, 192, 118, 105, 184, 10, 174, 123, 99, 11,
	10, 176, 115, 195, 114, 172, 100, 101, 162, 159,
	195, 138, 125, 201, 143, 78, 205, 99, 138, 165,
	165, 207, 208, 61, 154, 236, 139, 153, 213, 10,
	17, 231, 10, 139, 115, 230, 114, 107, 49, 179,
	107, 178, 229, 55, 54, 214, 212, 52, 53, 56,
	211, 191, 57, 141, 59, 10, 204, 199, 60, 58,
	63, 210, 62, 107, 49, 228, 160, 209, 202, 55,
	54, 169, 168, 52, 53, 56, 81, 80, 57, 182,
	59, 233, 77, 10, 60, 58, 10, 76, 115, 192,
	114, 75, 49, 74, 73, 71, 68, 55, 54, 67,
	22, 52, 53, 56, 188, 99, 57, 170, 59, 234,
	198, 10, 60, 58, 10, 232, 115, 197, 114, 107,
	49, 235, 107, 190, 65, 55, 54, 189, 180, 52,
	53, 56, 173, 167, 57, 141, 59, 49, 20, 19,
	60, 58, 55, 54, 40, 41, 52, 53, 56, 18,
	5, 57, 10, 59, 203, 14, 181, 60, 58, 90,
	107, 49, 89, 88, 87, 10, 55, 54, 86, 85,
	52, 53, 56, 84, 49, 57, 83, 59, 82, 55,
	54, 60, 58, 52, 53, 56, 79, 99, 57, 70,
	59, 49, 69, 21, 60, 58, 55, 54, 40, 41,
	52, 53, 56, 12, 49, 57, 47, 59, 46, 55,
	54, 60, 58, 52, 53, 56, 48, 126, 57, 121,
	59, 44, 194, 116, 60, 58, 10, 72, 115, 43,
	114, 102, 29, 156, 107, 155, 151, 150, 51, 98,
	97, 93, 92, 28, 45, 219, 217, 141, 177, 113,
	108, 42, 148, 146, 144, 50, 39, 38, 37, 36,
	35, 34, 33, 32, 31, 30, 164, 26, 25, 8,
	15, 23, 13, 4, 2, 1,
}
var yyPact = [...]int{

	-2, -1000, 52, 339, 100, 161, 284, -1000, -1000, 274,
	273, 329, 233, 310, 154, 193, 21, -1000, -1000, -1000,
	-1000, -1000, -1000, 256, -1000, -1000, -1000, -1000, 232, 229,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	328, 325, 228, 227, 226, 224, 220, 215, 146, 322,
	210, 209, 314, 312, 309, 305, 304, 300, 299, 298,
	295, -1000, -1000, 16, -1000, -1000, -1000, 123, 280, -1000,
	-1000, 214, -1000, 293, 39, 280, 160, 160, -1000, 145,
	14, 104, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 140, 198, 123, -1000, 139, -1000, 323, 323, 268,
	205, 204, 239, -1000, 136, -1000, -1000, 267, 128, -1000,
	-1000, 132, -1000, 172, 263, 292, 211, -1000, 125, -1000,
	-1000, 13, -1000, 112, 238, -1000, -1000, 262, -1000, 258,
	183, -1000, 108, -1000, -1000, 354, -1000, -1000, 94, -1000,
	-1000, 252, 242, -1000, 189, 89, -19, -1000, 201, 290,
	188, 104, -1000, 86, -1000, 323, 323, 200, 194, 182,
	-1000, -1000, -1000, 178, 323, -1000, 177, 79, -1000, -1000,
	-1000, -1000, -1000, 64, -1000, -1000, -1000, -1000, -1000, 19,
	59, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 48, 47,
	38, -1000, -1000, -1000, -1000, -1000, -1000, 37, -1000, -1000,
	-1000, -1000, 280, -1000, -1000, -1000, -1000, 174, 167, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 163, 250, 6,
	244, -1000, 257, -1000, -1000, -1000, -1000, -1000, 157, -1000,
	-1000, -1000, 33, -1000, 32, 29, -1000, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 415, 414, 413, 412, 411, 410, 15, 290, 409,
	29, 408, 407, 3, 5, 406, 405, 404, 403, 402,
	401, 400, 399, 398, 397, 396, 395, 394, 393, 6,
	392, 2, 391, 390, 14, 13, 10, 389, 388, 386,
	385, 384, 0, 30, 383, 382, 381, 16, 12, 380,
	379, 378, 377, 376, 4, 375, 373, 372, 371, 11,
	369, 367, 363, 9, 361, 359, 8, 357, 356, 348,
	22, 7, 346, 1,
}
var yyR1 = [...]int{

	0, 1, 2, 6, 4, 4, 7, 3, 3, 8,
	8, 8, 8, 9, 10, 10, 10, 5, 5, 14,
	14, 13, 13, 13, 13, 13, 13, 13, 13, 13,
	13, 15, 15, 24, 27, 27, 27, 26, 28, 28,
	29, 30, 16, 32, 33, 33, 34, 34, 34, 36,
	35, 37, 38, 38, 39, 39, 39, 19, 41, 31,
	31, 42, 42, 42, 23, 11, 44, 45, 45, 46,
	46, 47, 47, 47, 47, 49, 50, 25, 51, 52,
	52, 53, 53, 54, 54, 54, 54, 55, 56, 12,
	57, 58, 58, 59, 59, 59, 17, 61, 60, 62,
	62, 63, 63, 63, 18, 64, 65, 65, 66, 66,
	66, 66, 66, 66, 67, 22, 68, 20, 69, 70,
	70, 71, 71, 71, 71, 71, 43, 21, 72, 40,
	40, 73, 48,
}
var yyR2 = [...]int{

	0, 5, 3, 2, 2, 5, 2, 2, 3, 2,
	1, 1, 2, 2, 1, 1, 1, 1, 2, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 2, 4, 0, 2, 1, 2, 1, 2,
	4, 2, 4, 2, 1, 2, 1, 2, 1, 3,
	2, 2, 1, 3, 3, 1, 3, 4, 2, 1,
	2, 2, 1, 1, 3, 4, 2, 0, 1, 1,
	2, 2, 1, 3, 3, 2, 2, 4, 2, 0,
	1, 1, 2, 2, 1, 3, 3, 2, 2, 4,
	2, 1, 2, 2, 1, 1, 2, 3, 2, 1,
	2, 2, 1, 1, 4, 2, 1, 2, 2, 3,
	1, 1, 3, 1, 3, 2, 2, 4, 2, 1,
	2, 1, 2, 1, 1, 3, 3, 4, 2, 1,
	2, 3, 3,
}
var yyChk = [...]int{

	-1000, -1, -2, 25, -3, -8, 11, -7, -9, 15,
	12, 39, 4, -4, -8, -6, 13, 9, 5, 5,
	5, 4, 7, -5, -10, -11, -12, -13, -44, -57,
	-16, -17, -18, -19, -20, -21, -22, -23, -24, -25,
	28, 29, -32, -60, -64, -41, -69, -72, -68, 21,
	-26, -51, 30, 31, 27, 26, 32, 35, 42, 37,
	41, 9, 9, 7, 10, 8, -10, 7, 7, 4,
	4, 7, -61, 7, 7, 7, 7, 7, 9, 4,
	7, 7, 4, 4, 4, 4, 4, 4, 4, 4,
	4, -7, -45, -46, -47, -7, -48, -49, -50, 34,
	23, 24, -58, -59, -7, -43, -13, 20, -33, -34,
	-35, -7, -36, -37, 16, 14, -62, -63, -7, -48,
	-13, -65, -66, -7, 36, -43, -67, 22, -13, 19,
	-31, -42, -7, -43, -13, -70, -71, -35, -7, -43,
	-36, 33, -70, 9, -27, -7, -28, -29, -30, 38,
	-52, -53, -54, -7, -48, -55, -56, 23, 24, 9,
	8, -47, 9, -14, -15, -13, -14, 5, 7, 7,
	8, -59, 9, 5, 8, -34, 9, -38, 9, 7,
	5, 4, 8, -63, 9, 8, -66, 9, 6, 5,
	5, 8, -42, 9, 8, -71, 9, 5, 8, 8,
	9, -29, 7, 4, 8, -54, 9, -14, -14, 7,
	7, 8, 8, -13, 8, 9, 9, -39, 17, -40,
	43, -73, 18, 9, 9, 9, 9, 9, -31, 8,
	8, 8, 5, -73, 5, 4, 8, 9, 9, 9,
}
var yyDef = [...]int{

	0, -2, 0, 0, 0, 0, 0, 10, 11, 0,
	0, 0, 0, 0, 0, 0, 0, 7, 9, 12,
	6, 13, 2, 0, 17, 14, 15, 16, 0, 0,
	21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 8, 4, 0, 3, 1, 18, 67, 0, 66,
	90, 0, 96, 0, 0, 0, 0, 0, 115, 0,
	34, 79, 43, 98, 105, 58, 118, 128, 116, 37,
	78, 0, 0, 68, 69, 0, 72, 19, 19, 0,
	0, 0, 0, 91, 0, 94, 95, 0, 0, 44,
	46, 0, 48, 0, 0, 0, 0, 99, 0, 102,
	103, 0, 106, 0, 0, 110, 111, 0, 113, 0,
	0, 59, 0, 62, 63, 0, 119, 121, 0, 123,
	124, 0, 0, 64, 0, 0, 36, 38, 0, 0,
	0, 80, 81, 0, 84, 19, 19, 0, 0, 0,
	65, 70, 71, 0, 20, 31, 0, 0, 75, 76,
	89, 92, 93, 0, 42, 45, 47, 50, 52, 0,
	0, 51, 97, 100, 101, 104, 107, 108, 0, 0,
	0, 57, 60, 61, 117, 120, 122, 0, 127, 33,
	35, 39, 0, 41, 77, 82, 83, 0, 0, 87,
	88, 5, 73, 32, 74, 132, 126, 0, 0, 55,
	0, 129, 0, 49, 109, 112, 114, 125, 0, 85,
	86, 53, 0, 130, 0, 0, 40, 54, 56, 131,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43,
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
		//line parser.y:124
		{
			m := &schema.Module{Ident: yyDollar[2].token}
			yylval.stack.Push(m)
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:130
		{
			d := yylval.stack.Peek()
			r := &schema.Revision{Ident: yyDollar[2].token}
			d.(*schema.Module).Revision = r
			yylval.stack.Push(r)
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:138
		{
			yylval.stack.Pop()
		}
	case 5:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:141
		{
			yylval.stack.Pop()
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:145
		{
			yylval.stack.Peek().(schema.Describable).SetDescription(tokenString(yyDollar[2].token))
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:155
		{
			d := yylval.stack.Peek()
			d.(*schema.Module).Namespace = tokenString(yyDollar[2].token)
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:161
		{
			m := yylval.stack.Peek().(*schema.Module)
			m.Prefix = tokenString(yyDollar[2].token)
		}
	case 13:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:166
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
	case 33:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:212
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:224
		{
			yylval.stack.Push(&schema.Choice{Ident: yyDollar[2].token})
		}
	case 40:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:234
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 41:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:241
		{
			yylval.stack.Push(&schema.ChoiceCase{Ident: yyDollar[2].token})
		}
	case 42:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:249
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 43:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:256
		{
			yylval.stack.Push(&schema.Typedef{Ident: yyDollar[2].token})
		}
	case 49:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:270
		{
			if hasType, valid := yylval.stack.Peek().(schema.HasDataType); valid {
				hasType.GetDataType().SetDefault(tokenString(yyDollar[2].token))
			} else {
				yylex.Error("expected default statement on schema supporting details")
				goto ret1
			}
		}
	case 50:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:279
		{
			y := yylval.stack.Peek().(schema.HasDataType)
			y.SetDataType(yylval.dataType)
		}
	case 51:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:284
		{
			y := yylval.stack.Peek().(schema.HasDataType)
			yylval.dataType = schema.NewDataType(y, yyDollar[2].token)
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:294
		{
			var err error
			if err = yylval.dataType.DecodeLength(tokenString(yyDollar[2].token)); err != nil {
				yylex.Error(err.Error())
				goto ret1
			}
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:302
		{
			yylval.dataType.SetPath(tokenString(yyDollar[2].token))
		}
	case 57:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:310
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 58:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:317
		{
			yylval.stack.Push(&schema.Container{Ident: yyDollar[2].token})
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:331
		{
			yylval.stack.Push(&schema.Uses{Ident: yyDollar[2].token})
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 65:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:343
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 66:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:350
		{
			yylval.stack.Push(&schema.Rpc{Ident: yyDollar[2].token})
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:364
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:369
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 75:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:376
		{
			yylval.stack.Push(&schema.RpcInput{})
		}
	case 76:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:381
		{
			yylval.stack.Push(&schema.RpcOutput{})
		}
	case 77:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:389
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 78:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:396
		{
			yylval.stack.Push(&schema.Rpc{Ident: yyDollar[2].token})
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:410
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 86:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:415
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 87:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:422
		{
			yylval.stack.Push(&schema.RpcInput{})
		}
	case 88:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:427
		{
			yylval.stack.Push(&schema.RpcOutput{})
		}
	case 89:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:435
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 90:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:442
		{
			yylval.stack.Push(&schema.Notification{Ident: yyDollar[2].token})
		}
	case 96:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:458
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 98:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:470
		{
			yylval.stack.Push(&schema.Grouping{Ident: yyDollar[2].token})
		}
	case 104:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:486
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 105:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:493
		{
			yylval.stack.Push(&schema.List{Ident: yyDollar[2].token})
		}
	case 114:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:509
		{
			if list, valid := yylval.stack.Peek().(*schema.List); valid {
				list.Keys = strings.Split(tokenString(yyDollar[2].token), " ")
			} else {
				yylex.Error("expected a list for key statement")
				goto ret1
			}
		}
	case 115:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:519
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 116:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:526
		{
			yylval.stack.Push(schema.NewAny(yyDollar[2].token))
		}
	case 117:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:534
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 118:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:541
		{
			yylval.stack.Push(&schema.Leaf{Ident: yyDollar[2].token})
		}
	case 126:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:557
		{
			if hasDetails, valid := yylval.stack.Peek().(schema.HasDetails); valid {
				hasDetails.Details().SetConfig("true" == yyDollar[2].token)
			} else {
				yylex.Error("expected config statement on schema supporting details")
				goto ret1
			}
		}
	case 127:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:571
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 128:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:578
		{
			yylval.stack.Push(&schema.LeafList{Ident: yyDollar[2].token})
		}
	case 131:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:587
		{
			yylval.dataType.AddEnumeration(yyDollar[2].token)
		}
	}
	goto yystack /* stack new state and value */
}
