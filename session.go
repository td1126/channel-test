package main

import (
	"fmt"
	"io"
	"math/rand"
	"sync"
)

type session struct {
	done     chan struct{}
	doneOnce sync.Once
	data     chan int
}

func (sess *session) Serve() {
	go sess.loopRead()
	go sess.loopWrite()
}

func (sess *session) loopRead() { //消费者
	defer func() {
		if err := recover(); err != nil {
			sess.doneOnce.Do(func() { close(sess.done) })
		}
	}()

	var err error
	for {
		select {
		case <-sess.done:
			fmt.Printf("read receive sess.done\n")
			return
		case n := <-sess.data:
			fmt.Printf("read [%d]\n", n)
			if n == 0 {
				goto failed
			}
		default:
		}

		if err == io.ErrUnexpectedEOF || err == io.EOF {
			fmt.Printf("err:%s\n", err)
			goto failed
		}
	}
failed:
	sess.doneOnce.Do(func() { close(sess.done) }) //关闭done，确保仅关闭一次
}

func (sess *session) loopWrite() { //生产者
	defer func() {
		if err := recover(); err != nil {
			sess.doneOnce.Do(func() { close(sess.done) })
		}
	}()

	var err error
	for {
		select {
		case <-sess.done: //接收到关闭信号后，自动退出
			return
		case sess.data <- rand.Intn(100):
		}

		if err != nil {
			goto done
		}
	}
done:
	if err != nil {
		fmt.Printf("sess: loop write failed: %v, %s", err, sess)
	}
}

// ForceClose 强制关闭
func (sess *session) ForceClose() {
	sess.doneOnce.Do(func() { close(sess.done) })
}
