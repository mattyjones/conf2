package browse

import (
	"schema"
	"sort"
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
	Keys     []string
	Builder  StringIndexBuilder
}

func (i *StringIndex) CurrentKey() string {
	return i.Keys[i.Position]
}

func (i *StringIndex) OnNext(state *Selection, meta *schema.List, key []*Value, first bool) (hasMore bool, err error) {
	if len(key) > 0 {
		if first {
			i.Position = 0
			i.Keys = []string{key[0].Str}
			hasMore, err = i.Builder.Select(i.Keys[0]), nil
			state.SetKey(key)
		} else {
			hasMore = false
		}
	} else {
		if first {
			i.Keys = i.Builder.Build()
			sort.Strings(i.Keys)
		} else {
			i.Position++
		}
		if i.Position < len(i.Keys) {
			hasMore, err = i.Builder.Select(i.Keys[i.Position]), nil
			if hasMore {
				var positionKey []*Value
				positionKey, err = CoerseKeys(meta, []string{i.Keys[i.Position]})
				state.SetKey(positionKey)
			}
		} else {
			hasMore = false
		}
	}

	return
}
