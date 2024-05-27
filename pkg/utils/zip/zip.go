package zip

import (
	archive_tar "archive/tar"
	archive_zip "archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

func CompressZipFile(source, target string) error {
	zipFile, err := os.Create(target)
	if err != nil {
		fmt.Println("Error creating zip file")
		return err
	}
	defer zipFile.Close()

	zipWriter := archive_zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error walking through source directory")
			return err
		}

		relativePath, err := filepath.Rel(source, path)
		if err != nil {
			fmt.Println("Error getting relative path")
			return err
		}

		if info.IsDir() {
			_, err := zipWriter.Create(relativePath + "/")
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			fmt.Println("Error opening file")
			return err
		}
		defer file.Close()

		zipFileWriter, err := zipWriter.Create(relativePath)
		if err != nil {
			fmt.Println("Error creating file in zip file")
			return err
		}

		_, err = io.Copy(zipFileWriter, file)
		return err
	})

	if err != nil {
		fmt.Println("Error compressing zip file")
	}

	fmt.Println("Compressing zip file successfully")

	return err
}

func DecompressZipFile(zipFile, targetDir string) error {
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

func CompressTarFile(source, target string) error {
	tarFile, err := os.Create(target)
	if err != nil {
		fmt.Println("Error creating tar file")
		return err
	}
	defer tarFile.Close()

	tarWriter := archive_tar.NewWriter(tarFile)
	defer tarWriter.Close()

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error walking through source directory")
			return err
		}

		header, err := archive_tar.FileInfoHeader(info, "")
		if err != nil {
			fmt.Println("Error getting file info header")
			return err
		}

		relPath, err := filepath.Rel(source, path)
		if err != nil {
			fmt.Println("Error getting relative path")
			return err
		}
		header.Name = relPath

		if err := tarWriter.WriteHeader(header); err != nil {
			fmt.Println("Error writing header to tar file")
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			fmt.Println("Error opening file")
			return err
		}
		defer file.Close()

		_, err = io.Copy(tarWriter, file)
		if err != nil {
			fmt.Println("Error copying file to tar file")
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error compressing tar file")
		return err
	}

	fmt.Println("Compressing tar file successfully")

	return nil
}
