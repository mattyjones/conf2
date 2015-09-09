package yang

import (
	"io/ioutil"
	"fmt"
)

type ImportModule func (into *Module, name string) (e error)

func LoadModuleFromByteArray(data []byte, importer ImportModule) (*Module, error) {
	l := lex(string(data), importer)
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

func moduleCopy(dest *Module, src *Module) {
	iters := []MetaIterator {
		NewMetaListIterator(src.GetGroupings(), false),
		NewMetaListIterator(src.GetTypedefs(), false),
		NewMetaListIterator(src.DataDefs(), false),
		NewMetaListIterator(src.GetRpcs(), false),
		NewMetaListIterator(src.GetNotifications(), false),
	}
	for _, iter := range iters {
		for iter.HasNextMeta() {
			dest.AddMeta(iter.NextMeta())
		}
	}
}

func LoadModule(source StreamSource, yangfile string) (*Module, error) {
	importer := func(main *Module, submodName string) (suberr error) {
		var sub *Module
		// TODO: Performance - cache modules
		subFname := fmt.Sprint(submodName, ".yang")
		if sub, suberr = LoadModule(source, subFname); suberr != nil {
			return suberr
		}
		moduleCopy(main, sub)
		return nil
	}
	if res, err := source.OpenStream(yangfile); err != nil {
		return nil, err
	} else {
		defer CloseResource(res)
		if data, err := ioutil.ReadAll(res); err != nil {
			return nil, err
		} else {
			return LoadModuleFromByteArray(data, importer)
		}
	}
}
