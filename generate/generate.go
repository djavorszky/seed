package generate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"seed/consts"
	"seed/descriptor"
	"seed/files"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/go-yaml/yaml"
)

func CreateServiceFile(projectName string) error {
	f := jen.NewFilePath(projectName)

	f.Type().Id("Server").Struct()

	f.Func().Params(
		jen.Id("s").Op("*").Id("Server"),
	).Id("LoggerMw").Params(
		jen.Id("next").Qual("net/http", "Handler"),
	).Qual("net/http", "Handler").Block(
		jen.Comment("// Anything you add here will be executed once, during startup."),
		jen.Comment("// The returned http.Handler will be able to access these variables"),
		jen.Comment("// thanks to closure."),
		jen.Line(),
		jen.Id("prefix").Op(":=").Lit(fmt.Sprintf("[%s] - ", projectName)),
		jen.Return(
			jen.Qual("net/http", "HandlerFunc").Call(
				jen.Func().Params(
					jen.Id("w").Qual("net/http", "ResponseWriter"),
					jen.Id("r").Op("*").Qual("net/http", "Request"),
				).Block(
					jen.Qual("log", "Println").Call(
						jen.Id("prefix"),
						jen.Id("r").Dot("RemoteAddr"),
						jen.Id("r").Dot("Method"),
						jen.Id("r").Dot("RequestURI"),
					),
					jen.Line(),
					jen.Id("next").Dot("ServeHTTP").Call(
						jen.Id("w"), jen.Id("r"),
					),
				),
			),
		),
	)

	f.Line()

	f.Func().Params(
		jen.Id("s").Op("*").Id("Server"),
	).Id("Index").Params().Qual("net/http", "HandlerFunc").Block(
		jen.Comment("// Anything you add here will be executed once, during startup."),
		jen.Comment("// The returned http.HandlerFunc will be able to access these variables"),
		jen.Comment("// thanks to closure."),
		jen.Line(),
		jen.Id("defaultMsg").Op(":=").Index().Byte().Call(jen.Lit("I'm alive!")),
		jen.Return().Func().Params(
			jen.Id("w").Qual("net/http", "ResponseWriter"),
			jen.Id("r").Op("*").Qual("net/http", "Request"),
		).Block(
			jen.Id("w").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusOK")),
			jen.Line(),
			jen.List(
				jen.Id("_"), jen.Id("err"),
			).Op(":=").Id("w").Dot("Write").Call(jen.Id("defaultMsg")),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Panic(jen.Id("err")),
			),
		),
	)

	path := filepath.Join(files.Pwd, projectName, projectName+".go")

	err := f.Save(path)
	if err != nil {
		return fmt.Errorf("saving file: %v", err)
	}

	return nil
}

func CreateGeneratedFile(projectName string) error {
	f := jen.NewFilePath("gen")

	projectNameTitle := strings.Title(projectName)

	const mux = "github.com/gorilla/mux"

	f.Comment("// Service is the struct that will be exposed to serve HTTP traffic.")
	f.Type().Id("Service").Struct(
		jen.Id("router").Op("*").Qual(mux, "Router"),
		jen.Id("serviceImpl").Qual("gen", projectNameTitle+"Service"),
	)

	f.Comment("// ServeHTTP is what ultimately allows this service to be " +
		"used by the standard library's")
	f.Comment("// listen and serve functions")
	f.Func().Params(
		jen.Id("s").Op("*").Id("Service"),
	).Id("ServeHTTP").Params(
		jen.Id("w").Qual("net/http", "ResponseWriter"),
		jen.Id("r").Op("*").Qual("net/http", "Request"),
	).Block(
		jen.Id("s").Dot("router").Dot("ServeHTTP").Call(
			jen.Id("w"),
			jen.Id("r"),
		),
	)

	f.Comment("// Route is a struct that holds the path, the handler to " +
		"be called when that path is hit, and")
	f.Comment("// which list of methods it should serve.")

	f.Type().Id("route").Struct(
		jen.Id("path").String(),
		jen.Id("handler").Qual("net/http", "HandlerFunc"),
		jen.Id("methods").Index().String(),
	)

	f.Comment("// New returns a new service implementation, using the " +
		"service as a dependency. It also sets up the routes")
	f.Comment("// and the middlewares.")

	f.Func().Id("New").Params(
		jen.Id("service").Qual("gen", projectNameTitle+"Service"),
	).Op("*").Qual("gen", "Service").Block(
		jen.Id("s").Op(":=").Op("&").Qual("gen", "Service").Values(
			jen.Dict{
				jen.Id("router"):      jen.Qual(mux, "NewRouter").Call(),
				jen.Id("serviceImpl"): jen.Id("service"),
			},
		),
		jen.Empty(),
		jen.Id("s").Dot("routes").Call(),
		jen.Id("s").Dot("middlewares").Call(),
		jen.Empty(),
		jen.Return(jen.Id("s")),
	)

	f.Comment("// routes sets up the routes to be served by the service")
	f.Func().Params(
		jen.Id("s").Op("*").Id("Service"),
	).Id("routes").Params().Block(
		jen.Id("routes").Op(":=").Index().Id("route").Values(
			jen.Block(
				jen.Dict{
					jen.Id("handler"): jen.Id("s").Dot("serviceImpl").Dot("Index").Call(),
					jen.Id("methods"): jen.Index().String().Values(jen.Qual("net/http", "MethodGet")),
					jen.Id("path"):    jen.Lit("/"),
				},
			),
		),
		jen.Empty(),
		jen.For(
			jen.List(jen.Id("_"), jen.Id("route")).Op(":=").Range().Id("routes").Block(
				jen.Id("s").Dot("router").
					Dot("StrictSlash").Call(jen.Lit(true)).
					Dot("HandleFunc").
					Call(
						jen.Id("route").Dot("path"), jen.Id("route").Dot("handler")).
					Dot("Methods").
					Call(
						jen.Id("route").Dot("methods").Op("..."),
					),
			),
		),
	)

	f.Comment("// middlewares sets up the middlewares to be set up by the service")
	f.Func().Params(
		jen.Id("s").Op("*").Id("Service"),
	).Id("middlewares").Params().Block(
		jen.Id("mws").Op(":=").Index().Qual(mux, "MiddlewareFunc").Values(
			jen.Id("s").Dot("serviceImpl").Dot("LoggerMw"),
		),
		jen.Empty(),
		jen.For(
			jen.List(jen.Id("_"), jen.Id("mw")).Op(":=").Range().Id("mws").Block(
				jen.Id("s").Dot("router").Dot("Use").Call(jen.Id("mw")),
			),
		),
	)

	path := filepath.Join(files.Pwd, projectName, consts.GenFolder, consts.GeneratedFile)

	err := f.Save(path)
	if err != nil {
		return fmt.Errorf("saving file: %v", err)
	}

	return nil
}

func CreateServiceDescriptor(projectName string) error {
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

func CreateMainFile(projectName string) error {
	f := jen.NewFile("main")

	f.Func().Id("main").Params().Block(
		jen.Id("service").Op(":=").
			Qual(fmt.Sprintf("%s/gen", projectName), "New").
			Call(
				jen.Op("&").Qual(projectName, "Server").Values(),
			),
		jen.Qual("log", "Fatal").Call(
			jen.Qual("net/http", "ListenAndServe").Call(
				jen.Lit(":8080"), jen.Id("service"),
			),
		),
	)

	path := filepath.Join(files.Pwd, projectName, consts.CmdFolder, consts.MainFile)

	err := f.Save(path)
	if err != nil {
		return fmt.Errorf("saving file: %v", err)
	}

	return nil
}

func CreateInterfaceFile(projectName string) error {
	f := jen.NewFile("gen")

	title := strings.Title(projectName)
	service := fmt.Sprintf("%sService", title)
	handler := fmt.Sprintf("%sHandler", title)
	middleware := fmt.Sprintf("%sMiddleware", title)

	f.Commentf("// %s encapsulates the handler interface, which holds "+
		"all the methods to be called", service)
	f.Comment("// by the server, and middleware interface, which contains " +
		"all the middlewares to be added to the service.")

	f.Type().Id(service).Interface(
		jen.Id(handler),
		jen.Id(middleware),
	)

	f.Commentf("// %s is the interface for the handlers. Any new "+
		"endpoint added by seed will be added here as a", handler)
	f.Comment("// new method on the interface.")

	f.Type().Id(handler).Interface(
		jen.Id("Index").Params().Qual("net/http", "HandlerFunc"),
	)

	f.Commentf("// %s is the interface for all the middlewares that "+
		"will be added to all of the paths.", middleware)

	f.Type().Id(middleware).Interface(
		jen.Id("LoggerMw").Params(
			jen.Qual("net/http", "Handler"),
		).Qual("net/http", "Handler"),
	)

	path := filepath.Join(files.Pwd, projectName, consts.GenFolder, consts.InterfaceFile)

	err := f.Save(path)
	if err != nil {
		return fmt.Errorf("saving file: %v", err)
	}

	return nil
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

func GoModule(projectName string) error {
	goModContents := fmt.Sprintf(`module %s

go 1.12

require github.com/gorilla/mux v1.7.1
`, projectName)

	path := filepath.Join(files.Pwd, projectName, "go.mod")

	err := ioutil.WriteFile(path, []byte(goModContents), files.DefaultPerm)
	if err != nil {
		return fmt.Errorf("failed creating go.mod file: %v", err)
	}

	return nil
}
