package main

import (
	"log"
	"os"

	nats "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
)

func main() {
	nc, err := nats.Connect("localhost:4222")
	if err != nil {
		log.Printf("ERROR")
	} else {
		log.Printf("OK")
	}
	defer nc.Close()
	subj, _ := "test", []byte("hello")
	dat, err := os.ReadFile("model.json")
	if err != nil {
		log.Fatal("cant open")
	}
	// err = sc.Publish(subj, msg)
	// if err != nil {
	// 	log.Fatalf("Error during publish: %v\n", err)
	// }
	sc, err := stan.Connect("test-cluster", "sender", stan.NatsConn(nc))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: ", err)
	}
	defer sc.Close()
	err = sc.Publish(subj, dat)
	if err != nil {
		log.Fatalf("Error during publish: %v\n", err)
	}
}
