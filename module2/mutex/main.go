package main

import (
	"sync"
	"time"
)

func main() {
	go rLock() // 讀鎖不互斥
	go wLock()
	go Lock()
	time.Sleep(5 * time.Second)
}

func Lock() {
	lock := sync.Mutex{}
	for i := 0; i < 3; i++ {
		lock.Lock()
		defer lock.Unlock()
		println("Lock:", i)
	}
}

func rLock() {
	lock := sync.RWMutex{}
	for i := 0; i < 3; i++ {
		lock.RLock()
		defer lock.RUnlock()
		println("rLock:", i)
	}
}

func wLock() {
	lock := sync.RWMutex{}
	for i := 0; i < 3; i++ {
		lock.Lock()
		defer lock.Unlock()
		println("wLock:", i)
	}
}
