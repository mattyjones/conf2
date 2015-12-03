package data

import (
	"fmt"
	"net/http"
	"schema"
)

type HttpError interface {
	error
	HttpCode() int
}

type browseError struct {
	Code int
	Msg  string
}

func (err *browseError) Error() string {
	return err.Msg
}

func (err *browseError) HttpCode() int {
	return err.Code
}

func EditNotImplemented(meta schema.Meta) error {
	return &browseError{Code: http.StatusNotImplemented, Msg: fmt.Sprintf("editing of \"%s\" not implemented", meta.GetIdent())}
}

func NotImplementedByName(ident string) error {
	return &browseError{Code: http.StatusNotImplemented, Msg: fmt.Sprintf("browsing of \"%s\" not implemented", ident)}
}

func NotImplemented(meta schema.Meta) error {
	panic("STOP")
	return &browseError{Code: http.StatusNotImplemented, Msg: fmt.Sprintf("browsing of \"%s.%s\" not implemented",
		meta.GetParent().GetIdent(), meta.GetIdent())}
}

func PathNotFound(path string) error {
	return &browseError{Code: http.StatusNotFound, Msg: fmt.Sprintf("item identified with path \"%s\" not found", path)}
}

func ListItemNotFound(key string) error {
	return &browseError{Code: http.StatusNotFound, Msg: fmt.Sprintf("item identified with key \"%s\" not found", key)}
}

//func ListKeyRequired() error {
//	return &browseError{Code:http.StatusNotImplemented, Msg:fmt.Sprintf("List key required")}
//}
