package main

import (
	"bytes"
	"fmt"
	"log"

	"os"
	"regexp"
	"text/template"

	"github.com/fumiama/go-docx"
)

func main() {
	readFile, err := os.Open("test.docx")
	if err != nil {
		panic(err)
	}
	fileinfo, err := readFile.Stat()
	if err != nil {
		panic(err)
	}
	size := fileinfo.Size()
	doc, err := docx.Parse(readFile, size)
	if err != nil {
		panic(err)
	}

	for _, item := range doc.Document.Body.Items {
		switch item.(type) {
		case *docx.Paragraph, *docx.Table: // printable
			p, ok := item.(*docx.Paragraph)
			if ok {
				currentText := ""
				inIncompleteTag := false
				for _, pChild := range p.Children {
					run, ok := pChild.(*docx.Run)
					if ok {
						for _, runChild := range run.Children {
							text, ok := runChild.(*docx.Text)
							if ok {
								if inIncompleteTag {
									currentText += text.Text
								} else {
									currentText = text.Text
								}
								containsIncompleteTags, err := textContainsIncompleteTags(currentText)
								if err != nil {
									panic(err)
								}
								if containsIncompleteTags {
									text.Text = ""
									inIncompleteTag = true
								} else {
									inIncompleteTag = false
									containsTags, err := textContainsTags(currentText)
									if err != nil {
										panic(err)
									}
									if containsTags {
										newText, err := replaceTagsInText(currentText)
										if err != nil {
											panic(err)
										}
										fmt.Printf("Replacing %v with %v\n", text.Text, newText)
										text.Text = newText
									}
								}
							}
						}
					}
				}
			}
		}
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

type TemplateData struct {
	ProjectNumber string
	Status        string
}

func replaceTagsInText(text string) (string, error) {
	fmt.Println(text)
	tmpl, err := template.New("").Parse(text)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
		return "", err
	}
	templateData := TemplateData{"B-00001", "New"}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, templateData)
	return buf.String(), err
}

func textContainsTags(text string) (bool, error) {
	r, err := regexp.Compile("{{.+}}")
	if err != nil {
		return false, err
	}
	return r.MatchString(text), err
}

func textContainsIncompleteTags(text string) (bool, error) {
	r, err := regexp.Compile("{{[^}]*(}$|$)|{$|^}")
	if err != nil {
		return false, err
	}
	return r.MatchString(text), err
}
