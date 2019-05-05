package seed

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var pwd string

func init() {
	var err error

	pwd, err = os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("unable to get pwd: %v", err))
	}
}

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

	err = checkProjectDescriptor(actualName)
	if err != nil {
		t.Errorf("project descriptor: %v", err)
	}

	// check if foldername.go is correct
	err = checkServiceFile(actualName)
	if err != nil {
		t.Errorf("service file: %v", err)
	}

}

func checkProjectDescriptor(projectName string) error {
	filename := fmt.Sprintf("%s.yml", projectName)
	descriptor := filepath.Join(pwd, projectName, filename)

	f, err := os.Stat(descriptor)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s does not exist", descriptor)
		}

		return fmt.Errorf("checking %s: %v", descriptor, err)
	}

	err = checkFileIsCorrect(f)
	if err != nil {
		return fmt.Errorf("checking %s: %v", descriptor, err)
	}

	return nil
}

func checkFileIsCorrect(f os.FileInfo) error {
	fileMode := f.Mode()
	if fileMode.IsDir() {
		return fmt.Errorf("is a directory")
	}
	filePerm := fileMode.Perm()
	if filePerm != defaultPerm {
		return fmt.Errorf("expected fileperm %v, got: %v",
			defaultPerm, filePerm)
	}

	return nil
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
	if filePerm != defaultPerm {
		return fmt.Errorf("%s: expected fileperm %v, encountered: %v",
			projectName, defaultPerm, filePerm)
	}

	return nil
}

func checkServiceFile(projectName string) error {
	filename := fmt.Sprintf("%s.go", projectName)
	serviceFileName := filepath.Join(pwd, projectName, filename)

	f, err := os.Stat(serviceFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s does not exist", serviceFileName)
		}

		return fmt.Errorf("checking %s: %v", serviceFileName, err)
	}

	err = checkFileIsCorrect(f)
	if err != nil {
		return fmt.Errorf("checking %s: %v", serviceFileName, err)
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
