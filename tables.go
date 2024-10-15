package main

import (
	"regexp"
	"sync"

	"github.com/fumiama/go-docx"
)

func (d *DocxTmpl) handleTableRangeTags() error {
	var wg sync.WaitGroup
	errCh := make(chan error)

	for _, item := range d.Document.Body.Items {
		item := item
		table, ok := item.(*docx.Table)
		if ok {
			wg.Add(1)
			go func() {
				defer wg.Done()

				containsRangeTag := false
				rangeTag := ""
				rowsToDelete := make(map[int]struct{})

				for rowIndex, row := range table.TableRows {
					for cellIndex, cell := range row.TableCells {
						if cellIndex == 0 && containsRangeTag {
							addRangeTag(cell, rangeTag)
							containsRangeTag = false
							rangeTag = ""
							break
						}
						var err error
						containsRangeTag, rangeTag, err = cellContainsRangeTag(cell)
						if err != nil {
							errCh <- err
						}
						if containsRangeTag {
							rowsToDelete[rowIndex] = struct{}{}
							break
						}
					}
				}

				// Remove the rows with range/end in
				if len(rowsToDelete) > 0 {
					newTableRows := []*docx.WTableRow{}
					for rowIndex, row := range table.TableRows {
						if _, found := rowsToDelete[rowIndex]; !found {
							newTableRows = append(newTableRows, row)
						}
					}
					table.TableRows = newTableRows
				}
			}()
		}
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func cellContainsRangeTag(cell *docx.WTableCell) (bool, string, error) {
	r, err := regexp.Compile("{{range .*}}|{{ range .* }}")
	if err != nil {
		return false, "", err
	}

	for _, paragraph := range cell.Paragraphs {
		for _, pChild := range paragraph.Children {
			run, ok := pChild.(*docx.Run)
			if ok {
				for _, rChild := range run.Children {
					text, ok := rChild.(*docx.Text)
					if ok {
						containsRangeTag := r.MatchString(text.Text)
						if containsRangeTag {
							return true, text.Text, nil
						}
					}
				}
			}
		}
	}

	return false, "", err
}

func addRangeTag(cell *docx.WTableCell, rangeText string) {
	for _, paragraph := range cell.Paragraphs {
		for _, pChild := range paragraph.Children {
			run, ok := pChild.(*docx.Run)
			if ok {
				for _, rChild := range run.Children {
					text, ok := rChild.(*docx.Text)
					if ok {
						text.Text = rangeText + text.Text
						return
					}
				}
			}
		}
	}
}
