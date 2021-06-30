package util

import (
	"os"
	"path"
	"path/filepath"
	"strings"
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
