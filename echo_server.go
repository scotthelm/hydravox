package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

type EchoServer struct {
	Port     string
	listener *net.TCPListener
	done     bool
}

func (es *EchoServer) Start() {
	go func() {
		tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("localhost%s", es.Port))
		checkFatal(err)
		listener, err := net.ListenTCP("tcp", tcpAddr)
		es.listener = listener
		checkFatal(err)
		for es.done == false {
			conn, err := es.listener.Accept()
			if err != nil {
				continue
			}
			go handleClient(conn)
		}
	}()
}
func (es *EchoServer) Stop() {
	es.done = true
}

func checkFatal(err error) {
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Exit(1)
	}
}

func handleClient(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(4 * time.Second))
	request := make([]byte, 128)
	defer conn.Close()
	_, err := conn.Read(request)

	if err != nil {
		fmt.Println("---------------------------")
		fmt.Println(err)
	}

	line := string(request)
	conn.Write([]byte(fmt.Sprintf("You said: %s", line)))
}
