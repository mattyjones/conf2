package conf2

import (
	"net/http"
	"runtime"
	"bufio"
	"strconv"
	"bytes"
)

type Code int
const (
	NotFound = 404
	Conflict = http.StatusConflict
	InternalErr = http.StatusInternalServerError
	NotImplemented = http.StatusNotImplemented
	UserErr = http.StatusBadRequest
)

type Error interface {
	error
	Code() Code
	Stack() string
}

type codedErrorString struct {
	s string
	c Code
	stack string
}

func (e *codedErrorString) Error() string {
	return e.s
}

func (e *codedErrorString) Code() Code {
	return e.c
}

func (e *codedErrorString) Stack() string {
	return e.stack
}

func NewErr(msg string) error {
	return &codedErrorString{
		s : msg,
		c : InternalErr,
		stack: dumpStack(),
	}
}

func NewErrC(msg string, code Code) error {
	return &codedErrorString{
		s : msg,
		c : code,
		stack: dumpStack(),
	}
}

func trim(s string, max int) string {
	if len(s) > max {
		return "..." + s[len(s) - (max + 3):]
	}
	return s
}

func dumpStack() string {
	var buff bytes.Buffer
	w := bufio.NewWriter(&buff)
	var stack [25]uintptr
	len := runtime.Callers(2, stack[:])
	for i := 1; i < len; i++ {
		f := runtime.FuncForPC(stack[i])
		w.WriteRune(' ')
		w.WriteString(f.Name())
		w.WriteRune(' ')
		file, lineno := f.FileLine(stack[i - 1])
		w.WriteString(trim(file, 20))
		w.WriteRune(':')
		w.WriteString(strconv.Itoa(lineno))
		w.WriteString("\n")
	}
	w.Flush()
	return buff.String()
}