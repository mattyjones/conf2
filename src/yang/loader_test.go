package yang
import "testing"

func TestLoadLoader(t *testing.T) {
	//yyDebug = 4
	f := &FileResolver{}
	_, err:= LoadModule(f, "test_data/romancing-the-stone.yang")
	if err != nil {
		t.Error(err)
	}
}