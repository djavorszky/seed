package seed

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestInitProject(t *testing.T) {
	name := "testProject"

	defer teardownTestProject(t, name)

	err := InitProject(name)
	if err != nil {
		t.Errorf("InitProject(%q) failed = %v", name, err)
		return
	}

	actualName, err := getDirName(name)
	if err != nil {
		t.Errorf("InitProject(%q): checking result failed = %v", name, err)
		return
	}

	// check if project folder is created
	err = checkProjectCreated(actualName)
	if err != nil {
		t.Errorf("project: %v", err)
	}

	// check if foldername.go is correct
	err = checkServiceFile(actualName)
	if err != nil {
		t.Errorf("service file: %v", err)
	}

}

func checkProjectCreated(projectName string) error {
	f, err := os.Stat(projectName)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s folder does not exist", projectName)
		}

		return fmt.Errorf("checking %s: %v", projectName, err)
	}

	fileMode := f.Mode()
	if !fileMode.IsDir() {
		return fmt.Errorf("%s: is not a directory", projectName)
	}

	filePerm := fileMode.Perm()
	if filePerm != 0755 {
		return fmt.Errorf("%s: expected fileperm 0755, encountered: %v",
			projectName, filePerm)
	}

	return nil
}

func checkServiceFile(projectName string) error {
	serviceFileName := fmt.Sprintf("%s.go", projectName)

	f, err := os.Stat(serviceFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s does not exist", serviceFileName)
		}

		return fmt.Errorf("checking %s: %v", serviceFileName, err)
	}

	fileMode := f.Mode()
	if fileMode.IsDir() {
		return fmt.Errorf("%s: is a directory", serviceFileName)
	}

	filePerm := fileMode.Perm()
	if filePerm != 0755 {
		return fmt.Errorf("%s: expected fileperm 0755, encountered: %v",
			serviceFileName, filePerm)
	}

	return nil
}

func teardownTestProject(t *testing.T, name string) {
	dirName, err := getDirName(name)
	if err != nil {
		t.Errorf("teardown failed: %v", err)
		return
	}

	err = os.RemoveAll(dirName)
	if err != nil {
		panic(fmt.Sprintf("failed removing all: %v", err))
	}
}

func getDirName(name string) (string, error) {
	if name != "." {
		return name, nil
	}

	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed project setup: %v", err)
	}

	return filepath.Base(pwd), nil
}
