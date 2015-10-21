package mongo
import (
	"schema"
	"schema/browse"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoBrowser struct {
	schema schema.MetaList
	c *mgo.Collection
}

func NewMongoBrowser(schema schema.MetaList, c *mgo.Collection) *MongoBrowser {
	return &MongoBrowser{schema:schema, c:c}
}

func (self *MongoBrowser) Selector(path *browse.Path, strategy browse.Strategy) (browse.Selection, *browse.WalkState, error) {
	if strategy == browse.READ {
		return self.ReadSelector(path)
	}

	return self.WriteSelector(path, strategy)
}

func (self *MongoBrowser) Schema() schema.MetaList {
	return self.schema
}

type mongoBrowserReader struct{
}

type mongoBrowserWriter struct{
	c *mgo.Collection
	selector interface{}
	strategy browse.Strategy
}

func (self *MongoBrowser) WriteSelector(p *browse.Path, strategy browse.Strategy) (s browse.Selection, state *browse.WalkState, err error) {
	w := mongoBrowserWriter{c:self.c, strategy:strategy}
	if w.selector, state, err = PathToQuery(browse.NewWalkState(self.schema), p); err != nil {
		return nil, nil, err
	}
	data := make(bson.M, 10)
	s, err = w.writeResults(data, nil)
	return s, state, err
}

func (self *mongoBrowserWriter) writeResults(data bson.M, list []bson.M) (browse.Selection, error) {
	s := &browse.MySelection{}
	var created browse.Selection
	container := data
	s.OnNext = func(state *browse.WalkState, meta *schema.List, key []*browse.Value, first bool) (next browse.Selection, err error) {
		next, created = created, nil
		return
	}
	s.OnSelect = func(state *browse.WalkState, meta schema.MetaList) (child browse.Selection, err error) {
		child, created = created, nil
		return
	}
	s.OnWrite = func(state *browse.WalkState, meta schema.Meta, op browse.Operation, v *browse.Value) (err error) {
		switch op {
		case browse.CREATE_LIST:
			childList := make([]bson.M, 0, 1)
			data[meta.GetIdent()] = childList
			created, _ = self.writeResults(data, childList)
		case browse.CREATE_LIST_ITEM:
			container = make(bson.M, 10)
			list = append(list, container)
			// refresh parent reference in case list reference changed
			data[meta.GetIdent()] = list
			created, _ = self.writeResults(container, nil)
		case browse.CREATE_CHILD:
			child := make(bson.M, 10)
			container[meta.GetIdent()] = child
			created, _ = self.writeResults(child, nil)
		case browse.UPDATE_VALUE:
			container[meta.GetIdent()] = v.Value()
		case browse.END_EDIT:
			if self.selector == nil {
				err = self.c.Insert(data)
			} else {
				_, err = self.c.Upsert(self.selector, data)
			}
		}
		return
	}
	return s, nil
}

func (self *MongoBrowser) ReadSelector(p *browse.Path) (s browse.Selection, state *browse.WalkState, err error) {
	r := mongoBrowserReader{}
	var selector interface{}
	//var meta schema.MetaList

	var results bson.M

	if len(p.Segments) == 1 && len(p.Segments[0].Keys) == 1 {
		// TODO: Hack, specialized query for demo

		var recordResult bson.M
		err = self.c.FindId(bson.ObjectIdHex(p.Segments[0].Keys[0])).One(&recordResult)
		results = bson.M{"records": []interface{}{recordResult}}

	} else {
		if selector, _, err = PathToQuery(browse.NewWalkState(self.schema), p); err != nil {
			return nil, nil, err
		}
		err = self.c.Find(selector).One(&results)
	}

	// TODO: restrict results to only the data that was asked for.
	if err != nil {
		return nil, nil, nil
	}
	var root browse.Selection
	root, err = r.readContainer(results)

	// the result tree goes all the way back to the document root.  we need to navigate to
	// the point of the path and throw away results.
	s, state, err = browse.WalkPath(browse.NewWalkState(self.schema), root, p)
	return
}

func (self *mongoBrowserReader) readContainer(container bson.M) (browse.Selection, error) {
	s := &browse.MySelection{}
	s.OnSelect = func(state *browse.WalkState, meta schema.MetaList) (browse.Selection, error) {
		selectValue, found := container[meta.GetIdent()]
		if found {
			if schema.IsList(meta) {
				return self.readList(selectValue.([]interface{}))
			} else {
				return self.readContainer(selectValue.(bson.M))
			}
		}
		return nil, nil
	}
	s.OnRead = func(state *browse.WalkState, meta schema.HasDataType) (*browse.Value, error) {
		if meta == nil {
			return nil, nil
		}
		value := container[meta.GetIdent()]
		switch meta.GetDataType().Format {
		case schema.FMT_BOOLEAN:
			b := value.(int64) > 0
			return browse.SetValue(meta.GetDataType(), b)
		}
		return browse.SetValue(meta.GetDataType(), value)
	}
	return s, nil
}

func (self *mongoBrowserReader) readList(list []interface{}) (browse.Selection, error) {
	s := &browse.MySelection{}
	var i int
	s.OnNext = func(state *browse.WalkState, meta *schema.List, key []*browse.Value, first bool) (browse.Selection, error) {
		if len(key) > 0 {
			for _, candidate := range list {
				if m, ok := candidate.(bson.M); ok {
					if m[meta.Keys[0]] == key[0].Value() {
						return self.readContainer(m)
					}
				}
			}
		} else {
			if first {
				i = 0
			} else {
				i++
			}
			if i < len(list) {
				return self.readContainer(list[i].(bson.M))
			}
		}
		return nil, nil
	}
	return s, nil
}
