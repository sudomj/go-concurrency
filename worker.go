package main

import (
	"fmt"
	"sync"
	"time"
)

type Message struct {
	Data string
	Time time.Time
}

func worker(id int, msgCh chan Message, tunnel chan string) {

	for msg := range msgCh {
		fmt.Printf("worker %d processing the message\n", id)
		tunnel <- msg.Data
	}
}

func runWorker() {

	numberOfMessages := 100
	numberOfWorkers := 5
	msgCh := make(chan Message, numberOfMessages)
	tunnel := make(chan string, numberOfMessages)

	for index := 0; index < numberOfWorkers; index++ {
		go worker(index, msgCh, tunnel)
	}

	for index := 0; index < numberOfMessages; index++ {
		msgCh <- Message{
			Data: fmt.Sprintf("message %d", index),
		}
	}
	close(msgCh)

	wg := &sync.WaitGroup{}
	wg.Add(numberOfMessages)

	go func() {
		wg.Wait()
		close(tunnel)
	}()

	for msg := range tunnel {

		fmt.Println(msg)
		wg.Done()
	}
}
