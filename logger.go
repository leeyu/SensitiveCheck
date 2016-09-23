package sensitivecheck

import (
        "log"
        "os"
        "fmt"
       )

const (
    LevelTrace = iota
    LevelDebug
    LevelInfo
    LevelWarning
    LevelNotice
    LevelError
    LevelCritical
)


type Clog struct {
    flog *log.Logger
    level int
    }

var cLog *Clog
var log_root = "/home/q/logs/scheck/"

func (l *Clog) SetLevel(lev int) {
    l.level = lev
}

func Newlog(name string, lev int) *Clog {
    if("" == name) {
        name = "notice.log"
        }
    name = log_root+name
    fmt.Println(name)
    logFile,err  := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalln("open file error !")
    }

    cLog = &Clog{}
    cLog.flog = log.New(logFile, "logger: ", log.Lshortfile|log.LstdFlags)

    cLog.SetLevel(lev)

    return cLog
}


func (l *Clog) Notice(str string) {
    if LevelNotice >= l.level {
        l.flog.SetPrefix("[Info]")
        l.flog.Print(str)
    }
}

func (l *Clog) Fatalerr(str string) {
    if LevelNotice >= l.level {
        l.flog.SetPrefix("[Fatal]")
        l.flog.Fatal(str)
    }
}
/*

func main() {
    fileName := ""
    tlog := Newlog(fileName)

    tlog.Notice("this is a notice")
    tlog.Fatalerr("this is error")
    /*
//    var buf bytes.Buffer
    fileName := "scheck.log"
    logFile,err  := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
    defer logFile.Close()
    if err != nil {
        log.Fatalln("open file error !")
    }

    logger := log.New(logFile, "logger: ", log.Lshortfile|log.LstdFlags)
    logger.SetPrefix("[Info]")
    logger.Print("Hello, log file!")
    logger.Println("ln Hello, log file!")

//    fmt.Print(&buf)


    logger.SetPrefix("[Fatal]")
    logger.Fatal("Come with fatal,exit with 1 \n")
    logger.SetPrefix("[Fatal]")
    log.Fatal("Come with fatal,exit with 1 \n")
}
*/
