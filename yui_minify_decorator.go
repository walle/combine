package combine

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

// The types that can be minified
const (
	JS  = "js"
	CSS = "css"
)

// YuiMinifyCombiner is a combiner that minifies the output using yuicompressor.
// If the yuicompressor binary is not available in the $PATH an error is logged and the output is not minified.
type YuiMinifyDecorator struct {
	fileType string
}

// NewYuiMinifyDecorater creates a Decorater that minifies the output if the file type is js or css.
func NewYuiMinifyDecorator(fileType string) *YuiMinifyDecorator {
	return &YuiMinifyDecorator{fileType: fileType}
}

// Decorate minifies the result using yuicompressor.
// If the binary is not found in $PATH an error is returned and the output is not minified.
// If an error occurs while running the yuicompressor an error indicating this
// and the error message from yuicompressor is returned.
func (y *YuiMinifyDecorator) Decorate(combiner Combiner) error {
	_, err := exec.LookPath("yuicompressor")
	if err == nil {
		yuiCmd := exec.Command("yuicompressor", "--type", y.fileType)
		yuiCmd.Stdin = bytes.NewBuffer(combiner.Result())
		var buf []byte
		errBuf := bytes.NewBuffer(buf)
		yuiCmd.Stderr = errBuf
		res, err := yuiCmd.Output()
		if err != nil {
			return errors.New(fmt.Sprintf("Error running yuicompressor: %s\nError: %s\n", err, string(errBuf.Bytes())))
		}
		combiner.SetResult(res)
		return nil
	}
	return errors.New(fmt.Sprintf("No yuicompressor found in $PATH. Make sure you have it installed. The output will not be minified.\n"))
}
