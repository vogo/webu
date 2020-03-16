//author: wongoo
//date: 20190826

package cerror

import "net/http"

const (
	CodeOK              = 20000
	CodeServerErr       = 50000
	CodeBadErr          = 40000
	CodeNotFoundErr     = 40400
	CodeUnauthorizedErr = 40100
	CodeForbiddenErr    = 40300
)

var (
	ErrNotFound     = NewStatusCodeError(http.StatusNotFound, CodeNotFoundErr, "not found")
	ErrBadRequest   = NewStatusCodeError(http.StatusBadRequest, CodeBadErr, "forbidden")
	ErrArgRequired  = NewStatusCodeError(http.StatusBadRequest, CodeBadErr+1, "arg required")
	ErrValueInvalid = NewStatusCodeError(http.StatusBadRequest, CodeBadErr+2, "value invalid")
	ErrUnauthorized = NewStatusCodeError(http.StatusUnauthorized, CodeUnauthorizedErr, "unauthorized")
	ErrForbidden    = NewStatusCodeError(http.StatusForbidden, CodeForbiddenErr, "forbidden")
)

type Coder interface {
	Code() int
}

type StatusState interface {
	Status() int
}

type CodeError interface {
	error
	Coder
}
type StatusCodeError interface {
	CodeError
	StatusState
}

func NewCodeError(code int, err string) CodeError {
	return &codeError{c: code, m: err}
}

type codeError struct {
	c int
	m string
}

func (e *codeError) Code() int {
	return e.c
}
func (e *codeError) Error() string {
	return e.m
}

func NewStatusCodeError(status, code int, err string) CodeError {
	return &statusCodeError{s: status, c: code, m: err}
}

type statusCodeError struct {
	s int
	c int
	m string
}

func (e *statusCodeError) Code() int {
	return e.c
}

func (e *statusCodeError) Status() int {
	return e.s
}

func (e *statusCodeError) Error() string {
	return e.m
}
