package helpers

import (
	"fmt"
	"strings"
	"testing"
)

func TestMd5encode(t *testing.T) {
	fmt.Println(Md5encode("sdfsfd"))
}

func TestValidIdcard(t *testing.T)  {
	idcard := []string{
		"422801199301013819",
		"32048120160909009x",
		"432522199202294050",
		"500228199612193111",
		"500383199102288452",
	}

	for _,v := range idcard {
		v = strings.ToUpper(v)
		vv := []byte(v)
		result := IsValidCitizenNo(&vv)
		t.Log(result)
	}
}