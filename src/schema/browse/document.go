package browse
import "schema"


type Document interface {
	Selector(path *Path) (*Selection, error)
	Schema() (schema.MetaList)
}