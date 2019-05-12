package files

import (
	"fmt"
	"os"
	"runtime"
)

const (
	PermWindows = 0666
	PermLinux   = 0755
)

var (
	Pwd         string
	DefaultPerm os.FileMode
)

func init() {
	switch runtime.GOOS {
	case "windows":
		DefaultPerm = PermWindows
	default:
		DefaultPerm = PermLinux
	}

	var err error

	Pwd, err = os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("setting up pwd: %v", err))
	}
}
