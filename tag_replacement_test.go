package docxtpl

import (
	"regexp"
	"strings"
	"testing"
	"text/template"
)

func TestReplaceTagsInText(t *testing.T) {
	tests := []struct {
		name              string
		inputXml          string
		expectedOutputXml string
		data              map[string]any
		funcMap           *template.FuncMap
		expectError       bool
	}{
		{
			name: "Basic template",
			inputXml: `
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
			</w:document>`,
			expectedOutputXml: `
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
			</w:document>`,
			data: map[string]any{
				"Name": "Tom Watkins",
			},
			funcMap:     &defaultFuncMap,
			expectError: false,
		},
		{
			name: "Basic template with function call",
			inputXml: `
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
			</w:document>`,
			expectedOutputXml: `
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
			</w:document>`,
			data: map[string]any{
				"Name": "Tom Watkins",
			},
			funcMap:     &defaultFuncMap,
			expectError: false,
		},
		{
			name: "Template with table",
			inputXml: `
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
			</w:document>`,
			expectedOutputXml: `
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
			</w:document>`,
			data: map[string]any{
				"People": []map[string]any{
					{"Name": "Tom Watkins"},
					{"Name": "Evie Argyle"},
				},
			},
			funcMap:     &defaultFuncMap,
			expectError: false,
		},
		{
			name:              "Invalid template should result in an error",
			inputXml:          "{{ .IncompleteTag",
			expectedOutputXml: "",
			data: map[string]any{
				"Name": "Tom Watkins",
			},
			funcMap:     &defaultFuncMap,
			expectError: true,
		},
		{
			name:              "Error in the execution should result in an error being returned",
			inputXml:          "Function should cause an error: {{ fail }}",
			expectedOutputXml: "",
			data: map[string]any{
				"Name": "Tom Watkins",
			},
			funcMap: &template.FuncMap{
				"fail": func() string {
					panic("forced function error")
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputXml, err := replaceTagsInText(tt.inputXml, tt.data, tt.funcMap)
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}
			if removeXmlFormatting(outputXml) != removeXmlFormatting(tt.expectedOutputXml) {
				t.Fatalf("expected %v but got %v", tt.expectedOutputXml, outputXml)
			}
		})
	}
}

func removeXmlFormatting(originalXML string) string {
	newXml := strings.ReplaceAll(originalXML, "\n", "")
	newXml = strings.ReplaceAll(newXml, "\r", "")
	newXml = strings.ReplaceAll(newXml, "\t", "")

	newXml = regexp.MustCompile(`>\s+<`).ReplaceAllString(newXml, "><")

	return newXml
}
