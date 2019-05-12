package seed

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"seed/files"
	"seed/generate"
	"strings"
)

const (
	initFailed = "init failed: %v"
)

func InitProject(projectName string) error {
	err := generate.ProjectStructure(projectName)
	if err != nil {
		return fmt.Errorf(initFailed, err)
	}

	err = generate.CreateInterfaceFile(projectName)
	if err != nil {
		return fmt.Errorf(initFailed, err)
	}

	err = generate.CreateServiceFile(projectName)
	if err != nil {
		return fmt.Errorf(initFailed, err)
	}

	err = generate.CreateMainFile(projectName)
	if err != nil {
		return fmt.Errorf(initFailed, err)
	}

	err = generate.CreateGeneratedFile(projectName)
	if err != nil {
		return fmt.Errorf(initFailed, err)
	}

	err = generate.CreateServiceDescriptor(projectName)
	if err != nil {
		return fmt.Errorf(initFailed, err)
	}

	err = generate.GoModule(projectName)
	if err != nil {
		return fmt.Errorf(initFailed, err)
	}

	err = formatFiles(projectName)
	if err != nil {
		return fmt.Errorf(initFailed, err)
	}

	return nil
}

func formatFiles(projectName string) error {
	project := filepath.Join(files.Pwd, projectName)

	err := filepath.Walk(project, formatFile)
	if err != nil {
		return fmt.Errorf("source format: %v", err)
	}

	return nil
}

func formatFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
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

	err = ioutil.WriteFile(path, contents, files.DefaultPerm)
	if err != nil {
		return fmt.Errorf("failed writing file %s: %v",
			path, err)
	}

	return nil
}
