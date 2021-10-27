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
	ErrSyntax = errors.New("")
)

func main() {
	r, w := setup()
	defer r.Close()
	defer w.Close()

	idx := make(map[uint64]struct{}, 100_000)
	hash64 := fnv.New64()

	scan := bufio.NewScanner(r)
	for scan.Scan() {
		line := scan.Bytes()
		key := hasher(line, hash64)
		if _, exists := idx[key]; !exists {
			idx[key] = struct{}{}
			_, err := w.Write(line)
			fatal(err, "write error:")
		}
	}

	fatal(scan.Err(), "read error")
}

func fatal(err error, prefix string) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, prefix, err)
		os.Exit(1)
	}
}

func hasher(b []byte, hash hash.Hash64) uint64 {
	hash.Reset()
	_, err := hash.Write(b)
	fatal(err, "hasher:")
	return hash.Sum64()
}

func setup() (r *os.File, w io.WriteCloser) {
	w = os.Stdout
	switch len(os.Args) {
	case 1: //stdin
		r = os.Stdin
	case 2: //file
		path := os.Args[1]
		fh, err := os.Open(path)
		fatal(err, "read error:")
		r = fh
	default:
		fatal(ErrSyntax, "bad syntax:")
	}

	return
}
