package main

import (
	"log"

	"github.com/shoma07/go-mysql-kinesis/internal/producer"
)

func main() {
	// sr, _ := signal.NewSignalReceiver()
	p, err := producer.NewProducer()
	if err != nil {
		log.Printf("[ERROR]")
		return
	}

	p.Run()

	// select {
	// case n := <-sr.Receive():
	// 	log.Printf("[INFO] receive signal %s, closing\n", n)
	// }
}
