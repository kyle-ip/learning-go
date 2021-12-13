package log

import (
    "io/ioutil"
    "log"
    "os"
    "sync"
)

var (
    errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
    infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)
    loggers  = []*log.Logger{errorLog, infoLog}
    mu       sync.Mutex
)

// log methods
var (
    Error  = errorLog.Println
    Errorf = errorLog.Printf
    Info   = infoLog.Println
    Infof  = infoLog.Printf
)

// log levels
// 日志分级：不同层级以不同颜色区分（info 蓝，error 红）。
const (
    InfoLevel = iota
    ErrorLevel
    Disabled
)

// SetLevel controls log level
// 设置日志层级
func SetLevel(level int) {
    mu.Lock()
    defer mu.Unlock()

    // 三个层级声明为三个变量，通过设置 Output 来控制日志是否打印。
    for _, logger := range loggers {
        logger.SetOutput(os.Stdout)
    }

    if ErrorLevel < level {
        errorLog.SetOutput(ioutil.Discard)
    }
    if InfoLevel < level {
        infoLog.SetOutput(ioutil.Discard)
    }
}
