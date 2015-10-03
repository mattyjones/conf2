package mongo
import (
	"schema/browse"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"schema"
	"errors"
)

func PathToQuery(parent schema.MetaList, p *browse.Path) (q interface{}, meta schema.MetaList, err error) {
	nKeys := countKeys(p)
	var all []bson.M
	if all, meta, err = segmentsToQuearyParams(parent, p, nKeys); err == nil {
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

func segmentsToQuearyParams(parent schema.MetaList, p *browse.Path, expectedSegments int) (q []bson.M, meta schema.MetaList, err error) {
	q = make([]bson.M, expectedSegments)
	var ndx int
	var path string
	meta = parent
	pathAppend := func(a string, b string) string {
		if len(a) == 0 {
			return b
		}
		return fmt.Sprint(a, ".", b)
	}
	for _, segment := range p.Segments {
		propMeta := schema.FindByIdent2(meta, segment.Ident)
		if propMeta == nil {
			return nil, nil, errors.New(fmt.Sprintf("Schema \"%s\" not found", segment.Ident))
		} else if schema.IsLeaf(propMeta) {
			return nil, nil, errors.New(fmt.Sprintf("Cannot navigate into leaf node \"%s\"", segment.Ident))
		}
		meta = propMeta.(schema.MetaList)
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
