package combine

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// FileCombiner is a combiner that reads the input from file and writes it to file.
type FileCombiner struct {
	*StreamCombiner
	inputFile  string
	outputFile string
}

// NewFileCombiner creates a new FileCombiner.
func NewFileCombiner(inputFile, outputFile string) *FileCombiner {
	return &FileCombiner{inputFile: inputFile, outputFile: outputFile}
}

// Read initializes the Includer by reading the inputFile.
// Returns error if the inputFile is not valid or if the Includer cannot be parsed.
func (f *FileCombiner) Read(includer Includer) (Includer, error) {
	inputFileFileInfo, err := os.Stat(f.inputFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Not a valid input file: %s\n", f.inputFile))
	}

	if inputFileFileInfo.Mode().IsRegular() == false {
		return nil, errors.New(fmt.Sprintf("Not a valid input file: %s\n", f.inputFile))
	}

	file, err := os.Open(f.inputFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not read input file: %s\n", err))
	}
	defer file.Close()

	f.StreamCombiner = &StreamCombiner{input: file, baseDir: filepath.Dir(f.inputFile)}

	return f.StreamCombiner.Read(includer)
}

// Write writes the output of the Includer to outputFile.
// Returns error if the output could not be written.
func (f *FileCombiner) Write() error {
	file, err := os.Create(f.outputFile)
	if err != nil {
		return errors.New(fmt.Sprintf("Error creating output file: %s\n", err))
	}
	defer file.Close()

	f.StreamCombiner.output = file

	return f.StreamCombiner.Write()
}
