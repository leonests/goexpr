[![Build](https://github.com/leonests/rulengine/workflows/CI/badge.svg)](https://github.com/leonests/rulengine/actions?query=workflow)
[![Coverage](https://codecov.io/gh/leonests/rulengine/branch/main/graphs/badge.svg?branch=main)](https://codecov.io/gh/leonests/rulengine)
[![Go Report](https://goreportcard.com/badge/github.com/leonests/rulengine)](https://goreportcard.com/report/github.com/leonests/rulengine)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
------
# Goexpr

Goexpr provides support for evaluating expressions with parameters, arimethetic, logical, and string operations.

* basic expression: 1 > 0
* parameterized expression: x > 0
* nested parameterized expression: a.b > 0
* arithmetic expression: (x * y / 100) >= 50
* string expression: real == "expected"
* float64 expression: (part / total) * 100

## Installation
When used with Go modules, use the following import path:

    go get github.com/leonests/goexpr

## Quickstart

**Example 1: Simple Usage Without Parameters**
```go
	expr, err := goexpr.NewExpr("1 > 0")
	result, err := expr.Eval(nil)
	// result is true.
```

**Example 2: Simple Usage With Parameters**
```go
    param := map[string]interface{}{ "x": 100, "y": 50}
	expr, err := goexpr.NewExpr(`(x * y / 100) >= 50`)
	result, err := expr.Eval(parametes)
	// result is true.
```

## Advanced

### Bracket Accessor

### Dot Accessor

### Method