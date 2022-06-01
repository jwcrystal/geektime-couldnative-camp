package main

import (
	"fmt"
	"sync"
	"time"
)

type Queue struct {
	queue []string
	cond  *sync.Cond
}

func main() {
	q := Queue{
		queue: []string{},
		cond:  sync.NewCond(&sync.Mutex{}),
	}

	go func() {
		for {
			q.Enqueue("a")
			time.Sleep(time.Second * 2)
		}
	}()

	for {
		q.Dequeue()
		time.Sleep(time.Second)
	}
}

func (q *Queue) Enqueue(item string) {
	q.cond.L.Lock() // 指上面的mutex
	defer q.cond.L.Unlock()
	q.queue = append(q.queue, item)
	fmt.Print("putting ", item, " to queue, notify all\n")
	q.cond.Broadcast()
}

func (q *Queue) Dequeue() string {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if len(q.queue) == 0 { // sync.Cond 適合處理這種邊界問題
		fmt.Println("no data available, waiting...")
		q.cond.Wait() // wait until cond.Broadcast
	}
	result := q.queue[0]
	q.queue = q.queue[1:]
	return result
}
