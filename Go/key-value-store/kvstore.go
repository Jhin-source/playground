package main

import (
	"fmt"
	"io"
	"bufio"
	"net"
	"log"
	"strings"
	"time"
	"sync"
)

type Gstore interface {
	SET(key string, value string, id int) 
	GET(key string, id int) string
}

type Kvstore struct {
	store map[string]string
	mu sync.RWMutex
}

func main() {
	var S Gstore = &Kvstore{
		store : make(map[string]string),
	}

	l, err := net.Listen("tcp", ":6380")
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	id := 0
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn, s Gstore, id int) {
			reader := bufio.NewReader(conn)
			for {
				str, err := reader.ReadString('\n')
				if err == io.EOF {
					log.Printf("Client closed connection. Goroutine %d exiting!", id)
					return
				} else if err != nil {
					log.Fatal(err)
				}
				// log.Printf("Goroutine id: %d My first TCP connection: %s", id, str)
				cmd := strings.Fields(str)
				if strings.EqualFold(cmd[0], "SET"){
					S.SET(cmd[1],cmd[2],id)
				}
				if strings.EqualFold(cmd[0], "GET"){
					fmt.Fprintf(conn, "value: %v\n", S.GET(cmd[1],id))
				}
			}
		}(conn, S, id)
		id++
	}
}

func (s *Kvstore) SET(key string, value string, id int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	time.Sleep(15*time.Second)
	s.store[key] = value
	log.Printf("Goroutine %d stored Key: %s, Value %s", id, key, value)
}

func (s *Kvstore) GET(key string, id int) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if id % 2 == 1 {
		time.Sleep(15*time.Second)
	}
	
	ele, ok := s.store[key]
	if ok {
		log.Printf("Element is in Kvstore, Goroutine %d retrieved Key: %s, Value: %s", id, key, ele)
	} else {
		log.Printf("Element is not in Kvstore, Goroutine %d retrieved nothing", id)
	}

	return ele
}