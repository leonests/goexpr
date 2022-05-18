package rulengine

import (
	"encoding/json"
	"fmt"
	"testing"
)

type ParseAstTest struct {
	Name   string
	Input  string
	Params map[string]interface{}
	Wanted interface{}
}

type Param struct {
	String    string  `json:"string"`
	Int       int     `json:"int"`
	Int64     int64   `json:"int64"`
	Float     float64 `json:"float64"`
	Bool      bool    `json:"bool"`
	Interface interface{}
	Array     []interface{}
	Map       map[string]interface{}
}

var param = Param{
	String:    "string_test1",
	Int:       0,
	Int64:     1,
	Float:     1.0,
	Bool:      false,
	Interface: nil,
	Array: []interface{}{
		Param{
			String:    "string_test2",
			Int:       0,
			Int64:     1,
			Float:     1.0,
			Bool:      true,
			Interface: nil,
			Map: map[string]interface{}{
				"key_str":   "val",
				"key_int":   1,
				"key_bool":  true,
				"key_float": 1.0,
				"key_struct": Param{
					String:    "string_test2",
					Int:       -1,
					Bool:      true,
					Interface: nil,
				},
			},
		},
		"rulengine",
	},
	Map: map[string]interface{}{
		"key_str":   "val_string",
		"key_int":   1,
		"key_bool":  true,
		"key_float": 1.0,
		"key_struct": Param{
			String:    "string_test2",
			Int:       -1,
			Bool:      true,
			Interface: nil,
		},
		"key_interface": nil,
	},
}

func TestParseAst(t *testing.T) {
	parseAstTests := []ParseAstTest{
		// {
		// 	Name:   "Simple Char",
		// 	Input:  `'b' == 98`,
		// 	Wanted: true,
		// },
		// {
		// 	Name:   "Simple Char",
		// 	Input:  `'b' - 1`,
		// 	Wanted: 'a',
		// },
		{
			Name:   "Simple Char",
			Input:  `'a'`,
			Wanted: 'a',
		},
		{
			Name:   "Simple char",
			Input:  `'a' == "b"`,
			Wanted: false,
		},
		{
			Name:   "Simple string",
			Input:  `"a"`,
			Wanted: "a",
		},
		{
			Name:   "Simple char",
			Input:  `'a'`,
			Wanted: 'a',
		},

		{
			Name:   "Simple ADD",
			Input:  "1 + 2",
			Wanted: 3.0,
		},
		{
			Name:   "Simple ADD",
			Input:  `"abc" + "def"`,
			Wanted: "abcdef",
		},
		{
			Name:   "Simple ADD",
			Input:  `123 + "def"`,
			Wanted: "123def",
		},
		{
			Name:   "Simple ADD",
			Input:  `"abc" + 123`,
			Wanted: "abc123",
		},
		{
			Name:   "Simple ADD",
			Input:  `"abc" + true`,
			Wanted: "abctrue",
		},
		{
			Name:   "Simple ADD",
			Input:  `false + "abc" + true`,
			Wanted: "falseabctrue",
		},
		{
			Name:   "Simple SUB",
			Input:  "7 - 5",
			Wanted: 2.0,
		},
		{
			Name:   "Simple SUB",
			Input:  "55 - 77",
			Wanted: -22.0,
		},
		{
			Name:   "Simple MUL",
			Input:  "55 * 77",
			Wanted: 4235.0,
		},
		{
			Name:   "Simple QUO",
			Input:  "55 / 11",
			Wanted: 5.0,
		},
		{
			Name:   "Simple REM",
			Input:  "55 % 10",
			Wanted: 5.0,
		},
		{
			Name:   "Simple REM",
			Input:  "55 % 11",
			Wanted: 0.0,
		},
		{
			Name:   "Multi ADD",
			Input:  "1 + 2 + 30",
			Wanted: 33.0,
		},
		{
			Name:   "Multi ADD",
			Input:  `"1" + "2" + "30"`,
			Wanted: "1230",
		},
		{
			Name:   "Multi SUB",
			Input:  "1 - 2 - 30",
			Wanted: -31.0,
		},
		{
			Name:   "Multi OPERATOR",
			Input:  "2 * -2",
			Wanted: -4.0,
		},
		{
			Name:   "Multi OPERATOR",
			Input:  "1 + 2 - 3 + 4 - 5",
			Wanted: -1.0,
		},
		{
			Name:   "Multi OPERATOR",
			Input:  "1 - 2 + 3 - 4 + 5",
			Wanted: 3.0,
		},
		{
			Name:   "Multi OPERATOR",
			Input:  "1 + 2 * 3 + 4 - 5",
			Wanted: 6.0,
		},
		{
			Name:   "Multi OPERATOR",
			Input:  "1 / 2 * 3 / 4 * 5",
			Wanted: 1.875,
		},
		{
			Name:   "Multi OPERATOR",
			Input:  "1 / 2 * 3 + 4 * 5",
			Wanted: 21.5,
		},
		{
			Name:   "Multi OPERATOR",
			Input:  "1 / 2 * 3 + 4 * 5 + 6 % 7",
			Wanted: 27.5,
		},
		{
			Name:   "Multi OPERATOR",
			Input:  "1 * 2 / 4 * 3",
			Wanted: 1.5,
		},
		{
			Name:   "PAREN",
			Input:  "1 / 2 * (3 + 4) * 5",
			Wanted: 17.5,
		},
		{
			Name:   "PAREN RECURSIVE",
			Input:  "1 / ( 2 * (3 + 4)) * 7",
			Wanted: 0.5,
		},
		{
			Name:   "PAREN RECURSIVE",
			Input:  "1 * ( 2 * (3 + 4)) % 7",
			Wanted: 0.0,
		},
		{
			Name:   "Simple PREFIX",
			Input:  "-1",
			Wanted: -1.0,
		},
		{
			Name:   "NEG PREFIX",
			Input:  "-(1 * ( 2 * (3 + 4)) % 7)",
			Wanted: -0.0,
		},
		{
			Name:   "Simple SHL",
			Input:  "2 << 1",
			Wanted: 4.0,
		},
		{
			Name:   "Simple SHR",
			Input:  "2 >> 1",
			Wanted: 1.0,
		},
		{
			Name:   "Simple AND",
			Input:  "71 & 23",
			Wanted: 7.0,
		},
		{
			Name:   "Simple OR",
			Input:  "71 | 23",
			Wanted: 87.0,
		},
		{
			Name:   "Simple XOR",
			Input:  "71 ^ 23",
			Wanted: 80.0,
		},
		{
			Name:   "Multi BITWISE",
			Input:  "71 ^ (23 | (71 & 23))",
			Wanted: 80.0,
		},
		{
			Name:   "Multi BIT",
			Input:  "1 << 2 & 4",
			Wanted: 4.0,
		},
		{
			Name:   "Multi BIT",
			Input:  "1 << 2 & 15",
			Wanted: 4.0,
		},
		{
			Name:   "Multi BIT",
			Input:  "1 << 2 | 11",
			Wanted: 15.0,
		},
		{
			Name:   "Simple COMPARATOR",
			Input:  "3 > 5",
			Wanted: false,
		},
		{
			Name:   "Simple COMPARATOR",
			Input:  "3 < 5",
			Wanted: true,
		},
		{
			Name:   "Simple COMPARATOR",
			Input:  "3 < 5",
			Wanted: true,
		},
		{
			Name:   "Simple COMPARATOR",
			Input:  "3 == 5",
			Wanted: false,
		},
		{
			Name:   "Simple COMPARATOR",
			Input:  "3 != 5",
			Wanted: true,
		},
		{
			Name:   "Simple COMPARATOR",
			Input:  `"abc" != "adc"`,
			Wanted: true,
		},
		{
			Name:   "Simple COMPARATOR",
			Input:  `"abc" == "adc"`,
			Wanted: false,
		},
		{
			Name:   "Simple COMPARATOR",
			Input:  `"abc" == "abc"`,
			Wanted: true,
		},
		{
			Name:   "Simple COMPARATOR",
			Input:  `"ab" <= "abc"`,
			Wanted: true,
		},

		{
			Name:   "Simple COMPARATOR",
			Input:  `"ab" >= "abc"`,
			Wanted: false,
		},
		{
			Name:   "Simple COMPARATOR",
			Input:  `"aba" > "abc"`,
			Wanted: false,
		},
		{
			Name:   "Simple COMPARATOR",
			Input:  "-1 < 2",
			Wanted: true,
		},
		{
			Name:   "Simple COMPARATOR",
			Input:  `"ab" < "abc"`,
			Wanted: true,
		},
		{
			Name:   "Simple LOGICAL",
			Input:  "true && true || true && false",
			Wanted: true,
		},
		{
			Name:   "Simple LOGICAL",
			Input:  "false || true && true || false",
			Wanted: true,
		},
		{
			Name:   "Simple LOGICAL",
			Input:  "false && true || true",
			Wanted: true,
		},
		{
			Name:   "Simple LOGICAL",
			Input:  "true || false && true",
			Wanted: true,
		},
		{
			Name:   "Simple LOGICAL",
			Input:  "true && true || false && false",
			Wanted: true,
		},
		{
			Name:   "Simple LOGICAL",
			Input:  "(3 != 5) && (true == true)",
			Wanted: true,
		},
		{
			Name:   "Simple LOGICAL",
			Input:  "(3 == 5) || (false == false)",
			Wanted: true,
		},
		{
			Name:   "Simple LOGICAL",
			Input:  "7 > 5 || 2 <= 4",
			Wanted: true,
		},
		{
			Name:   "Simple LOGICAL",
			Input:  "7 < 5 && 2 <= 4",
			Wanted: false,
		},
		{
			Name:   "Multi LOGICAL",
			Input:  "7 < 5 && 2 <= 4 || 1 >= 1",
			Wanted: true,
		},
		{
			Name:   "Multi LOGICAL",
			Input:  "7 >= 5 && (2 <= 4 || 1 >= 3)",
			Wanted: true,
		},
		{
			Name:   "Multi LOGICAL",
			Input:  "7 >= 5 && (2 <= 4 && 1 != 3)",
			Wanted: true,
		},
		{
			Name:   "Multi OPERATOR",
			Input:  "4 + 3 << 2 * 2",
			Wanted: 112.0,
		},
		{
			Name:   "Multi OPERATOR",
			Input:  "!TRUE == FALSE",
			Wanted: true,
		},
		{
			Name:   "Multi OPERATOR",
			Input:  `!("x" < "y")`,
			Wanted: false,
		},
		{
			Name:   "Multi OPERATOR",
			Input:  "2 * -2",
			Wanted: -4.0,
		},
		{
			Name:   "Multi OPERATOR",
			Input:  `!("x" < "y")`,
			Wanted: false,
		},
		{
			Name:   "Simple TERNARY",
			Input:  "true ? -77.7",
			Wanted: -77.7,
		},
		{
			Name:   "Simple TERNARY",
			Input:  `true ? "string_test1"`,
			Wanted: "string_test1",
		},
		{
			Name:   "Simple TERNARY",
			Input:  "true ? false",
			Wanted: false,
		},
		{
			Name:   "Simple TERNARY",
			Input:  "false ? 1",
			Wanted: nil,
		},
		{
			Name:   "Multi TERNARY",
			Input:  `false ? "string_test1" : "string_test2"`,
			Wanted: "string_test2",
		},
		{
			Name:   "Multi TERNARY",
			Input:  `true ? "string_test1" : "string_test2"`,
			Wanted: "string_test1",
		},
		{
			Name:   "Simple TERNARY",
			Input:  "1 < 2 ? 3",
			Wanted: 3.0,
		},
		{
			Name:   "Simple TERNARY",
			Input:  "1 > 2 ? 3",
			Wanted: nil,
		},
		{
			Name:   "Multi TERNARY",
			Input:  "1 > 2 ? 3 : 4",
			Wanted: 4.0,
		},
		{
			Name:   "Multi TERNARY",
			Input:  "1 < 2 ? 3 : 4",
			Wanted: 3.0,
		},
		{
			Name:   "Multi TERNARY",
			Input:  "(3.0 * 2 - 3 % 2 > 4) ? (1010 / 5) : 4",
			Wanted: 202.0,
		},
		{
			Name:   "Multi TERNARY",
			Input:  "((3.0 * 2 - 3 % 2 > 4) ? (1010 / 5) : 4) > 201",
			Wanted: true,
		},
		{
			Name:   "Multi TERNARY",
			Input:  "((3.0 * 2 - 3 % 2 > 4) ? (1010 / 5) : 4) <= 202",
			Wanted: true,
		},
		{
			Name:   "Multi TERNARY",
			Input:  "1 != 2 ? 3 : (4 == 5 ? 6 : 7)",
			Wanted: 3.0,
		},
		{
			Name:   "Multi TERNARY",
			Input:  `1 >= 2 ? true : 4 > 5 ? 6 : 7 <= 8 ? "abc" : 10`,
			Wanted: "abc",
		},
		{
			Name:   "Multi TERNARY",
			Input:  "1 <= 2 ? true : 4 > 5 ? 6 : 7 == 8 ? 9 : 10",
			Wanted: true,
		},
		{
			Name:   "Multi TERNARY",
			Input:  "1 > 2 ? 3 : true ? 6 : 7 <= 8 ? 9 : 10",
			Wanted: 6.0,
		},
		{
			Name:   "Multi TERNARY",
			Input:  "1 > 2 ? 3 : 4 > 5 ? 6 : 7 <= 8 ? 9 : 10",
			Wanted: 9.0,
		},
		{
			Name:   "Multi TERNARY",
			Input:  "1 > 2 ? 3 : 4 > 5 ? 6 : 7 == 8 ? 9 : 10",
			Wanted: 10.0,
		},
		{
			Name:   "Multi TERNARY",
			Input:  "(1 != 2 ? 3 : 4 != 5 ? 6 : 7) != 0",
			Wanted: true,
		},
		{
			Name:   "Multi TERNARY",
			Input:  "1 == 2 ? 3 : 4 != 5 ? 6 : 7",
			Wanted: 6.0,
		},
		{
			Name:   "Multi TERNARY",
			Input:  "1 != 2 ? 3 : 4 == 5 ? 6 : 7",
			Wanted: 3.0,
		},
		{
			Name:   "Multi TERNARY",
			Input:  `(1 != 2 ? true : 4 == 5 ? 6 : 7) ? (1 == 2 ? 3 : 4 != 5 ? 6 : 7) : "abc"`,
			Wanted: 6.0,
		},
		{
			Name:   "Multi TERNARY",
			Input:  "(1 != 2 ? 3 : 4 == 5 ? 6 : 7) + (1 == 2 ? 3 : 4 != 5 ? 6 : 7)",
			Wanted: 9.0,
		},
		{
			Name:   "Multi TERNARY",
			Input:  `(1 != 2 ? false : 4 == 5 ? 6 : 7) ? (1 == 2 ? 3 : 4 != 5 ? 6 : 7) : "abc"`,
			Wanted: "abc",
		},
		{
			Name:   "Multi TERNARY",
			Input:  `(1 != 2 ? true : 4 == 5 ? 6 : 7) ? (1 == 2 ? 3 : 4 != 5 ? 6 : 7) : "abc"`,
			Wanted: 6.0,
		},
	}
	runParseAstTests(parseAstTests, t)
}
func TestParseAstWithParams(t *testing.T) {
	mapData := map[string]interface{}{
		"param": param,
		"a":     "b",
	}
	byteMap, _ := json.Marshal(mapData)
	var m map[string]interface{}
	json.Unmarshal(byteMap, &m)

	parseAstTests := []ParseAstTest{
		{
			Name:   "Simple Selector",
			Input:  "a.0",
			Params: m,
			Wanted: 'b',
		},
		{
			Name:   "parameter",
			Input:  `param["Array"][param.Array[0].int].float64`,
			Params: m,
			Wanted: 1.0,
		},
		{
			Name:  "parameter",
			Input: "param1 < param2",
			Params: map[string]interface{}{
				"param1": 1,
				"param2": 2,
			},
			Wanted: true,
		},
		{
			Name:   "Accessor Nested Bracket",
			Input:  `param["Array"][param.Array[0].int64]`,
			Params: m,
			Wanted: "rulengine",
		},
		{
			Name:   "Selector Parallel Bracket",
			Input:  `param.Array[param.Array[0]["int64"]]`,
			Params: m,
			Wanted: "rulengine",
		},
		{
			Name:   "Selector Parallel Bracket",
			Input:  `param.Array[param.Array[0]["Map"]["key_int"]]`,
			Params: m,
			Wanted: "rulengine",
		},
		{
			Name:   "Selector Parallel Bracket",
			Input:  `param.Array[0]["Map"]["key_int"]`,
			Params: m,
			Wanted: 1.0,
		},
		{
			Name:   "Parallel Bracket",
			Input:  `param["Map"]["key_interface"]`,
			Params: m,
			Wanted: nil,
		},
		{
			Name:   "Simple Bracket",
			Input:  `param["Interface"]`,
			Params: m,
			Wanted: nil,
		},
		{
			Name:   "Simple Bracket",
			Input:  `param["bool"]`,
			Params: m,
			Wanted: false,
		},
		{
			Name:   "Simple Bracket",
			Input:  `param["string"]`,
			Params: m,
			Wanted: "string_test1",
		},
		{
			Name:   "Simple Bracket",
			Input:  `param["string"]`,
			Params: m,
			Wanted: "string_test1",
		},
		{
			Name:   "Selector Bracket",
			Input:  `param.Array[param.Array.0.Map.key_int]`,
			Params: m,
			Wanted: "rulengine",
		},
		{
			Name:   "Selector Accessor Bracket",
			Input:  `param.Array[0].string`,
			Params: m,
			Wanted: "string_test2",
		},
		{
			Name:   "Selector Bracket",
			Input:  `param.Array[param.int64]`,
			Params: m,
			Wanted: "rulengine",
		},
		{
			Name:   "Selector Bracket",
			Input:  `param.Array[1]`,
			Params: m,
			Wanted: "rulengine",
		},
		{
			Name:   "Selector Bracket",
			Input:  `param.Map["key_str"]`,
			Params: m,
			Wanted: "val_string",
		},
		{
			Name:   "Nested Selector",
			Input:  `param.Array.1.1 == 'u'`,
			Params: m,
			Wanted: true,
		},
		{
			Name:   "Nested Selector",
			Input:  "param.Array.1.1",
			Params: m,
			Wanted: 'u',
		},
		{
			Name:   "Array Selector",
			Input:  "param.Array.0.bool",
			Params: m,
			Wanted: true,
		},
		{
			Name:   "Nested Selector",
			Input:  "param.Map.key_int",
			Params: m,
			Wanted: 1.0,
		},
		{
			Name:   "Simple Selector",
			Input:  "param.string",
			Params: m,
			Wanted: "string_test1",
		},
		{
			Name:   "Simple Selector",
			Input:  "param.Array.0.Interface",
			Params: m,
			Wanted: nil,
		},
	}

	runParseAstTests(parseAstTests, t)
}
func runParseAstTests(tests []ParseAstTest, test *testing.T) {

	var expr *Expr
	var res interface{}
	var err error

	fmt.Printf("Running %d eval ast test cases...\n", len(tests))

	// Run the test cases.
	for _, t := range tests {
		expr, err = NewExpr(t.Input)

		if err != nil {

			test.Logf("Test '%s' with input %s failed to eval: '%s'", t.Name, t.Input, err)
			test.Fail()
			continue
		}

		res, err = expr.Eval(t.Params)

		if err != nil {

			test.Logf("Test '%s' with input %s failed", t.Name, t.Input)
			test.Logf("Encountered error: %s", err.Error())
			test.Fail()
			continue
		}

		if res != t.Wanted {

			test.Logf("Test '%s' with input %s failed", t.Name, t.Input)
			test.Logf("Eval result '%v' does not match wanted: '%v'", res, t.Wanted)
			test.Fail()
		}
	}
}
