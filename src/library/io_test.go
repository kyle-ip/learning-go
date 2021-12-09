package std_lib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func GetProjectRoot() string {
	binDir, err := executableDir()
	if err != nil {
		return ""
	}
	return path.Dir(binDir)
}

func executableDir() (string, error) {
	pathAbs, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Dir(pathAbs), nil
}

func Utf8Index(str, substr string) int {
	asciiPos := strings.Index(str, substr)
	if asciiPos == -1 || asciiPos == 0 {
		return asciiPos
	}
	pos := 0
	totalSize := 0
	reader := strings.NewReader(str)
	for _, size, err := reader.ReadRune(); err == nil; _, size, err = reader.ReadRune() {
		totalSize += size
		pos++
		// 匹配到
		if totalSize == asciiPos {
			return pos
		}
	}
	return pos
}

func ReadFrom(reader io.Reader, num int) ([]byte, error) {
	p := make([]byte, num)
	n, err := reader.Read(p)
	if n > 0 {
		return p[:n], nil
	}
	return p, err
}

func TestRead(t *testing.T) {

	// read from standard input
	data, _ := ReadFrom(os.Stdin, 11)
	fmt.Println(data)

	// read from file (os.File)
	file, err := os.Open(GetProjectRoot() + "test.txt")
	if err != nil {
		fmt.Println("failed to open file!", err)
		return
	}
	data, err = ReadFrom(file, 9)
	_ = file.Close()
	fmt.Println(data)

	// read from string
	data, err = ReadFrom(strings.NewReader("from string"), 12)
	fmt.Println(data)

}

func TestWrite(t *testing.T) {
	_, err := fmt.Fprintln(os.Stderr, "Hello, World!")
	if err != nil {
		fmt.Println("failed write", err)
		return
	}

	_, err = fmt.Fprintln(os.Stdout, "Hello, World!")
}

func TestReadAt(t *testing.T) {
	reader := strings.NewReader("Learning Go")
	p := make([]byte, 6)
	// read from index 2, 6 bytes.
	n, err := reader.ReadAt(p, 2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s, %d\n", p, n)
}

func TestWriteAt(t *testing.T) {
	file, err := os.Create("test.txt")
	if err != nil {
		panic(err)
	}
	// close and delete file.
	defer func(file *os.File) {
		_ = file.Close()
		_ = os.Remove("test.txt")
	}(file)

	// return the number of characters that written to file.
	n, _ := file.WriteString("Hello, World!")
	fmt.Println(n)
	n, _ = file.WriteAt([]byte("Ywh"), 7)
	fmt.Println(n)

	// if file 'test.txt' not exists, file.Close() will cause panic.
	// defer file.Close() should be placed after error check.
	file, _ = os.Open("test.txt")
	writer := bufio.NewWriter(os.Stdout)
	_, err = writer.ReadFrom(file)
	_ = writer.Flush()
}

func TestSeeker(t *testing.T) {
	reader := strings.NewReader("Hello, World!")

	// return the index of the character.
	n, _ := reader.Seek(-6, io.SeekEnd)
	fmt.Println(n)
	r, _, _ := reader.ReadRune()
	fmt.Printf("%c\n", r)
}

func TestByte(t *testing.T) {
	var ch byte
	_, err := fmt.Scanf("%c\n", &ch)

	buffer := new(bytes.Buffer)
	err = buffer.WriteByte(ch)
	if err == nil {
		fmt.Println("write a byte successfully. then read it: ")
		newCh, _ := buffer.ReadByte()
		fmt.Printf("read byte: %c\n", newCh)
	} else {
		fmt.Println("write failed!", err)
	}
}
