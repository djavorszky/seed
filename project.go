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

	err = createDirs(
		[]string{pwd, projectName},
		[]string{pwd, projectName, cmdFolder},
		[]string{pwd, projectName, genFolder},
	)
	if err != nil {
		return fmt.Errorf("creating folders: %v", err)
	}
	err = createFiles(
		[]string{pwd, projectName, fmt.Sprintf("%v.go", projectName)},
		[]string{pwd, projectName, fmt.Sprintf("%v.yml", projectName)},
		[]string{pwd, projectName, cmdFolder, "main.go"},
		[]string{pwd, projectName, genFolder, "generated.go"},
		[]string{pwd, projectName, genFolder, "interface.go"},
	)
	if err != nil {
		return fmt.Errorf("creating files: %v", err)
	}

	return nil
}

func createDirs(paths ...[]string) error {
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

func createFiles(paths ...[]string) error {
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
