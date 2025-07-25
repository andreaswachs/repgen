package internal

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"time"
)

//go:embed template.html
var templateHTML string

//go:embed daisyui.css
var daisyuiCSS string

//go:embed jquery.js
var jqueryJS string

//go:embed tailwind.js
var tailwindJS string

//go:embed datatables.js
var datatablesJS string

//go:embed datatables.css
var datatablesCSS string

var dependencies = map[string]string{
	"daisyui.css":    daisyuiCSS,
	"jquery.js":      jqueryJS,
	"tailwind.js":    tailwindJS,
	"datatables.js":  datatablesJS,
	"datatables.css": datatablesCSS,
}

func Run(inReader io.Reader, outFilename string) {
	tree := NewTests()
	err := StreamTestOutput(inReader, func(f TestOutputLine) error {
		return tree.AddField(f)
	})
	if err != nil {
		fmt.Printf("An error occurred while parsing the test output: %s\n", err)
		os.Exit(1)
	}

	funcMap := template.FuncMap{
		"sumTests": func(a int, b int, c int) int {
			return a + b + c
		},
		"join": strings.Join,                        // Add join function for templates
		"add":  func(a, b int) int { return a + b }, // Add function for index math
		"len": func(x interface{}) int {
			switch v := x.(type) {
			case []interface{}:
				return len(v)
			case []string:
				return len(v)
			case []*TestDTO:
				return len(v)
			default:
				return 0
			}
		},
		"formatRFC3339": func(t time.Time) string {
			return t.Format(time.RFC3339)
		},
	}

	tmpl, err := template.New("template").Funcs(funcMap).Parse(templateHTML)
	if err != nil {
		fmt.Printf("An error occurred while parsing the report template file: %s\n", err)
		os.Exit(1)
	}

	file, err := os.Create(outFilename)
	if err != nil {
		fmt.Printf("An error occurred while creating the report file: %s\n", err)
		os.Exit(1)
	}

	fileWriter := bufio.NewWriter(file)
	templateStruct, err := tree.ToTemplateData()
	if err != nil {
		fmt.Printf("An error occurred while converting the internal test output representation to template data: %s\n", err)
		os.Exit(1)
	}

	templateStruct.Dependencies = dependencies
	err = tmpl.ExecuteTemplate(fileWriter, "template", templateStruct)
	if err != nil {
		fmt.Printf("An error occurred while assembling the test report: %s\n", err)
		os.Exit(1)
	}

	fileWriter.Flush()
}
