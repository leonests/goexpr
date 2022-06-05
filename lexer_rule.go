package goexpr

import (
	"fmt"
)

type lexerRule struct {
	isStartable   bool
	isTerminable  bool
	isNullable    bool
	nextAllowable map[TokenType]struct{}
}

var lexerRules = map[TokenType]lexerRule{
	ILLEGAL: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			CHAR:     {},
			STRING:   {},
			NUMBER:   {},
			BOOL:     {},
			VARIABLE: {},
			NOT:      {},
			NEG:      {},
			LPAREN:   {},
			SELECTOR: {},
		},
	},
	CHAR: {
		isStartable:  true,
		isTerminable: true,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			EQ:           {},
			NEQ:          {},
			LT:           {},
			GT:           {},
			LEQ:          {},
			GEQ:          {},
			LAND:         {},
			LOR:          {},
			TERNARY_IF:   {},
			TERNARY_ELSE: {},
			RPAREN:       {},
		},
	},
	STRING: {
		isStartable:  true,
		isTerminable: true,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			EQ:           {},
			NEQ:          {},
			LT:           {},
			GT:           {},
			LEQ:          {},
			GEQ:          {},
			ADD:          {},
			LAND:         {},
			LOR:          {},
			TERNARY_IF:   {},
			TERNARY_ELSE: {},
			RPAREN:       {},
		},
	},
	NUMBER: {
		isStartable:  true,
		isTerminable: true,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			EQ:           {},
			NEQ:          {},
			LT:           {},
			GT:           {},
			LEQ:          {},
			GEQ:          {},
			ADD:          {},
			SUB:          {},
			MUL:          {},
			QUO:          {},
			REM:          {},
			AND:          {},
			OR:           {},
			XOR:          {},
			SHL:          {},
			SHR:          {},
			LAND:         {},
			LOR:          {},
			TERNARY_IF:   {},
			TERNARY_ELSE: {},
			RPAREN:       {},
		},
	},
	BOOL: {
		isStartable:  true,
		isTerminable: true,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			EQ:           {},
			NEQ:          {},
			LAND:         {},
			LOR:          {},
			TERNARY_IF:   {},
			TERNARY_ELSE: {},
			RPAREN:       {},
		},
	},
	VARIABLE: {
		isStartable:  true,
		isTerminable: true,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			EQ:           {},
			NEQ:          {},
			LT:           {},
			GT:           {},
			LEQ:          {},
			GEQ:          {},
			ADD:          {},
			SUB:          {},
			MUL:          {},
			QUO:          {},
			REM:          {},
			AND:          {},
			OR:           {},
			XOR:          {},
			SHL:          {},
			SHR:          {},
			LAND:         {},
			LOR:          {},
			TERNARY_IF:   {},
			TERNARY_ELSE: {},
			RPAREN:       {},
			LBRACKET:     {},
			RBRACKET:     {},
		},
	},
	ACCESSOR: {
		isStartable:  false,
		isTerminable: true,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			EQ:           {},
			NEQ:          {},
			LT:           {},
			GT:           {},
			LEQ:          {},
			GEQ:          {},
			ADD:          {},
			SUB:          {},
			MUL:          {},
			QUO:          {},
			REM:          {},
			AND:          {},
			OR:           {},
			XOR:          {},
			SHL:          {},
			SHR:          {},
			LAND:         {},
			LOR:          {},
			TERNARY_IF:   {},
			TERNARY_ELSE: {},
			RPAREN:       {},
			LBRACKET:     {},
			RBRACKET:     {},
		},
	},
	SELECTOR: {
		isStartable:  true,
		isTerminable: true,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			EQ:           {},
			NEQ:          {},
			LT:           {},
			GT:           {},
			LEQ:          {},
			GEQ:          {},
			ADD:          {},
			SUB:          {},
			MUL:          {},
			QUO:          {},
			REM:          {},
			AND:          {},
			OR:           {},
			XOR:          {},
			SHL:          {},
			SHR:          {},
			LAND:         {},
			LOR:          {},
			TERNARY_IF:   {},
			TERNARY_ELSE: {},
			RPAREN:       {},
			LBRACKET:     {},
			RBRACKET:     {},
		},
	},
	LPAREN: {
		isStartable:  true,
		isTerminable: false,
		isNullable:   true,
		nextAllowable: map[TokenType]struct{}{
			CHAR:     {},
			STRING:   {},
			NUMBER:   {},
			BOOL:     {},
			VARIABLE: {},
			NOT:      {},
			NEG:      {},
			LPAREN:   {},
			RPAREN:   {},
			SELECTOR: {},
		},
	},
	RPAREN: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   true,
		nextAllowable: map[TokenType]struct{}{
			CHAR:     {},
			STRING:   {},
			NUMBER:   {},
			BOOL:     {},
			VARIABLE: {},
			NOT:      {},
			NEG:      {},
			SELECTOR: {},
			ACCESSOR: {},
			LPAREN:   {},
			RPAREN:   {},
			LBRACKET: {},
			RBRACKET: {},
		},
	},
	LBRACKET: {
		isStartable:  false,
		isTerminable: true,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			STRING:   {},
			NUMBER:   {},
			VARIABLE: {},
			LPAREN:   {},
			SELECTOR: {},
		},
	},
	RBRACKET: {
		isStartable:  false,
		isTerminable: true,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			EQ:           {},
			NEQ:          {},
			LT:           {},
			GT:           {},
			LEQ:          {},
			GEQ:          {},
			ADD:          {},
			SUB:          {},
			MUL:          {},
			QUO:          {},
			REM:          {},
			AND:          {},
			OR:           {},
			XOR:          {},
			SHL:          {},
			SHR:          {},
			LAND:         {},
			LOR:          {},
			TERNARY_IF:   {},
			TERNARY_ELSE: {},
			RPAREN:       {},
			SELECTOR:     {},
			ACCESSOR:     {},
		},
	},
	ADD: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			NEG:      {},
			CHAR:     {},
			STRING:   {},
			NUMBER:   {},
			VARIABLE: {},
			SELECTOR: {},
		},
	},
	SUB: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			NEG:      {},
			CHAR:     {},
			NUMBER:   {},
			VARIABLE: {},
			SELECTOR: {},
		},
	},
	MUL: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			NEG:      {},
			NUMBER:   {},
			VARIABLE: {},
			SELECTOR: {},
		},
	},
	QUO: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			NEG:      {},
			NUMBER:   {},
			VARIABLE: {},
			SELECTOR: {},
		},
	},
	REM: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			NEG:      {},
			NUMBER:   {},
			VARIABLE: {},
			SELECTOR: {},
		},
	},
	AND: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			NUMBER:   {},
			VARIABLE: {},
			SELECTOR: {},
		},
	},
	OR: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			NUMBER:   {},
			VARIABLE: {},
			SELECTOR: {},
		},
	},
	XOR: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			NUMBER:   {},
			VARIABLE: {},
			SELECTOR: {},
		},
	},
	SHL: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			NUMBER:   {},
			VARIABLE: {},
			SELECTOR: {},
		},
	},
	SHR: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			NUMBER:   {},
			VARIABLE: {},
			SELECTOR: {},
		},
	},
	EQ: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			CHAR:     {},
			STRING:   {},
			NUMBER:   {},
			BOOL:     {},
			VARIABLE: {},
			NOT:      {},
			NEG:      {},
			LPAREN:   {},
		},
	},
	NEQ: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			CHAR:     {},
			STRING:   {},
			NUMBER:   {},
			BOOL:     {},
			VARIABLE: {},
			NOT:      {},
			NEG:      {},
			LPAREN:   {},
		},
	},
	LT: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			CHAR:     {},
			STRING:   {},
			NUMBER:   {},
			BOOL:     {},
			VARIABLE: {},
			NOT:      {},
			NEG:      {},
			LPAREN:   {},
		},
	},
	GT: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			CHAR:     {},
			STRING:   {},
			NUMBER:   {},
			BOOL:     {},
			VARIABLE: {},
			NOT:      {},
			NEG:      {},
			LPAREN:   {},
		},
	},
	LEQ: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			CHAR:     {},
			STRING:   {},
			NUMBER:   {},
			BOOL:     {},
			VARIABLE: {},
			NOT:      {},
			NEG:      {},
			LPAREN:   {},
		},
	},
	GEQ: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			CHAR:     {},
			STRING:   {},
			NUMBER:   {},
			BOOL:     {},
			VARIABLE: {},
			NOT:      {},
			NEG:      {},
			LPAREN:   {},
		},
	},
	LAND: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			CHAR:     {},
			STRING:   {},
			NUMBER:   {},
			BOOL:     {},
			VARIABLE: {},
			NOT:      {},
			NEG:      {},
			LPAREN:   {},
		},
	},
	LOR: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			CHAR:     {},
			STRING:   {},
			NUMBER:   {},
			BOOL:     {},
			VARIABLE: {},
			NOT:      {},
			NEG:      {},
			LPAREN:   {},
		},
	},
	TERNARY_IF: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			CHAR:     {},
			STRING:   {},
			NUMBER:   {},
			BOOL:     {},
			VARIABLE: {},
			NOT:      {},
			NEG:      {},
			LPAREN:   {},
		},
	},
	TERNARY_ELSE: {
		isStartable:  false,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			CHAR:     {},
			STRING:   {},
			NUMBER:   {},
			BOOL:     {},
			VARIABLE: {},
			NOT:      {},
			NEG:      {},
			LPAREN:   {},
		},
	},
	NOT: {
		isStartable:  true,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			BOOL:     {},
			VARIABLE: {},
			LPAREN:   {},
			SELECTOR: {},
		},
	},
	NEG: {
		isStartable:  true,
		isTerminable: false,
		isNullable:   false,
		nextAllowable: map[TokenType]struct{}{
			NUMBER:   {},
			VARIABLE: {},
			LPAREN:   {},
			SELECTOR: {},
		},
	},
}

func (ls *lexerRule) hasNextAllowable(tokenType TokenType) bool {
	if _, ok := ls.nextAllowable[tokenType]; !ok {
		return false
	}
	return true
}

func getLexerRule(tokenType TokenType) (lexerRule, error) {
	if lc, ok := lexerRules[tokenType]; !ok {
		return lexerRules[ILLEGAL], fmt.Errorf("no lexer rule found for toke type %v", tokenType.String())
	} else {
		return lc, nil
	}
}

func checkLexerBalance(tokens []LexerToken) error {
	var parens, brackets int
	stream := newLexerStream(tokens)
	for stream.notEOF() {
		token := stream.flowForward()
		switch token.Type {
		case LPAREN:
			parens++
		case LBRACKET:
			brackets++
		case RPAREN:
			parens--
		case RBRACKET:
			brackets--
		}
	}
	if parens != 0 || brackets != 0 {
		return fmt.Errorf("unbalanced parenthesis or bracket")
	}
	return nil
}
