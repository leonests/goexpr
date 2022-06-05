package goexpr

import (
	"testing"
)

type ParseTokenTest struct {
	Name   string
	Input  string
	Wanted []LexerToken
}

func TestNumericParse(t *testing.T) {
	parseTokenTests := []ParseTokenTest{
		{
			Name:  "Just zero",
			Input: "0",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 0.0,
				},
			},
		},
		{
			Name:  "Single small number",
			Input: "35",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 35.0,
				},
			},
		},
		{
			Name:  "Single int number",
			Input: "3276842433",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 3276842433.0,
				},
			},
		},
		{
			Name:  "Single small float number",
			Input: "0.5",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 0.5,
				},
			},
		},
		{
			Name:  "Single large float number",
			Input: "32.7684237676",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 32.7684237676,
				},
			},
		},
	}
	runParseTokenTest(parseTokenTests, t)
}

func TestStringParse(t *testing.T) {
	parseTokenTests := []ParseTokenTest{
		{
			Name:  "Single string",
			Input: `"f**k"`,
			Wanted: []LexerToken{
				{
					Type:  STRING,
					Value: "f**k",
				},
			},
		},
		{
			Name:  "Single long string",
			Input: `"f**k-F**K_F**kf**K"`,
			Wanted: []LexerToken{
				{
					Type:  STRING,
					Value: "f**k-F**K_F**kf**K",
				},
			},
		},
	}
	runParseTokenTest(parseTokenTests, t)
}

func TestBooleanParse(t *testing.T) {
	parseTokenTests := []ParseTokenTest{
		{
			Name:  "Single boolean",
			Input: "true",
			Wanted: []LexerToken{
				{
					Type:  BOOL,
					Value: true,
				},
			},
		},
		{
			Name:  "Single boolean",
			Input: "True",
			Wanted: []LexerToken{
				{
					Type:  BOOL,
					Value: true,
				},
			},
		},
		{
			Name:  "Single boolean",
			Input: "TRUE",
			Wanted: []LexerToken{
				{
					Type:  BOOL,
					Value: true,
				},
			},
		},
		{
			Name:  "Single boolean",
			Input: "false",
			Wanted: []LexerToken{
				{
					Type:  BOOL,
					Value: false,
				},
			},
		},
		{
			Name:  "Single boolean",
			Input: "False",
			Wanted: []LexerToken{
				{
					Type:  BOOL,
					Value: false,
				},
			},
		},
		{
			Name:  "Single boolean",
			Input: "FALSE",
			Wanted: []LexerToken{
				{
					Type:  BOOL,
					Value: false,
				},
			},
		},
	}
	runParseTokenTest(parseTokenTests, t)
}

func TestModifierParse(t *testing.T) {
	parseTokenTests := []ParseTokenTest{
		{
			Name:  "ADD",
			Input: "1+2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  ADD,
					Value: "+",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "ADD",
			Input: "1 + 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  ADD,
					Value: "+",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "ADD",
			Input: "1+ 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  ADD,
					Value: "+",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "Sub",
			Input: "1-2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  SUB,
					Value: "-",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "Sub",
			Input: "1 - 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  SUB,
					Value: "-",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "Sub",
			Input: "1- 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  SUB,
					Value: "-",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "Mul",
			Input: "1*2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  MUL,
					Value: "*",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "Mul",
			Input: "1 * 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  MUL,
					Value: "*",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "Mul",
			Input: "1 * 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  MUL,
					Value: "*",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "Quotient",
			Input: "1 / 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  QUO,
					Value: "/",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "Quotient",
			Input: "1/ 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  QUO,
					Value: "/",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "Quotient",
			Input: "1/2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  QUO,
					Value: "/",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "Rem",
			Input: "1%2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  REM,
					Value: "%",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "Rem",
			Input: "1 %2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  REM,
					Value: "%",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "Rem",
			Input: "1 % 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  REM,
					Value: "%",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "And",
			Input: "1 & 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  AND,
					Value: "&",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "Or",
			Input: "1 | 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  OR,
					Value: "|",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "Xor",
			Input: "1 ^ 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  XOR,
					Value: "^",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "SHL",
			Input: "1 << 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  SHL,
					Value: "<<",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "SHR",
			Input: "1 >> 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  SHR,
					Value: ">>",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
	}
	runParseTokenTest(parseTokenTests, t)
}

func TestComparatorParse(t *testing.T) {
	parseTokenTests := []ParseTokenTest{
		{
			Name:  "EQ",
			Input: "1 == 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  EQ,
					Value: "==",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "EQ",
			Input: "1== 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  EQ,
					Value: "==",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "EQ",
			Input: "1==2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  EQ,
					Value: "==",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "NEQ",
			Input: "1!=2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  NEQ,
					Value: "!=",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "GT",
			Input: "1 > 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  GT,
					Value: ">",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "LT",
			Input: "1 < 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  LT,
					Value: "<",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "LEQ",
			Input: "1 <= 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  LEQ,
					Value: "<=",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "GEQ",
			Input: "1 >= 2",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  GEQ,
					Value: ">=",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
			},
		},
		{
			Name:  "String EQ",
			Input: `"F**K" == "f**k"`,
			Wanted: []LexerToken{
				{
					Type:  STRING,
					Value: "F**K",
				},
				{
					Type:  EQ,
					Value: "==",
				},
				{
					Type:  STRING,
					Value: "f**k",
				},
			},
		},
		{
			Name:  "String NEQ",
			Input: `"F**K.f**k" != "f**k.F**K"`,
			Wanted: []LexerToken{
				{
					Type:  STRING,
					Value: "F**K.f**k",
				},
				{
					Type:  NEQ,
					Value: "!=",
				},
				{
					Type:  STRING,
					Value: "f**k.F**K",
				},
			},
		},
		{
			Name:  "String NEQ",
			Input: `"F**K.f**k" > "f**k.F**K"`,
			Wanted: []LexerToken{
				{
					Type:  STRING,
					Value: "F**K.f**k",
				},
				{
					Type:  GT,
					Value: ">",
				},
				{
					Type:  STRING,
					Value: "f**k.F**K",
				},
			},
		},
		{
			Name:  "String NEQ",
			Input: `"F**K.f**k" < "f**k.F**K"`,
			Wanted: []LexerToken{
				{
					Type:  STRING,
					Value: "F**K.f**k",
				},
				{
					Type:  LT,
					Value: "<",
				},
				{
					Type:  STRING,
					Value: "f**k.F**K",
				},
			},
		},
		{
			Name:  "String NEQ",
			Input: `"F**K.f**k" <= "f**k.F**K"`,
			Wanted: []LexerToken{
				{
					Type:  STRING,
					Value: "F**K.f**k",
				},
				{
					Type:  LEQ,
					Value: "<=",
				},
				{
					Type:  STRING,
					Value: "f**k.F**K",
				},
			},
		},
		{
			Name:  "String NEQ",
			Input: `"F**K.f**k" >= "f**k.F**K"`,
			Wanted: []LexerToken{
				{
					Type:  STRING,
					Value: "F**K.f**k",
				},
				{
					Type:  GEQ,
					Value: ">=",
				},
				{
					Type:  STRING,
					Value: "f**k.F**K",
				},
			},
		},
	}
	runParseTokenTest(parseTokenTests, t)
}

func TestPrefixParse(t *testing.T) {
	parseTokenTests := []ParseTokenTest{
		{
			Name:  "Sub prefix before num",
			Input: "-1",
			Wanted: []LexerToken{
				{
					Type:  NEG,
					Value: "-",
				},
				{
					Type:  NUMBER,
					Value: 1.0,
				},
			},
		},
		{
			Name:  "Sub prefix before var",
			Input: "-fxxk",
			Wanted: []LexerToken{
				{
					Type:  NEG,
					Value: "-",
				},
				{
					Type:  VARIABLE,
					Value: "fxxk",
				},
			},
		},
		{
			Name:  "Bool prefix before var",
			Input: "!true",
			Wanted: []LexerToken{
				{
					Type:  NOT,
					Value: "!",
				},
				{
					Type:  BOOL,
					Value: true,
				},
			},
		},
		{
			Name:  "Bool prefix before var",
			Input: "!False",
			Wanted: []LexerToken{
				{
					Type:  NOT,
					Value: "!",
				},
				{
					Type:  BOOL,
					Value: false,
				},
			},
		},
	}
	runParseTokenTest(parseTokenTests, t)
}

func TestTernaryParse(t *testing.T) {
	parseTokenTests := []ParseTokenTest{
		{
			Name:  "Ternary with bool",
			Input: "true ? 1",
			Wanted: []LexerToken{
				{
					Type:  BOOL,
					Value: true,
				},
				{
					Type:  TERNARY_IF,
					Value: "?",
				},
				{
					Type:  NUMBER,
					Value: 1.0,
				},
			},
		},
		{
			Name:  "Ternary with bool",
			Input: "false ? a",
			Wanted: []LexerToken{
				{
					Type:  BOOL,
					Value: false,
				},
				{
					Type:  TERNARY_IF,
					Value: "?",
				},
				{
					Type:  VARIABLE,
					Value: "a",
				},
			},
		},
		{
			Name:  "Ternary with bool",
			Input: `false ? "a"`,
			Wanted: []LexerToken{
				{
					Type:  BOOL,
					Value: false,
				},
				{
					Type:  TERNARY_IF,
					Value: "?",
				},
				{
					Type:  STRING,
					Value: "a",
				},
			},
		},
		{
			Name:  "Ternary with bool",
			Input: "1 == 2 ? a",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  EQ,
					Value: "==",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
				{
					Type:  TERNARY_IF,
					Value: "?",
				},
				{
					Type:  VARIABLE,
					Value: "a",
				},
			},
		},
		{
			Name:  "Ternary with bool",
			Input: "1 == 2 ? a : b",
			Wanted: []LexerToken{
				{
					Type:  NUMBER,
					Value: 1.0,
				},
				{
					Type:  EQ,
					Value: "==",
				},
				{
					Type:  NUMBER,
					Value: 2.0,
				},
				{
					Type:  TERNARY_IF,
					Value: "?",
				},
				{
					Type:  VARIABLE,
					Value: "a",
				},
				{
					Type:  TERNARY_ELSE,
					Value: ":",
				},
				{
					Type:  VARIABLE,
					Value: "b",
				},
			},
		},
	}
	runParseTokenTest(parseTokenTests, t)
}

func runParseTokenTest(parseTokenTests []ParseTokenTest, t *testing.T) {
	var (
		wantedTokenLength, actualTokenLength int
		actualToken                          LexerToken
	)

	for _, test := range parseTokenTests {
		exprTokens, err := lexerScan(test.Input)
		if err != nil {
			t.Logf("Test '%s' failed:", test.Name)
			t.Logf("Expression: '%s' Error: %s", test.Input, err)
			t.Fail()
			continue
		}
		wantedTokenLength = len(test.Wanted)
		actualTokenLength = len(exprTokens)

		if wantedTokenLength != actualTokenLength {
			t.Logf("Test '%s' failed:", test.Name)
			t.Logf("Wanted: %d Actually: %d Error: %s", wantedTokenLength, actualTokenLength, "length not match")
			t.Fail()
			continue
		}

		for idx, wantedToken := range test.Wanted {
			actualToken = exprTokens[idx]
			if wantedToken.Type != actualToken.Type {
				t.Logf("Test '%s' failed:", test.Name)
				t.Logf("Wanted: %s Actually: %s Error: %s", wantedToken.Type.String(), actualToken.Type.String(), "token type not match")
				t.Fail()
				continue
			}

			if actualToken.Value != wantedToken.Value {
				t.Logf("Test '%s' failed:", test.Name)
				t.Logf("Wanted: %v Actually: %v Error: %v", wantedToken.Value, actualToken.Value, "token value not match")
				t.Fail()
				continue
			}
		}
	}
}
