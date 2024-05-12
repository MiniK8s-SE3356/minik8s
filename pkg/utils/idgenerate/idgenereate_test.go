package idgenerate_test

import (
	"fmt"
	"testing"

	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
)

func TestMain(m *testing.M) {
	a := idgenerate.GenerateID()
	b := idgenerate.GenerateID()
	fmt.Println(a, b)
	if a == b {
		panic("not random")
	}
}
