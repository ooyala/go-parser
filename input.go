package parser

import (
	"io"
)

// Basic string Input for parsing over a string input.
type StringInput struct {
	state State
	input []rune
	txn   []Position
	pos Position
}

func NewStringInput(s, filename string) *StringInput {
	return &StringInput{
		input: []rune(s),
		txn:   make([]Position, 0, 8),
		pos: Position{
			Name:   filename,
			Line:   1,
			Column: 1,
			Offset: 0,
		},
	}
}

func (s *StringInput) Begin() {
	s.txn = append(s.txn, s.pos)
}

func (s *StringInput) End(rollback bool) {
	p := s.txn[len(s.txn)-1]
	s.txn = s.txn[:len(s.txn)-1]
	if rollback {
		s.pos = p
	}
}

func (s *StringInput) Get(i int) (string, error) {
	if len(s.input) < s.pos.Offset+i {
		return "", io.EOF
	}

	return string(s.input[s.pos.Offset : s.pos.Offset+i]), nil
}

func (s *StringInput) Next() (rune, error) {
	if len(s.input) < s.pos.Offset+1 {
		return 0, io.EOF
	}

	return s.input[s.pos.Offset], nil
}

func (s *StringInput) Pop(i int) {
	for j := 0; j < i; j++ {
		if s.input[s.pos.Offset + j] == '\n' {
			s.pos.Line++
			s.pos.Column = 1
		} else {
			s.pos.Column++
		}
	}
	s.pos.Offset += i
}

func (s *StringInput) Position() Position { return s.pos }
