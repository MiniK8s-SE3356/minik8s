package zip_test

import (
	"testing"

	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/zip"
)

func TestMain(m *testing.M) {
	err := zip.DecompressZipFile("test.zip", "test")
	if err != nil {
		println("Error decompressing zip file")
		panic(err)
	}
}
