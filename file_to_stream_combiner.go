package combine

import (
	"io"
)

// FileToStreamCombiner is a combiner that reads the input from a stream and writes it to file.
type FileToStreamCombiner struct {
	*StreamCombiner
	*FileCombiner
}

// NewFileToStreamCombiner creates a new FileToStreamCombiner.
func NewFileToStreamCombiner(inputFile string, output io.Writer) *FileToStreamCombiner {
	return &FileToStreamCombiner{
		StreamCombiner: &StreamCombiner{
			output: output,
		}, FileCombiner: &FileCombiner{
			inputFile: inputFile,
		},
	}
}

// Read initializes the Includer by reading the inputFile.
// Returns error if the inputFile is not valid or if the Includer cannot be parsed.
func (f *FileToStreamCombiner) Read(includer Includer) (Includer, error) {
	return f.FileCombiner.Read(includer)
}

// Write writes the output of the Includer to the output stream.
// Returns error if the output could not be written.
func (f *FileToStreamCombiner) Write() error {
	return f.StreamCombiner.Write()
}
