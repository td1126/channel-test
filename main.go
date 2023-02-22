package main

import (
	"time"
)

func main() {

	sess := &session{data: make(chan int, 100),
		done: make(chan struct{}),
	}

	for i := 0; i < 10; i++ {
		sess.Serve()
	}

	time.Sleep(10000 * time.Millisecond)
	sess.ForceClose()
}
