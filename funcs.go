package goexpr

// ExprFunc represents a function that can be called within an expression
type ExprFunc func(args ...interface{})(interface{}, error)