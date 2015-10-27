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

//line parser.y:568

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 125
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 404

var yyAct = [...]int{

	125, 213, 124, 128, 145, 90, 156, 116, 130, 134,
	111, 97, 131, 141, 127, 103, 129, 27, 6, 10,
	16, 143, 9, 10, 3, 88, 24, 27, 214, 178,
	211, 214, 126, 10, 10, 7, 109, 7, 108, 10,
	123, 101, 47, 121, 61, 229, 11, 53, 52, 143,
	63, 50, 51, 54, 228, 10, 55, 118, 56, 6,
	10, 219, 57, 9, 167, 223, 150, 151, 10, 100,
	109, 172, 108, 171, 114, 122, 113, 93, 106, 218,
	99, 104, 60, 147, 59, 217, 119, 11, 133, 133,
	216, 136, 215, 86, 209, 158, 158, 89, 98, 159,
	100, 105, 208, 112, 117, 199, 132, 132, 164, 139,
	146, 99, 106, 154, 114, 104, 113, 193, 168, 122,
	89, 176, 189, 179, 186, 185, 180, 177, 169, 98,
	119, 165, 155, 152, 137, 105, 58, 17, 188, 222,
	221, 207, 197, 112, 133, 188, 10, 205, 117, 198,
	147, 133, 158, 158, 194, 200, 201, 150, 151, 204,
	203, 206, 132, 192, 10, 202, 195, 162, 93, 132,
	161, 123, 101, 47, 121, 77, 76, 146, 53, 52,
	74, 73, 50, 51, 54, 47, 72, 55, 118, 56,
	53, 52, 71, 57, 50, 51, 54, 70, 220, 55,
	227, 56, 181, 191, 10, 57, 68, 10, 65, 109,
	64, 108, 101, 47, 225, 101, 153, 22, 53, 52,
	10, 185, 50, 51, 54, 224, 184, 55, 135, 56,
	10, 94, 95, 57, 190, 183, 182, 173, 101, 47,
	46, 166, 93, 175, 53, 52, 10, 10, 50, 51,
	54, 160, 20, 55, 19, 56, 47, 94, 95, 57,
	18, 53, 52, 226, 196, 50, 51, 54, 93, 93,
	55, 163, 56, 174, 187, 10, 57, 85, 10, 84,
	109, 83, 108, 101, 47, 82, 101, 5, 62, 53,
	52, 81, 14, 50, 51, 54, 80, 45, 55, 135,
	56, 47, 79, 78, 57, 10, 53, 52, 39, 40,
	50, 51, 54, 101, 47, 55, 75, 56, 10, 53,
	52, 57, 67, 50, 51, 54, 66, 47, 55, 21,
	56, 12, 53, 52, 57, 120, 50, 51, 54, 115,
	93, 55, 43, 56, 47, 110, 69, 57, 42, 53,
	52, 39, 40, 50, 51, 54, 96, 10, 55, 109,
	56, 108, 29, 149, 57, 101, 148, 144, 49, 92,
	91, 87, 28, 44, 212, 210, 170, 107, 135, 102,
	41, 142, 140, 138, 48, 38, 37, 36, 35, 34,
	33, 32, 31, 30, 157, 26, 25, 8, 15, 23,
	13, 4, 2, 1,
}
var yyPact = [...]int{

	-1, -1000, 48, 327, 7, 128, 255, -1000, -1000, 249,
	247, 325, 210, 323, 127, 75, 34, -1000, -1000, -1000,
	-1000, -1000, -1000, 280, -1000, -1000, -1000, -1000, 203, 201,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 322,
	318, 199, 190, 185, 179, 174, 173, 312, 169, 168,
	299, 298, 292, 287, 281, 277, 275, 273, -1000, -1000,
	27, -1000, -1000, -1000, 234, 293, -1000, -1000, 22, -1000,
	306, 152, 293, 345, 345, 125, 11, 43, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 124, 208, -1000, 123,
	-1000, 164, 164, 246, 163, 160, 263, -1000, 122, -1000,
	-1000, 236, 56, -1000, -1000, 119, -1000, 64, 232, 269,
	235, -1000, 118, -1000, -1000, 21, -1000, 117, 196, -1000,
	-1000, 231, -1000, 230, 218, -1000, 115, -1000, -1000, 266,
	-1000, -1000, 113, -1000, -1000, 229, 195, -1000, 155, 108,
	-17, -1000, 159, 260, 134, -1000, 96, -1000, 164, 164,
	158, 153, 151, -1000, -1000, -1000, 139, 164, -1000, 133,
	93, -1000, -1000, -1000, -1000, -1000, 85, -1000, -1000, -1000,
	-1000, -1000, 13, 83, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 81, 76, 70, -1000, -1000, -1000, -1000, -1000, -1000,
	52, -1000, -1000, -1000, -1000, 293, -1000, -1000, -1000, -1000,
	132, 131, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	57, 220, 10, -1000, 259, -1000, -1000, -1000, -1000, -1000,
	192, -1000, -1000, -1000, 45, -1000, 36, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 403, 402, 401, 400, 399, 398, 32, 287, 397,
	26, 396, 395, 3, 6, 394, 393, 392, 391, 390,
	389, 388, 387, 386, 385, 384, 383, 382, 13, 381,
	2, 380, 379, 15, 12, 9, 377, 376, 375, 374,
	373, 0, 14, 372, 371, 25, 5, 370, 369, 368,
	367, 4, 366, 363, 362, 356, 11, 348, 346, 345,
	10, 342, 339, 7, 335, 297, 16, 8, 240, 1,
}
var yyR1 = [...]int{

	0, 1, 2, 6, 4, 4, 7, 3, 3, 8,
	8, 8, 8, 9, 10, 10, 10, 5, 5, 14,
	14, 13, 13, 13, 13, 13, 13, 13, 13, 13,
	15, 15, 23, 26, 26, 26, 25, 27, 27, 28,
	29, 16, 31, 32, 32, 33, 33, 33, 35, 34,
	36, 37, 37, 38, 38, 19, 40, 30, 30, 41,
	41, 41, 22, 11, 43, 44, 44, 45, 45, 45,
	45, 47, 48, 24, 49, 50, 50, 51, 51, 51,
	51, 52, 53, 12, 54, 55, 55, 56, 56, 56,
	17, 58, 57, 59, 59, 60, 60, 60, 18, 61,
	62, 62, 63, 63, 63, 63, 63, 63, 64, 20,
	65, 66, 66, 67, 67, 67, 67, 67, 42, 21,
	68, 39, 39, 69, 46,
}
var yyR2 = [...]int{

	0, 5, 3, 2, 2, 5, 2, 2, 3, 2,
	1, 1, 2, 2, 1, 1, 1, 1, 2, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 2, 4, 0, 2, 1, 2, 1, 2, 4,
	2, 4, 2, 1, 2, 1, 2, 1, 3, 2,
	2, 1, 3, 3, 1, 4, 2, 1, 2, 2,
	1, 1, 3, 4, 2, 1, 2, 2, 1, 3,
	3, 2, 2, 4, 2, 1, 2, 2, 1, 3,
	3, 2, 2, 4, 2, 1, 2, 2, 1, 1,
	2, 3, 2, 1, 2, 2, 1, 1, 4, 2,
	1, 2, 2, 3, 1, 1, 3, 1, 3, 4,
	2, 1, 2, 1, 2, 1, 1, 3, 3, 4,
	2, 1, 2, 3, 3,
}
var yyChk = [...]int{

	-1000, -1, -2, 25, -3, -8, 11, -7, -9, 15,
	12, 39, 4, -4, -8, -6, 13, 9, 5, 5,
	5, 4, 7, -5, -10, -11, -12, -13, -43, -54,
	-16, -17, -18, -19, -20, -21, -22, -23, -24, 28,
	29, -31, -57, -61, -40, -65, -68, 21, -25, -49,
	30, 31, 27, 26, 32, 35, 37, 41, 9, 9,
	7, 10, 8, -10, 7, 7, 4, 4, 7, -58,
	7, 7, 7, 7, 7, 4, 7, 7, 4, 4,
	4, 4, 4, 4, 4, 4, -7, -44, -45, -7,
	-46, -47, -48, 34, 23, 24, -55, -56, -7, -42,
	-13, 20, -32, -33, -34, -7, -35, -36, 16, 14,
	-59, -60, -7, -46, -13, -62, -63, -7, 36, -42,
	-64, 22, -13, 19, -30, -41, -7, -42, -13, -66,
	-67, -34, -7, -42, -35, 33, -66, 9, -26, -7,
	-27, -28, -29, 38, -50, -51, -7, -46, -52, -53,
	23, 24, 9, 8, -45, 9, -14, -15, -13, -14,
	5, 7, 7, 8, -56, 9, 5, 8, -33, 9,
	-37, 9, 7, 5, 4, 8, -60, 9, 8, -63,
	9, 6, 5, 5, 8, -41, 9, 8, -67, 9,
	5, 8, 8, 9, -28, 7, 4, 8, -51, 9,
	-14, -14, 7, 7, 8, 8, -13, 8, 9, 9,
	-38, 17, -39, -69, 18, 9, 9, 9, 9, 9,
	-30, 8, 8, 8, 5, -69, 4, 8, 9, 9,
}
var yyDef = [...]int{

	0, -2, 0, 0, 0, 0, 0, 10, 11, 0,
	0, 0, 0, 0, 0, 0, 0, 7, 9, 12,
	6, 13, 2, 0, 17, 14, 15, 16, 0, 0,
	21, 22, 23, 24, 25, 26, 27, 28, 29, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 8, 4,
	0, 3, 1, 18, 0, 0, 64, 84, 0, 90,
	0, 0, 0, 0, 0, 0, 33, 0, 42, 92,
	99, 56, 110, 120, 36, 74, 0, 0, 65, 0,
	68, 19, 19, 0, 0, 0, 0, 85, 0, 88,
	89, 0, 0, 43, 45, 0, 47, 0, 0, 0,
	0, 93, 0, 96, 97, 0, 100, 0, 0, 104,
	105, 0, 107, 0, 0, 57, 0, 60, 61, 0,
	111, 113, 0, 115, 116, 0, 0, 62, 0, 0,
	35, 37, 0, 0, 0, 75, 0, 78, 19, 19,
	0, 0, 0, 63, 66, 67, 0, 20, 30, 0,
	0, 71, 72, 83, 86, 87, 0, 41, 44, 46,
	49, 51, 0, 0, 50, 91, 94, 95, 98, 101,
	102, 0, 0, 0, 55, 58, 59, 109, 112, 114,
	0, 119, 32, 34, 38, 0, 40, 73, 76, 77,
	0, 0, 81, 82, 5, 69, 31, 70, 124, 118,
	0, 0, 54, 121, 0, 48, 103, 106, 108, 117,
	0, 79, 80, 52, 0, 122, 0, 39, 53, 123,
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
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:353
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:358
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 71:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:365
		{
			yylval.stack.Push(&schema.RpcInput{})
		}
	case 72:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:370
		{
			yylval.stack.Push(&schema.RpcOutput{})
		}
	case 73:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:378
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 74:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:385
		{
			yylval.stack.Push(&schema.Rpc{Ident: yyDollar[2].token})
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:395
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 80:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:400
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 81:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:407
		{
			yylval.stack.Push(&schema.RpcInput{})
		}
	case 82:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:412
		{
			yylval.stack.Push(&schema.RpcOutput{})
		}
	case 83:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:420
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 84:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:427
		{
			yylval.stack.Push(&schema.Notification{Ident: yyDollar[2].token})
		}
	case 90:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:443
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 92:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:455
		{
			yylval.stack.Push(&schema.Grouping{Ident: yyDollar[2].token})
		}
	case 98:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:471
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 99:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:478
		{
			yylval.stack.Push(&schema.List{Ident: yyDollar[2].token})
		}
	case 108:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:494
		{
			if list, valid := yylval.stack.Peek().(*schema.List); valid {
				list.Keys = strings.Split(tokenString(yyDollar[2].token), " ")
			} else {
				yylex.Error("expected a list for key statement")
				goto ret1
			}
		}
	case 109:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:508
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 110:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:515
		{
			yylval.stack.Push(&schema.Leaf{Ident: yyDollar[2].token})
		}
	case 118:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:531
		{
			if hasDetails, valid := yylval.stack.Peek().(schema.HasDetails); valid {
				hasDetails.Details().SetConfig("true" == yyDollar[2].token)
			} else {
				yylex.Error("expected config statement on schema supporting details")
				goto ret1
			}
		}
	case 119:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:545
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 120:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:552
		{
			yylval.stack.Push(&schema.LeafList{Ident: yyDollar[2].token})
		}
	case 123:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:561
		{
			yylval.dataType.Enumeration = append(yylval.dataType.Enumeration, yyDollar[2].token)
		}
	}
	goto yystack /* stack new state and value */
}
