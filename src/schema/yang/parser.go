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

//line parser.y:606

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 137
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 431

var yyAct = [...]int{

	134, 226, 133, 138, 156, 151, 167, 124, 145, 119,
	140, 103, 141, 111, 94, 10, 139, 27, 135, 153,
	190, 7, 3, 7, 10, 223, 227, 27, 24, 96,
	10, 132, 108, 49, 130, 227, 6, 10, 55, 54,
	9, 153, 52, 53, 56, 109, 64, 57, 126, 59,
	243, 225, 66, 60, 58, 6, 10, 16, 10, 9,
	117, 242, 116, 202, 11, 241, 137, 10, 184, 117,
	183, 116, 107, 231, 230, 108, 136, 122, 131, 63,
	114, 62, 91, 11, 112, 229, 95, 104, 109, 228,
	113, 221, 120, 125, 146, 142, 142, 220, 10, 149,
	157, 169, 169, 121, 10, 170, 107, 219, 165, 161,
	162, 158, 95, 210, 175, 100, 101, 204, 201, 114,
	99, 104, 122, 112, 180, 198, 99, 131, 188, 113,
	192, 191, 179, 235, 197, 106, 10, 120, 117, 189,
	116, 128, 125, 144, 144, 105, 181, 176, 121, 166,
	200, 127, 163, 143, 143, 147, 205, 200, 142, 78,
	209, 61, 17, 169, 169, 142, 211, 212, 10, 106,
	234, 233, 217, 218, 157, 132, 108, 49, 130, 105,
	216, 215, 55, 54, 208, 158, 52, 53, 56, 109,
	128, 57, 126, 59, 203, 199, 164, 60, 58, 10,
	127, 117, 214, 116, 213, 193, 144, 108, 206, 232,
	10, 173, 117, 144, 116, 172, 143, 240, 108, 81,
	109, 10, 80, 143, 77, 76, 237, 75, 74, 108,
	49, 109, 73, 197, 71, 55, 54, 68, 67, 52,
	53, 56, 109, 196, 57, 22, 59, 10, 238, 236,
	60, 58, 195, 194, 185, 108, 49, 178, 177, 171,
	20, 55, 54, 19, 18, 52, 53, 56, 109, 174,
	57, 239, 59, 10, 5, 207, 60, 58, 186, 14,
	90, 108, 49, 89, 88, 87, 86, 55, 54, 47,
	85, 52, 53, 56, 109, 84, 57, 10, 59, 83,
	82, 79, 60, 58, 70, 108, 49, 69, 21, 12,
	187, 55, 54, 46, 10, 52, 53, 56, 109, 48,
	57, 129, 59, 49, 123, 44, 60, 58, 55, 54,
	118, 72, 52, 53, 56, 65, 99, 57, 43, 59,
	102, 29, 160, 60, 58, 159, 155, 154, 49, 51,
	98, 97, 10, 55, 54, 40, 41, 52, 53, 56,
	93, 49, 57, 92, 59, 28, 55, 54, 60, 58,
	52, 53, 56, 45, 99, 57, 224, 59, 49, 222,
	182, 60, 58, 55, 54, 40, 41, 52, 53, 56,
	115, 49, 57, 110, 59, 42, 55, 54, 60, 58,
	52, 53, 56, 152, 150, 57, 148, 59, 50, 39,
	38, 60, 58, 37, 36, 35, 34, 33, 32, 31,
	30, 168, 26, 25, 8, 15, 23, 13, 4, 2,
	1,
}
var yyPact = [...]int{

	-3, -1000, 25, 305, 44, 153, 259, -1000, -1000, 258,
	255, 304, 238, 357, 152, 72, 36, -1000, -1000, -1000,
	-1000, -1000, -1000, 327, -1000, -1000, -1000, -1000, 231, 230,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	303, 300, 227, 225, 221, 220, 218, 217, 150, 297,
	215, 212, 296, 295, 291, 286, 282, 281, 280, 279,
	276, -1000, -1000, 18, -1000, -1000, -1000, 92, 285, -1000,
	-1000, 46, -1000, 340, 156, 285, 198, 198, -1000, 146,
	3, 86, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 143, 188, 92, -1000, 140, -1000, 370, 370, 254,
	208, 204, 261, -1000, 138, -1000, -1000, -1000, 253, 252,
	124, -1000, -1000, 137, -1000, 61, 249, 274, 302, -1000,
	130, -1000, -1000, 12, -1000, 121, 199, -1000, -1000, -1000,
	248, -1000, 247, 235, -1000, 116, -1000, -1000, -1000, 187,
	-1000, -1000, 109, -1000, -1000, -1000, 55, -1000, 186, 108,
	-19, -1000, 201, 271, 176, 86, -1000, 104, -1000, 370,
	370, 197, 195, 173, -1000, -1000, -1000, 172, 370, -1000,
	165, 98, -1000, -1000, -1000, -1000, -1000, 88, 82, -1000,
	-1000, -1000, -1000, -1000, 8, 80, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 76, 65, 64, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 285, -1000, -1000, -1000,
	-1000, 163, 162, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, 125, 244, 17, 243, -1000, 267, -1000, -1000,
	-1000, -1000, 209, -1000, -1000, -1000, 56, -1000, 52, 41,
	-1000, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 430, 429, 428, 427, 426, 425, 18, 274, 424,
	28, 423, 422, 3, 6, 421, 420, 419, 418, 417,
	416, 415, 414, 413, 410, 409, 408, 406, 404, 5,
	403, 2, 395, 393, 13, 12, 8, 390, 380, 379,
	376, 373, 0, 76, 66, 365, 363, 360, 14, 29,
	351, 350, 349, 347, 346, 4, 345, 342, 341, 340,
	11, 338, 331, 330, 9, 325, 324, 7, 321, 319,
	313, 16, 10, 289, 1,
}
var yyR1 = [...]int{

	0, 1, 2, 6, 4, 4, 7, 3, 3, 8,
	8, 8, 8, 9, 10, 10, 10, 5, 5, 14,
	14, 13, 13, 13, 13, 13, 13, 13, 13, 13,
	13, 15, 15, 24, 27, 27, 27, 26, 28, 28,
	29, 30, 16, 32, 33, 33, 34, 34, 34, 36,
	35, 37, 38, 38, 39, 39, 39, 19, 41, 31,
	31, 42, 42, 42, 42, 23, 11, 45, 46, 46,
	47, 47, 48, 48, 48, 48, 50, 51, 25, 52,
	53, 53, 54, 54, 55, 55, 55, 55, 56, 57,
	12, 58, 59, 59, 60, 60, 60, 60, 17, 62,
	61, 63, 63, 64, 64, 64, 18, 65, 66, 66,
	67, 67, 67, 67, 67, 67, 67, 68, 22, 69,
	20, 70, 71, 71, 72, 72, 72, 72, 72, 44,
	43, 21, 73, 40, 40, 74, 49,
}
var yyR2 = [...]int{

	0, 5, 3, 2, 2, 5, 2, 2, 3, 2,
	1, 1, 2, 2, 1, 1, 1, 1, 2, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 2, 4, 0, 2, 1, 2, 1, 2,
	4, 2, 4, 2, 1, 2, 1, 2, 1, 3,
	2, 2, 1, 3, 3, 1, 3, 4, 2, 1,
	2, 2, 1, 1, 1, 3, 4, 2, 0, 1,
	1, 2, 2, 1, 3, 3, 2, 2, 4, 2,
	0, 1, 1, 2, 2, 1, 3, 3, 2, 2,
	4, 2, 1, 2, 2, 1, 1, 1, 2, 3,
	2, 1, 2, 2, 1, 1, 4, 2, 1, 2,
	2, 3, 1, 1, 1, 3, 1, 3, 2, 2,
	4, 2, 1, 2, 1, 2, 1, 1, 1, 3,
	3, 4, 2, 1, 2, 3, 3,
}
var yyChk = [...]int{

	-1000, -1, -2, 25, -3, -8, 11, -7, -9, 15,
	12, 39, 4, -4, -8, -6, 13, 9, 5, 5,
	5, 4, 7, -5, -10, -11, -12, -13, -45, -58,
	-16, -17, -18, -19, -20, -21, -22, -23, -24, -25,
	28, 29, -32, -61, -65, -41, -70, -73, -69, 21,
	-26, -52, 30, 31, 27, 26, 32, 35, 42, 37,
	41, 9, 9, 7, 10, 8, -10, 7, 7, 4,
	4, 7, -62, 7, 7, 7, 7, 7, 9, 4,
	7, 7, 4, 4, 4, 4, 4, 4, 4, 4,
	4, -7, -46, -47, -48, -7, -49, -50, -51, 34,
	23, 24, -59, -60, -7, -43, -44, -13, 20, 33,
	-33, -34, -35, -7, -36, -37, 16, 14, -63, -64,
	-7, -49, -13, -66, -67, -7, 36, -43, -44, -68,
	22, -13, 19, -31, -42, -7, -43, -44, -13, -71,
	-72, -35, -7, -43, -44, -36, -71, 9, -27, -7,
	-28, -29, -30, 38, -53, -54, -55, -7, -49, -56,
	-57, 23, 24, 9, 8, -48, 9, -14, -15, -13,
	-14, 5, 7, 7, 8, -60, 9, 5, 5, 8,
	-34, 9, -38, 9, 7, 5, 4, 8, -64, 9,
	8, -67, 9, 6, 5, 5, 8, -42, 9, 8,
	-72, 9, 8, 8, 9, -29, 7, 4, 8, -55,
	9, -14, -14, 7, 7, 8, 8, -13, 8, 9,
	9, 9, -39, 17, -40, 43, -74, 18, 9, 9,
	9, 9, -31, 8, 8, 8, 5, -74, 5, 4,
	8, 9, 9, 9,
}
var yyDef = [...]int{

	0, -2, 0, 0, 0, 0, 0, 10, 11, 0,
	0, 0, 0, 0, 0, 0, 0, 7, 9, 12,
	6, 13, 2, 0, 17, 14, 15, 16, 0, 0,
	21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 8, 4, 0, 3, 1, 18, 68, 0, 67,
	91, 0, 98, 0, 0, 0, 0, 0, 118, 0,
	34, 80, 43, 100, 107, 58, 121, 132, 119, 37,
	79, 0, 0, 69, 70, 0, 73, 19, 19, 0,
	0, 0, 0, 92, 0, 95, 96, 97, 0, 0,
	0, 44, 46, 0, 48, 0, 0, 0, 0, 101,
	0, 104, 105, 0, 108, 0, 0, 112, 113, 114,
	0, 116, 0, 0, 59, 0, 62, 63, 64, 0,
	122, 124, 0, 126, 127, 128, 0, 65, 0, 0,
	36, 38, 0, 0, 0, 81, 82, 0, 85, 19,
	19, 0, 0, 0, 66, 71, 72, 0, 20, 31,
	0, 0, 76, 77, 90, 93, 94, 0, 0, 42,
	45, 47, 50, 52, 0, 0, 51, 99, 102, 103,
	106, 109, 110, 0, 0, 0, 57, 60, 61, 120,
	123, 125, 131, 33, 35, 39, 0, 41, 78, 83,
	84, 0, 0, 88, 89, 5, 74, 32, 75, 136,
	130, 129, 0, 0, 55, 0, 133, 0, 49, 111,
	115, 117, 0, 86, 87, 53, 0, 134, 0, 0,
	40, 54, 56, 135,
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
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:332
		{
			yylval.stack.Push(&schema.Uses{Ident: yyDollar[2].token})
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 66:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:344
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 67:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:351
		{
			yylval.stack.Push(&schema.Rpc{Ident: yyDollar[2].token})
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:365
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:370
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 76:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:377
		{
			yylval.stack.Push(&schema.RpcInput{})
		}
	case 77:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:382
		{
			yylval.stack.Push(&schema.RpcOutput{})
		}
	case 78:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:390
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 79:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:397
		{
			yylval.stack.Push(&schema.Rpc{Ident: yyDollar[2].token})
		}
	case 86:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:411
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:416
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 88:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:423
		{
			yylval.stack.Push(&schema.RpcInput{})
		}
	case 89:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:428
		{
			yylval.stack.Push(&schema.RpcOutput{})
		}
	case 90:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:436
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 91:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:443
		{
			yylval.stack.Push(&schema.Notification{Ident: yyDollar[2].token})
		}
	case 98:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:460
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 100:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:472
		{
			yylval.stack.Push(&schema.Grouping{Ident: yyDollar[2].token})
		}
	case 106:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:488
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 107:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:495
		{
			yylval.stack.Push(&schema.List{Ident: yyDollar[2].token})
		}
	case 117:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:512
		{
			if list, valid := yylval.stack.Peek().(*schema.List); valid {
				list.Key = strings.Split(tokenString(yyDollar[2].token), " ")
			} else {
				yylex.Error("expected a list for key statement")
				goto ret1
			}
		}
	case 118:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:522
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 119:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:529
		{
			yylval.stack.Push(schema.NewAny(yyDollar[2].token))
		}
	case 120:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:537
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 121:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:544
		{
			yylval.stack.Push(&schema.Leaf{Ident: yyDollar[2].token})
		}
	case 129:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:560
		{
			if hasDetails, valid := yylval.stack.Peek().(schema.HasDetails); valid {
				hasDetails.Details().SetMandatory("true" == yyDollar[2].token)
			} else {
				yylex.Error("expected mandatory statement on schema supporting details")
				goto ret1
			}
		}
	case 130:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:569
		{
			if hasDetails, valid := yylval.stack.Peek().(schema.HasDetails); valid {
				hasDetails.Details().SetConfig("true" == yyDollar[2].token)
			} else {
				yylex.Error("expected config statement on schema supporting details")
				goto ret1
			}
		}
	case 131:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:583
		{
			if HasError(yylex, popAndAddMeta(&yylval)) {
				goto ret1
			}
		}
	case 132:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:590
		{
			yylval.stack.Push(&schema.LeafList{Ident: yyDollar[2].token})
		}
	case 135:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:599
		{
			yylval.dataType.AddEnumeration(yyDollar[2].token)
		}
	}
	goto yystack /* stack new state and value */
}
