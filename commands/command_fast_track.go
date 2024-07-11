package commands

import (
	"fmt"
	"github.com/git-lfs/git-lfs/v3/tr"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	maxFileSize int64 = 10 * 1024 * 1024
	force             = false
) //10mb

type FileDetail struct {
	Path     string
	Size     int64
	FileName string
}

func backupAndReinitializeGit(workingDir string) error {
	// 获取当前 Git 仓库的远程仓库信息
	cmd := exec.Command("git", "remote", "-v")
	cmd.Dir = workingDir
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error getting git remotes: %s", err)
	}
	remotes := strings.TrimSpace(string(output))

	// 删除 .git 目录
	err = os.RemoveAll(filepath.Join(workingDir, ".git"))
	if err != nil {
		return fmt.Errorf("error removing .git directory: %s", err)
	}

	// 重新初始化 Git 仓库
	cmd = exec.Command("git", "init")
	cmd.Dir = workingDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error initializing git repository: %s", err)
	}

	// 恢复远程仓库信息
	remoteLines := strings.Split(remotes, "\n")
	for _, line := range remoteLines {
		parts := strings.Split(line, "\t")
		remoteName := strings.Split(parts[0], " ")[0]
		remoteUrl := strings.Split(parts[1], " ")[0]
		cmd = exec.Command("git", "remote", "add", remoteName, remoteUrl)
		cmd.Dir = workingDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error adding remote %s: %s", remoteName, err)
		}
	}
	return nil
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
	if force {
		Print(tr.Tr.Get("Forcing reinitialization of Git repository..."))
		workingDir := cfg.LocalWorkingDir()
		if err := backupAndReinitializeGit(workingDir); err != nil {
			Print(tr.Tr.Get("error: %s", err.Error()))
			return
		}
		Print(tr.Tr.Get("Repository reinitialized successfully."))
	}
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
		cmd.Flags().BoolVarP(&force, "f", "f", false, "")
	})
}
