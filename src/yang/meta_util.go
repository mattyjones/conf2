package yang
import (
	"strings"
	"unicode"
)

func FindByIdent(i MetaIterator, ident string) Meta {
	child := i.NextMeta()
	for child != nil {
		if child.GetIdent() == ident {
			return child
		}
		child = i.NextMeta()
	}
	return nil
}

func MetaNameToFieldName(in string) string {
	// assumes fix is always shorter because char can be dropped and not added
	fixed := make([]rune, len(in))
	cap := true
	j := 0
	for _, r := range in {
		if r == '-' {
			cap = true
		} else {
			if cap {
				fixed[j] = unicode.ToUpper(r)
			} else {
				fixed[j] = r
			}
			j += 1
			cap = false
		}
	}
	return string(fixed[:j])
}

func ListToArray(l MetaList) []Meta {
	// PERFORMANCE: is it better to iterate twice, pass 1 to find length?
	meta := make([]Meta, 0)
	i := NewMetaListIterator(l, true)
	for i.HasNextMeta() {
		m := i.NextMeta()
		meta = append(meta, m)
	}
	return meta
}

func FindByPathWithoutResolvingProxies(root MetaList, path string) Meta {
	c := find(root, path, false)
	return c
}

func FindByPath(root MetaList, path string) Meta {
	return find(root, path, true)
}

func find(root MetaList, path string, resolveProxies bool) (def Meta) {
	elems := strings.SplitN(path, "/", -1)
	lastLevel := len(elems) - 1
	var ok bool
	list := root
	i := NewMetaListIterator(list, resolveProxies)
	for level, elem := range elems {
		def = FindByIdent(i, elem)
		if def == nil {
			return nil
		}
		if level < lastLevel {
			if list, ok = def.(MetaList); ok {
				i = NewMetaListIterator(list, resolveProxies)
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