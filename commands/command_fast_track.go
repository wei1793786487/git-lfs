package commands

import (
	"github.com/git-lfs/git-lfs/v3/tr"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var maxFileSize int64 = 10 * 1024 * 1024 //10mb

type FileDetail struct {
	Path     string
	Size     int64
	FileName string
}

func getLfsTrackedPatterns() {

}

func getAllFileSizeMax(dir string) ([]FileDetail, error) {
	var files []FileDetail
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && strings.HasSuffix(path, "/.git") {
			return filepath.SkipDir
		}
		if !info.IsDir() && info.Size() > maxFileSize {
			files = append(files, FileDetail{
				Path:     path,
				Size:     info.Size(),
				FileName: filepath.Base(path),
			})
		}
		return nil

	})
	if err != nil {
		Print(tr.Tr.Get("error: %s", err.Error()))
	}
	return files, nil
}

func fastTrackCommand(cmd *cobra.Command, args []string) {
	Print(tr.Tr.Get("Start fast track"))
	workingDir := cfg.LocalWorkingDir()
	fileList, err := getAllFileSizeMax(workingDir)
	if err != nil {
		Print(tr.Tr.Get("error: %s", err.Error()))
	}
	if len(fileList) == 0 {
		Print(tr.Tr.Get("No files need track"))
		return
	}
	for _, file := range fileList {
		//print file list
		Print(tr.Tr.Get("file: %s, size:%.2f MB", file.Path, float64(file.Size)/(1024*1024)))

	}
}

func init() {
	RegisterCommand("fast track", fastTrackCommand, func(cmd *cobra.Command) {

	})
}
