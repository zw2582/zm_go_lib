package helpers

import (
	"testing"
)

func TestRemoveRepByLoop(t *testing.T) {
	params := []interface{}{
		[]interface{}{[]string{"a","b","a","b","c","a"}, 3},
		[]interface{}{[]int{1,2,2,3,4,5,5,4,2,3}, 5},
	}

	for _, val := range params {
		v := val.([]interface{})
		if tmp := RemoveRepByLoop(v[0]); tmp != nil {
			lenn := 0
			switch vv := tmp.(type) {
			case []string:
				lenn = len(vv)
			case []int:
				lenn = len(vv)
			}
			if lenn != v[1] {
				t.Error("去重失败")
			}
		}
	}
}

func TestRemoveRepByMap(t *testing.T) {
	params := []interface{}{
		[]interface{}{[]string{"a","b","a","b","c","a"}, 3},
		[]interface{}{[]int{1,2,2,3,4,5,5,4,2,3}, 5},
	}

	for _, val := range params {
		v := val.([]interface{})
		if tmp := RemoveRepByMap(v[0]); tmp != nil {
			lenn := 0
			switch vv := tmp.(type) {
			case []string:
				lenn = len(vv)
			case []int:
				lenn = len(vv)
			}
			if lenn != v[1] {
				t.Error("去重失败")
			}
		}
	}
}
