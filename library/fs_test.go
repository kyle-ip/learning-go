package library

import (
	"log"
	"os"
	"testing"
)

func TestOpenFile(t *testing.T) {
	f, err := os.Open("fs_test.go") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	defer func() { f.Close() }()
}
