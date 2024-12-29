package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

/*
 * Complete the 'TCPServer' function below.
 *
 * The function accepts chan bool ready as a parameter.
 */

func TCPServer(ready chan bool) {
	// Create a TCP address to listen on
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		panic(err)
	}

	// Create a TCP listener
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	// Signal that the server is ready to accept connections
	ready <- true

	// Accept and handle connections
	for {
		// Accept a new connection
		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}

		// Handle each connection in a goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn *net.TCPConn) {
	defer conn.Close()

	// Create buffer to read the message
	buffer := make([]byte, maxBufferSize)

	// Read from the connection
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}

	// Get the message and reverse it
	message := string(buffer[:n])
	reversed := reverseString(message)

	// Write the reversed string back to the connection
	_, err = conn.Write([]byte(reversed))
	if err != nil {
		return
	}
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

const maxBufferSize = 1024
const address = "127.0.0.1:3333"

func tcp_server() {
	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)
	writer := bufio.NewWriterSize(stdout, 16*1024*1024)

	messagesCount, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
	checkError(err)

	var messages []string

	for i := 0; i < int(messagesCount); i++ {
		messagesItem := readLine(reader)
		messages = append(messages, messagesItem)
	}

	ready := make(chan bool)
	go TCPServer(ready)
	<-ready
	reversed, err := tcpClient(messages)
	if err != nil {
		panic(err)
	}
	for _, msg := range reversed {
		fmt.Fprintf(writer, "%s\n", msg)
	}
	writer.Flush()
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func tcpClient(messages []string) ([]string, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return []string{}, err
	}

	reversed := []string{}

	for _, msg := range messages {

		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			return []string{}, err
		}
		_, err = conn.Write([]byte(msg))
		if err != nil {
			return []string{}, err
		}

		reply := make([]byte, maxBufferSize)

		n, err := conn.Read(reply)
		if err != nil {
			return []string{}, err
		}

		reversed = append(reversed, string(reply[:n]))
		conn.Close()
	}

	return reversed, nil
}
