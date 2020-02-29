package compile

import (
	"errors"
	"fmt"
)

const maxCells = 30000

type Compiler interface {
	Compile(src []byte) ([]byte, error)
}

type Error struct {
	Pos             int
	underlyingError error
}

func NewError(pos int, err error) *Error {
	return &Error{Pos: pos, underlyingError: err}
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (char: %d)", e.underlyingError.Error(), e.Pos)
}

func (e Error) Unwrap() error {
	return e.underlyingError
}

var ErrDPOutOfBounds = errors.New("data pointer will go out of bounds")
var ErrRParen = errors.New(`missing opening "[" before "]"`)
