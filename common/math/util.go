package math

func MakeArray(count int32, f func(int32) any) []any {
	var result = make([]any, count)
	for i := int32(0); i < count; i++ {
		result[i] = f(i)
	}
	return result
}
