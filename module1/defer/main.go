package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// defer 表示推延，相當於把指令加入stack
	// defer 最後執行
	defer fmt.Println("1")
	defer fmt.Println("2")
	defer fmt.Println("3")
	loopFunc()
	time.Sleep(time.Second)
}

func loopFunc() {
	lock := sync.Mutex{}
	for i := 0; i < 3; i++ {
		go func(i int) {
			lock.Lock()
			// !! defer是在"整個func退出執行"
			// 意味著，loop 1 加鎖 -> loop 2因還沒解鎖就要加鎖，進了等待，造成"deadlock"
			defer lock.Unlock()
			fmt.Println("loopFunc: ", i)
		}(i)
	}
}
