package main

import (
	"fmt"
	"sync"
	"time"
)

type FileInfo struct {
	buf         []byte
	size        uint64
	lastUseTime time.Time
}

type LogCache struct {
	files map[string]*FileInfo
	sync.RWMutex
}

const (
	CLEAN_USELESS_CACHE_DURATION = time.Duration(300 * time.Second)
	MAX_CACHE_LIVE_SECONS        = 1800
)

func NewLogCache() *LogCache {
	files := make(map[string]*FileInfo)
	return &LogCache{files: files}
}

func (this *LogCache) AutoCleanUselessCache() {
	timer := time.NewTicker(CLEAN_USELESS_CACHE_DURATION)
	for {
		select {
		case <-timer.C:
			go this.cleanUselessCache()
		}
	}
}

func (this *LogCache) cleanUselessCache() {
	this.Lock()
	defer this.Unlock()
	for name, info := range this.files {
		if time.Now().Sub(info.lastUseTime).Seconds() > MAX_CACHE_LIVE_SECONS {
			delete(this.files, name)
		}
	}
}

func (this *LogCache) GetLogCacheSize(fileName string) uint64 {
	this.RLock()
	defer this.RUnlock()
	if _, ok := this.files[fileName]; !ok {
		return 0
	}
	f := this.files[fileName]
	return f.size
}

func (this *LogCache) GetLog(fileName string) []byte {
	this.RLock()
	defer this.RUnlock()
	if _, ok := this.files[fileName]; !ok {
		return []byte{}
	}
	f := this.files[fileName]
	return f.buf
}

func (this *LogCache) AddLogCache(fileName string, buf []byte) (bool, error) {
	if len(buf) == 0 {
		return false, fmt.Errorf("AddLogCache failed: buf is empty")
	}
	if len(fileName) == 0 {
		return false, fmt.Errorf("AddLogCache failed: fileName is empty")
	}

	this.Lock()
	defer this.Unlock()
	if _, ok := this.files[fileName]; !ok || this.files[fileName].size == 0 {
		fileInfo := &FileInfo{
			buf:         buf,
			size:        uint64(len(buf)),
			lastUseTime: time.Now(),
		}
		this.files[fileName] = fileInfo
		return true, nil
	}

	orgFileInfo := this.files[fileName]
	orgFileInfo.buf = append(orgFileInfo.buf, buf...)
	orgFileInfo.lastUseTime = time.Now()
	orgFileInfo.size += uint64(len(buf))
	fmt.Printf("orgFileInfo:%v, this.files:%v\n, size:%d\n", orgFileInfo, this.files, this.files[fileName].size)
	return true, nil
}
