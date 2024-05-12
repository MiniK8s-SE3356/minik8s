package idgenerate_test

import (
	"fmt"
	"testing"

	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
)

func TestMain(m *testing.M) {
	a, err1 := idgenerate.GenerateID()
	b, err2 := idgenerate.GenerateID()
	if err1 != nil || err2 != nil {
		panic("error in generating")
	}
	fmt.Println(a, b)
	if a == b {
		panic("not random")
	}
}
