package browse
import "sort"

type StringIndexBuilder interface {
	Select(key string) bool
	Build() []string
}

type StringIndex struct {
	Position int
	Keys []string
	Builder StringIndexBuilder
}

func (i *StringIndex) OnNext(key []Value, first bool) (bool, error) {
	if (len(key) > 0) {
		if first {
			i.Position = 0
			i.Keys = []string { key[0].Str }
		} else {
			i.Position++
		}
	} else {
		if first {
			i.Keys = i.Builder.Build()
			sort.Strings(i.Keys)
		} else {
			i.Position++
		}
	}
	if i.Position < len(i.Keys) {
		return i.Builder.Select(i.Keys[i.Position]), nil
	}
	return false, nil
}
