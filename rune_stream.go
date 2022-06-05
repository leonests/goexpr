package goexpr

type runeStream struct {
	runes []rune
	pos   int
	len   int
}

func newRuneStream(expr string) *runeStream {
	var runes []rune
	for _, r := range expr {
		runes = append(runes, r)
	}
	return &runeStream{
		runes: runes,
		pos:   0,
		len:   len(runes),
	}
}

func (rs *runeStream) flowForward() rune {
	char := rs.runes[rs.pos]
	rs.pos += 1
	return char
}

func (rs *runeStream) flowBackward(step int) {
	rs.pos -= step
}

func (rs *runeStream) notEOF() bool {
	return rs.pos < rs.len
}
