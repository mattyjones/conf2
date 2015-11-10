package browse

import (
	"schema/yang"
	"testing"
)

const walkTestModule = `
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
		choice shipment {
			case water {
				container boat {
					leaf name {
						type string;
					}
				}
			}
			case air {
				container plane {
					leaf name {
						type string;
					}
				}
			}
		}
	}
}
`

func TestPathIntoListItemContainer(t *testing.T) {
	var err error
	var b *BucketBrowser
	if b, err = LoadPathTestData(); err != nil {
		t.Fatal(err)
	}
	var target *Selection
	if target, err = b.Selector(NewPath("fruits=apple/origin")); err != nil {
		t.Fatal(err)
	} else if target == nil {
		t.Fatal("Could not find target")
	}

	if target, err = b.Selector(NewPath("fruits=apple/boat")); err != nil {
		t.Fatal(err)
	} else if target == nil {
		t.Fatal("Could not find target")
	}
}

func LoadPathTestData() (*BucketBrowser, error) {
	m, err := yang.LoadModuleFromByteArray([]byte(walkTestModule), nil)
	if err != nil {
		return nil, err
	} else {
		bb := NewBucketBrowser(m)
		// avoid using json to load because that needs edit/INSERT and
		// we don't want to use code to load seed data that we're trying to test
		fruits := make([]map[string]interface{}, 2)
		fruits[0] = map[string]interface{}{
			"name": "banana",
		}
		fruits[0]["origin"] = map[string]interface{}{
			"country": "Brazil",
		}
		fruits[0]["plane"] = map[string]interface{}{
			"name": "747c",
		}
		fruits[1] = map[string]interface{}{
			"name": "apple",
		}
		fruits[1]["origin"] = map[string]interface{}{
			"country": "US",
		}
		fruits[1]["boat"] = map[string]interface{}{
			"name": "SS Hudson",
		}
		bb.Bucket["fruits"] = fruits
		return bb, nil
	}
}
