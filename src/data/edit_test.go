package data

import (
	"log"
	"schema/yang"
	"strings"
	"testing"
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
	var bd *BucketData
	if bd, err = LoadEditTestData(); err != nil {
		t.Fatal(err)
	}
	json := NewJsonReader(strings.NewReader(`{"origin":{"country":"Canada"}}`)).Node()

	// UPDATE
	// Here we're testing editing a specific list item. With FindTarget walk controller
	// needs to leave walkstate in a position for WalkTarget controller to make the edit
	// on the right item.
	log.Println("Testing edit\n")
	if err = bd.Select().Find("fruits=apple").Pull(json).Update().LastErr; err != nil {
		t.Fatal(err)
	}
	actual := MapValue(bd.Root, "fruits.1.origin.country")
	if actual != "Canada" {
		t.Error("Edit failed", actual)
	}

	// INSERT
	log.Println("Testing insert\n")
	json = NewJsonReader(strings.NewReader(`{"fruits":[{"name":"pear","origin":{"country":"Columbia"}}]}`)).Node()
	if err = bd.Select().Find("fruits").Pull(json).Insert().LastErr; err != nil {
		t.Fatal(err)
	}
	actual, found := bd.Root["fruits"]
	if !found {
		t.Error("fruits not found")
	} else {
		fruits := actual.([]map[string]interface{})
		if len(fruits) != 3 {
			t.Error("Expected 3 fruits but got ", len(fruits))
		}
	}
}

func LoadEditTestData() (*BucketData, error) {
	m, err := yang.LoadModuleFromByteArray([]byte(EDIT_TEST_MODULE), nil)
	if err != nil {
		return nil, err
	} else {
		// avoid using json to load because that needs edit/INSERT and
		// we don't want to use code to load seed data that we're trying to test
		fruits := map[string]interface{}{
			"fruits" : []map[string]interface{}{
				map[string]interface{}{
					"name" : "banana",
					"origin" : map[string]interface{}{
						"country" : "Brazil",
					},
				},
				map[string]interface{}{
					"name" : "apple",
					"origin" : map[string]interface{}{
						"country": "US",
					},
				},
			},
		}
		bb := &BucketData{Meta: m, Root:fruits}
		return bb, nil
	}
}
