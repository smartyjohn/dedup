package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestExec(t *testing.T) {
	r := new(bytes.Buffer)
	w := new(bytes.Buffer)

	line1 := []byte("line1 data\n")
	line2 := []byte("line2 data\n")
	line3 := []byte("line3 data\n")

	r.Write(line1)
	r.Write(line2)
	r.Write(line2)
	r.Write(line1)
	r.Write(line3)
	r.Write(line1)

	exec(r, w)

	actual := w.Bytes()
	expect := []byte(string(line1) + string(line2) + string(line3))
	if bytes.Compare(actual, expect) != 0 {
		t.Fatal(sprintCase(string(actual), string(expect)))
	}
}

func sprintCase(actual, expect interface{}) string {
	return fmt.Sprint("\n\t",
		"[actual]: ", actual, "\n\t",
		"[expect]:", expect, "\n")
}
