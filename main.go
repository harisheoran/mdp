package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

func main() {
	fileFlag := flag.String("f", "", "Provide Markdown flag")
	flag.Parse()
	data, err := Readfile(*fileFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := ParseFile(data)
		io.WriteString(w, string(html))
	})
	http.ListenAndServe(":3000", nil)

}

func Readfile(file string) ([]byte, error) {
	extension := path.Ext(file)
	if extension != ".md" {
		return nil, fmt.Errorf("Provide a markdown file.")
	}

	_, errOpen := os.Open(file)
	if errOpen != nil {
		return nil, errOpen
	}
	data, errRead := os.ReadFile(file)
	if errRead != nil {
		return nil, errRead
	}

	return data, nil
}

func ParseFile(markdown []byte) []byte {
	output := blackfriday.Run(markdown)
	santizedHTML := bluemonday.UGCPolicy().SanitizeBytes(output)
	return santizedHTML
}

func Templating(markdownHTML []byte) {
	var outputHTMLfilename = "output.html"
	template, err := template.New(outputHTMLfilename).Parse(outputHTMLfilename)
	if err != nil {
		panic(err)
	}
	errEx := template.Execute(os.Stdout, markdownHTML)
	if errEx != nil {
		panic(errEx)
	}
}
