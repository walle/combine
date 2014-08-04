package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/walle/combine"
)

const VERSION = "0.1.0"

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `combine is a utility for creating a new file from a template that includes other files.

Usage: combine [options...]

Options:
`)
		flag.PrintDefaults()
		os.Exit(1)
	}

	minify := flag.String("t", "", "Minify result, the type of file to minify, [js/css]")
	inputFile := flag.String("i", "", "Input file, the template to use as input, defaults to stdin")
	outputFile := flag.String("o", "", "Output file, the path to write the output to, defaults to stdout")
	baseDir := flag.String("d", "", "Base directory, the path files are relatively read from. Only used if reading input from stdin, if input file is given files are read relatively to the input file.")
	version := flag.Bool("version", false, "Version, print version information and exit.")

	flag.Parse()

	// No flags is not valid, show usage
	if flag.NFlag() == 0 {
		flag.Usage()
	}

	// If version is requested, print info and exit
	if *version {
		fmt.Fprintf(os.Stdout, "combine %s\n", VERSION)
		os.Exit(0)
	}

	var combiner combine.Combiner

	// Route the flags to the correct combiner
	if strings.ToLower(*minify) == combine.JS || strings.ToLower(*minify) == combine.CSS {
		if *inputFile != "" && *outputFile != "" {
			combiner = combine.NewYuiMinifyFileCombiner(*inputFile, *outputFile, *minify)
		} else if *inputFile != "" && *outputFile == "" {
			combiner = combine.NewYuiMinifyFileToStreamCombiner(*inputFile, os.Stdout, *minify)
		} else if *baseDir != "" {
			if *inputFile == "" && *outputFile != "" {
				combiner = combine.NewYuiMinifyStreamToFileCombiner(os.Stdin, *outputFile, *baseDir, *minify)
			} else {
				combiner = combine.NewYuiMinifyStreamCombiner(os.Stdin, os.Stdout, *baseDir, *minify)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Base directory (-d) is required when reading from stdin\n")
			flag.Usage()
		}
	} else {
		if *inputFile != "" && *outputFile != "" {
			combiner = combine.NewFileCombiner(*inputFile, *outputFile)
		} else if *inputFile != "" && *outputFile == "" {
			combiner = combine.NewFileToStreamCombiner(*inputFile, os.Stdout)
		} else if *baseDir != "" {
			if *inputFile == "" && *outputFile != "" {
				combiner = combine.NewStreamToFileCombiner(os.Stdin, *outputFile, *baseDir)
			} else {
				combiner = combine.NewStreamCombiner(os.Stdin, os.Stdout, *baseDir)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Base directory (-d) is required when reading from stdin\n")
			flag.Usage()
		}
	}

	templateIncluder := &combine.TemplateIncluder{}
	includer, err := combiner.Read(templateIncluder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		flag.Usage()
	}

	errors := combiner.Combine(includer)
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}

	err = combiner.Write()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		flag.Usage()
	}

	os.Exit(0)
}
