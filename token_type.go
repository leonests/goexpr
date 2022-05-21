package rulengine

import "strconv"

type TokenType int //enum

const (
	ILLEGAL TokenType = iota

	// literal operators
	CHAR     // 'a', '爱', ...
	STRING   // "abc"
	NUMBER   // 123, 123.456 treated as float64
	BOOL     // true, false
	VARIABLE // a1, b_2, c
	SELECTOR // a.b.c,
	ACCESSOR // .a.b

	// prefix operators
	NOT // !
	NEG // -

	// normal operators
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND // &
	OR  // |
	XOR // ^
	SHL // <<
	SHR // >>

	// ternary operators
	TERNARY_IF
	TERNARY_ELSE

	// logical operators
	LAND // &&
	LOR  // ||

	// comparer operators
	EQ  // ==
	NEQ // !=
	LT  // <
	GT  // >
	LEQ // <=
	GEQ // >=

	// clause operators
	LPAREN // (
	RPAREN // )

	LBRACKET // [
	RBRACKET // ]

	LITERAL // represent all literal operators
	CLAUSE  // represent all clause operator
	EOF
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	VARIABLE: "VARIABLE",
	SELECTOR: "SELECTOR",
	ACCESSOR: "ACCESSOR",

	CHAR:   "CHAR",
	STRING: "STRING",
	NUMBER: "NUMBER",
	BOOL:   "BOOL",

	NOT: "!",
	NEG: "-",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	AND: "&",
	OR:  "|",
	XOR: "^",
	SHL: "<<",
	SHR: ">>",

	LAND: "&&",
	LOR:  "||",

	EQ:  "==",
	NEQ: "!=",
	LT:  "<",
	GT:  ">",
	LEQ: "<=",
	GEQ: ">=",

	TERNARY_IF:   "?",
	TERNARY_ELSE: ":",

	LPAREN: "LPAREN",
	RPAREN: "RPAREN",

	LBRACKET: "LBRACKET",
	RBRACKET: "RBRACKET",

	LITERAL: "LITERAL",
	CLAUSE:  "CLAUSE",
	EOF:     "EOF",
}

func (t TokenType) String() string {
	s := ""
	if 0 <= t && t <= TokenType(len(tokens)) {
		s = tokens[t]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(t)) + ")"
	}
	return s
}

type opPriority int

const (
	priorityUNKNOWN opPriority = iota
	priorityTENARY
	priorityLOR
	priorityLAND
	priorityBIT
	priorityCOMPARER
	priorityBITSHIFT
	priorityADD
	priorityMUL
	priorityPREFIX
	priorityCLAUSE
	priorityLITERAL
)

func (op TokenType) Priority() opPriority {
	switch op {
	case TERNARY_IF, TERNARY_ELSE:
		return priorityTENARY
	case LOR:
		return priorityLOR
	case LAND:
		return priorityLAND
	case EQ, NEQ, GT, LT, GEQ, LEQ:
		return priorityCOMPARER
	case SHL, SHR:
		return priorityBITSHIFT
	case AND, OR, XOR:
		return priorityBIT
	case ADD, SUB:
		return priorityADD
	case MUL, QUO, REM:
		return priorityMUL
	case NOT, NEG:
		return priorityPREFIX
	case CLAUSE:
		return priorityCLAUSE
	case CHAR, STRING, NUMBER, BOOL, VARIABLE, SELECTOR, ACCESSOR, LITERAL:
		return priorityLITERAL
	}
	return priorityUNKNOWN
}

var tokenMap map[string]TokenType

func init() {
	tokenMap = make(map[string]TokenType)
	for i := ILLEGAL; i < EOF; i++ {
		tokenMap[tokens[i]] = i
	}
}

func (op TokenType) isTernary() bool {
	return op == TERNARY_IF || op == TERNARY_ELSE
}

func (op TokenType) isLogical() bool {
	return op == LAND || op == LOR
}

func (op TokenType) isShortCircuit() bool {
	return op.isTernary() || op.isLogical()
}

var tokenCOMPARER = map[TokenType]struct{}{
	EQ:  {},
	NEQ: {},
	GT:  {},
	GEQ: {},
	LT:  {},
	LEQ: {},
}

var tokenPREFIX = map[TokenType]struct{}{
	NEG: {},
	NOT: {},
}

var tokenBIT = map[TokenType]struct{}{
	XOR: {},
	AND: {},
	OR:  {},
}

var tokenBITSHIFT = map[TokenType]struct{}{
	SHR: {},
	SHL: {},
}

var tokenADD = map[TokenType]struct{}{
	ADD: {},
	SUB: {},
}

var tokenMUL = map[TokenType]struct{}{
	MUL: {},
	QUO: {},
	REM: {},
}

var tokenTERNARY = map[TokenType]struct{}{
	TERNARY_IF:   {},
	TERNARY_ELSE: {},
}
