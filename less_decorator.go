package combine

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

// LessDecorator is a decorator that converts the output using lessc.
type LessDecorator struct {
}

// NewLessDecorator returns a new less decorator. Should only be used for css files.
func NewLessDecorator() *LessDecorator {
	return &LessDecorator{}
}

// Decorate compiles the result using lessc.
// If the binary is not found in $PATH an error is returned and the output is not compiled.
// If an error occurs while running the lessc an error indicating this
// and the error message from lessc is returned.
func (l *LessDecorator) Decorate(combiner Combiner) error {
	_, err := exec.LookPath("lessc")
	if err == nil {
		yuiCmd := exec.Command("lessc", "-")
		yuiCmd.Stdin = bytes.NewBuffer(combiner.Result())
		var buf []byte
		errBuf := bytes.NewBuffer(buf)
		yuiCmd.Stderr = errBuf
		res, err := yuiCmd.Output()
		if err != nil {
			return errors.New(fmt.Sprintf("Error running lessc: %s\nError: %s\n", err, string(errBuf.Bytes())))
		}
		combiner.SetResult(res)
		return nil
	}
	return errors.New(fmt.Sprintf("No lessc found in $PATH. Make sure you have it installed. The output will not be compiled.\n"))
}
