package mongo
import (
	"schema/browse"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"schema"
	"errors"
)

func PathToQuery(initialState *browse.WalkState, p *browse.Path) (q interface{}, state *browse.WalkState, err error) {
	nKeys := countKeys(p)
	var all []bson.M
	if all, state, err = segmentsToQuearyParams(initialState, p, nKeys); err == nil {
		if nKeys == 0 {
			q = nil
		} else if nKeys == 1 {
			q = all[0]
		} else {
			q = bson.M{"$and" : all}
		}
	}
	return
}

func segmentsToQuearyParams(initialState *browse.WalkState, p *browse.Path, expectedSegments int) (q []bson.M, state *browse.WalkState, err error) {
	q = make([]bson.M, expectedSegments)
	var ndx int
	var path string
	state = initialState
	pathAppend := func(a string, b string) string {
		if len(a) == 0 {
			return b
		}
		return fmt.Sprint(a, ".", b)
	}
	for _, segment := range p.Segments {
		propMeta := schema.FindByIdent2(state.SelectedMeta(), segment.Ident)
		if propMeta == nil {
			return nil, nil, errors.New(fmt.Sprintf("Schema \"%s\" not found", segment.Ident))
		} else if schema.IsLeaf(propMeta) {
			return nil, nil, errors.New(fmt.Sprintf("Cannot navigate into leaf node \"%s\"", segment.Ident))
		}
		state.SetPosition(propMeta)
		state = state.Select()
		path = pathAppend(path, segment.Ident)
		if len(segment.Keys) > 0 {
			listMeta := propMeta.(*schema.List)
			values, err := browse.CoerseKeys(listMeta, segment.Keys)
			if err != nil {
				return nil, nil, err
			} else if len(values) > 1 {
				return nil, nil, errors.New("More than one key not supported in mongo db manager yet")
			}

			qpath := pathAppend(path, listMeta.Keys[0])
			q[ndx] = bson.M{qpath : values[0].Value()} // all?
			ndx++
			if ndx == expectedSegments {
				break
			}
		}
	}

	return
}

func countKeys(p *browse.Path) int {
	var nKeys int
	for _, segment := range p.Segments {
		if len(segment.Keys) > 0 {
			nKeys++
		}
	}
	return nKeys
}
