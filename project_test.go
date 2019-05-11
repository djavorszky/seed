package seed

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var pwd string

const name = "admiral"

func setup() {
	var err error

	pwd, err = os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("unable to get pwd: %v", err))
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

// TestInitProject___doInit needs to be the first, alphabetically. This was the
// cleanest way I could find to do an initialization within the test framework.
func TestInitProject___doInit(t *testing.T) {
	err := InitProject(name)
	if err != nil {
		t.Errorf("InitProject(%q) failed = %v", name, err)
	}
}

func TestInitProject_foldersAreCreated(t *testing.T) {
	err := checkIfFolderExists(name)
	if err != nil {
		t.Errorf("project: %v", err)
	}

	err = checkIfFolderExists(filepath.Join(name, cmdFolder))
	if err != nil {
		t.Errorf("project/%s: %v", cmdFolder, err)
	}

	err = checkIfFolderExists(filepath.Join(name, genFolder))
	if err != nil {
		t.Errorf("project/%s: %v", genFolder, err)
	}
}

func TestInitProject_projectDescriptor(t *testing.T) {
	filename := fmt.Sprintf("%s.yml", name)
	descriptor := filepath.Join(pwd, name, filename)

	f, err := os.Stat(descriptor)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("%s does not exist", descriptor)
		}

		t.Errorf("checking %s: %v", descriptor, err)
	}

	err = checkFileIsCorrect(f)
	if err != nil {
		t.Errorf("checking %s: %v", descriptor, err)
	}
}

func TestInitProject_serviceFile(t *testing.T) {
	filename := fmt.Sprintf("%s.go", name)
	serviceFileName := filepath.Join(pwd, name, filename)

	f, err := os.Stat(serviceFileName)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("%s does not exist", serviceFileName)
		}

		t.Errorf("checking %s: %v", serviceFileName, err)
	}

	err = checkFileIsCorrect(f)
	if err != nil {
		t.Errorf("checking %s: %v", serviceFileName, err)
	}
}

func TestInitProject_interfaceFile(t *testing.T) {
	filename := "interface.go"
	interfaceFile := filepath.Join(pwd, name, genFolder, filename)

	f, err := os.Stat(interfaceFile)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("%s does not exist", interfaceFile)
		}

		t.Errorf("checking %s: %v", interfaceFile, err)
	}

	err = checkFileIsCorrect(f)
	if err != nil {
		t.Errorf("checking %s: %v", interfaceFile, err)
	}
}

func TestInitProject_generatedFile(t *testing.T) {
	filename := "generated.go"
	interfaceFile := filepath.Join(pwd, name, genFolder, filename)

	f, err := os.Stat(interfaceFile)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("%s does not exist", interfaceFile)
		}

		t.Errorf("checking %s: %v", interfaceFile, err)
	}

	err = checkFileIsCorrect(f)
	if err != nil {
		t.Errorf("checking %s: %v", interfaceFile, err)
	}
}

func TestInitProject_interfaceContents(t *testing.T) {
	fileName := "interface.go"

	generatedBytes, err :=
		ioutil.ReadFile(filepath.Join(pwd, name, genFolder, fileName))
	if err != nil {
		t.Errorf("reading generated interface file: %v", err)
	}

	exampleBytes, err :=
		ioutil.ReadFile(
			filepath.Join(pwd, "example", "admiral", genFolder, fileName))
	if err != nil {
		t.Errorf("reading example interface file: %v", err)
	}

	if !bytes.Equal(generatedBytes, exampleBytes) {
		t.Errorf(
			"generated != example:\nGENERATED:\n%v\n\nvs EXAMPLE:\n\n%v",
			string(generatedBytes), string(exampleBytes))
	}
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

func checkIfFolderExists(path string) error {
	f, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s directory does not exist", path)
		}

		return fmt.Errorf("checking %s: %v", path, err)
	}

	if !f.IsDir() {
		return fmt.Errorf("%s: not a directory", path)
	}

	return nil
}

func teardown() {
	dirName, err := getDirName(name)
	if err != nil {
		panic(fmt.Sprintf("teardown failed: %v", err))
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
