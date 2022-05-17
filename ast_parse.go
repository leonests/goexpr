package rulengine

import "fmt"

func parseAST(tokens []LexerToken) (*astNode, error) {
	stream := newLexerStream(tokens)

	ast, err := parseAst(stream)
	if err != nil {
		return nil, err
	}

	adjustAst(ast)

	return ast, nil
}

func parseAst(stream *lexerStream) (*astNode, error) {
	if !stream.notEOF() {
		return nil, nil
	}
	priority := stream.getLowestPriority()
	switch priority {
	case priorityLITERAL, priorityPREFIX, priorityCLAUSE:
		return parseSelectorAndVariable(stream)
	case priorityMUL:
		return parseMul(stream)
	case priorityADD:
		return parseAdd(stream)
	case priorityBITSHIFT:
		return parseBitShift(stream)
	case priorityCOMPARATOR:
		return parseComparator(stream)
	case priorityBIT:
		return parseBit(stream)
	case priorityLAND:
		return parseLogicalAnd(stream)
	case priorityLOR:
		return parseLogicalOr(stream)
	default:
		return parseTernary(stream)
	}
}

type parser func(stream *lexerStream) (*astNode, error)

var (
	parsePrefix     parser
	parseMul        parser
	parseAdd        parser
	parseBit        parser
	parseBitShift   parser
	parseComparator parser
	parseLogicalAnd parser
	parseLogicalOr  parser
	parseTernary    parser
)

type parserPkg struct {
	validToken   map[TokenType]struct{}
	nextPriority parser // higher priority
	rightParser  parser // right part
	errFormat    string
}

func init() {
	// different operators have different priorities
	// ref: https://en.cppreference.com/w/c/language/operator_precedence
	parsePrefix = buildParserWithPkg(&parserPkg{
		validToken:  tokenPREFIX,
		rightParser: parseSelectorAndVariable,
		errFormat:   errPrefixFormat,
	})
	parseMul = buildParserWithPkg(&parserPkg{
		validToken:   tokenMUL,
		nextPriority: parseSelectorAndVariable,
		errFormat:    errNumericFormat,
	})
	parseAdd = buildParserWithPkg(&parserPkg{
		validToken:   tokenADD,
		nextPriority: parseMul,
		errFormat:    errNumericFormat,
	})
	parseBitShift = buildParserWithPkg(&parserPkg{
		validToken:   tokenBITSHIFT,
		nextPriority: parseAdd,
		errFormat:    errNumericFormat,
	})
	parseComparator = buildParserWithPkg(&parserPkg{
		validToken:   tokenCOMPARATOR,
		nextPriority: parseBitShift,
		errFormat:    errComparatorFormat,
	})
	parseBit = buildParserWithPkg(&parserPkg{
		validToken:   tokenBIT,
		nextPriority: parseComparator,
		errFormat:    errNumericFormat,
	})
	parseLogicalAnd = buildParserWithPkg(&parserPkg{
		validToken:   map[TokenType]struct{}{LAND: {}},
		nextPriority: parseBit,
		errFormat:    errLogicalFormat,
	})
	parseLogicalOr = buildParserWithPkg(&parserPkg{
		validToken:   map[TokenType]struct{}{LOR: {}},
		nextPriority: parseLogicalAnd,
		errFormat:    errLogicalFormat,
	})
	parseTernary = buildParserWithPkg(&parserPkg{
		validToken:   tokenTERNARY,
		nextPriority: parseLogicalOr,
		errFormat:    errTernaryFormat,
	})
}

func buildParserWithPkg(pkg *parserPkg) parser {
	var (
		parsed    parser
		nextRight parser
	)
	parsed = func(stream *lexerStream) (*astNode, error) {
		return parseAstNode(stream, pkg.validToken, pkg.nextPriority, nextRight, pkg.errFormat)
	}

	if pkg.rightParser != nil {
		nextRight = pkg.rightParser
	} else {
		nextRight = parsed
	}
	return parsed
}

func parseAstNode(stream *lexerStream, validToken map[TokenType]struct{}, leftParser, rightParser parser, errFormat string) (*astNode, error) {
	var (
		token LexerToken
		op    TokenType
		left  *astNode
		right *astNode
		check typeChecks
		err   error
	)
	if leftParser != nil {
		left, err = leftParser(stream)
		if err != nil {
			return nil, err
		}
	}
	for stream.notEOF() {
		token = stream.flowForward()

		// check if it is a valid operator
		if validToken != nil {
			if _, ok := validToken[token.Type]; !ok {
				break
			} else {
				op = token.Type
			}
		}

		if rightParser != nil {
			right, err = rightParser(stream)
			if err != nil {
				return nil, err
			}
		}

		check = getTypeChecks(op)

		return &astNode{
			operator:   op,
			left:       left,
			right:      right,
			leftCheck:  check.left,
			rightCheck: check.right,
			bothCheck:  check.both,
			calculator: opCalculator[op],
			err:        errFormat,
		}, nil
	}
	stream.flowBackward()
	return left, nil
}

func parseSelectorAndVariable(stream *lexerStream) (*astNode, error) {
	var (
		rightNode *astNode
		nextToken LexerToken
		err       error
		rightList []*astNode
		cal       calculator
	)

	if !stream.notEOF() {
		return nil, nil
	}
	token := stream.flowForward()
	if token.Type != VARIABLE && token.Type != SELECTOR && token.Type != ACCESSOR {
		stream.flowBackward()
		return parseValue(stream)
	}

	if (token.Type == SELECTOR || token.Type == ACCESSOR) && stream.notEOF() {
		if stream.flowForward().Type == LPAREN {
			rightNode, err = parseAst(stream)
			if err != nil {
				return nil, err
			}
			return &astNode{
				operator:   LITERAL,
				right:      rightNode,
				calculator: calculatorVARIABLE(token.Value.(string)),
				err:        errSelectorFormat,
			}, nil
		} else {
			stream.flowBackward()
		}
	}

	rightList = make([]*astNode, 0)
	for stream.notEOF() {
		nextToken = stream.flowForward()
		if nextToken.Type == LBRACKET {
			stream.flowBackward()
			rightNode, err = parseAst(stream)
			if err != nil {
				return nil, err
			}
			rightList = append(rightList, rightNode)
		} else if nextToken.Type == SELECTOR {
			rightList = append(rightList, buildSelectorNode(nextToken))
		} else if nextToken.Type == ACCESSOR {
			rightList = append(rightList, buildAccessorNode(nextToken))
		} else {
			stream.flowBackward()
			break
		}
	}

	rightNode, rightList = resetRightAndRightList(rightNode, rightList)

	if token.Type == SELECTOR {
		cal = calculatorSELECTOR(token.Value.([]string))
	} else if token.Type == VARIABLE {
		cal = calculatorVARIABLE(token.Value.(string))
	}

	return &astNode{
		operator:   LITERAL,
		right:      rightNode,
		calculator: cal,
		err:        errSelectorFormat,
		rightList:  rightList,
	}, nil
}

func parseValue(stream *lexerStream) (*astNode, error) {
	var (
		cal calculator
		op  TokenType
	)
	if !stream.notEOF() {
		return nil, nil
	}

	token := stream.flowForward()

	switch token.Type {
	case LPAREN, LBRACKET:
		node, err := parseAst(stream)
		if err != nil {
			return nil, err
		}

		stream.flowForward() // jump over the RPAREN
		node = &astNode{
			operator:   CLAUSE,
			right:      node,
			calculator: calculatorCLAUSE,
		}
		return node, nil
	case NEG, NOT:
		stream.flowBackward()
		return parsePrefix(stream)
	case NUMBER, STRING, CHAR, BOOL:
		op = LITERAL
		cal = calculatorLITERAL(token.Value)
	}
	if cal == nil {
		return nil, fmt.Errorf("unable to deal with token type: %s, value: %v", token.Type.String(), token.Value)
	}
	return &astNode{
		operator:   op,
		calculator: cal,
	}, nil
}

func resetRightAndRightList(right *astNode, rightList []*astNode) (*astNode, []*astNode) {
	if rightList == nil {
		return right, rightList
	}
	switch {
	case len(rightList) > 1:
		right = nil
	case len(rightList) == 1:
		right = rightList[0]
		rightList = nil
	case len(rightList) == 0:
		rightList = nil
	}
	return right, rightList
}

// need to adjust operator seq, because nodes with equal priority are parsed into a reverse order to evaluate
// this will get the wrong result if using the original AST
func adjustAst(root *astNode) {
	var (
		samePriority []*astNode
		tmp          *astNode
		priority     opPriority
	)
	tmp = root
	priority = root.operator.Priority() //root priority

	for tmp != nil {
		if tmp.left != nil {
			adjustAst(tmp.left)
		}

		if tmp.operator.Priority() != priority {
			// if tmp has different priorities, swap the same priority nodes
			if len(samePriority) > 1 {
				swapTrees(samePriority)
			}
			// reset priority and same priority list
			samePriority = []*astNode{tmp}
			priority = tmp.operator.Priority()
			tmp = tmp.right
			continue
		}
		// if tmp has the same priority, add into list and keep searching
		samePriority = append(samePriority, tmp)
		tmp = tmp.right
	}
	// if same priority nodes are still more than 1, deal with it
	if len(samePriority) > 1 {
		swapTrees(samePriority)
	}
}
func swapTrees(nodes []*astNode) {
	var tmp *astNode
	length := len(nodes)
	// reverse left and right of all the nodes
	for _, node := range nodes {
		tmp = node.right
		node.right = node.left
		node.left = tmp
	}

	// end left swap with root right
	tmp = nodes[0].right
	nodes[0].right = nodes[length-1].left
	nodes[length-1].left = tmp

	// swap right tree of counterparts in the list except the root
	for i := 0; i < length/2; i++ {
		tmp = nodes[i+1].right
		nodes[i+1].right = nodes[length-i-1].right
		nodes[length-i-1].right = tmp
	}

	// swap all other information of all counterpart nodes
	for i := 0; i < length/2; i++ {
		swap(nodes[i], nodes[length-i-1])
	}
}

func swap(x *astNode, y *astNode) {
	tmp := *y
	y.operator = x.operator
	y.calculator = x.calculator
	y.leftCheck = x.leftCheck
	y.rightCheck = x.rightCheck
	y.bothCheck = x.bothCheck
	y.err = x.err

	x.operator = tmp.operator
	x.calculator = tmp.calculator
	x.leftCheck = tmp.leftCheck
	x.rightCheck = tmp.rightCheck
	x.bothCheck = tmp.bothCheck
	x.err = tmp.err
}
