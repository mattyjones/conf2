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

//line parser.y:576

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 129
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 398

var yyAct = [...]int{

	126, 215, 125, 129, 147, 158, 142, 131, 117, 112,
	130, 98, 10, 144, 135, 127, 89, 27, 7, 132,
	7, 3, 180, 24, 91, 104, 10, 27, 6, 10,
	16, 216, 9, 124, 102, 47, 122, 10, 144, 177,
	53, 52, 128, 10, 50, 51, 54, 63, 10, 55,
	119, 56, 47, 61, 231, 57, 11, 53, 52, 152,
	153, 50, 51, 54, 230, 94, 55, 221, 56, 101,
	94, 174, 57, 173, 115, 123, 86, 220, 6, 10,
	90, 99, 9, 107, 106, 137, 113, 118, 105, 133,
	133, 219, 140, 148, 218, 114, 160, 160, 217, 161,
	10, 101, 149, 211, 90, 156, 11, 210, 100, 166,
	201, 95, 96, 99, 120, 115, 134, 134, 107, 106,
	123, 178, 94, 105, 195, 181, 187, 113, 191, 170,
	213, 216, 118, 188, 182, 193, 114, 179, 190, 10,
	100, 110, 60, 109, 59, 190, 133, 102, 196, 171,
	167, 200, 157, 133, 160, 160, 202, 203, 154, 120,
	136, 138, 148, 208, 169, 205, 58, 10, 10, 17,
	110, 149, 109, 134, 124, 102, 47, 122, 225, 224,
	134, 53, 52, 223, 209, 50, 51, 54, 207, 206,
	55, 119, 56, 199, 194, 229, 57, 183, 155, 10,
	222, 10, 204, 110, 197, 109, 164, 102, 47, 163,
	77, 76, 74, 53, 52, 73, 227, 50, 51, 54,
	226, 186, 55, 187, 56, 10, 72, 71, 57, 70,
	68, 65, 64, 102, 47, 22, 192, 185, 184, 53,
	52, 175, 168, 50, 51, 54, 162, 165, 55, 20,
	56, 10, 19, 18, 57, 228, 5, 198, 176, 102,
	47, 14, 85, 84, 62, 53, 52, 83, 82, 50,
	51, 54, 81, 46, 55, 80, 56, 47, 79, 78,
	57, 10, 53, 52, 39, 40, 50, 51, 54, 102,
	47, 55, 75, 56, 10, 53, 52, 57, 67, 50,
	51, 54, 66, 47, 55, 21, 56, 12, 53, 52,
	57, 45, 50, 51, 54, 121, 94, 55, 116, 56,
	47, 43, 111, 57, 69, 53, 52, 39, 40, 50,
	51, 54, 47, 42, 55, 97, 56, 53, 52, 29,
	57, 50, 51, 54, 151, 189, 55, 150, 56, 10,
	146, 110, 57, 109, 10, 145, 110, 102, 109, 49,
	93, 92, 102, 88, 87, 28, 44, 214, 212, 172,
	136, 108, 103, 41, 143, 136, 141, 139, 48, 38,
	37, 36, 35, 34, 33, 32, 31, 30, 159, 26,
	25, 8, 15, 23, 13, 4, 2, 1,
}
var yyPact = [...]int{

	-4, -1000, 67, 303, 17, 160, 248, -1000, -1000, 247,
	244, 301, 228, 299, 157, 135, 43, -1000, -1000, -1000,
	-1000, -1000, -1000, 256, -1000, -1000, -1000, -1000, 225, 224,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 298,
	294, 223, 222, 220, 219, 208, 205, 288, 204, 203,
	275, 274, 271, 268, 264, 263, 259, 258, -1000, -1000,
	25, -1000, -1000, -1000, 88, 269, -1000, -1000, 189, -1000,
	282, 155, 269, 342, 342, 152, 0, 36, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 149, 190, 88, -1000,
	143, -1000, 311, 311, 241, 202, 199, 239, -1000, 141,
	-1000, -1000, 237, 156, -1000, -1000, 140, -1000, 64, 236,
	254, 31, -1000, 128, -1000, -1000, 14, -1000, 125, 191,
	-1000, -1000, 233, -1000, 232, 213, -1000, 124, -1000, -1000,
	337, -1000, -1000, 119, -1000, -1000, 231, 127, -1000, 186,
	115, -25, -1000, 197, 253, 185, 36, -1000, 101, -1000,
	311, 311, 195, 158, 181, -1000, -1000, -1000, 180, 311,
	-1000, 176, 98, -1000, -1000, -1000, -1000, -1000, 94, -1000,
	-1000, -1000, -1000, -1000, 113, 89, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 85, 82, 68, -1000, -1000, -1000, -1000,
	-1000, -1000, 58, -1000, -1000, -1000, -1000, 269, -1000, -1000,
	-1000, -1000, 175, 171, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, 170, 215, 13, -1000, 251, -1000, -1000, -1000,
	-1000, -1000, 187, -1000, -1000, -1000, 55, -1000, 45, -1000,
	-1000, -1000,
}
var yyPgo = [...]int{

	0, 397, 396, 395, 394, 393, 392, 15, 256, 391,
	23, 390, 389, 3, 5, 388, 387, 386, 385, 384,
	383, 382, 381, 380, 379, 378, 377, 376, 6, 374,
	2, 373, 372, 25, 19, 14, 371, 369, 368, 367,
	366, 0, 42, 365, 364, 363, 16, 24, 361, 360,
	359, 355, 350, 4, 347, 344, 339, 335, 11, 333,
	324, 322, 9, 321, 318, 8, 315, 311, 10, 7,
	273, 1,
}
var yyR1 = [...]int{

	0, 1, 2, 6, 4, 4, 7, 3, 3, 8,
	8, 8, 8, 9, 10, 10, 10, 5, 5, 14,
	14, 13, 13, 13, 13, 13, 13, 13, 13, 13,
	15, 15, 23, 26, 26, 26, 25, 27, 27, 28,
	29, 16, 31, 32, 32, 33, 33, 33, 35, 34,
	36, 37, 37, 38, 38, 19, 40, 30, 30, 41,
	41, 41, 22, 11, 43, 44, 44, 45, 45, 46,
	46, 46, 46, 48, 49, 24, 50, 51, 51, 52,
	52, 53, 53, 53, 53, 54, 55, 12, 56, 57,
	57, 58, 58, 58, 17, 60, 59, 61, 61, 62,
	62, 62, 18, 63, 64, 64, 65, 65, 65, 65,
	65, 65, 66, 20, 67, 68, 68, 69, 69, 69,
	69, 69, 42, 21, 70, 39, 39, 71, 47,
}
var yyR2 = [...]int{

	0, 5, 3, 2, 2, 5, 2, 2, 3, 2,
	1, 1, 2, 2, 1, 1, 1, 1, 2, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 2, 4, 0, 2, 1, 2, 1, 2, 4,
	2, 4, 2, 1, 2, 1, 2, 1, 3, 2,
	2, 1, 3, 3, 1, 4, 2, 1, 2, 2,
	1, 1, 3, 4, 2, 0, 1, 1, 2, 2,
	1, 3, 3, 2, 2, 4, 2, 0, 1, 1,
	2, 2, 1, 3, 3, 2, 2, 4, 2, 1,
	2, 2, 1, 1, 2, 3, 2, 1, 2, 2,
	1, 1, 4, 2, 1, 2, 2, 3, 1, 1,
	3, 1, 3, 4, 2, 1, 2, 1, 2, 1,
	1, 3, 3, 4, 2, 1, 2, 3, 3,
}
var yyChk = [...]int{

	-1000, -1, -2, 25, -3, -8, 11, -7, -9, 15,
	12, 39, 4, -4, -8, -6, 13, 9, 5, 5,
	5, 4, 7, -5, -10, -11, -12, -13, -43, -56,
	-16, -17, -18, -19, -20, -21, -22, -23, -24, 28,
	29, -31, -59, -63, -40, -67, -70, 21, -25, -50,
	30, 31, 27, 26, 32, 35, 37, 41, 9, 9,
	7, 10, 8, -10, 7, 7, 4, 4, 7, -60,
	7, 7, 7, 7, 7, 4, 7, 7, 4, 4,
	4, 4, 4, 4, 4, 4, -7, -44, -45, -46,
	-7, -47, -48, -49, 34, 23, 24, -57, -58, -7,
	-42, -13, 20, -32, -33, -34, -7, -35, -36, 16,
	14, -61, -62, -7, -47, -13, -64, -65, -7, 36,
	-42, -66, 22, -13, 19, -30, -41, -7, -42, -13,
	-68, -69, -34, -7, -42, -35, 33, -68, 9, -26,
	-7, -27, -28, -29, 38, -51, -52, -53, -7, -47,
	-54, -55, 23, 24, 9, 8, -46, 9, -14, -15,
	-13, -14, 5, 7, 7, 8, -58, 9, 5, 8,
	-33, 9, -37, 9, 7, 5, 4, 8, -62, 9,
	8, -65, 9, 6, 5, 5, 8, -41, 9, 8,
	-69, 9, 5, 8, 8, 9, -28, 7, 4, 8,
	-53, 9, -14, -14, 7, 7, 8, 8, -13, 8,
	9, 9, -38, 17, -39, -71, 18, 9, 9, 9,
	9, 9, -30, 8, 8, 8, 5, -71, 4, 8,
	9, 9,
}
var yyDef = [...]int{

	0, -2, 0, 0, 0, 0, 0, 10, 11, 0,
	0, 0, 0, 0, 0, 0, 0, 7, 9, 12,
	6, 13, 2, 0, 17, 14, 15, 16, 0, 0,
	21, 22, 23, 24, 25, 26, 27, 28, 29, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 8, 4,
	0, 3, 1, 18, 65, 0, 64, 88, 0, 94,
	0, 0, 0, 0, 0, 0, 33, 77, 42, 96,
	103, 56, 114, 124, 36, 76, 0, 0, 66, 67,
	0, 70, 19, 19, 0, 0, 0, 0, 89, 0,
	92, 93, 0, 0, 43, 45, 0, 47, 0, 0,
	0, 0, 97, 0, 100, 101, 0, 104, 0, 0,
	108, 109, 0, 111, 0, 0, 57, 0, 60, 61,
	0, 115, 117, 0, 119, 120, 0, 0, 62, 0,
	0, 35, 37, 0, 0, 0, 78, 79, 0, 82,
	19, 19, 0, 0, 0, 63, 68, 69, 0, 20,
	30, 0, 0, 73, 74, 87, 90, 91, 0, 41,
	44, 46, 49, 51, 0, 0, 50, 95, 98, 99,
	102, 105, 106, 0, 0, 0, 55, 58, 59, 113,
	116, 118, 0, 123, 32, 34, 38, 0, 40, 75,
	80, 81, 0, 0, 85, 86, 5, 71, 31, 72,
	128, 122, 0, 0, 54, 125, 0, 48, 107, 110,
	112, 121, 0, 83, 84, 52, 0, 126, 0, 39,
	53, 127,
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
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:267
		{
			if hasType, valid := yylval.stack.Peek().(schema.HasDataType); valid {
				hasType.GetDataType().Default = tokenString(yyDollar[2].token)
			} else {
				yylex.Error("expected default statement on schema supporting details")
				goto ret1
			}
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:276
		{
			y := yylval.stack.Peek().(schema.HasDataType)
			y.SetDataType(yylval.dataType)
		}
	case 50:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:281
		{
			yylval.dataType = schema.NewDataType(yyDollar[2].token)
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:290
		{
			var err error
			if err = yylval.dataType.DecodeLength(tokenString(yyDollar[2].token)); err != nil {
				yylex.Error(err.Error())
				goto ret1
			}
		}
	case 55:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:303
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 56:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:310
		{
			yylval.stack.Push(&schema.Container{Ident: yyDollar[2].token})
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:324
		{
			yylval.stack.Push(&schema.Uses{Ident: yyDollar[2].token})
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 63:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:336
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 64:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:343
		{
			yylval.stack.Push(&schema.Rpc{Ident: yyDollar[2].token})
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:357
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 72:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:362
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 73:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:369
		{
			yylval.stack.Push(&schema.RpcInput{})
		}
	case 74:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:374
		{
			yylval.stack.Push(&schema.RpcOutput{})
		}
	case 75:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:382
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 76:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:389
		{
			yylval.stack.Push(&schema.Rpc{Ident: yyDollar[2].token})
		}
	case 83:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:403
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 84:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:408
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 85:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:415
		{
			yylval.stack.Push(&schema.RpcInput{})
		}
	case 86:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:420
		{
			yylval.stack.Push(&schema.RpcOutput{})
		}
	case 87:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:428
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 88:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:435
		{
			yylval.stack.Push(&schema.Notification{Ident: yyDollar[2].token})
		}
	case 94:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:451
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 96:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:463
		{
			yylval.stack.Push(&schema.Grouping{Ident: yyDollar[2].token})
		}
	case 102:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:479
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 103:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:486
		{
			yylval.stack.Push(&schema.List{Ident: yyDollar[2].token})
		}
	case 112:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:502
		{
			if list, valid := yylval.stack.Peek().(*schema.List); valid {
				list.Keys = strings.Split(tokenString(yyDollar[2].token), " ")
			} else {
				yylex.Error("expected a list for key statement")
				goto ret1
			}
		}
	case 113:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:516
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 114:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:523
		{
			yylval.stack.Push(&schema.Leaf{Ident: yyDollar[2].token})
		}
	case 122:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:539
		{
			if hasDetails, valid := yylval.stack.Peek().(schema.HasDetails); valid {
				hasDetails.Details().SetConfig("true" == yyDollar[2].token)
			} else {
				yylex.Error("expected config statement on schema supporting details")
				goto ret1
			}
		}
	case 123:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:553
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 124:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:560
		{
			yylval.stack.Push(&schema.LeafList{Ident: yyDollar[2].token})
		}
	case 127:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:569
		{
			yylval.dataType.Enumeration = append(yylval.dataType.Enumeration, yyDollar[2].token)
		}
	}
	goto yystack /* stack new state and value */
}
