package browse
import (
	"sort"
	"schema"
)

// Example:
//   s := &MySelection{}
//   index := newMappingIndex(data)
//   s.OnNext = index.Index.OnNext
//   ...
//
//	type mappingIndex struct {
//		Index browse.StringIndex
//		Data map[string]*BridgeMapping
//		Selected *BridgeMapping
//	}
//
//	func newMappingIndex(data map[string]*BridgeMapping) *mappingIndex {
//		ndx := &mappingIndex{Data:data}
//		ndx.Index.Builder = ndx
//		return ndx
//	}
//
//	func (impl *mappingIndex) Select(key string) (found bool) {
//		impl.Selected, found = impl.Data[key]
//		return
//	}
//
//	func (impl *mappingIndex) Build() []string {
//		index := make([]string, len(impl.Data))
//		j := 0
//		for key, _ := range impl.Data {
//			index[j] = key
//			j++
//		}
//		return index
//	}

type StringIndexBuilder interface {
	Select(key string) bool
	Build() []string
}

type StringIndex struct {
	Position int
	Keys []string
	Builder StringIndexBuilder
}

func (i *StringIndex) CurrentKey() string {
	return i.Keys[i.Position]
}

func (i *StringIndex) OnNext(state *WalkState, meta *schema.List, key []*Value, first bool) (bool, error) {
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
