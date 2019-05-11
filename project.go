package seed

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	cmdFolder = "cmd"
	genFolder = "gen"

	initFailed = "init failed: %v"

	permWindows = 0666
	permLinux   = 0755
)

var defaultPerm os.FileMode

func init() {
	switch runtime.GOOS {
	case "windows":
		defaultPerm = permWindows
	default:
		defaultPerm = permLinux
	}
}

func InitProject(projectName string) error {
	err := createProjectStructure(projectName)
	if err != nil {
		return fmt.Errorf(initFailed, err)
	}

	return nil
}

func createProjectStructure(projectName string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = createDirs(pwd, projectName)
	if err != nil {
		return fmt.Errorf("creating folders: %v", err)
	}
	err = createFiles(pwd, projectName)
	if err != nil {
		return fmt.Errorf("creating files: %v", err)
	}

	return nil
}

func createDirs(pwd, projectName string) error {
	paths := [][]string{
		{pwd, projectName},
		{pwd, projectName, cmdFolder},
		{pwd, projectName, genFolder},
	}

	for _, path := range paths {
		fullPath := filepath.Join(path...)

		err := os.MkdirAll(fullPath, defaultPerm)
		if err != nil {
			return fmt.Errorf("creating %v failed: %v", path, err)
		}

		err = os.Chmod(fullPath, defaultPerm)
		if err != nil {
			return fmt.Errorf("changing %v permissions failed: %v", path, err)
		}
	}

	return nil
}

func createFiles(pwd, projectName string) error {
	paths := [][]string{
		{pwd, projectName, fmt.Sprintf("%v.go", projectName)},
		{pwd, projectName, fmt.Sprintf("%v.yml", projectName)},
		{pwd, projectName, cmdFolder, "main.go"},
		{pwd, projectName, genFolder, "generated.go"},
		{pwd, projectName, genFolder, "interface.go"},
	}
	for _, path := range paths {
		fullPath := filepath.Join(path...)

		f, err := os.Create(fullPath)
		if err != nil {
			return fmt.Errorf("creating %v failed: %v", path, err)
		}

		err = f.Close()
		if err != nil {
			return fmt.Errorf("closing %v failed: %v", path, err)
		}

		err = os.Chmod(fullPath, defaultPerm)
		if err != nil {
			return fmt.Errorf("updating %v permissions failed: %v", path, err)
		}

	}

	return nil
}
