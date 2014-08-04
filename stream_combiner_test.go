package combine_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/walle/combine"
)

func Test_Read(t *testing.T) {
	input := bytes.NewBufferString("input for read")
	output := bytes.NewBufferString("")
	baseDir := "/tmp"
	includer := &combine.TemplateIncluder{}

	streamCombiner := combine.NewStreamCombiner(input, output, baseDir)
	includeFile, err := streamCombiner.Read(includer)
	if err != nil {
		t.Errorf("Error reading on StreamCombiner: %s", err)
	}

	if includeFile.BaseDir() != "/tmp" {
		t.Errorf("Error assigning base directory on IncludeFile")
	}
}

func Test_ReadError(t *testing.T) {
	input := bytes.NewBufferString("input for read {{ me.Read 'test.txt }")
	output := bytes.NewBufferString("")
	baseDir := "/tmp"
	includer := &combine.TemplateIncluder{}

	streamCombiner := combine.NewStreamCombiner(input, output, baseDir)
	_, err := streamCombiner.Read(includer)
	if err == nil {
		t.Errorf("No error when reading broken template")
	}
}

func Test_Combine(t *testing.T) {
	expected, err := ioutil.ReadFile("README.md")
	if err != nil {
		t.Errorf("Error reading expected value: %s", err)
	}
	input := bytes.NewBufferString("{{ .Read \"README.md\"}}")
	output := bytes.NewBufferString("")
	baseDir := ""
	includer := &combine.TemplateIncluder{}

	streamCombiner := combine.NewStreamCombiner(input, output, baseDir)
	includeFile, err := streamCombiner.Read(includer)
	if err != nil {
		t.Errorf("Error reading on StreamCombiner: %s", err)
	}

	errors := streamCombiner.Combine(includeFile)
	if errors != nil {
		t.Errorf("Errors occured while combining: %+v", errors)
	}

	if !bytes.Equal(expected, streamCombiner.Result()) {
		t.Errorf("Did not get expected result after Combine")
	}
}

func Test_CombineErrors(t *testing.T) {
	input := bytes.NewBufferString("{{ .Read \"nonexisting\"}}")
	output := bytes.NewBufferString("")
	baseDir := ""
	includer := &combine.TemplateIncluder{}

	streamCombiner := combine.NewStreamCombiner(input, output, baseDir)
	includeFile, err := streamCombiner.Read(includer)
	if err != nil {
		t.Errorf("Error reading on StreamCombiner: %s", err)
	}

	errors := streamCombiner.Combine(includeFile)
	if errors == nil {
		t.Errorf("No errors occured while combining: %+v", errors)
	}
}

func Test_Write(t *testing.T) {
	expected, err := ioutil.ReadFile("README.md")
	if err != nil {
		t.Errorf("Error reading expected value: %s", err)
	}
	input := bytes.NewBufferString("{{ .Read \"README.md\"}}")
	output := bytes.NewBufferString("")
	baseDir := ""
	includer := &combine.TemplateIncluder{}

	streamCombiner := combine.NewStreamCombiner(input, output, baseDir)
	includeFile, err := streamCombiner.Read(includer)
	if err != nil {
		t.Errorf("Error reading on StreamCombiner: %s", err)
	}

	errors := streamCombiner.Combine(includeFile)
	if errors != nil {
		t.Errorf("Errors occured while combining: %+v", errors)
	}

	err = streamCombiner.Write()
	if err != nil {
		t.Errorf("Error writing StreamCombiner: %s", err)
	}

	if !bytes.Equal(expected, output.Bytes()) {
		t.Errorf("Did not get expected result after Write")
	}
}
