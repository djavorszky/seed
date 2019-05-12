package seed

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"seed/consts"
	"seed/descriptor"
	"seed/files"
	"strings"
	"testing"

	"github.com/go-yaml/yaml"

	"github.com/stretchr/testify/assert"
)

const name = "example2"

func TestMain(m *testing.M) {
	//setup()
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

	err = checkIfFolderExists(filepath.Join(name, consts.CmdFolder))
	if err != nil {
		t.Errorf("project/%s: %v", consts.CmdFolder, err)
	}

	err = checkIfFolderExists(filepath.Join(name, consts.GenFolder))
	if err != nil {
		t.Errorf("project/%s: %v", consts.GenFolder, err)
	}
}

func TestInitProject_executableFile(t *testing.T) {
	mainPath := filepath.Join(files.Pwd, name, consts.CmdFolder, consts.MainFile)

	f, err := os.Stat(mainPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("%s does not exist", mainPath)
			return
		}

		t.Errorf("checking %s: %v", mainPath, err)
	}

	err = checkFileIsCorrect(f)
	if err != nil {
		t.Errorf("checking %s: %v", mainPath, err)
	}
}

func TestInitProject_executableContents(t *testing.T) {
	path := filepath.Join(files.Pwd, name, consts.CmdFolder, consts.MainFile)

	actual, err := readFile(path)
	if err != nil {
		t.Errorf("reading result file for %q: %v", "go.mod", err)
	}

	expected, err := parseExpected("main.expected", name)
	if err != nil {
		t.Errorf("parsing expected file: %v", err)
	}

	assert.Equal(t, expected, actual)
}

func TestInitProject_generatedFile(t *testing.T) {
	filename := "generated.go"
	interfaceFile := filepath.Join(files.Pwd, name, consts.GenFolder, filename)

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

func TestInitProject_generatedContents(t *testing.T) {
	path := filepath.Join(files.Pwd, name, consts.GenFolder, consts.GeneratedFile)

	actual, err := readFile(path)
	if err != nil {
		t.Errorf("reading result file for %q: %v", "go.mod", err)
	}

	expected, err := parseExpected("generated.expected", name)
	if err != nil {
		t.Errorf("parsing expected file: %v", err)
	}

	assert.Equal(t, expected, actual)
}

func TestInitProject_projectDescriptor(t *testing.T) {
	filename := fmt.Sprintf("%s.yml", name)
	d := filepath.Join(files.Pwd, name, filename)

	f, err := os.Stat(d)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("%s does not exist", d)
			return
		}

		t.Errorf("checking %s: %v", d, err)
		return
	}

	err = checkFileIsCorrect(f)
	if err != nil {
		t.Errorf("checking %s: %v", d, err)
	}
}

func TestInitProject_projectDescriptorContents(t *testing.T) {
	filename := fmt.Sprintf("%s.yml", name)
	d := filepath.Join(files.Pwd, name, filename)

	b, err := ioutil.ReadFile(d)
	if err != nil {
		t.Errorf("failed reading descriptor: %v", err)
	}

	var sd descriptor.ServiceDescriptor

	err = yaml.Unmarshal(b, &sd)
	if err != nil {
		t.Errorf("failed unmarshalling descriptor: %v", err)
	}

	info := descriptor.Info{
		Name:    name,
		Summary: "just a test for now",
	}

	expected := descriptor.Base(info)

	assert.Equal(t, expected, sd)
}

func TestInitProject_goModFile(t *testing.T) {
	d := filepath.Join(files.Pwd, name, "go.mod")

	f, err := os.Stat(d)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("%s does not exist", d)
			return
		}

		t.Errorf("checking %s: %v", d, err)
		return
	}

	err = checkFileIsCorrect(f)
	if err != nil {
		t.Errorf("checking %s: %v", d, err)
	}
}
func TestInitProject_goModFileContents(t *testing.T) {
	path := filepath.Join(files.Pwd, name, "go.mod")

	actual, err := readFile(path)
	if err != nil {
		t.Errorf("reading result file for %q: %v", "go.mod", err)
	}

	expected, err := parseExpected("gomod.expected", name)
	if err != nil {
		t.Errorf("parsing expected file: %v", err)
	}

	assert.Equal(t, expected, actual)
}

func TestInitProject_serviceFile(t *testing.T) {
	filename := fmt.Sprintf("%s.go", name)
	serviceFile := filepath.Join(files.Pwd, name, filename)

	f, err := os.Stat(serviceFile)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("%s does not exist", serviceFile)
		}

		t.Errorf("checking %s: %v", serviceFile, err)
	}

	err = checkFileIsCorrect(f)
	if err != nil {
		t.Errorf("checking %s: %v", serviceFile, err)
	}
}

func TestInitProject_serviceFileContents(t *testing.T) {
	path := filepath.Join(files.Pwd, name, name+".go")

	actual, err := readFile(path)
	if err != nil {
		t.Errorf("reading result file for %q: %v", "go.mod", err)
	}

	expected, err := parseExpected("service.expected", name)
	if err != nil {
		t.Errorf("parsing expected file: %v", err)
	}

	assert.Equal(t, expected, actual)
}

func TestInitProject_interfaceFile(t *testing.T) {
	interfaceFile := filepath.Join(files.Pwd, name, consts.GenFolder, consts.InterfaceFile)

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
	path := filepath.Join(files.Pwd, name, consts.GenFolder, consts.InterfaceFile)

	actual, err := readFile(path)
	if err != nil {
		t.Errorf("reading result file for %q: %v", "go.mod", err)
	}

	expected, err := parseExpected("interface.expected", name)
	if err != nil {
		t.Errorf("parsing expected file: %v", err)
	}

	assert.Equal(t, expected, actual)
}

func checkFileIsCorrect(f os.FileInfo) error {
	fileMode := f.Mode()
	if fileMode.IsDir() {
		return fmt.Errorf("is a directory")
	}

	filePerm := fileMode.Perm()
	if filePerm != files.DefaultPerm {
		return fmt.Errorf("expected fileperm %v, got: %v",
			files.DefaultPerm, filePerm)
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

func parseExpected(filename, projectName string) (string, error) {
	path := filepath.Join(files.Pwd, "testdata", filename)

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return "", fmt.Errorf("failed parsing template %q: %v", filename, err)
	}

	details := struct {
		Package string
		Title   string
	}{
		Package: projectName,
		Title:   strings.Title(projectName),
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, details)
	if err != nil {
		return "", fmt.Errorf("failed parsing template %q: %v", filename, err)
	}

	result := strings.ReplaceAll(buf.String(), "\r\n", "\n")

	return result, nil
}

func readFile(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed reading %q: %v", path, err)
	}

	return strings.ReplaceAll(string(b), "\r\n", "\n"), nil
}
