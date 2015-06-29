package yang
import (
	"io/ioutil"
	"fmt"
)

func LoadModule(yangfile string) (*Module, error) {
	data, err := ioutil.ReadFile(yangfile)
	if err == nil {
		l := lex(string(data))
		err_code := yyParse(l)
		if err_code == 0 {
			d := l.stack.Peek()
			return d.(*Module), nil
		}
		return nil, yangError{fmt.Sprintf("Error %d loading yang file", err_code)}
	}
	return nil, err
}