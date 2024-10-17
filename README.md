# Go Docx Template

## Introduction

A simple library for merging docx files with data.

This is a wrapper library around the [go-docx](https://github.com/fumiama/go-docx) library by [fumiama](https://github.com/fumiama). All of the methods on the `docx.Docx` struct from go-docx are inherited by the `docxtpl.DocxTmpl` struct with an additional `Render` method.

## Usage

```go
package main

import (
  "os"

  "github.com/tomwatkins1994/go-docx-template"
)

func main() {
  // Parse the document
  doc, err := docxtpl.Parse("template.docx")
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
  _, err = doc.WriteTo(f)
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
