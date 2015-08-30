package yang

import (
	"io/ioutil"
	"fmt"
)

func LoadModuleFromByteArray(data []byte) (*Module, error) {
	l := lex(string(data))
	err_code := yyParse(l)
	if err_code != 0 || l.lastError != nil {
		if l.lastError == nil {
			// Developer - Find out why there's no error
			msg := fmt.Sprint("Error parsing, code ", string(err_code))
			l.lastError = &yangError{msg}

		}
		return nil, l.lastError
	}

	d := l.stack.Peek()
	return d.(*Module), nil
}

func LoadModule(source StreamSource, yangfile string) (*Module, error) {
	if res, err := source.OpenStream(yangfile); err != nil {
		return nil, err
	} else {
		defer res.Close()
		if data, err := ioutil.ReadAll(res); err != nil {
			return nil, err
		} else {
			return LoadModuleFromByteArray(data)
		}
	}
}
