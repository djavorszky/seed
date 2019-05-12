package generate

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"seed/consts"
	"seed/descriptor"
	"seed/files"
	"strings"

	"github.com/go-yaml/yaml"

	. "github.com/dave/jennifer/jen"
)

func ServiceFile(projectName string) ([]byte, error) {
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
				Add(httpHandlerFunc()).Block(
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
		Return().Add(
			httpHandlerFunc().Block(
				Id("w").Dot("WriteHeader").Call(Qual("net/http", "StatusOK")),
				Line(),
				List(
					Id("_"), Id("err"),
				).Op(":=").Id("w").Dot("Write").Call(Id("defaultMsg")),
				If(Id("err").Op("!=").Nil()).Block(
					Panic(Id("err")),
				),
			),
		),
	)

	var buf bytes.Buffer

	err := f.Render(&buf)
	if err != nil {
		return nil, fmt.Errorf("rendering file: %v", err)
	}

	return buf.Bytes(), nil
}

func BootstrapFile(projectName string) ([]byte, error) {
	f := NewFilePath("gen")

	projectNameTitle := strings.Title(projectName)

	const mux = "github.com/gorilla/mux"

	f.Comment("// Service is the struct that will be exposed to serve HTTP traffic.")
	f.Type().Id("Service").Struct(
		Id("router").Op("*").Qual(mux, "Router"),
		Id("serviceImpl").Qual("gen", projectNameTitle+"Service"),
	)

	f.Comment("// ServeHTTP is what ultimately allows this service to be " +
		"used by the standard library's")
	f.Comment("// listen and serve functions")
	f.Func().Params(
		Id("s").Op("*").Id("Service"),
	).Id("ServeHTTP").Add(httpMethodParams()).Block(
		Id("s").Dot("router").Dot("ServeHTTP").Call(
			Id("w"),
			Id("r"),
		),
	)

	f.Comment("// Route is a struct that holds the path, the handler to " +
		"be called when that path is hit, and")
	f.Comment("// which list of methods it should serve.")

	f.Type().Id("route").Struct(
		Id("path").String(),
		Id("handler").Qual("net/http", "HandlerFunc"),
		Id("methods").Index().String(),
	)

	f.Comment("// New returns a new service implementation, using the " +
		"service as a dependency. It also sets up the routes")
	f.Comment("// and the middlewares.")

	f.Func().Id("New").Params(
		Id("service").Qual("gen", projectNameTitle+"Service"),
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

	var buf bytes.Buffer

	err := f.Render(&buf)
	if err != nil {
		return nil, fmt.Errorf("rendering file: %v", err)
	}

	return buf.Bytes(), nil
}

func MainFile(projectName string) ([]byte, error) {
	f := NewFile("main")

	f.Func().Id("main").Params().Block(
		Id("service").Op(":=").
			Qual(fmt.Sprintf("%s/gen", projectName), "New").
			Call(
				Op("&").Qual(projectName, "Server").Values(),
			),
		Qual("log", "Fatal").Call(
			Qual("net/http", "ListenAndServe").Call(
				Lit(":8080"), Id("service"),
			),
		),
	)

	var buf bytes.Buffer

	err := f.Render(&buf)
	if err != nil {
		return nil, fmt.Errorf("rendering file: %v", err)
	}

	return buf.Bytes(), nil
}

func InterfaceFile(projectName string) ([]byte, error) {
	f := NewFile("gen")

	title := strings.Title(projectName)
	service := fmt.Sprintf("%sService", title)
	handler := fmt.Sprintf("%sHandler", title)
	middleware := fmt.Sprintf("%sMiddleware", title)

	f.Commentf("// %s encapsulates the handler interface, which holds "+
		"all the methods to be called", service)
	f.Comment("// by the server, and middleware interface, which contains " +
		"all the middlewares to be added to the service.")

	f.Type().Id(service).Interface(
		Id(handler),
		Id(middleware),
	)

	f.Commentf("// %s is the interface for the handlers. Any new "+
		"endpoint added by seed will be added here as a", handler)
	f.Comment("// new method on the interface.")

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

	var buf bytes.Buffer

	err := f.Render(&buf)
	if err != nil {
		return nil, fmt.Errorf("rendering file: %v", err)
	}

	return buf.Bytes(), nil
}

func ProjectStructure(projectName string) error {
	err := createDirs(files.Pwd, projectName)
	if err != nil {
		return fmt.Errorf("creating folders: %v", err)
	}

	return nil
}

func createDirs(pwd, projectName string) error {
	paths := [][]string{
		{pwd, projectName},
		{pwd, projectName, consts.CmdFolder},
		{pwd, projectName, consts.GenFolder},
	}

	for _, path := range paths {
		fullPath := filepath.Join(path...)

		err := os.MkdirAll(fullPath, files.DefaultPerm)
		if err != nil {
			return fmt.Errorf("creating %v failed: %v", path, err)
		}

		err = os.Chmod(fullPath, files.DefaultPerm)
		if err != nil {
			return fmt.Errorf("changing %v permissions failed: %v", path, err)
		}
	}

	return nil
}

func ServiceDescriptor(projectName string) error {
	desc := descriptor.Base(descriptor.Info{
		Name:    projectName,
		Summary: "just a test for now",
	})

	path := filepath.Join(files.Pwd, projectName, projectName+".yml")

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed creating project descriptor: %v", err)
	}
	defer f.Close()

	err = yaml.NewEncoder(f).Encode(&desc)
	if err != nil {
		return fmt.Errorf("failed creating descriptor: %v", err)
	}

	return nil
}

func GoModule(projectName string) ([]byte, error) {
	goModContents := fmt.Sprintf(`module %s

go 1.12

require github.com/gorilla/mux v1.7.1
`, projectName)

	return []byte(goModContents), nil
}

func httpHandlerFunc() *Statement {
	return Func().Add(httpMethodParams())
}

func httpMethodParams() *Statement {
	return Params(
		Id("w").Qual("net/http", "ResponseWriter"),
		Id("r").Op("*").Qual("net/http", "Request"),
	)
}
