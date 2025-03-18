package docxtpl

import (
	"regexp"
	"strings"
	"testing"
)

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

	data := map[string]any{
		"Name": "Tom Watkins",
	}
	outputXml, err := replaceTagsInText(xmlString, data, &defaultFuncMap)
	if err != nil {
		t.Fatalf("Error in basic template %v", err)
	}

	expectedOutputXml := `
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

	if removeXmlFormatting(outputXml) != removeXmlFormatting(expectedOutputXml) {
		t.Fatalf("Output does not match expected: %v", outputXml)
	}
}

func TestBasicTemplateWithFunctions(t *testing.T) {
	xmlString := `
	<w:document>  
		<w:body>  
			<w:p>  
			<w:r>  
				<w:t>Text</w:t>  
			</w:r>  
			<w:fldSimple w:instr="AUTHOR">  
				<w:r>  
					<w:t>Author Name (Upper): {{upper .Name}}</w:t>  
				</w:r> 
				<w:r>  
					<w:t>Author Name (Lower): {{lower .Name}}</w:t>  
				</w:r> 
				<w:r>  
					<w:t>Author Name (Title): {{title .Name}}</w:t>  
				</w:r>  
			</w:fldSimple>  
			</w:p>  
		</w:body>  
	</w:document>`

	data := map[string]any{
		"Name": "Tom watkins",
	}
	outputXml, err := replaceTagsInText(xmlString, data, &defaultFuncMap)
	if err != nil {
		t.Fatalf("Error in basic template %v", err)
	}

	expectedOutputXml := `
	<w:document>  
		<w:body>  
			<w:p>  
			<w:r>  
				<w:t>Text</w:t>  
			</w:r>  
			<w:fldSimple w:instr="AUTHOR">  
				<w:r>  
					<w:t>Author Name (Upper): TOM WATKINS</w:t>  
				</w:r> 
				<w:r>  
					<w:t>Author Name (Lower): tom watkins</w:t>  
				</w:r>  
				<w:r>  
					<w:t>Author Name (Title): Tom Watkins</w:t>  
				</w:r>   
			</w:fldSimple>  
			</w:p>  
		</w:body>  
	</w:document>`

	if removeXmlFormatting(outputXml) != removeXmlFormatting(expectedOutputXml) {
		t.Fatalf("Output does not match expected: %v", outputXml)
	}
}

func TestTableTemplate(t *testing.T) {
	originalXml := `
	<w:document>  
		<w:body>  
			<w:tbl>  
				<w:tblPr>  
					<w:tblW w:w="5000" w:type="pct"/>  
					<w:tblBorders>  
						<w:top w:val="single" w:sz="4" w:space="0" w:color="auto"/>  
						<w:left w:val="single" w:sz="4" w:space="0" w:color="auto"/>  
						<w:bottom w:val="single" w:sz="4" w:space="0" w:color="auto"/>  
						<w:right w:val="single" w:sz="4" w:space="0" w:color="auto"/>  
					</w:tblBorders>  
				</w:tblPr>  
				<w:tblGrid>  
					<w:gridCol w:w="10296"/>  
				</w:tblGrid>  
				{{ range .People }}<w:tr>  
					<w:tc>  
						<w:tcPr>  
							<w:tcW w:w="0" w:type="auto"/>  
						</w:tcPr>  
						<w:p>
							<w:r> 
								<w:t>{{ .Name }}</w:t>
							</w:r>  
						</w:p>
					</w:tc>  
				</w:tr>
				{{ end }}
			</w:tbl>
		</w:body>  
	</w:document>`

	data := map[string]any{
		"People": []map[string]any{
			{"Name": "Tom Watkins"},
			{"Name": "Evie Argyle"},
		},
	}
	outputXml, err := replaceTagsInText(originalXml, data, &defaultFuncMap)
	if err != nil {
		t.Fatalf("Error in basic template %v", err)
	}

	expectedOutputXml := `
	<w:document>  
		<w:body>  
			<w:tbl>  
				<w:tblPr>  
					<w:tblW w:w="5000" w:type="pct"/>  
					<w:tblBorders>  
					<w:top w:val="single" w:sz="4" w:space="0" w:color="auto"/>  
					<w:left w:val="single" w:sz="4" w:space="0" w:color="auto"/>  
					<w:bottom w:val="single" w:sz="4" w:space="0" w:color="auto"/>  
					<w:right w:val="single" w:sz="4" w:space="0" w:color="auto"/>  
					</w:tblBorders>  
				</w:tblPr>  
				<w:tblGrid>  
					<w:gridCol w:w="10296"/>  
				</w:tblGrid>  
				<w:tr>  
					<w:tc>  
						<w:tcPr>  
							<w:tcW w:w="0" w:type="auto"/>  
						</w:tcPr>  
						<w:p>
							<w:r> 
								<w:t>Tom Watkins</w:t>
							</w:r>  
						</w:p> 
					</w:tc>  
				</w:tr>  
				<w:tr>  
					<w:tc>  
						<w:tcPr>  
							<w:tcW w:w="0" w:type="auto"/>  
						</w:tcPr>  
						<w:p>
							<w:r> 
								<w:t>Evie Argyle</w:t>
							</w:r>  
						</w:p> 
					</w:tc>  
				</w:tr>  
			</w:tbl>
		</w:body>  
	</w:document>`

	if removeXmlFormatting(outputXml) != removeXmlFormatting(expectedOutputXml) {
		t.Fatalf("Output does not match expected: %v", outputXml)
	}
}

func removeXmlFormatting(originalXML string) string {
	newXml := strings.ReplaceAll(originalXML, "\n", "")
	newXml = strings.ReplaceAll(newXml, "\r", "")
	newXml = strings.ReplaceAll(newXml, "\t", "")

	newXml = regexp.MustCompile(`>\s+<`).ReplaceAllString(newXml, "><")

	return newXml
}
