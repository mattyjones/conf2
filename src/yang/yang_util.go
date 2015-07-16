package yang
import (
	"strings"
)

func FindByIdent(i DefIterator, ident string) Def {
	child := i.NextDef()
	for child != nil {
		if child.GetIdent() == ident {
			return child
		}
		child = i.NextDef()
	}
	return nil
}

func FindByPathWithoutResolvingProxies(root DefList, path string) Def {
	c := find(root, path, false)
	return c
}

func FindByPath(root DefList, path string) Def {
	return find(root, path, true)
}

func find(root DefList, path string, resolveProxies bool) (def Def) {
	elems := strings.SplitN(path, "/", -1)
	lastLevel := len(elems) - 1
	var ok bool
	list := root
	i := NewDefListIterator(list, resolveProxies)
	for level, elem := range elems {
		def = FindByIdent(i, elem)
		if def == nil {
			return nil
		}
		if level < lastLevel {
			if list, ok = def.(DefList); ok {
				i = NewDefListIterator(list, resolveProxies)
			} else {
				return nil
			}
		}
	}
	return
}


type yangError struct {
	s string
}

func (err *yangError) Error() string {
	return err.s
}