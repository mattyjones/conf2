//line parser.y:2
package yang

import __yyfmt__ "fmt"

//line parser.y:2
import (
	"fmt"
)

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
	fmt.Println(fmt.Sprintf("%s at line %d, col %d", e, line, col))
}

func popAndAddDef(yylval *yySymType) bool {
	i := yylval.stack.Pop()
	def, ok := i.(Def)
	if ok {
		parent := yylval.stack.Peek()
		parentList, ok := parent.(DefList)
		if ok {
			err := parentList.AddDef(def)
			if err == nil {
				return true
			}
			if yyDebug > 1 {
				__yyfmt__.Printf(err.Error())
			}
		} else if yyDebug > 1 {
			__yyfmt__.Printf("Internal Error: %s doesn't implement DefList.", parent.GetIdent())
		}
	} else {
		__yyfmt__.Printf("Internal Error: %s doesn't implement Def.", i.GetIdent())
	}

	return false
}

//line parser.y:49
type yySymType struct {
	yys   int
	ident string
	token string
	stack *yangDefStack
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

//line parser.y:408

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

const yyNprod = 101
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 302

var yyAct = [...]int{

	180, 157, 70, 156, 104, 96, 139, 90, 81, 67,
	61, 158, 9, 3, 159, 181, 69, 7, 123, 7,
	178, 181, 9, 101, 102, 9, 43, 199, 73, 86,
	85, 84, 87, 21, 71, 189, 37, 36, 198, 9,
	25, 65, 78, 162, 197, 79, 83, 161, 160, 25,
	37, 36, 88, 94, 45, 108, 78, 59, 62, 79,
	163, 63, 196, 68, 82, 91, 97, 105, 98, 195,
	9, 111, 62, 194, 184, 63, 116, 86, 85, 84,
	87, 68, 176, 88, 37, 36, 175, 6, 9, 124,
	78, 8, 94, 79, 83, 82, 9, 132, 65, 171,
	64, 137, 141, 141, 91, 170, 108, 142, 146, 145,
	97, 110, 98, 9, 169, 9, 136, 65, 105, 64,
	9, 107, 106, 153, 164, 154, 131, 37, 36, 168,
	9, 101, 102, 78, 167, 166, 79, 115, 93, 92,
	165, 9, 71, 173, 37, 36, 155, 6, 9, 14,
	78, 8, 42, 79, 41, 37, 36, 151, 183, 9,
	72, 78, 147, 71, 79, 138, 183, 107, 106, 9,
	133, 9, 24, 37, 36, 125, 117, 93, 92, 78,
	192, 24, 79, 37, 36, 37, 36, 112, 109, 78,
	40, 78, 79, 71, 79, 182, 9, 15, 65, 9,
	162, 65, 190, 162, 161, 160, 44, 161, 160, 37,
	36, 38, 39, 34, 35, 174, 172, 163, 150, 144,
	163, 143, 120, 119, 37, 36, 38, 39, 34, 35,
	52, 51, 50, 49, 48, 46, 19, 126, 191, 188,
	187, 186, 149, 135, 130, 129, 128, 127, 118, 113,
	18, 17, 16, 5, 193, 185, 148, 134, 12, 122,
	121, 114, 58, 57, 56, 55, 54, 53, 10, 77,
	76, 80, 30, 66, 47, 29, 103, 33, 100, 99,
	95, 32, 89, 31, 179, 177, 152, 60, 28, 75,
	74, 140, 27, 26, 23, 22, 13, 20, 11, 4,
	2, 1,
}
var yyPact = [...]int{

	-12, -1000, 76, 264, 136, 188, 247, -1000, 246, 245,
	229, 183, 181, 145, 16, -1000, -1000, -1000, -1000, -1000,
	198, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 228, 227,
	226, 225, 224, 223, 263, 262, 261, 260, 259, 258,
	-1000, -1000, 13, -1000, -1000, -1000, 84, -1000, 159, 58,
	157, 0, 147, -1000, -1000, -1000, -1000, -1000, -1000, 179,
	103, -1000, -1000, 178, 244, 257, 129, -1000, 167, -1000,
	-1000, 243, -1000, -1000, -1000, -1000, 216, 215, 256, 255,
	10, -1000, 166, 231, 242, 241, 240, 239, -1000, 118,
	-1000, 161, 253, 238, -1000, 108, -1000, 156, -1000, 24,
	24, 214, 212, 101, -1000, 153, 252, 237, -1000, 210,
	-1000, -1000, -1000, 148, 116, -1000, -1000, -1000, 137, 184,
	184, -1000, -1000, -1000, -1000, -1000, 131, 126, 125, 120,
	105, -1000, -1000, -1000, 96, 90, -1000, -1000, -1000, 208,
	24, -1000, 207, -1000, -1000, -1000, -1000, -1000, 77, 73,
	-1000, -1000, -1000, 3, -1000, -1000, 187, -1000, -1000, 65,
	251, 236, 235, 234, 27, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 194, 233, -3,
	-1000, 250, -1000, -1000, -1000, 64, 60, 53, 35, -1000,
	-1000, 29, -1000, 18, -1000, -1000, -1000, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 301, 300, 299, 298, 297, 296, 14, 253, 33,
	295, 294, 160, 28, 293, 292, 6, 291, 2, 290,
	289, 288, 287, 10, 11, 286, 285, 284, 283, 282,
	7, 281, 280, 5, 16, 279, 278, 277, 276, 4,
	275, 274, 273, 9, 272, 271, 8, 270, 3, 1,
	269, 0,
}
var yyR1 = [...]int{

	0, 1, 2, 6, 4, 4, 7, 3, 3, 8,
	8, 8, 9, 9, 9, 9, 9, 9, 5, 5,
	16, 16, 18, 18, 18, 18, 17, 17, 10, 21,
	22, 22, 23, 23, 23, 24, 25, 25, 26, 26,
	13, 28, 29, 29, 30, 30, 30, 30, 14, 31,
	32, 32, 33, 33, 33, 33, 35, 36, 15, 37,
	38, 38, 39, 39, 39, 39, 11, 41, 40, 42,
	42, 43, 43, 43, 12, 44, 45, 45, 46, 46,
	46, 46, 46, 46, 46, 19, 47, 48, 48, 49,
	49, 49, 49, 49, 49, 20, 50, 27, 27, 51,
	34,
}
var yyR2 = [...]int{

	0, 5, 3, 2, 2, 5, 2, 2, 3, 2,
	1, 2, 1, 1, 1, 1, 1, 1, 1, 2,
	0, 1, 1, 1, 1, 1, 1, 2, 4, 2,
	1, 2, 1, 2, 3, 3, 3, 1, 3, 1,
	4, 2, 1, 2, 2, 3, 3, 1, 4, 2,
	1, 2, 2, 1, 3, 3, 2, 2, 4, 2,
	1, 2, 2, 3, 3, 1, 2, 3, 2, 1,
	2, 2, 1, 1, 4, 2, 1, 2, 2, 3,
	3, 3, 3, 3, 1, 4, 2, 1, 2, 1,
	2, 3, 3, 3, 3, 4, 2, 1, 2, 3,
	3,
}
var yyChk = [...]int{

	-1000, -1, -2, 25, -3, -8, 11, -7, 15, 12,
	4, -4, -8, -6, 13, 9, 5, 5, 5, 7,
	-5, -9, -10, -11, -12, -13, -14, -15, -21, -40,
	-44, -28, -31, -37, 30, 31, 27, 26, 28, 29,
	9, 9, 7, 10, 8, -9, 7, -41, 7, 7,
	7, 7, 7, 4, 4, 4, 4, 4, 4, -7,
	-22, -23, -24, -7, 16, 14, -42, -43, -7, -34,
	-18, 34, -12, -13, -19, -20, -47, -50, 32, 35,
	-45, -46, -7, 36, 21, 20, 19, 22, -18, -29,
	-30, -7, 21, 20, -18, -32, -33, -7, -34, -35,
	-36, 23, 24, -38, -39, -7, 21, 20, -18, 9,
	8, -23, 9, 5, 4, 8, -43, 9, 5, 7,
	7, 4, 4, 8, -46, 9, 6, 5, 5, 5,
	5, 8, -30, 9, 4, 5, 8, -33, 9, -16,
	-17, -18, -16, 7, 7, 8, -39, 9, 4, 5,
	8, 9, -25, 7, 9, 9, -48, -49, -24, -7,
	21, 20, 16, 33, -48, 9, 9, 9, 9, 9,
	9, 9, 8, -18, 8, 9, 9, -26, 17, -27,
	-51, 18, 8, -49, 9, 4, 5, 5, 5, 8,
	8, 5, -51, 4, 9, 9, 9, 9, 9, 9,
}
var yyDef = [...]int{

	0, -2, 0, 0, 0, 0, 0, 10, 0, 0,
	0, 0, 0, 0, 0, 7, 9, 11, 6, 2,
	0, 18, 12, 13, 14, 15, 16, 17, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	8, 4, 0, 3, 1, 19, 0, 66, 0, 0,
	0, 0, 0, 29, 68, 75, 41, 49, 59, 0,
	0, 30, 32, 0, 0, 0, 0, 69, 0, 72,
	73, 0, 22, 23, 24, 25, 0, 0, 0, 0,
	0, 76, 0, 0, 0, 0, 0, 0, 84, 0,
	42, 0, 0, 0, 47, 0, 50, 0, 53, 20,
	20, 0, 0, 0, 60, 0, 0, 0, 65, 0,
	28, 31, 33, 0, 0, 67, 70, 71, 0, 0,
	0, 86, 96, 74, 77, 78, 0, 0, 0, 0,
	0, 40, 43, 44, 0, 0, 48, 51, 52, 0,
	21, 26, 0, 56, 57, 58, 61, 62, 0, 0,
	5, 34, 35, 0, 37, 100, 0, 87, 89, 0,
	0, 0, 0, 0, 0, 79, 80, 81, 82, 83,
	45, 46, 54, 27, 55, 63, 64, 0, 0, 39,
	97, 0, 85, 88, 90, 0, 0, 0, 0, 95,
	36, 0, 98, 0, 91, 92, 93, 94, 38, 99,
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
		//line parser.y:102
		{
			m := &Module{Ident: yyDollar[2].token}
			yylval.stack.Push(m)
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:108
		{
			d := yylval.stack.Peek()
			r := &Revision{Ident: yyDollar[2].token}
			d.(*Module).Revision = r
			yylval.stack.Push(r)
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:116
		{
			yylval.stack.Pop()
		}
	case 5:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:119
		{
			yylval.stack.Pop()
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:123
		{
			yylval.stack.Peek().(Describable).SetDescription(yyDollar[2].token)
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:133
		{
			d := yylval.stack.Peek()
			d.(*Module).Namespace = yyDollar[2].token
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:138
		{
			m := yylval.stack.Peek().(*Module)
			m.Prefix = yyDollar[2].token
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:152
		{
			if !popAndAddDef(&yylval) {
				goto ret1
			}
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:157
		{
			if !popAndAddDef(&yylval) {
				goto ret1
			}
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:174
		{
			if !popAndAddDef(&yylval) {
				goto ret1
			}
		}
	case 27:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:179
		{
			if !popAndAddDef(&yylval) {
				goto ret1
			}
		}
	case 29:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:192
		{
			yylval.stack.Push(&Typedef{Ident: yyDollar[2].token})
		}
	case 41:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:226
		{
			yylval.stack.Push(&Container{Ident: yyDollar[2].token})
		}
	case 47:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:238
		{
			if !popAndAddDef(&yylval) {
				goto ret1
			}
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:251
		{
			yylval.stack.Push(&Rpc{Ident: yyDollar[2].token})
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:262
		{
			input := yylval.stack.Pop().(*RpcInput)
			rpc := yylval.stack.Peek().(*Rpc)
			rpc.Input = input
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:267
		{
			output := yylval.stack.Pop().(*RpcOutput)
			rpc := yylval.stack.Peek().(*Rpc)
			rpc.Output = output
		}
	case 56:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:274
		{
			yylval.stack.Push(&RpcInput{})
		}
	case 57:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:279
		{
			yylval.stack.Push(&RpcOutput{})
		}
	case 59:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:290
		{
			yylval.stack.Push(&Notification{Ident: yyDollar[2].token})
		}
	case 65:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:303
		{
			if !popAndAddDef(&yylval) {
				goto ret1
			}
		}
	case 68:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:319
		{
			yylval.stack.Push(&Grouping{Ident: yyDollar[2].token})
		}
	case 73:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:330
		{
			if !popAndAddDef(&yylval) {
				goto ret1
			}
		}
	case 75:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:342
		{
			yylval.stack.Push(&List{Ident: yyDollar[2].token})
		}
	case 84:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:357
		{
			if !popAndAddDef(&yylval) {
				goto ret1
			}
		}
	case 86:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:370
		{
			yylval.stack.Push(&Leaf{Ident: yyDollar[2].token})
		}
	case 96:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:395
		{
			yylval.stack.Push(&LeafList{Ident: yyDollar[2].token})
		}
	}
	goto yystack /* stack new state and value */
}
