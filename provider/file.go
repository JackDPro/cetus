package provider

import (
	"io"
	"os"
	"path/filepath"
)

func CopyFile(srcFile, dstFile string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return dst.Sync()
}

func CopyDir(srcDir, dstDir string) error {
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		err := os.MkdirAll(dstDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())

		if entry.IsDir() {
			// 如果是目录，递归复制
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// 如果是文件，复制文件
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}
