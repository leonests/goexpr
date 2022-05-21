package rulengine

import "fmt"

const rightShortCircuit int = 0

var ternaryShortCircuit interface{}

type Expr struct {
	tokens  []LexerToken
	astNode *astNode
	input   string
}

func NewExpr(expr string) (res *Expr, err error) {
	res = &Expr{
		input: expr,
	}
	res.tokens, err = lexerScan(expr)
	if err != nil {
		return nil, err
	}
	res.astNode, err = parseAST(res.tokens)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (expr *Expr) Eval(params map[string]interface{}) (interface{}, error) {
	if expr.astNode == nil {
		return nil, nil
	}
	value, err := expr.eval(expr.astNode, params)
	ternaryShortCircuit = nil
	return value, err
}

func (expr *Expr) eval(node *astNode, params map[string]interface{}) (interface{}, error) {
	var (
		left, right interface{}
		rightList   []interface{}
		err         error
	)

	if node.left != nil {
		left, err = expr.eval(node.left, params)
		if err != nil {
			return nil, err
		}
	}
	if node.operator.isTernary() && ternaryShortCircuit != nil {
		return ternaryShortCircuit, nil
	}
	if !node.operator.isTernary() && ternaryShortCircuit != nil {
		ternaryShortCircuit = nil
	}

	if node.operator.isShortCircuit() {
		switch node.operator {
		case AND:
			if left == false {
				return false, nil
			}
		case OR:
			if left == true {
				return true, nil
			}
		case TERNARY_IF:
			if left == false {
				right = rightShortCircuit
			}
		case TERNARY_ELSE:
			if left != nil {
				ternaryShortCircuit = left
				right = rightShortCircuit
			}
		}
	}

	if right != rightShortCircuit {
		if node.right != nil {
			right, err = expr.eval(node.right, params)
			if err != nil {
				return nil, err
			}
		} else if node.rightList != nil {
			rightList = make([]interface{}, len(node.rightList))
			for i, r := range node.rightList {
				right, err = expr.eval(r, params)
				if err != nil {
					return nil, err
				}
				rightList[i] = right
			}
		}
	}

	if !node.operator.isTernary() && ternaryShortCircuit != nil {
		ternaryShortCircuit = nil
	}

	if err = typeCheck(node, left, right); err != nil {
		return nil, err
	}

	if rightList != nil {
		return node.calculator(left, rightList, params)
	}
	return node.calculator(left, right, params)
}

func typeCheck(node *astNode, left, right interface{}) error {
	if node.bothCheck == nil {
		if node.leftCheck != nil && !node.leftCheck(left) {
			return fmt.Errorf(node.err, left, node.operator)
		}
		if node.rightCheck != nil && !node.rightCheck(right) {
			return fmt.Errorf(node.err, right, node.operator)
		}
	} else {
		if !node.bothCheck(left, right) {
			return fmt.Errorf(node.err, left, node.operator)
		}
	}
	return nil
}
