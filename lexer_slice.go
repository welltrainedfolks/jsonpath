package jsonpath

import (
	"errors"
	"fmt"
)

type SliceLexer struct {
	lex
	input []byte // the []byte being scanned.
	start Pos    // start position of this Item.
	pos   Pos    // current position in the input
}

func NewSliceLexer(input []byte, initial stateFn) *SliceLexer {
	l := &SliceLexer{
		lex:   newLex(initial),
		input: input,
	}

	return l
}

func (l *SliceLexer) take() int {
	if int(l.pos) >= len(l.input) {
		return eof
	}

	r := int(l.input[l.pos])
	l.pos++

	return r
}

func (l *SliceLexer) takeString() error {
	curPos := l.pos
	inputLen := len(l.input)

	if int(curPos) >= inputLen {
		return errors.New("end of file where string expected")
	}

	cur := int(l.input[curPos])
	curPos++

	if cur != '"' {
		l.pos = curPos
		return fmt.Errorf("expected \" as start of string instead of %#U", cur)
	}

	var previous int
looper:
	for {
		if int(curPos) >= inputLen {
			l.pos = curPos
			return errors.New("end of file where string expected")
		}
		cur := int(l.input[curPos])
		curPos++
		if cur == '"' {
			if previous == noValue || previous != '\\' {
				break looper
			} else {
				l.take()
			}
		}

		previous = cur
	}

	l.pos = curPos

	return nil
}

func (l *SliceLexer) peek() int {
	if int(l.pos) >= len(l.input) {
		return eof
	}

	return int(l.input[l.pos])
}

func (l *SliceLexer) emit(t int) {
	l.setItem(t, l.start, l.input[l.start:l.pos])

	l.hasItem = true

	// Ignore whitespace after this token
	for int(l.pos) < len(l.input) {
		r := l.input[l.pos]
		if r == ' ' || r == '\t' || r == '\r' || r == '\n' {
			l.pos++
		} else {
			break
		}
	}

	l.start = l.pos
}

func (l *SliceLexer) setItem(typ int, pos Pos, val []byte) {
	l.item.typ = typ
	l.item.pos = pos
	l.item.val = val
}

func (l *SliceLexer) ignore() {
	l.start = l.pos
}

func (l *SliceLexer) next() (*Item, bool) {
	for {
		if l.currentStateFn == nil {
			break
		}

		l.currentStateFn = l.currentStateFn(l, &l.stack)

		if l.hasItem {
			l.hasItem = false
			return &l.item, true
		}
	}

	return &l.item, false
}

func (l *SliceLexer) errorf(format string, args ...interface{}) stateFn {
	l.setItem(lexError, l.start, []byte(fmt.Sprintf(format, args...)))

	l.start = l.pos
	l.hasItem = true

	return nil
}

func (l *SliceLexer) reset() {
	l.start = 0
	l.pos = 0
	l.lex = newLex(l.initialState)
}
