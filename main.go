package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	defaultTemplate = `<!DOCTYPE html>​
​ 	​<html>​
​ 	​  <head>​
​ 	​    <meta http-equiv="content-type" content="text/html; charset=utf-8">​
​ 	​    <title>{{ .Title }}</title>​
​ 	​  </head>​
​ 	​  <body>​
​ 	​{{ .Body }}​
​ 	​  </body>​
​ 	​</html>
`
)

type content struct {
	Title string
	Body  template.HTML
}

func main() {
	fileFlag := flag.String("f", "", "Provide Markdown flag")
	skipFlag := flag.Bool("s", false, "Skip the preview")
	tempFile := flag.String("t", "", "Provide HTML template file")
	flag.Parse()

	if *fileFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	err := Readfile(*fileFlag, os.Stdout, *skipFlag, *tempFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func Readfile(file string, output io.Writer, skipPreview bool, tName string) error {
	extension := path.Ext(file)
	if extension != ".md" {
		return fmt.Errorf("Provide a markdown file.")
	}

	_, errOpen := os.Open(file)
	if errOpen != nil {
		return errOpen
	}
	data, errRead := os.ReadFile(file)
	if errRead != nil {
		return errRead
	}

	parsedData, errParse := ParseFile(data, tName)
	if errParse != nil {
		return errParse
	}
	tempFile, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}

	errClose := tempFile.Close()
	if errClose != nil {
		return errClose
	}

	fmt.Fprintln(output, tempFile.Name())

	errSave := SaveHTML(parsedData, tempFile.Name())
	if errSave != nil {
		return errSave
	}
	if skipPreview {
		return nil
	}

	defer os.Remove(tempFile.Name())

	return Preview(tempFile.Name())
}

func ParseFile(markdown []byte, tName string) ([]byte, error) {
	output := blackfriday.Run(markdown)
	santizedHTML := bluemonday.UGCPolicy().SanitizeBytes(output)

	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	if tName != "" {
		t, err = template.ParseFiles(tName)
		if err != nil {
			return nil, err
		}
	}

	content := content{
		Title: "Markdown Preview Tool",
		Body:  template.HTML(santizedHTML),
	}

	var buffer bytes.Buffer
	errExe := t.Execute(&buffer, content)
	if errExe != nil {
		return nil, errExe

	}

	return buffer.Bytes(), nil
}

func SaveHTML(data []byte, filename string) error {
	return os.WriteFile(filename, []byte(data), 0644)
}

func Preview(filename string) error {
	cName := ""

	cParams := []string{}

	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}

	cParams = append(cParams, filename)

	cPath, err := exec.LookPath(cName)

	if err != nil {
		return err
	}

	errRun := exec.Command(cPath, cParams...).Run()

	time.Sleep(2 * time.Second)

	return errRun

}
