package combine

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"text/template"
)

// TemplateIncluder is a text/template implementation of Includer.
type TemplateIncluder struct {
	template *template.Template
	baseDir  string
	errors   []error
}

// NewTemplateIncluder creates a new include file ready for processing.
//  Returns error if reading or parsing fails.
func NewIncludeFile(input io.Reader, baseDir string) (*TemplateIncluder, error) {
	t := &TemplateIncluder{}
	err := t.Initialize(input, baseDir)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *TemplateIncluder) Initialize(input io.Reader, baseDir string) error {
	contents, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}

	tmpl := template.New("input")
	template, err := tmpl.Parse(string(contents))
	if err != nil {
		return err
	}

	t.template = template
	t.baseDir = baseDir

	return nil
}

// Precess executes the template and returns the bytes generated.
func (t *TemplateIncluder) Process() []byte {
	buffer := new(bytes.Buffer)
	t.template.Execute(buffer, t)
	return buffer.Bytes()
}

// Read reads the file content of filename based on the path in baseDir and returns it as a string
// Errors are logged internally and can be accessed by the Errors method.
func (t *TemplateIncluder) Read(filename string) string {
	contents, err := ioutil.ReadFile(filepath.Join(t.baseDir, filename))
	if err != nil {
		if t.errors == nil {
			t.errors = make([]error, 0, 5)
		}
		t.errors = append(t.errors, errors.New(fmt.Sprintf("Could not read included file: %s\n", err)))
	}

	return string(contents)
}

// AnyErrors returns true if any errors are logged, false otherwise.
func (t *TemplateIncluder) AnyErrors() bool {
	return len(t.errors) > 0
}

// Errors returns a copy of all errors that were logged.
func (t *TemplateIncluder) Errors() []error {
	errors := make([]error, len(t.errors))
	copy(errors, t.errors)

	return errors
}

// BaseDir returns the base directory.
func (t *TemplateIncluder) BaseDir() string {
	return t.baseDir
}
