package main

import (
	"fmt"
	"os"

	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/zip"
)

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("Usage: compress-construct <source> <destination>")
		return
	}

	sourceDir := args[0]
	destFile := args[1]

	err := zip.CompressZipFile(sourceDir, destFile)
	if err != nil {
		fmt.Println("Error compressing directory:", err)
		panic(err)
	}
	fmt.Println("Directory compressed successfully")
}
