package db
import "schema/browse"

type Store interface {
	Load(p *browse.Path) (map[string]*browse.Value, error)
	Clear(p *browse.Path) error
	Upsert(p *browse.Path, vals map[string]*browse.Value) error
}