package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

const internalCommunicationChannel = "127.0.0.1:20202"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("start <My\\Path\\To\\Config>|listen|check|end")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		if len(os.Args) < 3 {
			fmt.Println("Usage: MQTTToDB <start My\\Path\\To\\Config>")
		}
		start(os.Args[2])
	case "listen":
		listen()
	case "check":
		check()
	case "stop":
		stop()
	default:
		fmt.Println("Usage: MQTTToDB <start <My\\Path\\To\\Config>|listen|check|end>")
		os.Exit(1)
	}
}

func listen() {
	fmt.Println("Listening...")
	c, err := net.Dial("tcp", internalCommunicationChannel)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	c.Write([]byte{1})
	buf := make([]byte, 1024)
	for {
		n, err := c.Read(buf[:])
		if err != nil {
			return
		}
		println("Client got:", string(buf[0:n]))
	}
}
func check() bool {
	fmt.Println("Checking...")
	c, err := net.Dial("tcp", internalCommunicationChannel)
	if err != nil {
		fmt.Println("Service not running")
		return false
	}
	defer c.Close()
	c.Write([]byte{2})
	buf := make([]byte, 1)
	c.Read(buf)
	if buf[0] != 1 {
		fmt.Println("Service answered incorrectly")
		return false
	}
	fmt.Println("Service is running")
	return true
}
func stop() {
	fmt.Println("Stopping...")
	c, err := net.Dial("tcp", internalCommunicationChannel)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	c.Write([]byte{3})
}

func start(configPath string) {
	fmt.Println("Starting...")
	l, err := net.Listen("tcp", internalCommunicationChannel)
	if err != nil {
		panic(err)
	}
	var m sync.Mutex
	var cons []*net.Conn
	go core(&m, &cons, configPath)
	for {
		fd, err := l.Accept()
		if err != nil {
			panic(err)
		}
		buf := make([]byte, 1)
		_, err = fd.Read(buf)
		if err != nil {
			return
		}
		switch buf[0] {
		case 1:
			fmt.Println("Add Listen")
			m.Lock()
			cons = append(cons, &fd)
			m.Unlock()
		case 2:
			fmt.Println("Got Checked")
			fd.Write([]byte{1})
		case 3:
			fmt.Println("Exit by command")
			os.Exit(1)
		default:
			fmt.Println("Unknown command")
		}

	}
}

func actionLoop() string {
	time.Sleep(100 * time.Millisecond)
	return fmt.Sprintf("%d\n", rand.Intn(1000))
}
