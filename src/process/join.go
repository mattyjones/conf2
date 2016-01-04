package process

type Join struct {
	On   Table
	Into Table
}

func (j *Join) Next() (error) {
	// TODO: link thru key
	var err error
	if err = j.On.Next(); err != nil {
		return err
	}
	if err = j.Into.Next(); err != nil {
		return err
	}
	return nil
}

func (j *Join) HasNext() bool {
	return j.On.HasNext() && j.Into.HasNext()
}

func (j *Join) Select(identPath string, autocreate bool) (t Table, err error) {
	if autocreate {
		return j.Into.Select(identPath, true)
	}
	return j.On.Select(identPath, false)
}

func (j *Join) Get(key string) (v interface{}, err error) {
	if v, err = j.On.Get(key); v == nil && err == nil {
		v, err = j.Into.Get(key)
	}
	return v, err
}

func (j *Join) Set(key string, v interface{}) (err error) {
	return j.Into.Set(key, v)
}
