package browse
import "schema"


type Data interface {
	Selector(path *Path) (*Selection, error)
	Schema() (schema.MetaList)
}