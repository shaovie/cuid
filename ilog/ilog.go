package ilog

import (
	"log"
	"os"
	"sync"
)

var iLog *log.Logger
var lmtx sync.Mutex

func InitLog(logPath string) error {
	lmtx.Lock()
	defer lmtx.Unlock()

	os.Mkdir(logPath, 0755)
	lf, err := os.OpenFile(logPath+"/"+"log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	iLog = log.New(lf, "", log.Ldate|log.Lmicroseconds)
	return nil
}
func Log(format string, v ...interface{}) {
	lmtx.Lock()
	defer lmtx.Unlock()

	if iLog != nil {
		iLog.Printf(format, v...)
	}
}
