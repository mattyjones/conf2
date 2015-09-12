package yang

import (
	"io/ioutil"
	"fmt"
	"schema"
)

type ImportModule func (into *schema.Module, name string) (e error)

func LoadModuleFromByteArray(data []byte, importer ImportModule) (*schema.Module, error) {
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
	return d.(*schema.Module), nil
}

func moduleCopy(dest *schema.Module, src *schema.Module) {
	iters := []schema.MetaIterator {
		schema.NewMetaListIterator(src.GetGroupings(), false),
		schema.NewMetaListIterator(src.GetTypedefs(), false),
		schema.NewMetaListIterator(src.DataDefs(), false),
		schema.NewMetaListIterator(src.GetRpcs(), false),
		schema.NewMetaListIterator(src.GetNotifications(), false),
	}
	for _, iter := range iters {
		for iter.HasNextMeta() {
			dest.AddMeta(iter.NextMeta())
		}
	}
}

func LoadModule(source schema.StreamSource, yangfile string) (*schema.Module, error) {
	importer := func(main *schema.Module, submodName string) (suberr error) {
		var sub *schema.Module
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
		defer schema.CloseResource(res)
		if data, err := ioutil.ReadAll(res); err != nil {
			return nil, err
		} else {
			return LoadModuleFromByteArray(data, importer)
		}
	}
}
