# Combine

A library and binary for concatenating files using a template that includes other files.

Usage example is for a JavaScript file, and a template reading in the different parts to create a single file can be used. Also used for stylesheets(CSS). Combine works with yuicompressor to minify the result if requested.

Another usage example is to generate a book using [gimli](https:///github.com/walle/gimli) and want to create a custom front page and TOC and order the chapters correctly.

See [Examples](#examples) below for more details.

[Documentation on godoc.org](http://godoc.org/github.com/walle/combine)

## Installation

Installing using go get is the easiest.

To be able to install with `go get` requires you to have your `$GOPATH` setup and your `$GOPATH/bin` added to path as described here http://golang.org/doc/code.html#GOPATH.

### Installing the binary

    go get github.com/walle/combine/cmd/combine

### Installing the library

    go get github.com/walle/combine

## Examples

### Running the binary

To use the binary first [install it](#installation).

The options available for the binary.

    combine is a utility for creating a new file from a template that includes other files.

    Usage: combine [options...]

    Options:

      -d="": Base directory, the path files are relatively read from. Only used if reading input from stdin, if input file is given files are read relatively to the input file.
      -i="": Input file, the template to use as input, defaults to stdin
      -o="": Output file, the path to write the output to, defaults to stdout
      -t="": Minify result, the type of file to minify, [js/css]
      -version=false: Version, print version information and exit.

#### Combining a JavaScript file

Lets say you have two JavaScript files `functions.js` and `app.js`. `app.js` uses functions from `functions.js`. Now you want to combine these files to one when adding it to a website so the user doesn't have to download two files. You can do that with combine by creating a template like this called `application.template`

    (function(window, document) {
      'use strict';

      {{ .Read "functions.js" }}
      {{ .Read "app.js" }}

    })(window, document);

The anonymous function around your code is to protect from leaking variables and functions to the global window object. [More information](http://stackoverflow.com/questions/2421911/what-is-the-purpose-of-wrapping-whole-javascript-files-in-anonymous-functions-li)

You can then use this command to save your combined file as `application.js`

    $ combine -i application.template -o application.js

If you want to minify the `application.js` file you can add the `-t js` flag. This only works if you have `yuicompressor` installed. If you are using Homebrew on OSX you can install it with `$brew install yuicompressor`. This command minifies the output.

    $ combine -i application.template -o application.js -t js

#### Combining CSS files

Lets say you have separated two CSS files so it's easier to work with. `typography.css` and `layout.css` but want to combine them before using them so the user doesn't have to download two files. You can do that by creating a template like this called `style.template`

    {{ .Read "typography.css" }}
    {{ .Read "layout.css" }}

You can then use this command to save your combined file as `style.css`

    $ combine -i style.template -o style.css

If you want to minify the `style.css` file you can add the `-t css` flag. This only works if you have `yuicompressor` installed. If you are using Homebrew on OSX you can install it with `$brew install yuicompressor`. This command minifies the output.

    $ combine -i style.template -o style.css -t css

#### Combining a gimli document

If you are writing a book using markdown and want to use [gimli](https:///github.com/walle/gimli) to create a PDF of it you can use combine to do that.
See [gimli#installation](https://github.com/walle/gimli#installation) on how to get it if you don't have it installed.
Lets say that you have the files `front.md`, `chapter1.md` and `chapter2.md` and want to create a PDF with the files in that order. Setup a template like this called `my_book.template`

    {{ .Read "front.md" }}
    <div class="page-break"></div>
    {{ .Read "chapter1.md" }}
    <div class="page-break"></div>
    {{ .Read "chapter2.md" }}

You can then use this command to save your combined file as `my_book.md`

    $ combine -i my_book.template -o my_book.md

Then you can use `gimli` to generate the `my_book.pdf`.

    $ gimli -f my_book.md

### Using the library

If you want to use the library to combine for an example CSS files from your code you can use it like this.

```go
incl := &combine.TemplateIncluder{}
combiner := combine.NewFileToStreamCombiner(filename, writer)
includer, err := combiner.Read(incl)
if err != nil {
  log.Println(err)
}

errors := combiner.Combine(includer)
if len(errors) > 0 {
  log.Println(err)
}

err = combiner.Write()
if err != nil {
  log.Println(err)
}
```

Where `filename` points to the template and `writer` is the stream e.g. `http.ResponseWriter` you want to write to.

## Contributing

All contributions are welcome! See [CONTRIBUTING](CONTRIBUTING.md) for more info.

A suggested entry point for reading the code is [cmd/combine/main.go](cmd/combine/main.go) which uses the whole library.

## License

Licensed under MIT license. See [LICENSE](LICENSE) for more information.

## Authors

* Fredrik Wallgren