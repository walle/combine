package combine

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

// The types that can be minified
const (
	JS  = "js"
	CSS = "css"
)

// YuiMinifyCombiner is a combiner that minifies the output using yuicompressor.
// If the yuicompressor binary is not available in the $PATH an error is logged and the output is not minified.
type YuiMinifyCombiner struct {
	Combiner
	fileType string
}

// NewYuiMinifyFileCombiner creates a Combiner that minifies the output and reads and writes to file.
func NewYuiMinifyFileCombiner(inputFile, outputFile, fileType string) *YuiMinifyCombiner {
	return &YuiMinifyCombiner{Combiner: NewFileCombiner(inputFile, outputFile), fileType: fileType}
}

// NewYuiMinifyFileToStreamCombiner creates a Combiner that minifies the output and reads from file and writes to stream.
func NewYuiMinifyFileToStreamCombiner(inputFile string, output io.Writer, fileType string) *YuiMinifyCombiner {
	return &YuiMinifyCombiner{Combiner: NewFileToStreamCombiner(inputFile, output), fileType: fileType}
}

// NewYuiMinifyFileToStreamCombiner creates a Combiner that minifies the output and reads and writes to stream.
func NewYuiMinifyStreamCombiner(input io.Reader, output io.Writer, baseDir, fileType string) *YuiMinifyCombiner {
	return &YuiMinifyCombiner{Combiner: NewStreamCombiner(input, output, baseDir), fileType: fileType}
}

// NewYuiMinifyFileToStreamCombiner creates a Combiner that minifies the output and reads from stream and writes to file.
func NewYuiMinifyStreamToFileCombiner(input io.Reader, outputFile, baseDir, fileType string) *YuiMinifyCombiner {
	return &YuiMinifyCombiner{Combiner: NewStreamToFileCombiner(input, outputFile, baseDir), fileType: fileType}
}

// Combine minifies the result using yuicompressor.
// If the binary is not found in $PATH an error is returned and the output is not minified.
// If an error occurs while running the yuicompressor an error indicating this
// and the error message from yuicompressor is returned.
func (y *YuiMinifyCombiner) Combine(includer Includer) []error {
	errorsList := y.Combiner.Combine(includer)
	if len(errorsList) > 0 {
		return errorsList
	}

	_, err := exec.LookPath("yuicompressor")
	if err == nil {
		yuiCmd := exec.Command("yuicompressor", "--type", y.fileType)
		yuiCmd.Stdin = bytes.NewBuffer(y.Combiner.Result())
		var buf []byte
		errBuf := bytes.NewBuffer(buf)
		yuiCmd.Stderr = errBuf
		res, err := yuiCmd.Output()
		if err != nil {
			errorsList := []error{
				errors.New(fmt.Sprintf("Error running yuicompressor: %s\n", err)),
				errors.New(fmt.Sprintf("Error: %s\n", string(errBuf.Bytes()))),
			}
			return errorsList
		}
		y.Combiner.SetResult(res)
		return nil
	}
	return []error{errors.New(fmt.Sprintf("No yuicompressor found in $PATH. Make sure you have it installed. The output will not be minified.\n"))}
}
