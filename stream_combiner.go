package combine

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// StreamCombiner is a combiner that reads an input stream and writes to an output stream.
type StreamCombiner struct {
	baseDir string
	input   io.Reader
	output  io.Writer
	result  []byte
}

// NewStreamCombiner creates a new StreamCombiner.
func NewStreamCombiner(input io.Reader, output io.Writer, baseDir string) *StreamCombiner {
	return &StreamCombiner{input: input, output: output, baseDir: baseDir}
}

// Read initializes the Includer by reading the input stream.
// Returns error if the input cannot be read or if the Includer cannot be parsed.
func (s *StreamCombiner) Read(includer Includer) (Includer, error) {
	err := includer.Initialize(s.input, s.baseDir)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing template: %s\n", err))
	}

	return includer, nil
}

// Combine processes the includer.
// Returns a list of errors if any errors occur in the processing.
func (s *StreamCombiner) Combine(includer Includer) []error {
	s.result = includer.Process()
	if includer.AnyErrors() {
		return includer.Errors()
	}

	return nil
}

// SetResult sets the output of the combiner. Usable for extending a combiner.
// See YuiMinifyCombiner#Combine for usage.
func (s *StreamCombiner) SetResult(r []byte) {
	s.result = r
}

// Result gets the output of the combiner.
func (s *StreamCombiner) Result() []byte {
	return s.result
}

// Write writes the output of the Includer to the output stream.
// Returns error if the output could not be written.
func (s *StreamCombiner) Write() error {
	_, err := io.Copy(s.output, bytes.NewBuffer(s.Result()))
	if err != nil {
		return errors.New(fmt.Sprintf("Error writing output: %s\n", err))
	}
	return nil
}
