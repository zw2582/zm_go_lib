package helpers

import (
	"math/rand"
	"reflect"
	"time"
)

//IndexOf 判断元素是否存在于slice中
func IndexOf(larr interface{}, a interface{}) int {
	v := reflect.ValueOf(a)
	arr := reflect.ValueOf(larr)

	var t = arr.Kind()

	if t != reflect.Slice && t != reflect.Array {
		panic("Type Error! Second argument must be an array or a slice.")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == v.Interface() {
			return i
		}
	}
	return -1
}

func randSlice(slice interface{}) interface{} {
	//slice1 := slice.([]interface{})
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic("type error")
	}

	rand.Seed(time.Now().Unix())
	for i:=v.Len() -1;i>0;i-- {
		num := rand.Intn(i +1)
		numTmp := reflect.ValueOf(v.Index(num).Interface())
		//numTmp := v.Index(num)
		v.Index(num).Set(v.Index(i))
		v.Index(i).Set(numTmp)
	}
	return v.Interface()
}