package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Filter struct {
	LogCache *LogCache
	sync.RWMutex
}

const (
	LOG_ALL   = 0 //log all message
	LOG_INFO  = 1 //log info message
	LOG_WARN  = 2
	LOG_DEBUG = 3
	LOG_ERROR = 4
)

func NewFilter() *Filter {
	logCache := NewLogCache()
	return &Filter{LogCache: logCache}
}

func (this *Filter) StartCleanCache() {
	go this.LogCache.AutoCleanUselessCache()
}

func (this *Filter) GetFileContString(nIndex, logType uint8, more int) (string, error) {
	this.RLock()
	defer this.RUnlock()
	buf, err := this.GetFileContByNode(nIndex)
	if err != nil {
		return "", err
	}
	if logType == 0 {
		return string(buf), nil
	}

	lineArray := strings.Split(string(buf), "\n")

	var ret []string
	var lineNums []int

	keywords := []string{"", "[INFO ]", "[WARN ]", "[DEBUG]", "[ERROR]"}
	keyword := keywords[logType]

	for i, line := range lineArray {
		if strings.Index(line, keyword) >= 0 {
			// more: append from (i - more)  to i line to result
			if logType != LOG_ALL && logType != LOG_INFO && more > 0 {
				var start int
				if len(lineNums) > 0 {
					// has appended lines
					lastAddLine := lineNums[len(lineNums)-1]
					if i-lastAddLine <= more {
						start = lastAddLine + 1
					} else {
						start = i - more
					}
				} else {
					if i-more > 0 {
						start = i - more
					} else {
						start = 0
					}
				}
				for j := start; j <= i; j++ {
					ret = append(ret, lineArray[j])
				}
				lineNums = append(lineNums, i)
			} else {
				ret = append(ret, line)
			}
		}
	}
	return strings.Join(ret, "\n"), nil
}

func (this *Filter) GetFileContByNode(i uint8) ([]byte, error) {
	this.RLock()
	defer this.RUnlock()
	path := "../../node" + strconv.Itoa(int(i)) + "/Log"
	f, err := getNodeLatestFile(i)
	if err != nil {
		return []byte{}, err
	}
	fileName := f.Name()
	fPath := path + "/" + fileName

	fmt.Printf("fPath :%s\n", fPath)

	file, err := os.Open(fPath)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	cacheFileSize := this.LogCache.GetLogCacheSize(fileName)
	if cacheFileSize == (uint64)(f.Size()) {
		return this.LogCache.GetLog(fileName), nil
	}
	buf := make([]byte, f.Size()-(int64)(cacheFileSize))
	_, err = file.ReadAt(buf, (int64)(cacheFileSize))
	if err != nil {
		log.Fatal("file readat failed", file.Name(), err)
	}
	_, err = this.LogCache.AddLogCache(fileName, buf)
	if err != nil {
		return []byte{}, err
	}
	return this.LogCache.GetLog(fileName), nil
}

func (this *Filter) LogTypeFromString(logType string) uint8 {
	logType = strings.Replace(logType, " ", "", -1)
	lowType := strings.ToLower(logType)
	switch lowType {
	case "all":
		return LOG_ALL
	case "info":
		return LOG_INFO
	case "warn":
		return LOG_WARN
	case "error":
		return LOG_ERROR
	case "debug":
		return LOG_DEBUG
	default:
		return LOG_ALL
	}
}

// getNodeLatestFile get the lastest file info by node index
func getNodeLatestFile(i uint8) (os.FileInfo, error) {
	path := "../../node" + strconv.Itoa(int(i)) + "/Log"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var f os.FileInfo
	for i := len(files) - 1; i > 0; i-- {
		f = files[i]
		if !f.IsDir() && f.Size() > 0 {
			break
		}
	}
	return f, nil
}
