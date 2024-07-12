package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

const (
	inputfile  = "./testdata/test1.md"
	goldenFile = "./testdata/test1.md.html"
)

func Testparse(t *testing.T) {
	data, err := os.ReadFile(inputfile)
	if err != nil {
		t.Errorf(err.Error())
	}

	result := ParseFile(data)

	goldenFileData, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Errorf(err.Error())
	}

	if bytes.Equal(result, goldenFileData) {
		t.Errorf("result is %s and expected was %s", result, goldenFileData)
	}

}

func TestRead(t *testing.T) {
	var mockStdOut bytes.Buffer
	err := Readfile(inputfile, &mockStdOut)
	if err != nil {
		t.Errorf("THIS %s", err.Error())
	}

	resultFile := strings.TrimSpace(mockStdOut.String())
	result, errReadResult := os.ReadFile(resultFile)
	if errReadResult != nil {
		t.Errorf(errReadResult.Error())
	}

	expectedResult, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Errorf(err.Error())
	}

	if bytes.Equal(expectedResult, result) {
		t.Errorf("expected result %s and actual result is %s", expectedResult, string(result))
	}

	os.Remove(resultFile)
}
