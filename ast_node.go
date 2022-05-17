package rulengine

type astNode struct {
	operator   TokenType
	left       *astNode
	right      *astNode
	rightList  []*astNode
	leftCheck  nodeTypeCheck
	rightCheck nodeTypeCheck
	bothCheck  bothTypeCheck
	calculator calculator
	err        string
}

type nodeTypeCheck func(value interface{}) bool
type bothTypeCheck func(left, right interface{}) bool
type calculator func(left, right interface{}, params map[string]interface{}) (interface{}, error)

var opCalculator = map[TokenType]calculator{
	EQ:           calculatorEQ,
	NEQ:          calculatorNEQ,
	GT:           calculatorGT,
	GEQ:          calculatorGEQ,
	LT:           calculatorLT,
	LEQ:          calculatorLEQ,
	ADD:          calculatorADD,
	SUB:          calculatorSUB,
	MUL:          calculatorMUL,
	QUO:          calculatorQUO,
	REM:          calculatorREM,
	AND:          calculatorAND,
	OR:           calculatorOR,
	XOR:          calculatorXOR,
	SHL:          calculatorSHL,
	SHR:          calculatorSHR,
	LAND:         calculatorLAND,
	LOR:          calculatorLOR,
	TERNARY_IF:   calculatorTernaryIf,
	TERNARY_ELSE: calculatorTernaryElse,
	NOT:          calculatorNOT,
	NEG:          calculatorNEG,
}

func buildSelectorNode(token LexerToken) *astNode {
	return &astNode{
		operator:   LITERAL,
		right:      nil,
		calculator: calculatorSELECTOR(token.Value.([]string)),
		err:        errSelectorFormat,
	}
}

func buildAccessorNode(token LexerToken) *astNode {
	return &astNode{
		operator:   LITERAL,
		right:      nil,
		calculator: calculatorACCESSOR(token.Value.([]string)),
		err:        errAccessorFormat,
	}
}
