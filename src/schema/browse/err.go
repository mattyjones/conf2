package browse
import (
	"schema"
	"fmt"
)

type browseError struct {
	Code ResponseCode
	Msg string
}

func (err *browseError) Error() string {
	return err.Msg
}

type ResponseCode int
const (
	UNSPECIFIED ResponseCode = iota
	NOT_IMPLEMENTED
	NOT_FOUND
	MISSING_KEY
)

func EditNotImplemented(meta schema.Meta) error {
	return &browseError{Code:NOT_IMPLEMENTED, Msg:fmt.Sprintf("editing of \"%s\" not implemented", meta.GetIdent())}
}

func NotImplementedByName(ident string) error {
	return &browseError{Code:NOT_IMPLEMENTED, Msg:fmt.Sprintf("browsing of \"%s\" not implemented", ident)}
}

func NotImplemented(meta schema.Meta) error {
	return &browseError{Code:NOT_IMPLEMENTED, Msg:fmt.Sprintf("browsing of \"%s.%s\" not implemented",
		meta.GetParent().GetIdent(), meta.GetIdent())}
}

func NotFound(key string) error {
	return &browseError{Code:NOT_IMPLEMENTED, Msg:fmt.Sprintf("item identified with key \"%s\" not found", key)}
}

func ListKeyRequired() error {
	return &browseError{Code:MISSING_KEY, Msg:fmt.Sprintf("List key required")}
}

