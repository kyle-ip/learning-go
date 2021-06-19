package main

import (
    "bufio"
    "errors"
    "fmt"
    "io/ioutil"
    "math"
    "os"
    "os/exec"
    "strings"
)

func rgb(i int) (int, int, int) {
    var f = 0.1
    return int(math.Sin(f*float64(i)+0)*127 + 128),
        int(math.Sin(f*float64(i)+2*math.Pi/3)*127 + 128),
        int(math.Sin(f*float64(i)+4*math.Pi/3)*127 + 128)
}

func printLn(str string) {

}

// ErrNoPath is returned when 'cd' was called without a second argument.
var ErrNoPath = errors.New("path required")

// cd changes work directory
func cd(dir string) error {
    // Change the directory and return the error.
    return os.Chdir(dir)
}

// pwd returns a rooted path name corresponding to the
func pwd() error {
    dir, err := os.Getwd()
    fmt.Println(dir)
    return err
}

func ls(args []string) error {
    dir := "./"
    if len(args) > 1 {
        dir = args[1]
    }
    files, err := ioutil.ReadDir(dir)
    for _, f := range files {
        fmt.Printf("%s ", f.Name())
    }
    fmt.Println()
    return err
}

// execStr handle string command by calling operating system shell interface
func execStr(input string) error {

    // Remove the newline character.
    input = strings.TrimSuffix(input, "\n")

    // Split the input separate the command and the arguments.
    args := strings.Split(input, " ")

    // Check for built-in commands.
    switch args[0] {
    case "cd":
        // 'cd' to home with empty path not yet supported.
        if len(args) < 2 {
            return ErrNoPath
        }
        // Change the directory and return the error.
        return cd(args[1])
    case "pwd":
        return pwd()
    case "ls":
        return ls(args)
    case "exit":
        os.Exit(0)
    }

    // Prepare the command to execute.
    cmd := exec.Command(args[0], args[1:]...)

    // Set the correct output device.
    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout

    // Execute the command and return the error.
    return cmd.Run()
}

func main() {

    // Create a reader for standard input stream
    reader := bufio.NewReader(os.Stdin)
    for {
        dir, err := os.Getwd()
        fmt.Printf("%s> ", dir)

        // Read the keyboard input.
        // This operation would block until the first occurrence of delim in the input.
        input, err := reader.ReadString('\n')
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
        }

        // Handle the execution of the input.
        if err = execStr(input); err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
    }
}
