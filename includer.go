package combine

import (
	"io"
)

// Includer is responsible for reading files into it's output.
type Includer interface {
	Initialize(io.Reader, string) error
	Process() []byte
	Read(string) string
	AnyErrors() bool
	Errors() []error
	BaseDir() string
}
