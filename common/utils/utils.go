package utils

func MakeArray(count int32, f func(int32) interface{}) []interface{} {
	var result = make([]interface{}, count)
	for i := int32(0); i < count; i++ {
		result[i] = f(i)
	}
	return result
}
