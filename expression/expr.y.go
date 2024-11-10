// Code generated by goyacc -o expr.y.go expr.y. DO NOT EDIT.

//line expr.y:2
package expression

import __yyfmt__ "fmt"

//line expr.y:2
import ()

//line expr.y:7
type yySymType struct {
	yys        int
	exprVal    Expression
	sVal       string
	iVal       int64
	funcArgs   []Expression
	structArgs map[string]Expression
}

const NUM = 57346
const STRING = 57347
const IDENTIFIER = 57348
const VARIABLE = 57349
const FUNC = 57350
const RETURN = 57351
const IF = 57352
const ELSE = 57353
const SWITCH = 57354
const CASE = 57355
const DEFAULT = 57356
const NEWLINE = 57357
const COMPARE = 57358

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"NUM",
	"STRING",
	"IDENTIFIER",
	"VARIABLE",
	"FUNC",
	"RETURN",
	"IF",
	"ELSE",
	"SWITCH",
	"CASE",
	"DEFAULT",
	"NEWLINE",
	"COMPARE",
	"'='",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'%'",
	"'('",
	"')'",
	"'{'",
	"'}'",
	"':'",
	"','",
	"'$'",
	"';'",
}

var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line expr.y:143

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 105

var yyAct = [...]int{
	2, 60, 71, 14, 8, 6, 17, 58, 46, 67,
	47, 59, 68, 67, 28, 84, 68, 15, 16, 86,
	14, 8, 13, 17, 19, 85, 93, 81, 18, 57,
	92, 38, 39, 40, 15, 16, 48, 79, 80, 13,
	67, 19, 5, 68, 89, 18, 88, 74, 56, 53,
	77, 55, 73, 51, 49, 32, 41, 63, 61, 70,
	62, 25, 26, 27, 4, 50, 36, 37, 69, 45,
	31, 72, 23, 24, 22, 21, 52, 44, 34, 30,
	83, 29, 75, 76, 78, 87, 82, 35, 66, 65,
	90, 91, 64, 54, 43, 42, 33, 12, 11, 10,
	9, 7, 3, 20, 1,
}

var yyPact = [...]int{
	16, 60, -1000, 58, 54, 41, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 16, -1000, 77, 75, -1000, 47, 72,
	-1000, -1000, 16, 16, 16, 16, 16, 16, 32, -1000,
	-1000, 71, 46, -18, 9, 54, 41, 41, -1000, -1000,
	-1000, -1000, 30, 42, -1000, 29, -1000, 70, 16, -1000,
	16, 23, 2, -1000, -17, -1000, -1000, 16, -1000, 16,
	31, -1000, -1000, -1000, -1000, -1000, -1000, 16, 36, -28,
	16, -1000, 28, 22, -1000, 24, -1000, -1000, 1, -1,
	-2, -1000, -8, -1000, 16, 21, 19, -1000, -1000, -1000,
	4, 0, -1000, -1000,
}

var yyPgo = [...]int{
	0, 104, 0, 103, 102, 64, 42, 5, 101, 100,
	99, 98, 97, 96, 95, 94, 93, 1, 92, 89,
	88, 86, 84, 83, 82,
}

var yyR1 = [...]int{
	0, 1, 1, 2, 4, 4, 5, 5, 5, 6,
	6, 6, 6, 7, 7, 7, 7, 7, 7, 7,
	12, 13, 13, 8, 8, 8, 9, 10, 14, 16,
	16, 16, 11, 18, 18, 17, 17, 21, 21, 22,
	23, 24, 24, 20, 20, 19, 15, 3,
}

var yyR2 = [...]int{
	0, 1, 2, 1, 1, 3, 1, 3, 3, 1,
	3, 3, 3, 1, 1, 1, 1, 1, 1, 3,
	3, 3, 5, 1, 2, 2, 1, 4, 4, 0,
	1, 3, 7, 1, 1, 0, 2, 1, 2, 5,
	6, 0, 2, 7, 8, 3, 1, 1,
}

var yyChk = [...]int{
	-1000, -1, -2, -4, -5, -6, -7, -8, 5, -9,
	-10, -11, -12, 23, 4, 18, 19, 7, 29, 25,
	-3, 15, 16, 18, 19, 20, 21, 22, -2, 4,
	4, 23, 8, -13, 6, -5, -6, -6, -7, -7,
	-7, 24, -14, -15, 6, 23, 26, 28, 27, 24,
	23, 24, 6, -2, -16, -2, 25, 27, 24, 28,
	-17, -2, -2, 26, -18, -19, -20, 9, 12, -2,
	23, 30, -2, 24, 25, -24, -23, 26, -22, 13,
	14, 26, -21, -2, 16, 27, 27, -2, 25, 25,
	-17, -17, 26, 26,
}

var yyDef = [...]int{
	0, -2, 1, 3, 4, 6, 9, 13, 14, 15,
	16, 17, 18, 0, 23, 0, 0, 26, 0, 0,
	2, 47, 0, 0, 0, 0, 0, 0, 0, 24,
	25, 0, 0, 0, 0, 5, 7, 8, 10, 11,
	12, 19, 0, 0, 46, 0, 20, 0, 0, 27,
	29, 0, 0, 21, 0, 30, 35, 0, 28, 0,
	0, 22, 31, 32, 36, 33, 34, 0, 0, 0,
	0, 45, 0, 0, 41, 0, 42, 43, 0, 0,
	0, 44, 0, 37, 0, 0, 0, 38, 35, 35,
	0, 0, 39, 40,
}

var yyTok1 = [...]int{
	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 29, 22, 3, 3,
	23, 24, 20, 18, 28, 19, 3, 21, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 27, 30,
	3, 17, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 25, 3, 26,
}

var yyTok2 = [...]int{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16,
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
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
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
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
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
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
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
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
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
			yyrcvr.char = -1
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

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
//line expr.y:24
		{
			currentExecutor.Push(yyDollar[1].exprVal)
		}
	case 3:
		yyDollar = yyS[yypt-1 : yypt+1]
//line expr.y:27
		{
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
//line expr.y:30
		{
		}
	case 5:
		yyDollar = yyS[yypt-3 : yypt+1]
//line expr.y:31
		{
			yyVAL.exprVal = NewMathStatement([]Expression{yyDollar[1].exprVal, yyDollar[3].exprVal}, yyDollar[2].sVal)
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
//line expr.y:35
		{
		}
	case 7:
		yyDollar = yyS[yypt-3 : yypt+1]
//line expr.y:36
		{
			yyVAL.exprVal = NewMathStatement([]Expression{yyDollar[1].exprVal, yyDollar[3].exprVal}, "+")
		}
	case 8:
		yyDollar = yyS[yypt-3 : yypt+1]
//line expr.y:39
		{
			yyVAL.exprVal = NewMathStatement([]Expression{yyDollar[1].exprVal, yyDollar[3].exprVal}, "-")
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
//line expr.y:43
		{
		}
	case 10:
		yyDollar = yyS[yypt-3 : yypt+1]
//line expr.y:44
		{
			yyVAL.exprVal = NewMathStatement([]Expression{yyDollar[1].exprVal, yyDollar[3].exprVal}, "*")
		}
	case 11:
		yyDollar = yyS[yypt-3 : yypt+1]
//line expr.y:47
		{
			yyVAL.exprVal = NewMathStatement([]Expression{yyDollar[1].exprVal, yyDollar[3].exprVal}, "/")
		}
	case 12:
		yyDollar = yyS[yypt-3 : yypt+1]
//line expr.y:50
		{
			yyVAL.exprVal = NewMathStatement([]Expression{yyDollar[1].exprVal, yyDollar[3].exprVal}, "%")
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
//line expr.y:54
		{
			yyVAL.exprVal = yyDollar[1].exprVal
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
//line expr.y:55
		{
			yyVAL.exprVal = NewString(yyDollar[1].sVal)
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
//line expr.y:56
		{
			yyVAL.exprVal = yyDollar[1].exprVal
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
//line expr.y:57
		{
			yyVAL.exprVal = yyDollar[1].exprVal
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
//line expr.y:58
		{
			yyVAL.exprVal = yyDollar[1].exprVal
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
//line expr.y:59
		{
			yyVAL.exprVal = yyDollar[1].exprVal
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
//line expr.y:60
		{
			yyVAL.exprVal = yyDollar[2].exprVal
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
//line expr.y:62
		{
			yyVAL.exprVal = NewStructStatement(yyDollar[2].structArgs)
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
//line expr.y:67
		{
			yyVAL.structArgs = make(map[string]Expression)
			yyVAL.structArgs[yyDollar[1].sVal] = yyDollar[3].exprVal
		}
	case 22:
		yyDollar = yyS[yypt-5 : yypt+1]
//line expr.y:72
		{
			yyVAL.structArgs[yyDollar[3].sVal] = yyDollar[5].exprVal
		}
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
//line expr.y:76
		{
			yyVAL.exprVal = NewInt(yyDollar[1].iVal)
		}
	case 24:
		yyDollar = yyS[yypt-2 : yypt+1]
//line expr.y:77
		{
			yyVAL.exprVal = NewInt(yyDollar[1].iVal)
		}
	case 25:
		yyDollar = yyS[yypt-2 : yypt+1]
//line expr.y:78
		{
			yyVAL.exprVal = NewMathStatement([]Expression{NewInt(yyDollar[1].iVal)}, "neg")
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
//line expr.y:82
		{
			yyVAL.exprVal = NewVariableStatement(yyDollar[1].sVal)
		}
	case 27:
		yyDollar = yyS[yypt-4 : yypt+1]
//line expr.y:87
		{
			yyVAL.exprVal = yyDollar[3].exprVal
		}
	case 28:
		yyDollar = yyS[yypt-4 : yypt+1]
//line expr.y:91
		{
			yyVAL.exprVal = NewCallStatement(yyDollar[3].funcArgs, yyDollar[1].sVal)
		}
	case 29:
		yyDollar = yyS[yypt-0 : yypt+1]
//line expr.y:95
		{
			yyVAL.funcArgs = nil
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
//line expr.y:96
		{
			yyVAL.funcArgs = append(yyVAL.funcArgs, yyDollar[1].exprVal)
		}
	case 31:
		yyDollar = yyS[yypt-3 : yypt+1]
//line expr.y:99
		{
			yyVAL.funcArgs = append(yyDollar[1].funcArgs, yyDollar[3].exprVal)
		}
	case 32:
		yyDollar = yyS[yypt-7 : yypt+1]
//line expr.y:104
		{
			yyVAL.exprVal = NewFunctionStatement(yyDollar[6].funcArgs)
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
//line expr.y:108
		{
			yyVAL.exprVal = yyDollar[1].exprVal
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
//line expr.y:109
		{
			yyVAL.exprVal = yyDollar[1].exprVal
		}
	case 35:
		yyDollar = yyS[yypt-0 : yypt+1]
//line expr.y:111
		{
			yyVAL.funcArgs = []Expression{}
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
//line expr.y:112
		{
			yyVAL.funcArgs = append(yyDollar[1].funcArgs, yyDollar[2].exprVal)
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
//line expr.y:114
		{
			yyVAL.exprVal = NewCaseExpression("", yyDollar[1].exprVal)
		}
	case 38:
		yyDollar = yyS[yypt-2 : yypt+1]
//line expr.y:115
		{
			yyVAL.exprVal = NewCaseExpression(yyDollar[1].sVal, yyDollar[2].exprVal)
		}
	case 39:
		yyDollar = yyS[yypt-5 : yypt+1]
//line expr.y:118
		{
			yyVAL.exprVal = NewDefaultStatement(yyDollar[4].funcArgs)
		}
	case 40:
		yyDollar = yyS[yypt-6 : yypt+1]
//line expr.y:121
		{
			yyVAL.exprVal = NewCaseStatement(yyDollar[2].exprVal, yyDollar[5].funcArgs)
		}
	case 41:
		yyDollar = yyS[yypt-0 : yypt+1]
//line expr.y:123
		{
			yyVAL.funcArgs = []Expression{}
		}
	case 42:
		yyDollar = yyS[yypt-2 : yypt+1]
//line expr.y:125
		{
			yyVAL.funcArgs = append(yyDollar[1].funcArgs, yyDollar[2].exprVal)
		}
	case 43:
		yyDollar = yyS[yypt-7 : yypt+1]
//line expr.y:128
		{
			yyVAL.exprVal = NewSwitchStatement(yyDollar[3].exprVal, yyDollar[6].funcArgs)
		}
	case 44:
		yyDollar = yyS[yypt-8 : yypt+1]
//line expr.y:130
		{
			tempArgs := append(yyDollar[6].funcArgs, yyDollar[7].exprVal)
			yyVAL.exprVal = NewSwitchStatement(yyDollar[3].exprVal, tempArgs)
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
//line expr.y:136
		{
			yyVAL.exprVal = NewReturnStatement(yyDollar[2].exprVal)
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
//line expr.y:140
		{
		}
	case 47:
		yyDollar = yyS[yypt-1 : yypt+1]
//line expr.y:141
		{
		}
	}
	goto yystack /* stack new state and value */
}
