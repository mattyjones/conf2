package mongo
import "gopkg.in/mgo.v2/bson"

func DeepCopy(original bson.M) (copy bson.M) {
	if len(original) > 0 {
		copy = make(bson.M, len(original))
		for k, v := range original {
			switch val := v.(type) {
			case bson.M:
				copy[k] = DeepCopy(val)
			default:
				copy[k] = val
			}
		}
	}
	return copy
}
