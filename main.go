package main

import (
	"fmt"
	"log"
	"syscall"
)

func main() {
	socket, sockErr := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if sockErr != nil {
		log.Fatalln("Error Creating socket:", sockErr)
	}

	defer syscall.Close(socket)

	// allows the program to reuse a port that is in TIME_WAIT state
	reuseErr := syscall.SetsockoptInt(socket, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if reuseErr != nil {
		log.Fatalln("Error setting SO_REUSEADDR:", reuseErr)
	}

	sockAddr := syscall.SockaddrInet4{Port: 8080}
	copy(sockAddr.Addr[:], []byte{127, 0, 0, 1})

	bindErr := syscall.Bind(socket, &sockAddr)
	if bindErr != nil {
		log.Fatalln("Error Binding Address to Socket:", bindErr)
	}

	listenErr := syscall.Listen(socket, 5)
	if listenErr != nil {
		log.Fatalln("Error Listening on socket:", listenErr)
	}

	log.Println("Server Listening")

	for {
		conn, _, connErr := syscall.Accept(socket)
		if connErr != nil {
			log.Println("Error Accepting Connection:", connErr)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(handle int) {
	buffer := make([]byte, 1024)

	n, err := syscall.Read(handle, buffer)
	if err != nil {
		log.Println("Error reading from connection:", err)
		syscall.Close(handle)
		return
	}

	if n > 0 {

		fmt.Println(string(buffer[:n]))

		message := "{some:\"Hello, world!\"}"

		length := len(message)

		response := "HTTP/1.1 200 OK\r\n" +
			"Content-Type: application/json\r\n" +
			"Content-Length: " + string(length) + "\r\n" +
			"\r\n" +
			message

		_, err := syscall.Write(handle, []byte(response))
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			syscall.Close(handle)
			return
		}

		syscall.Close(handle)
		return
	}
}
