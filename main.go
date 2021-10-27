package main

import (
	"bufio"
	"errors"
	"fmt"
	"hash"
	"hash/fnv"
	"io"
	"os"
)

var (
	ErrSyntax = errors.New("\nDe-duplicates lines from input, writing to standard out.  If an argument is given, input is read from that path; otherwise data is read from stdin.")
)

//goland:noinspection GoUnhandledErrorResult
func main() {
	r, w := setup()
	defer r.Close()
	defer w.Close()

	exec(r, w)
}

func exec(r io.Reader, w io.Writer) {
	idx := make(map[uint64]struct{}, 100_000)
	hash64 := fnv.New64()
	eol := []byte("\n") //TODO, use same as input

	out := bufio.NewWriter(w)
	defer func(out *bufio.Writer) {
		fatal(out.Flush(), "write error:")
	}(out)

	scan := bufio.NewScanner(r)
	dupes := 0
	for scan.Scan() {
		line := scan.Bytes()
		key := hasher(line, hash64)
		if _, exists := idx[key]; !exists {
			idx[key] = struct{}{}
			fatalWrite(out.Write(line))
			fatalWrite(out.Write(eol))
		} else {
			dupes++
		}
	}

	fatal(scan.Err(), "read error:")
	fatalWrite(fmt.Fprintln(os.Stderr, "Duplicates removed:", dupes))
}

func fatal(err error, prefix string) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, prefix, err)
		os.Exit(1)
	}
}

func fatalWrite(n int, err error) {
	fatal(err, "write error:")
}

func hasher(b []byte, hash hash.Hash64) uint64 {
	hash.Reset()
	_, err := hash.Write(b)
	fatal(err, "hasher error:")
	return hash.Sum64()
}

func setup() (r io.ReadCloser, w io.WriteCloser) {
	w = os.Stdout
	switch len(os.Args) {
	case 1: //stdin
		r = os.Stdin
	case 2: //file
		path := os.Args[1]
		fh, err := os.Open(path)
		fatal(err, "input file error:")
		r = fh
	default:
		fatal(ErrSyntax, "bad syntax:")
	}

	return
}
