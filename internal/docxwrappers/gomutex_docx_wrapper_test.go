package docxwrappers

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGomutexGetDocumentXml(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	docx, err := NewGomutexDocxFromFilename("../../test_templates/test_basic.docx")
	require.NoError(err)

	xmlString, err := docx.GetDocumentXml()
	require.NoError(err)
	assert.NotEmpty(xmlString)
}

func TestGomutexSetDocumentXml(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	docx, err := NewGomutexDocxFromFilename("../../test_templates/test_basic.docx")
	require.NoError(err)

	documentSectPtrXml, err := xml.Marshal(docx.Document.Body.SectPr)
	require.NoError(err)
	newXmlString := "<w:body><w:p><w:r><w:t>Hello, World!</w:t></w:r></w:p>" + string(documentSectPtrXml) + "</w:body>"
	err = docx.ReplaceDocumentXml(newXmlString)
	require.NoError(err)

	xmlString, err := docx.GetDocumentXml()
	println(xmlString)
	require.NoError(err)
	assert.Equal(newXmlString, xmlString)
}

// func TestMergeTags(t *testing.T) {
// 	assert := assert.New(t)

// 	pStartText := docx.Text{
// 		Text: "{{ .tag ",
// 	}
// 	pEndText := docx.Text{
// 		Text: "}}",
// 	}
// 	p := docx.Paragraph{
// 		Children: []any{
// 			&docx.Run{
// 				Children: []any{
// 					&pStartText,
// 					&pEndText,
// 				},
// 			},
// 		},
// 	}

// 	tblStartText := docx.Text{
// 		Text: "{{ .tag ",
// 	}
// 	tblEndText := docx.Text{
// 		Text: "}}",
// 	}
// 	tbl := docx.Table{
// 		TableRows: []*docx.WTableRow{
// 			{
// 				TableCells: []*docx.WTableCell{
// 					{
// 						Paragraphs: []*docx.Paragraph{
// 							{
// 								Children: []any{
// 									&docx.Run{
// 										Children: []any{
// 											&tblStartText,
// 											&tblEndText,
// 										},
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	items := []any{
// 		&p,
// 		&tbl,
// 	}

// 	mergeTags(items)

// 	assert.Equal(pStartText.Text, "")
// 	assert.Equal(pEndText.Text, "{{ .tag }}")
// 	assert.Equal(tblStartText.Text, "")
// 	assert.Equal(tblEndText.Text, "{{ .tag }}")
// }

// func TestMergeTagsInParagraph(t *testing.T) {
// 	t.Run("Tags in text nodes in same run should get merged", func(t *testing.T) {
// 		assert := assert.New(t)

// 		startText := docx.Text{
// 			Text: "{{ .tag ",
// 		}
// 		endText := docx.Text{
// 			Text: "}}",
// 		}
// 		p := docx.Paragraph{
// 			Children: []any{
// 				&docx.Run{
// 					Children: []any{
// 						&startText,
// 						&endText,
// 					},
// 				},
// 			},
// 		}

// 		mergeTagsInParagraph(&p)

// 		assert.Equal(startText.Text, "")
// 		assert.Equal(endText.Text, "{{ .tag }}")
// 	})

// 	t.Run("Tags in text nodes in different runs should get merged", func(t *testing.T) {
// 		assert := assert.New(t)

// 		startText := docx.Text{
// 			Text: "{{ .tag ",
// 		}
// 		endText := docx.Text{
// 			Text: "}}",
// 		}
// 		p := docx.Paragraph{
// 			Children: []any{
// 				&docx.Run{
// 					Children: []any{
// 						&startText,
// 					},
// 				},
// 				&docx.Run{
// 					Children: []any{
// 						&endText,
// 					},
// 				},
// 			},
// 		}

// 		mergeTagsInParagraph(&p)

// 		assert.Equal(startText.Text, "")
// 		assert.Equal(endText.Text, "{{ .tag }}")
// 	})
// }

// func TestMergeTagsInTable(t *testing.T) {
// 	assert := assert.New(t)

// 	p1StartText := docx.Text{
// 		Text: "{{ .tag1 ",
// 	}
// 	p1EndText := docx.Text{
// 		Text: "}}",
// 	}
// 	p2StartText := docx.Text{
// 		Text: "{{ .tag2 ",
// 	}
// 	p2EndText := docx.Text{
// 		Text: "}}",
// 	}
// 	tbl := docx.Table{
// 		TableRows: []*docx.WTableRow{
// 			{
// 				TableCells: []*docx.WTableCell{
// 					{
// 						Paragraphs: []*docx.Paragraph{
// 							{
// 								Children: []any{
// 									&docx.Run{
// 										Children: []any{
// 											&p1StartText,
// 											&p1EndText,
// 										},
// 									},
// 								},
// 							},
// 							{
// 								Children: []any{
// 									&docx.Run{
// 										Children: []any{
// 											&p2StartText,
// 											&p2EndText,
// 										},
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	mergeTagsInTable(&tbl)

// 	assert.Equal(p1StartText.Text, "")
// 	assert.Equal(p1EndText.Text, "{{ .tag1 }}")
// 	assert.Equal(p2StartText.Text, "")
// 	assert.Equal(p2EndText.Text, "{{ .tag2 }}")
// }

// func TestAddInlineImage(t *testing.T) {
// 	assert := assert.New(t)
// 	require := require.New(t)

// 	t.Run("Should return the XML string for the image", func(t *testing.T) {
// 		reader, err := os.Open("../../test_templates/test_basic.docx")
// 		require.NoError(err)

// 		fileinfo, err := reader.Stat()
// 		require.NoError(err)
// 		size := fileinfo.Size()

// 		docx, err := NewFumiamaDocx(reader, size)
// 		require.NoError(err)

// 		image, err := images.CreateInlineImage("../../test_templates/test_image.png")
// 		require.NoError(err)

// 		imageXml, err := docx.AddInlineImage(image)
// 		assert.NoError(err)
// 		assert.NotEmpty(imageXml)
// 	})

// 	t.Run("Should add the PNG content type to the documents content types", func(t *testing.T) {
// 		reader, err := os.Open("../../test_templates/test_basic.docx")
// 		require.NoError(err)

// 		fileinfo, err := reader.Stat()
// 		require.NoError(err)
// 		size := fileinfo.Size()

// 		docx, err := NewFumiamaDocx(reader, size)
// 		require.NoError(err)

// 		image, err := images.CreateInlineImage("../../test_templates/test_image.png")
// 		require.NoError(err)

// 		_, err = docx.AddInlineImage(image)
// 		assert.NoError(err)
// 		assert.Contains(docx.contentTypes.Defaults, contenttypes.PNG_CONTENT_TYPE)
// 	})
// }

// func TestSave(t *testing.T) {
// 	assert := assert.New(t)
// 	require := require.New(t)

// 	reader, err := os.Open("../../test_templates/test_basic.docx")
// 	require.NoError(err)

// 	fileinfo, err := reader.Stat()
// 	require.NoError(err)
// 	size := fileinfo.Size()

// 	docx, err := NewFumiamaDocx(reader, size)
// 	require.NoError(err)

// 	f, err := os.Create("../../test_templates/generated_fumiama_test_basic.docx")
// 	require.Nil(err, "Error creating document")

// 	err = docx.Save(f)
// 	assert.Nil(err, "Error saving document")
// }
