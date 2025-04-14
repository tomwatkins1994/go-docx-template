package docxtpl

import (
	"fmt"
	"maps"
	"text/template"

	"github.com/tomwatkins1994/go-docx-template/internal/functions"
)

// Register a function which can then be used within your template
//
//	d.RegisterFunction("sayHello", func(text string) string {
//		return "Hello " + text
//	})
func (d *DocxTmpl) RegisterFunction(name string, fn any) error {
	if !functions.FunctionNameValid(name) {
		return fmt.Errorf("function name %q is not a valid identifier", name)
	}

	// Check the function signature
	err := functions.FunctionValid(fn)
	if err != nil {
		return fmt.Errorf("error registering function (%s): %s", name, err.Error())
	}

	// Add to the function map
	d.funcMap[name] = fn

	return nil
}

// Get a pointer to the documents function map. This will include built-in functions.
func (d *DocxTmpl) GetRegisteredFunctions() *template.FuncMap {
	copiedFuncMap := make(template.FuncMap)
	maps.Copy(copiedFuncMap, d.funcMap)
	return &copiedFuncMap
}
