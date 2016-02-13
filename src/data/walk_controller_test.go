package data

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
	var b *BucketData
	if b, err = LoadPathTestData(); err != nil {
		t.Fatal(err)
	}
	tests := []string {
		"fruits=apple/origin",
		"fruits=apple/boat",
	}
	for _, test := range tests {
		target := b.Select().Find(test)
		if target.LastErr != nil {
			t.Fatal(target.LastErr)
		} else if target.Selection == nil {
			t.Fatal("Could not find target " + test)
		}
	}
}

func LoadPathTestData() (*BucketData, error) {
	m, err := yang.LoadModuleFromByteArray([]byte(walkTestModule), nil)
	if err != nil {
		return nil, err
	} else {
		// avoid using json to load because that needs edit/INSERT and
		// we don't want to use code to load seed data that we're trying to test
		data := map[string]interface{}{
			"fruits" : []map[string]interface{}{
				map[string]interface{}{
					"name" : "banana",
					"origin" : map[string]interface{}{
						"country" : "Brazil",
					},
					"plane" : map[string]interface{}{
						"name" : "747c",
					},
				},
				map[string]interface{}{
					"name" : "apple",
					"origin" : map[string]interface{}{
						"country" : "US",
					},
					"boat" : map[string]interface{}{
						"name" : "SS Hudson",
					},
				},
			},
		}
		bb := &BucketData{Meta :m, Root: data}
		return bb, nil
	}
}
