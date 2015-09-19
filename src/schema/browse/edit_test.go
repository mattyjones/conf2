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
	if b, err = LoadTestData(); err != nil {
		t.Fatal(err)
	}
	var s Selection
	if s, err = b.RootSelector(); err != nil {
		t.Fatal(err)
	}
	p, _ := NewPath("fruits=apple")
	var target Selection
	log.Printf("Walk path to find apple in list\n")
	target, err = WalkPath(s, p)
	if target == nil {
		t.Fatal("Could not find target");
	}
	var edit Selection
	log.Println("Testing edit\n")
	edit, err = LoadTestEdit(target, `{"origin":{"country":"Canada"}}`)
	if err != nil {
		t.Fatal(err)
	}

	// UPDATE
	var wc WalkController
	wc, err = p.WalkTargetController()
	if err != nil {
		t.Error(err)
	} else {
		err = Update(edit, target, wc)
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
	}

	// INSERT
	p, _ = NewPath("fruits")
	target, err = WalkPath(s, p)
	if target == nil {
		t.Fatal("Could not find target");
	}
	wc, err = p.WalkTargetController()
	if err != nil {
		t.Fatal(err)
	}
	edit, err = LoadTestEdit(target, `{"fruits":[{"name":"pear","origin":{"country":"Columbia"}}]}`)
	err = Insert(edit, target, wc)
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

func LoadTestEdit(target Selection, edit string) (Selection, error) {
	r := NewJsonReader(strings.NewReader(edit))
	return r.GetSelector(target.WalkState().Meta, target.WalkState().InsideList)
}

func LoadTestData() (*BucketBrowser, error) {
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