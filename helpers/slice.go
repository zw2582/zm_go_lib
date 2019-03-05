package helpers

import "reflect"

//IndexOf 判断元素是否存在于slice中
func IndexOf(larr interface{}, a interface{}) int {
	v := reflect.ValueOf(a)
	arr := reflect.ValueOf(larr)

	var t = arr.Kind()

	if t != reflect.Slice && t != reflect.Array {
		panic("Type Error! Second argument must be an array or a slice.")
	}

	for i := 0; i < arr.Len()-1; i++ {
		if arr.Index(i).Interface() == v.Interface() {
			return i
		}
	}
	return -1
}