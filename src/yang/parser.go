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

const yyLast = 333

var yyAct = [...]int{

	188, 150, 152, 131, 120, 122, 115, 106, 7, 101,
	7, 93, 119, 87, 26, 80, 121, 160, 133, 89,
	3, 9, 189, 26, 186, 189, 9, 55, 113, 109,
	45, 111, 141, 21, 140, 48, 47, 54, 210, 53,
	208, 49, 45, 9, 50, 108, 51, 48, 47, 39,
	40, 37, 38, 49, 57, 9, 50, 85, 51, 83,
	78, 209, 202, 90, 82, 104, 88, 94, 102, 133,
	112, 118, 201, 107, 116, 81, 177, 204, 129, 200,
	9, 95, 85, 126, 124, 82, 6, 9, 123, 90,
	8, 199, 88, 198, 197, 136, 81, 196, 94, 153,
	144, 125, 195, 104, 148, 194, 102, 190, 112, 184,
	157, 107, 95, 161, 179, 181, 173, 118, 9, 169,
	116, 168, 162, 158, 172, 113, 109, 45, 111, 149,
	145, 172, 48, 47, 180, 137, 134, 56, 49, 127,
	171, 50, 108, 51, 9, 135, 85, 52, 124, 9,
	45, 85, 123, 83, 192, 48, 47, 39, 40, 37,
	38, 49, 15, 167, 50, 125, 51, 9, 193, 191,
	9, 183, 85, 178, 124, 117, 45, 155, 123, 154,
	72, 48, 47, 203, 163, 156, 70, 49, 206, 9,
	50, 125, 51, 69, 68, 67, 62, 103, 45, 61,
	205, 147, 143, 48, 47, 9, 9, 60, 58, 49,
	9, 19, 50, 176, 51, 45, 98, 99, 117, 45,
	48, 47, 175, 174, 48, 47, 49, 91, 91, 50,
	49, 51, 9, 50, 170, 51, 9, 166, 165, 164,
	103, 45, 9, 159, 207, 45, 48, 47, 146, 138,
	48, 47, 49, 98, 99, 50, 49, 51, 91, 50,
	45, 51, 44, 18, 91, 48, 47, 17, 6, 9,
	14, 49, 8, 16, 50, 5, 51, 182, 142, 77,
	12, 76, 75, 74, 73, 71, 66, 65, 64, 63,
	10, 43, 110, 105, 41, 86, 59, 28, 100, 30,
	97, 96, 92, 29, 114, 42, 187, 185, 139, 84,
	79, 27, 132, 130, 128, 46, 36, 35, 34, 33,
	32, 31, 151, 25, 24, 23, 22, 13, 20, 11,
	4, 2, 1,
}
var yyPact = [...]int{

	-5, -1000, 75, 286, 257, 153, 268, -1000, 262, 258,
	204, 21, 138, 30, 17, -1000, -1000, -1000, -1000, -1000,
	129, -1000, -1000, -1000, -1000, -1000, -1000, 201, 200, 192,
	189, -1000, -1000, -1000, -1000, -1000, -1000, 285, 284, 283,
	282, 188, 187, 186, 179, 281, 173, 280, 279, 278,
	277, 275, -1000, -1000, 14, -1000, -1000, -1000, 43, -1000,
	224, 230, 220, -1000, -1000, -1000, -1000, 106, 198, 158,
	158, 130, 31, -1000, -1000, -1000, -1000, -1000, 127, 137,
	-1000, -1000, 126, 244, 25, 274, 194, -1000, 121, -1000,
	-1000, 243, 193, -1000, 120, -1000, 239, 239, 172, 170,
	177, -1000, 114, 238, -1000, 9, -1000, 113, 178, 234,
	-1000, 233, -1000, 232, 155, -1000, 110, 229, -1000, 132,
	-1000, -1000, 107, 218, 217, 208, 68, -1000, 165, 105,
	-20, -1000, 108, 273, 163, -1000, -1000, -1000, 100, -1000,
	-1000, 7, -1000, -1000, -1000, -1000, 98, -1000, -1000, -1000,
	161, 239, -1000, 160, -1000, -1000, -1000, -1000, -1000, 96,
	-1000, -1000, -1000, 93, 88, 85, 84, -1000, -1000, -1000,
	82, -1000, -1000, -1000, 70, 63, 53, -1000, -1000, -1000,
	-1000, 239, -1000, -1000, -1000, 69, 195, 4, -1000, 240,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 32, -1000, 52, -1000, 29, -1000, -1000,
	-1000,
}
var yyPgo = [...]int{

	0, 332, 331, 330, 329, 328, 327, 5, 275, 33,
	326, 325, 324, 323, 2, 1, 322, 321, 320, 319,
	318, 317, 316, 315, 314, 313, 3, 312, 311, 310,
	15, 16, 309, 308, 307, 306, 305, 304, 6, 303,
	302, 11, 19, 301, 300, 299, 298, 9, 297, 296,
	295, 13, 294, 293, 7, 292, 291, 12, 4, 262,
	0,
}
var yyR1 = [...]int{

	0, 1, 2, 6, 4, 4, 7, 3, 3, 8,
	8, 8, 9, 9, 9, 9, 9, 5, 5, 15,
	15, 14, 14, 14, 14, 14, 14, 16, 16, 22,
	24, 24, 24, 23, 25, 25, 26, 27, 10, 28,
	29, 29, 30, 30, 30, 31, 32, 33, 33, 34,
	34, 18, 36, 37, 37, 38, 38, 38, 21, 12,
	39, 40, 40, 41, 41, 41, 41, 43, 44, 13,
	45, 46, 46, 47, 47, 47, 11, 49, 48, 50,
	50, 51, 51, 51, 17, 52, 53, 53, 54, 54,
	54, 54, 54, 54, 55, 19, 56, 57, 57, 58,
	58, 58, 58, 58, 20, 59, 35, 35, 60, 42,
}
var yyR2 = [...]int{

	0, 5, 3, 2, 2, 5, 2, 2, 3, 2,
	1, 2, 1, 1, 1, 1, 1, 1, 2, 0,
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
	-5, -9, -10, -11, -12, -13, -14, -28, -48, -39,
	-45, -17, -18, -19, -20, -21, -22, 30, 31, 28,
	29, -52, -36, -56, -59, 21, -23, 27, 26, 32,
	35, 37, 9, 9, 7, 10, 8, -9, 7, -49,
	7, 7, 7, 4, 4, 4, 4, 7, 7, 7,
	7, 4, 7, 4, 4, 4, 4, 4, -7, -29,
	-30, -31, -7, 16, -32, 14, -50, -51, -7, -42,
	-14, 34, -40, -41, -7, -42, -43, -44, 23, 24,
	-46, -47, -7, 20, -14, -53, -54, -7, 36, 20,
	-55, 22, -14, 19, -37, -38, -7, 20, -14, -57,
	-58, -31, -7, 20, 16, 33, -57, 9, -24, -7,
	-25, -26, -27, 38, 9, 8, -30, 9, 5, -33,
	9, 7, 4, 8, -51, 9, 5, 8, -41, 9,
	-15, -16, -14, -15, 7, 7, 8, -47, 9, 5,
	8, -54, 9, 6, 5, 5, 5, 8, -38, 9,
	5, 8, -58, 9, 5, 5, 5, 8, 8, 9,
	-26, 7, 4, 8, 9, -34, 17, -35, -60, 18,
	9, 8, -14, 8, 9, 9, 9, 9, 9, 9,
	9, 9, 9, -15, 8, 5, -60, 4, 8, 9,
	9,
}
var yyDef = [...]int{

	0, -2, 0, 0, 0, 0, 0, 10, 0, 0,
	0, 0, 0, 0, 0, 7, 9, 11, 6, 2,
	0, 17, 12, 13, 14, 15, 16, 0, 0, 0,
	0, 21, 22, 23, 24, 25, 26, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 8, 4, 0, 3, 1, 18, 0, 76,
	0, 0, 0, 39, 78, 60, 70, 0, 0, 0,
	0, 0, 30, 85, 52, 96, 105, 33, 0, 0,
	40, 42, 0, 0, 0, 0, 0, 79, 0, 82,
	83, 0, 0, 61, 0, 64, 19, 19, 0, 0,
	0, 71, 0, 0, 75, 0, 86, 0, 0, 0,
	91, 0, 93, 0, 0, 53, 0, 0, 57, 0,
	97, 99, 0, 0, 0, 0, 0, 58, 0, 0,
	32, 34, 0, 0, 0, 38, 41, 43, 0, 45,
	47, 0, 46, 77, 80, 81, 0, 59, 62, 63,
	0, 20, 27, 0, 67, 68, 69, 72, 73, 0,
	84, 87, 88, 0, 0, 0, 0, 51, 54, 55,
	0, 95, 98, 100, 0, 0, 0, 104, 29, 31,
	35, 19, 37, 5, 44, 0, 0, 50, 106, 0,
	109, 65, 28, 66, 74, 89, 90, 92, 94, 56,
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
