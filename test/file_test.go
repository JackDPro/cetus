package test

import (
	"github.com/JackDPro/cetus/provider"
	"testing"
)

func TestFileCopyFile(t *testing.T) {
	err := provider.CopyFile("./data/1.txt", "./data/2.txt")
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestFileCopyDir(t *testing.T) {
	err := provider.CopyDir("./data", "./data1")
	if err != nil {
		t.Fatal(err)
		return
	}
}
