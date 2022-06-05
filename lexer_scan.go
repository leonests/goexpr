package goexpr

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type LexerToken struct {
	Type  TokenType
	Value interface{}
}

func lexerScan(expr string) (tokens []LexerToken, err error) {
	var (
		token LexerToken
		exist bool
	)

	stream := newRuneStream(expr)
	tokenRule := lexerRules[ILLEGAL]

	for stream.notEOF() {
		token, exist, err = tokenScan(stream, tokenRule)
		if err != nil {
			return tokens, err
		}

		if !exist {
			break
		}

		tokenRule, err = getLexerRule(token.Type)
		if err != nil {
			return tokens, err
		}

		tokens = append(tokens, token)
	}
	err = checkLexerBalance(tokens)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func tokenScan(stream *runeStream, rule lexerRule) (LexerToken, bool, error) {
	var (
		char      rune
		tokenType TokenType
		tokenStr  string
		tokenVal  interface{}
		completed bool
		err       error
	)
	for stream.notEOF() {
		char = stream.flowForward()
		if unicode.IsSpace(char) {
			continue
		}

		tokenType = ILLEGAL
		if unicode.IsDigit(char) {
			tokenStr = readWithCond(stream, isNumeric)
			tokenVal, err = strconv.ParseFloat(tokenStr, 64)
			if err != nil {
				return LexerToken{}, false, fmt.Errorf("unable to parse numeric value '%v' to float64", tokenStr)
			}
			tokenType = NUMBER
			break
		}

		if unicode.IsLetter(char) {
			tokenStr = readWithCond(stream, isVariable)

			tokenVal = tokenStr
			tokenType = VARIABLE

			if strings.ToUpper(tokenStr) == "TRUE" {
				tokenType = BOOL
				tokenVal = true
			} else if strings.ToUpper(tokenStr) == "FALSE" {
				tokenType = BOOL
				tokenVal = false
			}

			if strings.Contains(tokenStr, ".") {
				//can not be the last one
				if tokenStr[len(tokenStr)-1] == '.' {
					return LexerToken{}, false, fmt.Errorf("selector at tail of token %v", tokenStr)
				}
				tokenType = SELECTOR
				tokenVal = strings.Split(tokenStr, ".")
			}
			break
		}
		// start with .
		if isDot(char) {
			tokenStr = readWithCond(stream, isVariable)
			if tokenStr[len(tokenStr)-1] == '.' {
				return LexerToken{}, false, fmt.Errorf("accessor at tail of token %v", tokenStr)
			}
			tokenType = ACCESSOR
			tokenVal = strings.Split(tokenStr, ".")[1:]
			break
		}

		if isDoubleQuote(char) {
			tokenStr, completed = readWithFlagAndCond(stream, false, true, isNotDoubleQuote)
			if !completed {
				return LexerToken{}, false, fmt.Errorf("literal string unclosed")
			}

			stream.flowBackward(-1) //jump over "
			tokenVal = tokenStr
			tokenType = STRING
			break
		}

		if isSingleQuote(char) {
			tokenVal = stream.flowForward()
			tokenType = CHAR
			//jump over '
			if stream.flowForward() != '\'' {
				return LexerToken{}, false, fmt.Errorf("more than 1 charactor for char type")
			}
			break
		}

		if char == '(' {
			tokenVal = char
			tokenType = LPAREN
			break
		}

		if char == ')' {
			tokenVal = char
			tokenType = RPAREN
			break
		}

		if char == '[' {
			tokenVal = char
			tokenType = LBRACKET
			break
		}
		if char == ']' {
			tokenVal = char
			tokenType = RBRACKET
			break
		}

		//then it must be an operator
		tokenStr = readWithCond(stream, isNotAlphanumeric)
		tokenVal = tokenStr

		//'-' may be PREFIX or SUB, determined by current rule
		if rule.hasNextAllowable(NEG) {
			if tokenStr == "-" {
				tokenType = NEG
				break
			}
		}

		if tok, ok := tokenMap[tokenStr]; ok {
			tokenType = tok
			break
		}
		return LexerToken{}, false, fmt.Errorf("invalid token %v", tokenStr)
	}
	res := LexerToken{
		Type:  tokenType,
		Value: tokenVal,
	}

	return res, tokenType != ILLEGAL, nil
}

func isNumeric(char rune) bool {
	return unicode.IsDigit(char) || char == '.'
}

// a_b1.c2_d3
func isVariable(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsDigit(char) ||
		char == '_' || char == '.'
}

func isDoubleQuote(char rune) bool {
	return char == '"'
}

func isNotDoubleQuote(char rune) bool {
	return !isDoubleQuote(char)
}

func isSingleQuote(char rune) bool {
	return char == '\''
}

func isDot(char rune) bool {
	return char == '.'
}

func isNotAlphanumeric(char rune) bool {
	return !(unicode.IsDigit(char) || unicode.IsLetter(char) ||
		char == '(' || char == ')' || char == '[' || char == ']')
}

func readWithCond(stream *runeStream, cond func(rune) bool) string {
	stream.flowBackward(1)
	res, _ := readWithFlagAndCond(stream, true, false, cond)
	return res
}

// return string when the given condition was false, or broken with white space
// return false if stream ended before condition was met or broken with space
func readWithFlagAndCond(stream *runeStream, breakWhenSpace, includeWhenSpace bool, cond func(rune) bool) (res string, matched bool) {
	var (
		buffer bytes.Buffer
		char   rune
	)
	for stream.notEOF() {
		char = stream.flowForward()
		if unicode.IsSpace(char) {
			if breakWhenSpace && buffer.Len() > 0 {
				matched = true
				break
			}
			if !includeWhenSpace {
				continue
			}
		}
		if cond(char) {
			buffer.WriteString(string(char))
		} else {
			matched = true
			stream.flowBackward(1)
			break
		}
	}
	res = buffer.String()
	return
}
