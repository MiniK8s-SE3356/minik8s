package zip_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/zip"
)

func TestMain(m *testing.M) {

	dirPath := "exampleDir"
	filePath := "exampleFile.txt"
	message := "Hello, World!"

	err := os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		println("Error creating directory")
		panic(err)
	}

	err = os.Mkdir(filepath.Join(dirPath, "subDir"), os.ModePerm)
	if err != nil {
		println("Error creating sub directory")
		panic(err)
	}

	file, err := os.Create(filepath.Join(dirPath, filePath))
	if err != nil {
		println("Error creating file")
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(message)
	if err != nil {
		println("Error writing to file")
		panic(err)
	}

	file2, err := os.Create(filepath.Join(dirPath, "subDir", filePath))
	if err != nil {
		println("Error creating file")
		panic(err)
	}
	defer file2.Close()

	_, err = file2.WriteString(message)
	if err != nil {
		println("Error writing to file")
		panic(err)
	}

	err = zip.CompressZipFile(dirPath, "example.zip")
	if err != nil {
		println("Error compressing zip file")
		panic(err)
	}

	err = os.RemoveAll(dirPath)
	if err != nil {
		println("Error removing directory")
		panic(err)
	}

	err = zip.DecompressZipFile("example.zip", "exampleDir")
	if err != nil {
		println("Error decompressing zip file")
		panic(err)
	}

	err = os.RemoveAll("example.zip")
	if err != nil {
		println("Error removing zip file")
		panic(err)
	}

	err = os.RemoveAll("exampleDir")
	if err != nil {
		println("Error removing directory")
		panic(err)
	}
}
