package utils

import (
	"os"
	"path/filepath"

	"github.com/peterbourgon/diskv/v3"
)

var PresistentStorage *diskv.Diskv = NewStorage(".zerok")

func NewStorage(dirPath string) *diskv.Diskv {
	var err error

	var baseDir string
	if baseDir, err = os.UserHomeDir(); err != nil {
		baseDir = os.TempDir()
	}

	diskv := diskv.New(diskv.Options{
		BasePath:  filepath.Join(baseDir, dirPath),
		Transform: func(s string) []string { return []string{} },
	})

	return diskv
}
