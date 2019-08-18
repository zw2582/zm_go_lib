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
	for i:=v.Len() -1;i>0;i-- {
		num := rand.Intn(i +1)
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
	if len(params) == 1{
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
func RemoveRepByLoop(slc []interface{}) []interface{} {
var result = make([]interface{}, 0)  // 存放结果
	for i := range slc{
		flag := true
		for j := range result{
			if slc[i] == result[j] {
				flag = false  // 存在重复元素，标识为false
				break
			}
		}
		if flag {  // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

// 通过map主键唯一的特性过滤重复元素
func RemoveRepByMap(slc []interface{}) []interface{} {
	var result = make([]interface{}, 0)
	tempMap := map[interface{}]byte{}  // 存放不重复主键
	for _, e := range slc{
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l{  // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

// 元素去重
func RemoveRep(slc []interface{}) []interface{} {
	if len(slc) < 1024 {
		// 切片长度小于1024的时候，循环来过滤
		return RemoveRepByLoop(slc)
	}else{
		// 大于的时候，通过map来过滤
		return RemoveRepByMap(slc)
	}
}