# Go Docx Template

> [!IMPORTANT] 
> This library is currently in active development and is therefore not recommended for production use

## Introduction

A simple Go library for merging docx files with data.

This is a wrapper library around the [go-docx](https://github.com/fumiama/go-docx) library by [fumiama](https://github.com/fumiama). All of the methods on the `docx.Docx` struct from go-docx are inherited by the `docxtpl.DocxTmpl` struct with an additional `Render` method.

## Usage

### Installation

```sh
go get github.com/tomwatkins1994/go-docx-template@latest 
```

### Render a document 

```go
package main

import (
  "os"

  "github.com/tomwatkins1994/go-docx-template"
)

func main() {
  // Parse the document 
  // If using a reader, use docxtpl.Parse instead
  doc, err := docxtpl.ParseFromFilename("template.docx")
  if err != nil {
    panic(err)
  }

  // Render the document
  data := struct {
    FirstName     string
    LastName      string
    Gender        string
  }{
    FirstName: "Tom",
    LastName:  "Watkins",
    Gender:    "Male",
  }
  err = doc.Render(data)
  if err != nil {
    panic(err)
  }

  // Create a new file for the output
  f, err := os.Create("generated.docx")
  if err != nil {
    panic(err)
  }
  err = doc.Save(f)
  if err != nil {
    panic(err)
  }
  err = f.Close()
  if err != nil {
    panic(err)
  }
}
```
Examples of docx files can be found in the [tests](https://github.com/tomwatkins1994/go-docx-template/tree/main/test_templates) directory of this repository.

## Acknowledgements

A lof of the heavy lifting such as XML parsing is done by the [go-docx](https://github.com/fumiama/go-docx) library by [fumiama](https://github.com/fumiama).

This library was also heavily inspired by the excellent [python-docx-template](https://github.com/elapouya/python-docx-template) library for Python written by [elapouya](https://github.com/elapouya).

## License

AGPL-3.0. See [LICENSE](LICENSE)
