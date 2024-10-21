package docxtpl

import "testing"

func TestBasicTemplate(t *testing.T) {
	xmlString := `
	<w:document>  
		<w:body>  
			<w:p>  
			<w:r>  
				<w:t>Text</w:t>  
			</w:r>  
			<w:fldSimple w:instr="AUTHOR">  
				<w:r>  
				<w:t>Author Name: {{.Name}}</w:t>  
				</w:r>  
			</w:fldSimple>  
			</w:p>  
		</w:body>  
	</w:document>`
	data := map[string]interface{}{
		"Name": "Tom Watkins",
	}
	output, err := replaceTagsInText(xmlString, data)
	if err != nil {
		t.Fatalf("Error in basic template %v", err)
	}
	expectedOutput := `
	<w:document>  
		<w:body>  
			<w:p>  
			<w:r>  
				<w:t>Text</w:t>  
			</w:r>  
			<w:fldSimple w:instr="AUTHOR">  
				<w:r>  
				<w:t>Author Name: Tom Watkins</w:t>  
				</w:r>  
			</w:fldSimple>  
			</w:p>  
		</w:body>  
	</w:document>`
	if output != expectedOutput {
		t.Fatalf("Output does not match expected: %v", output)
	}
}
