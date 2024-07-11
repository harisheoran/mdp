package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

const (
	data = `
## Markdown Preview CLI
A simple markdown preview CLI tool built in GO language.

## Usage
- mdp -f=Readme.md
`
)

var (
	binary = "mdp"
	mdfile = "test.md"
)

func TestMain(m *testing.M) {
	fmt.Println("Building the tool")

	if runtime.GOOS == "windows" {
		binary += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", binary)
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot build tool %s:%s", binary, err)
		os.Exit(1)
	}

	result := m.Run()
	os.Remove(binary)
	os.Remove(mdfile)
	fmt.Println("Removing the files")
	fmt.Println(result)
}

func TestCLI(t *testing.T) {
	// write data to md file
	var buffer bytes.Buffer
	buffer.WriteString(data)
	err := os.WriteFile(mdfile, buffer.Bytes(), 0644)
	if err != nil {
		t.Errorf(err.Error())
	}

	dir, dirErr := os.Getwd()
	if dirErr != nil {
		t.Fatal(dirErr)
	}

	cmdPath := filepath.Join(dir, binary)
	cmd2 := exec.Command(cmdPath, "-f=test.md")
	cmd2.Run()

	outputHTMLfilename := fmt.Sprintf("%s.html", mdfile)
	result, errREAD := os.ReadFile(outputHTMLfilename)
	if errREAD != nil {
		t.Errorf(errREAD.Error())
	}

	expectedResult := `
<!DOCTYPE html>

<html lang="en">

<head>
  <meta charset="UTF-8" />
  <title>Markdown Preview</title>
  <meta name="viewport" content="width=device-width,initial-scale=1" />
  <meta name="description" content="" />
  <link rel="icon" href="favicon.png">
</head>
<body><h2>Markdown Preview CLI</h2>

<p>A simple markdown preview CLI tool built in GO language.</p>

<h2>Usage</h2>

<ul>
<li>mdp -f=Readme.md</li>
</ul>

</body>
</html>`

	if expectedResult != string(result) {
		t.Errorf("expected result %s and actual result is %s", expectedResult, string(result))
	}

}
