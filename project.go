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

	permWindows = 0777
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

var pwd string

func InitProject(projectName string) error {
	var err error

	pwd, err = os.Getwd()
	if err != nil {
		return fmt.Errorf(initFailed, err)
	}

	err = createDirs(defaultPerm,
		[]string{pwd, projectName, cmdFolder},
		[]string{pwd, projectName, genFolder},
	)
	if err != nil {
		return fmt.Errorf(initFailed, err)
	}

	return nil
}

func createDirs(mode os.FileMode, paths ...[]string) error {
	for _, path := range paths {
		fullPath := filepath.Join(path...)

		err := os.MkdirAll(fullPath, defaultPerm)
		if err != nil {
			return fmt.Errorf("creating %v failed: %v", path, err)
		}
	}

	return nil
}
