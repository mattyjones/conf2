package mongo
import (
	"gopkg.in/mgo.v2/bson"
)

func DeepCopy(o interface{}) (copy interface{}) {
	switch val := o.(type) {
	case bson.M: {
		copy := make(bson.M, len(val))
		for key, mapObj := range val {
			switch mapObj.(type) {
			case bson.M, []interface{}:
				copy[key] = DeepCopy(mapObj)
			default:
				copy[key] = mapObj
			}
		}
		return copy
	}
	case []interface{}:
		copy := make([]interface{}, len(val))
		for i, arrayObj := range val {
			switch arrayObj.(type) {
			case bson.M, []interface{}:
				copy[i] = DeepCopy(arrayObj)
			default:
				copy[i] = arrayObj
			}
		}
		return copy
	default:
		return o
	}
}
