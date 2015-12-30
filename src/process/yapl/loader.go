package yapl
import (
	"process"
	"fmt"
	"errors"
)

func Load(data string) (map[string]*process.Script, error) {
	l := lex(data)
	err_code := yyParse(l)
	if err_code != 0 || l.lastError != nil {
		if l.lastError == nil {
			// Developer - Find out why there's no error
			msg := fmt.Sprint("Error parsing, code ", string(err_code))
			l.lastError = errors.New(msg)

		}
		return nil, l.lastError
	}

	return l.stack.scripts, nil
}
