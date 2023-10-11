package internal

import (
	"bufio"
	_ "embed"
	"html/template"
	"io"
	"os"
)

//go:embed template.html
var templateHTML string

func Run(inReader io.Reader, outFilename string) {
	tree := NewTests()
	err := StreamTestOutput(inReader, func(f TestOutputLine) error {
		return tree.AddField(f)
	})
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New("template").Parse(templateHTML)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(outFilename)
	if err != nil {
		panic(err)
	}

	fileWriter := bufio.NewWriter(file)

	templateStruct, err := tree.ToTemplateData()
	if err != nil {
		panic(err)
	}

	err = tmpl.ExecuteTemplate(fileWriter, "template", templateStruct)
	if err != nil {
		panic(err)
	}

	fileWriter.Flush()
}
