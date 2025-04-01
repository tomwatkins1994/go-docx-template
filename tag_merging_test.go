package docxtpl

import (
	"testing"

	"github.com/fumiama/go-docx"
	"github.com/stretchr/testify/assert"
)

func TestMergeTagsInParagraph(t *testing.T) {
	t.Run("Tags in text nodes in same run should get merged", func(t *testing.T) {
		assert := assert.New(t)

		startText := docx.Text{
			Text: "{{ .tag ",
		}
		endText := docx.Text{
			Text: "}}",
		}
		r := docx.Run{
			Children: []any{
				&startText,
				&endText,
			},
		}
		p := docx.Paragraph{
			Children: []any{
				&r,
			},
		}

		mergeTagsInParagraph(&p)

		assert.Equal(startText.Text, "")
		assert.Equal(endText.Text, "{{ .tag }}")
	})
}
