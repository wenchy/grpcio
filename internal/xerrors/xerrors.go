package xerrors

import (
	"fmt"

	"github.com/wenchy/grpcio/internal/corepb"
)

type Error struct {
	code corepb.Ecode
	// msg  string
}

func (e *Error) Code() corepb.Ecode {
	return e.code
}

func (e *Error) Error() string {
	return fmt.Sprintf("error:%v", e.code)
}

func New(code corepb.Ecode) error {
	return &Error{
		code: code,
	}
}

func Is(err error, code corepb.Ecode) bool {
	if pberr, ok := err.(*Error); ok {
		return pberr.Code() == code
	} else {
		return false
	}
}
