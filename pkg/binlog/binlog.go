package binlog

import (
	"os"
	"path/filepath"
	"sync"
)

// TODO: add bufio writer
var (
	f       *os.File
	mu      sync.Mutex
	logDir  = "data/log"
	logName = filepath.Join(logDir, "log.bin")
)

func NewContext() (err error) {
	err = os.MkdirAll(logDir, 0777)
	if err != nil {
		return err
	}
	f, err = os.OpenFile(logName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}

	return
}

func Append(data []byte) error {
	mu.Lock()
	defer mu.Unlock()
	_, err := f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func Track() {

}
