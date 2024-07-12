package main

import (
	"bufio"
	"log"
	"net"
	"sync"
)

func main() {
	addr := "localhost:9000"
	sock, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer sock.Close()

	log.Println("Started listening on:", addr)

	stream := make(chan []byte)
	clients := make(map[*net.Conn]bool)
	var mutex sync.Mutex

	go func() { // broadcast messages
		for {
			msg := <-stream
			mutex.Lock()
			for c := range clients {
				(*c).Write(msg)
			}
			mutex.Unlock()
		}
	}()

	for {
		conn, err := sock.Accept()
		if err != nil {
			continue
		}

		mutex.Lock()
		clients[&conn] = true
		mutex.Unlock()

		go func(conn net.Conn) {
			defer conn.Close()

			id := conn.RemoteAddr().String()

			reader := bufio.NewReader(conn)

			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					continue
				}

				stream <- []byte(id + ": " + line)
			}
		}(conn)
	}
}
