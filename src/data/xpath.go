package data

import (
	"strings"
	"conf2"
	"schema"
)

type XPath struct {
	AutoCreate bool
}

func (x XPath) IsFwdSlash(r rune) bool {
	return r == '/'
}

func (x XPath) Get(cwd *Selection, xpath string) (interface{}, error) {
	var err error
	sel, prop, err := x.SelectProperty(cwd, xpath)
	if err != nil {
		return nil, err
	}
	if sel == nil || prop == nil {
		return nil, err
	}
	if schema.IsLeaf(prop) {
		v, err := sel.Node.Read(sel, prop.(schema.HasDataType))
		if err != nil {
			return nil, err
		}
		return v.Value(), nil
	} else {
		return sel.Peek(), nil
	}
}

func (x XPath) SelectProperty(cwd *Selection, xpath string) (sel *Selection, meta schema.Meta, err error) {
	slash := strings.LastIndexFunc(xpath, x.IsFwdSlash)
	sel = cwd
	ident := xpath
	if slash > 0 {
		if sel, err = x.Select(cwd, xpath[:slash]); err != nil {
			return
		}
		ident = xpath[slash + 1:]
	}
	meta = schema.FindByIdent2(sel.State.SelectedMeta(), ident)
	return
}

func (x XPath) Select(cwd *Selection, xpath string) (*Selection, error) {
	if strings.HasPrefix(xpath, "../") {
		if cwd.parent != nil {
			return x.Select(cwd.parent, xpath[3:])
		} else {
			return nil, conf2.NewErrC("No parent path to resolve " + xpath, conf2.NotFound)
		}
	}

	p, err := ParsePath(xpath, cwd.State.SelectedMeta())
	if err != nil {
		return nil, err
	}
	return WalkPath(cwd, p)
}
