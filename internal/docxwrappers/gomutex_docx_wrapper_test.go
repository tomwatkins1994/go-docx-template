package docxwrappers

import (
	"encoding/xml"
	"testing"

	"github.com/gomutex/godocx"
	"github.com/gomutex/godocx/wml/ctypes"
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

	tblCell := docx.AddTable().AddRow().AddCell()
	tblP1 := tblCell.AddParagraph("").GetCT()
	tblStartText := ctypes.Text{
		Text: "{{ .tbltag ",
	}
	tblEndText := ctypes.Text{
		Text: "}}",
	}
	tblP1.Children = []ctypes.ParagraphChild{
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
	// tblP2 := tblCell.AddParagraph("").GetCT()
	// tblEndText := ctypes.Text{
	// 	Text: "}}",
	// }
	// tblP2.Children = []ctypes.ParagraphChild{
	// 	{
	// 		Run: &ctypes.Run{
	// 			Children: []ctypes.RunChild{
	// 				{
	// 					Text: &tblEndText,
	// 				},
	// 			},
	// 		},
	// 	},
	// }

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
