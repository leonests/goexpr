package goexpr

type lexerStream struct {
	tokens []LexerToken
	pos    int
	len    int
}

func newLexerStream(tokens []LexerToken) *lexerStream {
	return &lexerStream{
		tokens: tokens,
		pos:    0,
		len:    len(tokens),
	}
}

func (rs *lexerStream) flowForward() LexerToken {
	token := rs.tokens[rs.pos]
	rs.pos += 1
	return token
}

func (rs *lexerStream) flowBackward() {
	rs.pos -= 1
}

func (rs *lexerStream) notEOF() bool {
	return rs.pos < rs.len
}

func (rs *lexerStream) getLowestPriority() opPriority {
	p := priorityLITERAL
	for i := rs.pos; i < rs.len; i++ {
		newP := rs.tokens[i].Type.Priority()
		if newP < p {
			p = newP
		}
	}
	return p
}
