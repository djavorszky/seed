package seed

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/dave/jennifer/jen"
)

const (
	cmdFolder     = "cmd"
	genFolder     = "gen"
	interfaceFile = "interface.go"
	generatedFile = "generated.go"
	mainFile      = "main.go"

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

	err = createInterfaceFile(projectName)
	if err != nil {
		return fmt.Errorf(initFailed, err)
	}

	err = createMainFile(projectName)
	if err != nil {
		return fmt.Errorf(initFailed, err)
	}

	err = createGeneratedFile(projectName)
	if err != nil {
		return fmt.Errorf(initFailed, err)
	}

	err = createServiceDescriptor(projectName)
	if err != nil {
		return fmt.Errorf(initFailed, err)
	}

	err = formatFiles()
	if err != nil {
		return fmt.Errorf(initFailed, err)
	}

	return nil
}

func createGeneratedFile(projectName string) error {
	f := NewFile("gen")

	path := filepath.Join(pwd, projectName, genFolder, generatedFile)

	err := f.Save(path)
	if err != nil {
		return fmt.Errorf("saving file: %v", err)
	}

	return nil
}

func createServiceDescriptor(projectName string) error {
	// for now, only create the file

	path := filepath.Join(pwd, projectName, projectName+".yml")

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed creating project descriptor: %v", err)
	}
	f.Close()

	return nil
}

func createMainFile(projectName string) error {
	f := NewFile("main")

	f.Func().Id("main").Params().Block(
		Id("service").Op(":=").
			Qual(fmt.Sprintf("%s/gen", projectName), "New").
			Call(
				Op("&").Qual("admiral", "Server").Values(),
			),
		Qual("log", "Fatal").Call(
			Qual("net/http", "ListenAndServe").Call(
				Lit(":8080"), Id("service"),
			),
		),
	)

	path := filepath.Join(pwd, projectName, cmdFolder, mainFile)

	err := f.Save(path)
	if err != nil {
		return fmt.Errorf("saving file: %v", err)
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

	err = ioutil.WriteFile(path, contents, defaultPerm)
	if err != nil {
		return fmt.Errorf("failed writing file %s: %v",
			path, err)
	}

	return nil
}

func createInterfaceFile(projectName string) error {
	f := NewFile("gen")

	title := strings.Title(projectName)
	service := fmt.Sprintf("%sService", title)
	handler := fmt.Sprintf("%sHandler", title)
	middleware := fmt.Sprintf("%sMiddleware", title)

	f.Commentf("// %s encapsulates the handler interface, which "+
		"holds all the methods to be called\n// by the server, and middleware "+
		"interface, which contains all the middlewares to be added to the service.",
		service)

	f.Type().Id(service).Interface(
		Id(handler),
		Id(middleware),
	)

	f.Commentf("// %s is the interface for the handlers. Any new "+
		"endpoint added by seed will be added here as a\n// new method on "+
		"the interface.", handler)

	f.Type().Id(handler).Interface(
		Id("Index").Params().Qual("net/http", "HandlerFunc"),
	)

	f.Commentf("// %s is the interface for all the middlewares that "+
		"will be added to all of the paths.", middleware)

	f.Type().Id(middleware).Interface(
		Id("LoggerMw").Params(
			Qual("net/http", "Handler"),
		).Qual("net/http", "Handler"),
	)

	path := filepath.Join(pwd, projectName, genFolder, interfaceFile)

	err := f.Save(path)
	if err != nil {
		return fmt.Errorf("saving file: %v", err)
	}

	return nil
}

func createProjectStructure(projectName string) error {
	err := createDirs(pwd, projectName)
	if err != nil {
		return fmt.Errorf("creating folders: %v", err)
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
