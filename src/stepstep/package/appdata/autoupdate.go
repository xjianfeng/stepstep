package appdata

import (
	"bufio"
	"io"
	"os"
	"strings"
	"time"
)

var fileName = "autoupdate.txt"
var fileInfo = make(map[string]interface{})

func updateFile() {
	f, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		line = strings.Replace(line, " ", "", -1)
		line = strings.Replace(line, "\n", "", -1)
		if line == "" {
			return
		}
		if err != nil || io.EOF == err {
			return
		}
		UpDateGameData(line)
	}
}

func checkFile() {
	time.AfterFunc(10*time.Second, checkFile)

	stat, err := os.Stat(fileName)
	if err != nil {
		return
	}

	oldTime, ok := fileInfo[fileName]
	t := stat.ModTime()
	if !ok {
		updateFile()
	} else if oldTime != t {
		updateFile()
	}
	fileInfo[fileName] = t
}

func init() {
	time.AfterFunc(10*time.Second, checkFile)
}
