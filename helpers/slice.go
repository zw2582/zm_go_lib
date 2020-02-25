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

func RandSlice(slice interface{}) interface{} {
	//slice1 := slice.([]interface{})
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic("type error")
	}

	rand.Seed(time.Now().Unix())
	for i := v.Len() - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		numTmp := reflect.ValueOf(v.Index(num).Interface())
		//numTmp := v.Index(num)
		v.Index(num).Set(v.Index(i))
		v.Index(i).Set(numTmp)
	}
	return v.Interface()
}

type MapFunc func(i int, val interface{}) interface{}

//SliceMap map slice
func SliceMap(slice interface{}, fn MapFunc) interface{} {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic("type error")
	}
	res := make([]interface{}, 0, v.Cap())
	for k := 0; k < v.Len(); k++ {
		tmp := fn(k, v.Index(k).Interface())
		res = append(res, reflect.ValueOf(tmp))
	}
	return res
}

func SliceColumn(slice interface{}, keystr string, params ...string) map[interface{}]interface{} {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic("type error")
	}
	valstr := ""
	if len(params) == 1 {
		valstr = params[0]
	}
	if keystr == "" {
		panic("give keystr is empty")
	}
	res := make(map[interface{}]interface{})
	for k := 0; k < v.Len(); k++ {
		obj := v.Index(k)
		if obj.Kind() != reflect.Struct {
			panic("slice value must stuct")
		}
		key := obj.FieldByName(keystr).Interface()
		if valstr == "" {
			res[key] = obj.Interface()
		} else {
			res[key] = obj.FieldByName(valstr).Interface()
		}
	}
	return res
}

// 通过两重循环过滤重复元素
func RemoveRepByLoop(slice interface{}) interface{} {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic("type error")
	}
	// 存放结果
	result := reflect.MakeSlice(reflect.TypeOf(slice),0, v.Len())

	for i:=0; i< v.Len(); i++ {
		flag := true
		for j :=0; j < result.Len(); j++ {
			if v.Index(i).Interface() == result.Index(j).Interface() {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag {
			result = reflect.Append(result, v.Index(i))
			//reflect.AppendSlice(result, v.Index(i))
		}
	}
	return result.Interface()
}

// 通过map主键唯一的特性过滤重复元素
func RemoveRepByMap(slc interface{}) interface{} {
	v := reflect.ValueOf(slc)
	if v.Kind() != reflect.Slice {
		panic("type error")
	}
	// 存放结果
	result := reflect.MakeSlice(reflect.TypeOf(slc),0, v.Len())
	// 存放不重复主键
	tempMap := reflect.MakeMap(reflect.TypeOf(map[interface{}]int{}))

	for i:=0;i<v.Len();i++ {
		t := v.Index(i)
		if tempMap.MapIndex(t) == reflect.ValueOf(nil) {
			result = reflect.Append(result, t)
			tempMap.SetMapIndex(t, reflect.ValueOf(1))
		}
	}
	return result.Interface()
}

// 元素去重
func RemoveRep(slc interface{}) interface{} {
	v := reflect.ValueOf(slc)
	if v.Kind() != reflect.Slice {
		panic("type error")
	}
	if v.Len() < 1024 {
		// 切片长度小于1024的时候，循环来过滤
		return RemoveRepByLoop(slc)
	} else {
		// 大于的时候，通过map来过滤
		return RemoveRepByMap(slc)
	}
}

//求切片的差集
func DiffSlice(one, two interface{}) interface{} {
	v1 := reflect.ValueOf(one)
	if v1.Kind() != reflect.Slice {
		panic("type error")
	}
	v2 := reflect.ValueOf(two)
	if v2.Kind() != reflect.Slice {
		panic("type error")
	}

	// 存放结果
	result := reflect.MakeSlice(reflect.TypeOf(one),0, v1.Len())

	for i:=0; i< v1.Len(); i++ {
		flag := true
		for j :=0; j < v2.Len(); j++ {
			if v1.Index(i).Interface() == v2.Index(j).Interface() {
				flag = false // 该值在v2集合中存在，不是差集，结束循环
				break
			}
		}
		if flag {
			result = reflect.Append(result, v1.Index(i))
		}
	}
	return result.Interface()
}