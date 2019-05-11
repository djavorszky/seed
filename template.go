package seed

// InterfaceStruct is used to generate the interface.go file
type InterfaceStruct struct {
	PackageName string
	Imports     []string
	ServiceName string
	Services    []string
	Middlewares []string
}

// MainStruct is used to generate the main.go file
type MainStruct struct {
	Imports        []string
	ServicePackage string
	ServiceAddress string
}
