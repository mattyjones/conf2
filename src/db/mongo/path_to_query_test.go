package mongo
import (
	"testing"
	"schema/yang"
	"schema/browse"
	"reflect"
	"gopkg.in/mgo.v2/bson"
)

func TestPathToMongoQuery(t *testing.T) {
	mstr := `
module birding {
	prefix "";
	namespace "";
	revision 0000-00-00 {
		description "";
	}
	list birds {
		key "name";
		leaf name {
			type string;
		}
		container details {
			leaf wingSpan {
				type int32;
			}
			list regions {
				key "region";
				leaf region {
					type string;
				}
			}
		}
	}
}
	`
	tests := [] struct {
		xpath string
		expected interface{}
		valid bool
	} {
		{"birds", nil, true},
		{"birds=sparrow", bson.M{"birds.name": "sparrow"}, true},
		{"birds=sparrow/details/regions=US", bson.M{
			"$and" : []bson.M{
				bson.M{"birds.name": "sparrow"},
				bson.M{"birds.details.regions.region": "US"},
			}}, true,
		},
		{"birds/bogus=xx", nil, false},
	}
	m, err := yang.LoadModuleFromByteArray([]byte(mstr), nil)
	if err != nil {
		t.Fatal(err)
	}
	var actual interface{}
	for _, test := range tests {
		actual, _, err = PathToQuery(browse.NewWalkState(m), browse.NewPath(test.xpath))
		if err != nil {
			if test.valid {
				t.Error(err)
			}
		} else {
			if ! reflect.DeepEqual(actual, test.expected) {
				t.Errorf("\nExpected : %s\n  Actual : %s", test.expected, actual)
			}
		}
	}
}
