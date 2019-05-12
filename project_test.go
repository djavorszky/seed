package seed

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			return
		}

		t.Errorf("checking %s: %v", descriptor, err)
		return
	}

	err = checkFileIsCorrect(f)
	if err != nil {
		t.Errorf("checking %s: %v", descriptor, err)
	}

	t.Errorf("test is not complete yet.")
}

func TestInitProject_serviceFile(t *testing.T) {
	filename := fmt.Sprintf("%s.go", name)

	generatedBytes, err :=
		ioutil.ReadFile(filepath.Join(pwd, name, filename))
	if err != nil {
		t.Errorf("reading generated service file: %v", err)
	}

	exampleBytes, err := ioutil.ReadFile(
		filepath.Join(pwd, "testdata", "service.expected"))
	if err != nil {
		t.Errorf("reading example service file: %v", err)
	}

	genString := string(generatedBytes)
	exampleString := string(exampleBytes)

	genString = strings.ReplaceAll(genString, "\r\n", "\n")
	exampleString = strings.ReplaceAll(exampleString, "\r\n", "\n")

	assert.Equal(t, exampleString, genString)
}

func TestInitProject_interfaceFile(t *testing.T) {
	interfaceFile := filepath.Join(pwd, name, genFolder, interfaceFile)

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

	t.Errorf("test is not complete yet.")

}

func TestInitProject_generatedFile(t *testing.T) {
	filename := "generated.go"
	interfaceFile := filepath.Join(pwd, name, genFolder, filename)

	f, err := os.Stat(interfaceFile)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("%s does not exist", interfaceFile)
			return
		}

		t.Errorf("checking %s: %v", interfaceFile, err)
	}

	err = checkFileIsCorrect(f)
	if err != nil {
		t.Errorf("checking %s: %v", interfaceFile, err)
	}
}

func TestInitProject_interfaceContents(t *testing.T) {
	generatedBytes, err :=
		ioutil.ReadFile(filepath.Join(pwd, name, genFolder, interfaceFile))
	if err != nil {
		t.Errorf("reading generated interface file: %v", err)
	}

	exampleBytes, err := ioutil.ReadFile(
		filepath.Join(pwd, "testdata", "interface.expected"))
	if err != nil {
		t.Errorf("reading example interface file: %v", err)
	}

	genString := string(generatedBytes)
	exampleString := string(exampleBytes)

	genString = strings.ReplaceAll(genString, "\r\n", "\n")
	exampleString = strings.ReplaceAll(exampleString, "\r\n", "\n")

	assert.Equal(t, exampleString, genString)
}

func TestInitProject_executableContents(t *testing.T) {
	generatedBytes, err := ioutil.ReadFile(
		filepath.Join(pwd, name, cmdFolder, mainFile))
	if err != nil {
		t.Errorf("reading generated main.go: %v", err)
	}

	exampleBytes, err := ioutil.ReadFile(
		filepath.Join(pwd, "testdata", "main.expected"))
	if err != nil {
		t.Errorf("reading example main.go: %v", err)
	}

	genString := string(generatedBytes)
	exampleString := string(exampleBytes)

	genString = strings.ReplaceAll(genString, "\r\n", "\n")
	exampleString = strings.ReplaceAll(exampleString, "\r\n", "\n")

	assert.Equal(t, exampleString, genString)
}

func TestInitProject_generatedContents(t *testing.T) {
	generatedBytes, err :=
		ioutil.ReadFile(filepath.Join(pwd, name, genFolder, generatedFile))
	if err != nil {
		t.Errorf("reading generated interface file: %v", err)
	}

	exampleBytes, err := ioutil.ReadFile(
		filepath.Join(pwd, "testdata", "generated.expected"))
	if err != nil {
		t.Errorf("reading example interface file: %v", err)
	}

	genString := string(generatedBytes)
	exampleString := string(exampleBytes)

	genString = strings.ReplaceAll(genString, "\r\n", "\n")
	exampleString = strings.ReplaceAll(exampleString, "\r\n", "\n")

	assert.Equal(t, exampleString, genString)
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
