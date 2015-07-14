package yang

import (
	"C"
	"io/ioutil"
	"fmt"
)

//export LoadModuleFromCByteArray
func LoadModuleFromCByteArray(cdata *C.char, len C.int) {
	// TODO: improve performance by not copying
	gdata := []byte(C.GoStringN(cdata, len))
	loadModuleFromByteArray(gdata)
}

func LoadModule(resolver ResourceResolver, yangfile string) (*Module, error) {
	data, err := resolver.LoadResource(yangfile)
	if err == nil {
		return loadModuleFromByteArray(data)
	}
	return nil, err
}

func loadModuleFromByteArray(data []byte) (*Module, error) {
	l := lex(string(data))
	err_code := yyParse(l)
	if err_code == 0 {
		d := l.stack.Peek()
		return d.(*Module), nil
	}
	return nil, &yangError{fmt.Sprintf("Error %d loading yang file", err_code)}
}

type ResourceResolver interface {
	LoadResource(resource string) ([]byte, error)
}

type FileResolver struct {
}

func (fs *FileResolver) LoadResource(fname string) ([]byte, error) {
	return ioutil.ReadFile(fname)
}
