package dns

import (
	"reflect"
)

func GetTagValue(myStruct interface{}, myField string, myTag string) string {
	t := reflect.TypeOf(myStruct)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == myField {
			tag := field.Tag.Get(myTag)
			return tag
		}
	}

	return ""
}
