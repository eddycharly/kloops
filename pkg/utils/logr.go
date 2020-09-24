package utils

// MapValues maps a map to an array of key values as expected by logr
func MapValues(values map[string]string) []interface{} {
	var ret []interface{}
	for k, v := range values {
		ret = append(ret, k, v)
	}
	return ret
}
