package browse
import (
	"testing"
	"schema/yang"
	"strings"
	"log"
)

const EDIT_TEST_MODULE = `
module food {
	prefix "x";
	namespace "y";
	revision 0000-00-00 {
		description "";
	}
	list fruits  {
		key "name";
		leaf name {
			type string;
		}
		container origin {
			leaf country {
				type string;
			}
		}
	}
}
`

func TestEditListItem(t *testing.T) {
	var err error
	var b *BucketBrowser
	if b, err = LoadEditTestData(); err != nil {
		t.Fatal(err)
	}
	in := NewJsonFragmentReader(strings.NewReader(`{"origin":{"country":"Canada"}}`))

	// UPDATE
	// Here we're testing editing a specific list item. With FindTarget walk controller
	// needs to leave walkstate in a position for WalkTarget controller to make the edit
	// on the right item.
	log.Println("Testing edit\n")
	err = Update(NewPath("fruits=apple"), in, b)
	if err != nil {
		t.Error(err)
	} else {
		var actual interface{}
		if actual, err = b.Read("fruits.1.origin.country"); err != nil {
			t.Error(err)
		} else if actual != "Canada" {
			t.Error("Edit failed", actual)
		}
	}

	// INSERT
	in = NewJsonFragmentReader(strings.NewReader(`{"fruits":[{"name":"pear","origin":{"country":"Columbia"}}]}`))
	log.Println("Testing insert\n")
	err = Insert(NewPath("fruits"), in, b)
	if err != nil {
		t.Error(err)
	} else {
		var actual interface{}
		if actual, err = b.Read("fruits"); err != nil {
			t.Error(err)
		} else {
			fruits := actual.([]map[string]interface{})
			if len(fruits) != 3 {
				t.Error("Expected 3 fruits but got ", len(fruits))
			}
		}
	}
}

func LoadEditTestData() (*BucketBrowser, error) {
	m, err := yang.LoadModuleFromByteArray([]byte(EDIT_TEST_MODULE), nil)
	if err != nil {
		return nil, err
	} else {
		bb := NewBucketBrowser(m)
		// avoid using json to load because that needs edit/INSERT and
		// we don't want to use code to load seed data that we're trying to test
		fruits := make([]map[string]interface{}, 2)
		fruits[0] = map[string]interface{} {
			"name" : "banana",
		}
		fruits[0]["origin"] = map[string]interface{} {
			"country" : "Brazil",
		}
		fruits[1] = map[string]interface{} {
			"name" : "apple",
		}
		fruits[1]["origin"] = map[string]interface{} {
			"country" : "US",
		}
		bb.Bucket["fruits"] = fruits
		return bb, nil
	}
}