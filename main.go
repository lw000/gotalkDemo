package main

import (
	"github.com/rsms/gotalk"
	"log"
	"os"
	"os/signal"
	"time"
)

type GreetIn struct {
	Name string `json:"name"`
}
type GreetOut struct {
	Greeting string `json:"greeting"`
}

func main() {
	go server()
	time.AfterFunc(time.Second, func() {
		go client()
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Kill, os.Interrupt)
	<-c
}

func server() {
	gotalk.Handle("greet", func(in GreetIn) (GreetOut, error) {
		return GreetOut{"Hello " + in.Name}, nil
	})
	if err := gotalk.Serve("tcp", "localhost:1234", func(sock *gotalk.Sock) {

	}); err != nil {
		log.Fatalln(err)
	}
}

func client() {
	s, err := gotalk.Connect("tcp", "localhost:1234")
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer func() {
		s.Close()
	}()
	i := 0
	for i < 100 {
		i++
		greeting := &GreetOut{}
		if err := s.Request("greet", GreetIn{"Rasmus"}, greeting); err != nil {
			log.Fatalln(err)
		}
		log.Printf("greeting: %+v\n", greeting)
	}
}
