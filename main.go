package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
	"udp-forward-tcp/queue"
)

var (
	fListenPort        = flag.Uint("p", 0, "udp port on localhost to listen")
	fListenUnixgram    = flag.String("u", "", "unix datagram socket address to listen")
	fMaxDatagramLength = flag.Uint("m", 1450, "max datagram length to receive")

	fTcpAddr = flag.String("s", "", "tcp addr:port to send data")

	q *queue.Queue
)

func main() {
	flag.Parse()

	var conn net.Conn

	if *fListenPort > 0 {
		laddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", *fListenPort))
		if err != nil {
			log.Fatal("Failed to resolve udp address:", err)
		}
		conn, err = net.ListenUDP("udp", laddr)
		if err != nil {
			log.Fatal("Failed to listen:", err)
		}
	} else if *fListenUnixgram != "" {
		laddr, err := net.ResolveUnixAddr("unixgram", *fListenUnixgram)
		if err != nil {
			log.Fatal("Failed to resolve udp datagram address:", err)
		}
		conn, err = net.ListenUnixgram("unixgram", laddr)
		if err != nil {
			log.Fatal("Failed to listen:", err)
		}
	} else {
		log.Fatal("You must listen on an udp port or an unix datagram socket")
	}

	log.Println("Listening...")

	q = queue.NewQueue()

	go sender()

	buf := make([]byte, *fMaxDatagramLength)
	for {
		l, err := conn.Read(buf)
		if err != nil {
			log.Fatal("read error:", err)
		}

		err = q.Push(buf[:l])
		if err != nil {
			log.Println("Failed to write queue:", err)
		}
	}
}

func sender() {
	for {
		conn, err := net.Dial("tcp", *fTcpAddr)
		if err != nil {
			log.Println("Failed to connect tcp endpoint:", err)
			time.Sleep(3 * time.Second)
			continue
		}
		log.Println("tcp connected")

		for {
			item := q.Pop().([]byte)
			_, err = conn.Write(item)
			if err != nil {
				log.Println("Failed to write tcp:", err)
				break
			}
		}
		time.Sleep(time.Second)
	}
}
