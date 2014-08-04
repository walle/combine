package combine

import (
	"io"
)

// StreamCombiner is a combiner that reads an input stream and writes to an output file.
type StreamToFileCombiner struct {
	*StreamCombiner
	*FileCombiner
}

// NewStreamToFileCombiner creates a new StreamToFileCombiner.
func NewStreamToFileCombiner(input io.Reader, outputFile, baseDir string) *StreamToFileCombiner {
	return &StreamToFileCombiner{
		StreamCombiner: &StreamCombiner{
			input:   input,
			baseDir: baseDir,
		}, FileCombiner: &FileCombiner{
			outputFile: outputFile,
		},
	}
}

// Read initializes the Includer by reading the input stream.
// Returns error if the input cannot be read or if the Includer cannot be parsed.
func (s *StreamToFileCombiner) Read(includer Includer) (Includer, error) {
	return s.StreamCombiner.Read(includer)
}

// Write writes the output of the Includer to file.
// Returns error if the output could not be written.
func (s *StreamToFileCombiner) Write() error {
	s.FileCombiner.StreamCombiner = s.StreamCombiner
	return s.FileCombiner.Write()
}
