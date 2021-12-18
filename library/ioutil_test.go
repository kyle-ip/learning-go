package library

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

// NopCloser
func TestNopCloser(t *testing.T) {
	resp, err := http.Get("https://www.baidu.com")
	if err != nil {
		panic(err)
	}
	body := resp.Body
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	s, _ := ioutil.ReadAll(rc)
	fmt.Println(string(s))
}

// ReadAll
func TestReadDir(t *testing.T) {
	dir, _ := os.Getwd()
	listAll(dir, 0)
}

func listAll(path string, curHier int) {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, info := range fileInfos {
		if info.IsDir() {
			for tmpHier := curHier; tmpHier > 0; tmpHier-- {
				fmt.Printf("|\t")
			}
			fmt.Println(info.Name(), "\\")
			listAll(path+"/"+info.Name(), curHier+1)
		} else {
			for tmpHier := curHier; tmpHier > 0; tmpHier-- {
				fmt.Printf("|\t")
			}
			fmt.Println(info.Name())
		}
	}
}

// ReadFile
func TestReadFile(t *testing.T) {
	dir, _ := os.Getwd()
	ret, _ := ioutil.ReadFile(fmt.Sprintf("%s%cioutil_test.go", dir, os.PathSeparator))
	fmt.Println(string(ret))
}

// WriteFile
func TestWriteFile(t *testing.T) {
	dir, _ := os.Getwd()
	file := fmt.Sprintf("%s%cioutil_test.go", dir, os.PathSeparator)
	ret, _ := ioutil.ReadFile(file)
	ioutil.WriteFile(file, ret, 0666)
}

// TempDir
func TestTempDir(t *testing.T) {
	dir, _ := os.Getwd()
	f, _ := ioutil.TempDir(dir, "tmp")
	defer func() {
		_ = os.Remove(f)
	}()
	fmt.Println(f)
	time.Sleep(10 * time.Duration(time.Second))
}

// TempFile
func TestTempFile(t *testing.T) {
	dir, _ := os.Getwd()
	f, _ := ioutil.TempFile(dir, "tmp")
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()
	fmt.Println(f)
	time.Sleep(10 * time.Duration(time.Second))
}

// Discard
