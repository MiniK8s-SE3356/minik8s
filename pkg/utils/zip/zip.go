package zip

import (
	archive_zip "archive/zip"
	"fmt"
	"io"
	"os"
	"path"
)

func DecompressZipFile(zipFile string, targetDir string) error {
	r, err := archive_zip.OpenReader(zipFile)
	if err != nil {
		fmt.Println("Error opening zip file")
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := path.Join(targetDir, f.Name)

		if f.FileInfo().IsDir() {
			err := os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				fmt.Println("Error creating directory for decompressing zip file")
				return err
			}
			continue
		}

		if err := os.MkdirAll(path.Dir(fpath), os.ModePerm); err != nil {
			fmt.Println("Error creating directory for decompressing zip file")
			return err
		}

		inFile, err := f.Open()
		if err != nil {
			fmt.Println("Error opening file in zip file")
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			fmt.Println("Error opening file for writing in zip file")
			return err
		}

		_, err = io.Copy(outFile, inFile)
		if err != nil {
			fmt.Println("Error copying file in zip file")
			return err
		}

		inFile.Close()
		outFile.Close()
	}

	fmt.Println("Decompressing zip file successfully")

	return nil
}
