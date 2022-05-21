package rulengine

const (
	errNumericFormat  string = "value '%v' cannot be used with the numeric operator '%v', it is not a number"
	errLogicalFormat  string = "value '%v' cannot be used with the logical operator '%v', it is not a bool"
	errComparerFormat string = "value '%v' cannot be used with the COMPARER operator '%v', it is not a number"
	errTernaryFormat  string = "value '%v' cannot be used with the ternary operator '%v', it is not a bool"
	errPrefixFormat   string = "value '%v' cannot be used with the prefix operator '%v'"
	errSelectorFormat string = "fail to select parameter '%v'"
	errAccessorFormat string = "fail to access parameter '%v'"
)

type typeChecks struct {
	left  nodeTypeCheck
	right nodeTypeCheck
	both  bothTypeCheck
}

func getTypeChecks(op TokenType) typeChecks {
	switch op {
	case ADD:
		return typeChecks{
			both: addTypeCheck,
		}
	case SUB, MUL, QUO, REM:
		return typeChecks{
			left:  isFloat64,
			right: isFloat64,
		}
	case GT, LT, GEQ, LEQ:
		return typeChecks{
			both: comparerTypeCheck,
		}
	case AND, OR, XOR, SHL, SHR:
		return typeChecks{
			left:  isFloat64,
			right: isFloat64,
		}
	case LAND, LOR:
		return typeChecks{
			left:  isBool,
			right: isBool,
		}
	case NOT:
		return typeChecks{
			right: isBool,
		}
	case NEG:
		return typeChecks{
			right: isFloat64,
		}
	case TERNARY_IF:
		return typeChecks{
			left: isBool,
		}
	default:
		return typeChecks{}
	}
}

func addTypeCheck(left, right interface{}) bool {
	// both number
	if isFloat64(left) && isFloat64(right) {
		return true
	}
	// or either is string, string concat
	if isString(left) || isString(right) {
		return true
	}
	//TODO: support char(rune), type system need to improve
	return false
}

func comparerTypeCheck(left, right interface{}) bool {
	if isFloat64(left) && isFloat64(right) {
		return true
	}
	if isString(left) && isString(right) {
		return true
	}
	return false
}
