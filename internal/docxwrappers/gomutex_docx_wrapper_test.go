package docxwrappers

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/gomutex/godocx"
	"github.com/gomutex/godocx/wml/ctypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tomwatkins1994/go-docx-template/internal/images"
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
	require.NoError(err)
	assert.Equal(newXmlString, xmlString)
}

// Commented out as currently not full working due to <a:stretch> element issues
// func TestGomutexSetDocumentXmlWithImage(t *testing.T) {
// 	assert := assert.New(t)
// 	require := require.New(t)

// 	docx, err := NewGomutexDocxFromFilename("../../test_templates/test_basic.docx")
// 	require.NoError(err)

// 	image, err := images.CreateInlineImage("../../test_templates/test_image.png")
// 	require.NoError(err)

// 	imageXml, err := docx.AddInlineImage(image)
// 	require.NoError(err)

// 	documentSectPtrXml, err := xml.Marshal(docx.Document.Body.SectPr)
// 	require.NoError(err)

// 	newXmlString := "<w:body><w:p><w:r>" + imageXml + "</w:r></w:p>" + string(documentSectPtrXml) + "</w:body>"
// 	err = docx.ReplaceDocumentXml(newXmlString)
// 	require.NoError(err)

// 	xmlString, err := docx.GetDocumentXml()
// 	require.NoError(err)
// 	assert.Equal(newXmlString, xmlString)
// }

func TestGomutexMergeTags(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	docx, err := godocx.NewDocument()
	require.NoError(err)

	p := docx.AddParagraph("").GetCT()
	pStartText := ctypes.Text{
		Text: "{{ .tag ",
	}
	pEndText := ctypes.Text{
		Text: "}}",
	}
	p.Children = []ctypes.ParagraphChild{
		{
			Run: &ctypes.Run{
				Children: []ctypes.RunChild{
					{
						Text: &pStartText,
					},
					{
						Text: &pEndText,
					},
				},
			},
		},
	}

	tblPara := docx.AddTable().AddRow().AddCell().AddParagraph("").GetCT()
	tblStartText := ctypes.Text{
		Text: "{{ .tbltag ",
	}
	tblEndText := ctypes.Text{
		Text: "}}",
	}
	tblPara.Children = []ctypes.ParagraphChild{
		{
			Run: &ctypes.Run{
				Children: []ctypes.RunChild{
					{
						Text: &tblStartText,
					},
					{
						Text: &tblEndText,
					},
				},
			},
		},
	}

	mergeGomutexTags(docx.Document.Body.Children)

	assert.Equal("", pStartText.Text)
	assert.Equal("{{ .tag }}", pEndText.Text)
	assert.Equal("", tblStartText.Text)
	assert.Equal("{{ .tbltag }}", tblEndText.Text)
}

func TestGomutexMergeTagsInParagraph(t *testing.T) {
	t.Run("Tags in text nodes in same run should get merged", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		docx, err := godocx.NewDocument()
		require.NoError(err)

		p := docx.AddParagraph("").GetCT()

		startText := ctypes.Text{
			Text: "{{ .tag ",
		}
		endText := ctypes.Text{
			Text: "}}",
		}

		p.Children = []ctypes.ParagraphChild{
			{
				Run: &ctypes.Run{
					Children: []ctypes.RunChild{
						{
							Text: &startText,
						},
						{
							Text: &endText,
						},
					},
				},
			},
		}

		mergeGomutexTagsInParagraph(p)

		assert.Equal("", startText.Text)
		assert.Equal("{{ .tag }}", endText.Text)
	})

	t.Run("Tags in text nodes in different runs should get merged", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		docx, err := godocx.NewDocument()
		require.NoError(err)

		p := docx.AddParagraph("").GetCT()

		startText := ctypes.Text{
			Text: "{{ .tag ",
		}
		endText := ctypes.Text{
			Text: "}}",
		}

		p.Children = []ctypes.ParagraphChild{
			{
				Run: &ctypes.Run{
					Children: []ctypes.RunChild{
						{
							Text: &startText,
						},
					},
				},
			},
			{
				Run: &ctypes.Run{
					Children: []ctypes.RunChild{
						{
							Text: &endText,
						},
					},
				},
			},
		}

		mergeGomutexTagsInParagraph(p)

		assert.Equal("", startText.Text)
		assert.Equal("{{ .tag }}", endText.Text)
	})
}

func TestGomutexMergeTagsInTable(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	docx, err := godocx.NewDocument()
	require.NoError(err)

	tbl := docx.AddTable().GetCT()

	p1StartText := ctypes.Text{
		Text: "{{ .tag1 ",
	}
	p1EndText := ctypes.Text{
		Text: "}}",
	}
	p1 := &ctypes.Paragraph{
		Children: []ctypes.ParagraphChild{
			{
				Run: &ctypes.Run{
					Children: []ctypes.RunChild{
						{
							Text: &p1StartText,
						},
						{
							Text: &p1EndText,
						},
					},
				},
			},
		},
	}
	p2StartText := ctypes.Text{
		Text: "{{ .tag2 ",
	}
	p2EndText := ctypes.Text{
		Text: "}}",
	}
	p2 := &ctypes.Paragraph{
		Children: []ctypes.ParagraphChild{
			{
				Run: &ctypes.Run{
					Children: []ctypes.RunChild{
						{
							Text: &p2StartText,
						},
						{
							Text: &p2EndText,
						},
					},
				},
			},
		},
	}

	tbl.RowContents = []ctypes.RowContent{
		{
			Row: &ctypes.Row{
				Contents: []ctypes.TRCellContent{
					{
						Cell: &ctypes.Cell{
							Contents: []ctypes.TCBlockContent{
								{
									Paragraph: p1,
								},
								{
									Paragraph: p2,
								},
							},
						},
					},
				},
			},
		},
	}

	mergeGomutexTagsInTable(tbl)

	assert.Equal("", p1StartText.Text)
	assert.Equal("{{ .tag1 }}", p1EndText.Text)
	assert.Equal("", p2StartText.Text)
	assert.Equal("{{ .tag2 }}", p2EndText.Text)
}

func TestGomutexAddInlineImage(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	t.Run("Should return the XML string for the image", func(t *testing.T) {
		docx, err := NewGomutexDocxFromFilename("../../test_templates/test_basic.docx")
		require.NoError(err)

		image, err := images.CreateInlineImage("../../test_templates/test_image.png")
		require.NoError(err)

		imageXml, err := docx.AddInlineImage(image)

		assert.NoError(err)
		assert.NotEmpty(imageXml)
	})
}

func TestGomutexSave(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	docx, err := NewGomutexDocxFromFilename("../../test_templates/test_basic.docx")
	require.NoError(err)

	f, err := os.Create("../../test_templates/generated_gomutex_test_basic.docx")
	require.Nil(err, "Error creating document")

	err = docx.Save(f)
	assert.Nil(err, "Error saving document")
}
