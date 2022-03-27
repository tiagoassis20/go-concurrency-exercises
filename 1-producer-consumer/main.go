//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func producer(ctx context.Context, stream Stream, tweets chan *Tweet) {
	log.Println("PRODUCER")

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in PRODUCER", r)
		}
	}()
	defer log.Println("PRODUCER Close")
	defer close(tweets)

	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			return
		}

		select {
		case <-ctx.Done():
			return
		case tweets <- tweet:
		}

	}
}

func consumer(tweets <-chan *Tweet, done chan<- bool) {
	log.Println("CONSUMER")
	defer func() { done <- true }()
	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}

	log.Println("CONSUMER Close")

}

func main() {
	start := time.Now()
	stream := GetMockStream()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 40*time.Second)
	tweets := make(chan *Tweet, 3)
	done := make(chan bool)
	// Producer
	go producer(ctx, stream, tweets)
	// Consumer
	go consumer(tweets, done)
	go consumer(tweets, done)
	go consumer(tweets, done)

	<-done
	cancel()
	fmt.Printf("Process took %s\n", time.Since(start))
}
