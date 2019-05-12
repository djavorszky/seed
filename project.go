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

	err = createServiceFile(projectName)
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

func createServiceFile(projectName string) error {
	f := NewFilePath(projectName)

	f.Type().Id("Server").Struct()

	f.Func().Params(
		Id("s").Op("*").Id("Server"),
	).Id("LoggerMw").Params(
		Id("next").Qual("net/http", "Handler"),
	).Qual("net/http", "Handler").Block(
		Comment("// Anything you add here will be executed once, during startup."),
		Comment("// The returned http.Handler will be able to access these variables"),
		Comment("// thanks to closure."),
		Line(),
		Id("prefix").Op(":=").Lit(fmt.Sprintf("[%s] - ", projectName)),
		Return(
			Qual("net/http", "HandlerFunc").Call(
				Func().Params(
					Id("w").Qual("net/http", "ResponseWriter"),
					Id("r").Op("*").Qual("net/http", "Request"),
				).Block(
					Qual("log", "Println").Call(
						Id("prefix"),
						Id("r").Dot("RemoteAddr"),
						Id("r").Dot("Method"),
						Id("r").Dot("RequestURI"),
					),
					Line(),
					Id("next").Dot("ServeHTTP").Call(
						Id("w"), Id("r"),
					),
				),
			),
		),
	)

	f.Line()

	f.Func().Params(
		Id("s").Op("*").Id("Server"),
	).Id("Index").Params().Qual("net/http", "HandlerFunc").Block(
		Comment("// Anything you add here will be executed once, during startup."),
		Comment("// The returned http.HandlerFunc will be able to access these variables"),
		Comment("// thanks to closure."),
		Line(),
		Id("defaultMsg").Op(":=").Index().Byte().Call(Lit("I'm alive!")),
		Return().Func().Params(
			Id("w").Qual("net/http", "ResponseWriter"),
			Id("r").Op("*").Qual("net/http", "Request"),
		).Block(
			Id("w").Dot("WriteHeader").Call(Qual("net/http", "StatusOK")),
			Line(),
			List(
				Id("_"), Id("err"),
			).Op(":=").Id("w").Dot("Write").Call(Id("defaultMsg")),
			If(Id("err").Op("!=").Nil()).Block(
				Panic(Id("err")),
			),
		),
	)

	path := filepath.Join(pwd, projectName, projectName+".go")

	err := f.Save(path)
	if err != nil {
		return fmt.Errorf("saving file: %v", err)
	}

	return nil
}

func createGeneratedFile(projectName string) error {
	f := NewFilePath("gen")

	const mux = "github.com/gorilla/mux"

	f.Comment("// Service is the struct that will be exposed to serve HTTP traffic.")
	f.Type().Id("Service").Struct(
		Id("router").Op("*").Qual(mux, "Router"),
		Id("serviceImpl").Qual("gen", "AdmiralService"),
	)

	f.Comment("// ServeHTTP is what ultimately allows this service to be " +
		"used by the standard library's\n// listen and serve functions")
	f.Func().Params(
		Id("s").Op("*").Id("Service"),
	).Id("ServeHTTP").Params(
		Id("w").Qual("net/http", "ResponseWriter"),
		Id("r").Op("*").Qual("net/http", "Request"),
	).Block(
		Id("s").Dot("router").Dot("ServeHTTP").Call(
			Id("w"),
			Id("r"),
		),
	)

	f.Comment("// Route is a struct that holds the path, the handler to " +
		"be called when that path is hit, and\n// which list of methods it should serve.")

	f.Type().Id("route").Struct(
		Id("path").String(),
		Id("handler").Qual("net/http", "HandlerFunc"),
		Id("methods").Index().String(),
	)

	f.Comment("// New returns a new service implementation, using the " +
		"service as a dependency. It also sets up the routes\n// and the middlewares.")

	f.Func().Id("New").Params(
		Id("service").Qual("gen", "AdmiralService"),
	).Op("*").Qual("gen", "Service").Block(
		Id("s").Op(":=").Op("&").Qual("gen", "Service").Values(
			Dict{
				Id("router"):      Qual(mux, "NewRouter").Call(),
				Id("serviceImpl"): Id("service"),
			},
		),
		Empty(),
		Id("s").Dot("routes").Call(),
		Id("s").Dot("middlewares").Call(),
		Empty(),
		Return(Id("s")),
	)

	f.Comment("// routes sets up the routes to be served by the service")
	f.Func().Params(
		Id("s").Op("*").Id("Service"),
	).Id("routes").Params().Block(
		Id("routes").Op(":=").Index().Id("route").Values(
			Block(
				Dict{
					Id("handler"): Id("s").Dot("serviceImpl").Dot("Index").Call(),
					Id("methods"): Index().String().Values(Qual("net/http", "MethodGet")),
					Id("path"):    Lit("/"),
				},
			),
		),
		Empty(),
		For(
			List(Id("_"), Id("route")).Op(":=").Range().Id("routes").Block(
				Id("s").Dot("router").
					Dot("StrictSlash").Call(Lit(true)).
					Dot("HandleFunc").
					Call(
						Id("route").Dot("path"), Id("route").Dot("handler")).
					Dot("Methods").
					Call(
						Id("route").Dot("methods").Op("..."),
					),
			),
		),
	)

	f.Comment("// middlewares sets up the middlewares to be set up by the service")
	f.Func().Params(
		Id("s").Op("*").Id("Service"),
	).Id("middlewares").Params().Block(
		Id("mws").Op(":=").Index().Qual(mux, "MiddlewareFunc").Values(
			Id("s").Dot("serviceImpl").Dot("LoggerMw"),
		),
		Empty(),
		For(
			List(Id("_"), Id("mw")).Op(":=").Range().Id("mws").Block(
				Id("s").Dot("router").Dot("Use").Call(Id("mw")),
			),
		),
	)

	path := filepath.Join(pwd, projectName, genFolder, generatedFile)

	err := f.Save(path)
	if err != nil {
		return fmt.Errorf("saving file: %v", err)
	}

	return nil
}

func createServiceDescriptor(projectName string) error {
	// TODO: for now, only create the file

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
