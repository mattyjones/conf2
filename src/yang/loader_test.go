package yang
import "testing"

func TestLoadBson(t *testing.T) {
	//yyDebug = 4
	err, _ := LoadModule("test_data/romancing-the-stone.yang")
	if err != nil {
		t.Error(err)
	}
}