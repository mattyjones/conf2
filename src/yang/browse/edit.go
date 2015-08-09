package browse
import (
	"yang"
)

type Operation int
const (
	CREATE_CHILD Operation = 1 + iota // 1
	POST_CREATE_CHILD                 // 2
	CREATE_LIST                       // 3
	POST_CREATE_LIST                  // 4
	UPDATE_VALUE                      // 5
	DELETE_CHILD                      // 6
	DELETE_LIST                       // 7
	BEGIN_EDIT                        // 8
	END_EDIT                          // 9
)

type strategy int
const (
	UPSERT strategy = iota + 1
	INSERT
	UPDATE
	DELETE
	CLEAR
)

type editor struct {
	target *Path
	strategy strategy
	from *Selection
}

func InsertIntoPath(from *Selection, to *Selection, p *Path) error {
	return edit(from, to, p, to.Meta, INSERT)
}

func Insert(from *Selection, to *Selection) error {
	return edit(from, to, nil, from.Meta, INSERT)
}

func Upsert(from *Selection, to *Selection) error {
	return edit(from, to, nil, from.Meta, UPSERT)
}

func UpsertIntoPath(from *Selection, to *Selection, p *Path) error {
	return edit(from, to, p, to.Meta, UPSERT)
}

func Delete(from *Selection, to *Selection, p *Path) error {
	return edit(from, to, p, to.Meta, DELETE)
}

func Update(from *Selection, to *Selection) error {
	return edit(from, to, nil, from.Meta, UPDATE)
}

func UpdateIntoPath(from *Selection, to *Selection, p *Path) error {
	return edit(from, to, p, to.Meta, UPDATE)
}

func edit(from *Selection, dest *Selection, p *Path, meta yang.MetaList, strategy strategy) (err error) {
	e := editor{target:p, strategy:strategy, from:from}
	var s *Selection
	s, err = e.findTarget(dest, []string{})
	if err == nil {
		s.Meta = meta
		if err = dest.Edit(BEGIN_EDIT, nil); err == nil {
			if err = Walk(s, p); err == nil {
				err = dest.Edit(END_EDIT, nil)
			} else {
			}
		}
	}
	return
}

func isEditTarget(target *Path, candidate []string) bool {
	if target == nil {
		return true
	}
	if len(target.Segments) != len(candidate) {
		return false
	}
	for i, _ := range target.Segments {
		if target.Segments[i].Ident != candidate[i] {
			return false
		}
	}
	return true
}

func (e *editor) findTarget(dest *Selection, path []string) (*Selection, error) {
	if isEditTarget(e.target, path) {
		return e.editTarget(e.from, dest, e.strategy)
	}
	s := &Selection{}
	s.Select = func(ident string) (bSelection *Selection, err error) {
		dest.Meta = s.Meta
		bSelection, err = dest.Select(ident)
		s.Position = dest.Position
		if err == nil && bSelection != nil {
			return e.findTarget(dest, append(path, ident))
		}
		return
	}
	s.Read = dest.Read
	s.Iterate = dest.Iterate
	return s, nil
}

func (e *editor) editTarget(from *Selection, to *Selection, strategy strategy) (*Selection, error) {
	var createdChild bool
	var createdList bool
	s := &Selection{}
	s.Select = func(ident string) (c *Selection, err error) {
		from.Meta = s.Meta
		to.Meta = s.Meta

		fromChild, err := from.Select(ident)
		s.Position = from.Position
		if err != nil {
			return
		}

		toChild, err := to.Select(ident)
		if err != nil {
			return
		} else if toChild == nil {
			return nil, &browseError{Msg:"source could not be selected"}
		}

		if fromChild == nil || (!from.Found && !to.Found) {
			s.Found = from.Found
			return
		}

		s.Found = from.Found
		nextStrategy := strategy
		if from.Found && to.Found {
			switch strategy {
			case INSERT:
				err = &browseError{Msg:"Duplicate object found"}
			case UPDATE:
				strategy = CLEAR
			case DELETE:
				if yang.IsList(s.Position) {
					err = to.DeleteList()
				} else {
					err = to.DeleteChild()
				}
				s.Found = false
			}
		} else if from.Found && !to.Found {
			switch strategy {
			case UPSERT, INSERT, CLEAR:
				if yang.IsList(s.Position) {
					err = to.CreateList()
					createdList = true
				} else {
					if err = to.CreateChild(); err == nil {
						createdChild = true
						toChild, err = to.Select(ident)
						if toChild == nil {
							err = &browseError{Msg:"Could not select object that was just created"}
						}
					}
				}
			case UPDATE, DELETE:
				err = &browseError{Msg:"No such object"}
			}
		} else if !from.Found && to.Found {
			switch strategy {
			case DELETE, CLEAR:
				if yang.IsList(s.Position) {
					err = to.DeleteList()
				} else {
					err = to.DeleteChild()
				}
				s.Found = false
			}
		}

		if err == nil && s.Found {
			return e.editTarget(fromChild, toChild, nextStrategy)
		}

		return
	}
	s.Exit = func() (err error) {
		if createdChild {
			if err = to.FinishCreateChild(); err != nil {
				return
			}
			createdChild = false
		}
		if createdList {
			if err = to.FinishCreateList(); err != nil {
				return
			}
			createdList = false
		}
		if from.Exit != nil {
			if err = from.Exit(); err != nil {
				return
			}
		}
		if to.Exit != nil {
			if err = to.Exit(); err != nil {
				return
			}
		}
		return
	}
	s.Read = func(v *Value) (err error) {
		var copy bool
		var clear bool
		if from.Found && to.Found {
			switch strategy {
			case UPSERT, UPDATE, CLEAR:
				copy = true
			}
		} else if from.Found && !to.Found {
			switch strategy {
			case INSERT, UPSERT, CLEAR:
				copy = true
			}
		} else if !from.Found && to.Found {
			switch strategy {
			case UPDATE, CLEAR:
				clear = true
			}
		}
		if copy {
			if err = from.Read(v); err != nil {
				return
			}
		}
		if copy || clear {
			if err = to.SetValue(v); err != nil {
				return
			}
		}
		return
	}
	if from.Iterate != nil {
		s.Iterate = func(fromKeys []string, first bool) (fromMore bool, err error) {
			from.Meta = s.Meta
			to.Meta = s.Meta
			fromMore, err = from.Iterate(fromKeys, first)
			if err != nil || !fromMore {
				return
			}

// TODO
//			toKeys := fromKeys
//			if len(toKeys) == 0 {
//				keyIdents := s.Meta.(*yang.List).Keys
//				toKeys = make([]string, len(keyIdents))
//				for i, keyIdent := range keyIdents {
//					v := &Value{}
//					if _, err = from.Select(keyIdent); err != nil {
//						return
//					}
//					if err = from.Read(v); err != nil {
//						return
//					}
//					// TODO: don't assume key is a string
//					toKeys[i] = v.Str
//				}
//			}
//
//			// ignore if exists or not, next Select will detect existance for lists and container
//			// selections.
//			_, err = to.Iterate(toKeys, true)
//			if err != nil {
//				return
//			}
			return fromMore, err
		}
	}
	return s, nil
}
