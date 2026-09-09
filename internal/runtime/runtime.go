package runtime

import "fmt"

// Runtime defines the behavior required to execute a program
// using a specific runtime environment.
type Runtime interface {
	Image() string
	Extension() string
	Command(fileName string) []string
}

// GetRuntime is a function that takes a string name
// and returns a specific Runtime implementation.
func GetRuntime(name string) (Runtime, error) {
	switch name {
	case "python":
		return &PythonRuntime{
			"python:3.11-alpine",
			".py",
		}, nil
	case "go":
		return &GoRuntime{
			"golang:1.21-alpine",
			".go",
		}, nil
	default:
		return nil, fmt.Errorf("unsupported runtime: %s", name)
	}
}
