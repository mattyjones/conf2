package data
import (
	"schema"
	"schema/yang"
	"bytes"
)

type Testing interface {
	Fatal(args ...interface{})
}

type ModuleTestSetup struct {
	Module *schema.Module
	Store *BufferStore
	Data *StoreData
}

func ModuleSetup(mstr string, t Testing) (setup *ModuleTestSetup) {
	setup = &ModuleTestSetup{}
	var err error
	setup.Module, err = yang.LoadModuleFromByteArray([]byte(mstr), nil)
	if err != nil {
		t.Fatal(err)
	}
	setup.Store = NewBufferStore()
	setup.Data = NewStoreData(setup.Module, setup.Store)
	return
}

func (setup *ModuleTestSetup) ToString(t Testing) string {
	var actualBuff bytes.Buffer
	out := NewJsonWriter(&actualBuff).Node()
	err := NodeToNode(setup.Data.Node(), out, setup.Data.Schema()).Insert()
	if err != nil {
		t.Fatal(err)
	}
	return actualBuff.String()
}


