package seed

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

const (
	cmdFolder     = "cmd"
	genFolder     = "gen"
	interfaceFile = "interface.go"

	initFailed = "init failed: %v"

	permWindows = 0666
	permLinux   = 0755
)

var (
	defaultPerm os.FileMode
	pwd         string
)

func init() {
	switch runtime.GOOS {
	case "windows":
		defaultPerm = permWindows
	default:
		defaultPerm = permLinux
	}

	var err error

	pwd, err = os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("setting up pwd: %v", err))
	}
}

func InitProject(projectName string) error {
	err := createProjectStructure(projectName)
	if err != nil {
		return fmt.Errorf(initFailed, err)
	}

	err = setupInterfaceFile(projectName)
	if err != nil {
		return fmt.Errorf(initFailed, err)
	}

	err = formatFiles()
	if err != nil {
		return fmt.Errorf(initFailed, err)
	}

	return nil
}

func formatFiles() error {
	err := filepath.Walk(pwd, formatFile)
	if err != nil {
		return fmt.Errorf("source format: %v", err)
	}

	return nil
}

func formatFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	if !strings.HasSuffix(info.Name(), ".go") {
		return nil
	}

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed reading file %s: %v",
			path, err)
	}

	contents, err = format.Source(contents)
	if err != nil {
		return fmt.Errorf("failed formatting file %s: %v",
			path, err)
	}

	err = ioutil.WriteFile(path, contents, defaultPerm)
	if err != nil {
		return fmt.Errorf("failed writing file %s: %v",
			path, err)
	}

	return nil
}

func setupInterfaceFile(projectName string) error {
	path := filepath.Join(pwd, projectName, genFolder, interfaceFile)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, defaultPerm)
	if err != nil {
		return fmt.Errorf("failed opening %s: %v", interfaceFile, err)
	}
	defer f.Close()

	templateFile := "interface.tmpl"
	templatePath := filepath.Join(pwd, "tmpl", templateFile)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed parsing template %s: %v", templateFile, err)
	}

	info := InterfaceStruct{
		PackageName: genFolder,
		ServiceName: strings.Title(projectName),
		Imports:     []string{"net/http"},
		Middlewares: []string{"LoggerMw(http.Handler) http.Handler"},
		Services:    []string{"Index() http.HandlerFunc"},
	}

	err = tmpl.Execute(f, info)
	if err != nil {
		return fmt.Errorf("failed executing template: %v", err)
	}

	return nil
}

func createProjectStructure(projectName string) error {
	err := createDirs(pwd, projectName)
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
