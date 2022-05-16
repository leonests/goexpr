package rulengine

import (
	"fmt"
	"math"
	"reflect"
)

// bool to interface, predefined to avoid cost
// since boolean to interface needs allocation
var (
	_true  = interface{}(true)
	_false = interface{}(false)
)

func calculatorEQ(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	return convertBool2Interface(reflect.DeepEqual(left, right)), nil
}
func calculatorNEQ(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	return convertBool2Interface(!reflect.DeepEqual(left, right)), nil
}
func calculatorGT(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	if isString(left) && isString(right) {
		return convertBool2Interface(left.(string) > right.(string)), nil
	}
	return convertBool2Interface(left.(float64) > right.(float64)), nil
}
func calculatorGEQ(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	if isString(left) && isString(right) {
		return convertBool2Interface(left.(string) >= right.(string)), nil
	}
	return convertBool2Interface(left.(float64) >= right.(float64)), nil
}
func calculatorLT(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	if isString(left) && isString(right) {
		return convertBool2Interface(left.(string) < right.(string)), nil
	}
	return convertBool2Interface(left.(float64) < right.(float64)), nil
}
func calculatorLEQ(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	if isString(left) && isString(right) {
		return convertBool2Interface(left.(string) <= right.(string)), nil
	}
	return convertBool2Interface(left.(float64) <= right.(float64)), nil
}
func calculatorADD(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	if isString(left) || isString(right) {
		return fmt.Sprintf("%v%v", left, right), nil
	}
	return left.(float64) + right.(float64), nil
}
func calculatorSUB(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	return left.(float64) - right.(float64), nil
}
func calculatorMUL(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	return left.(float64) * right.(float64), nil
}
func calculatorQUO(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	return left.(float64) / right.(float64), nil
}
func calculatorREM(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	return math.Mod(left.(float64), right.(float64)), nil
}
func calculatorNEG(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	return -right.(float64), nil
}
func calculatorNOT(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	return convertBool2Interface(!right.(bool)), nil
}
func calculatorLAND(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	return convertBool2Interface(left.(bool) && right.(bool)), nil
}
func calculatorLOR(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	return convertBool2Interface(left.(bool) || right.(bool)), nil
}
func calculatorAND(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	return float64(int64(left.(float64)) & int64(right.(float64))), nil
}
func calculatorOR(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	return float64(int64(left.(float64)) | int64(right.(float64))), nil
}
func calculatorXOR(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	return float64(int64(left.(float64)) ^ int64(right.(float64))), nil
}
func calculatorSHL(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	return float64(int64(left.(float64)) << int64(right.(float64))), nil
}
func calculatorSHR(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	return float64(int64(left.(float64)) >> int64(right.(float64))), nil
}
func calculatorTernaryIf(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	if left.(bool) { // left ? right : ...
		return right, nil
	}
	return nil, nil
}
func calculatorTernaryElse(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	if left != nil { // ... ? left : right
		return left, nil
	}
	return right, nil
}
func calculatorCLAUSE(left, right interface{}, params map[string]interface{}) (interface{}, error) {
	return right, nil
}
func calculatorVARIABLE(paramName string) calculator {
	return func(left, right interface{}, params map[string]interface{}) (interface{}, error) {
		path, err := buildPathFromRight(right, []string{paramName})
		if err != nil {
			return nil, err
		}
		value, err := extractValueFromParams(params, path)
		if err != nil {
			return nil, err
		}
		return value, nil
	}
}
func calculatorLITERAL(literal interface{}) calculator {
	return func(left, right interface{}, params map[string]interface{}) (interface{}, error) {
		return literal, nil
	}
}
func calculatorACCESSOR(parts []string) calculator {
	return func(left, right interface{}, params map[string]interface{}) (interface{}, error) {
		return parts, nil
	}
}
func calculatorSELECTOR(parts []string) calculator {
	return func(left, right interface{}, params map[string]interface{}) (res interface{}, err error) {
		parts, err = buildPathFromRight(right, parts)
		if err != nil {
			return nil, err
		}

		value, err := extractValueFromParams(params, parts)
		if err != nil {
			return nil, err
		}
		return value, nil
	}
}

func isString(value interface{}) bool {
	switch value.(type) {
	case string:
		return true
	}
	return false
}

func isFloat64(value interface{}) bool {
	switch value.(type) {
	case float64:
		return true
	}
	return false
}

func isBool(value interface{}) bool {
	switch value.(type) {
	case bool:
		return true
	}
	return false
}

func convertBool2Interface(b bool) interface{} {
	if b {
		return _true
	}
	return _false
}