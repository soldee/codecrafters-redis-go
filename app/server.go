package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/commands"
)

func main() {

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer listener.Close()

	Run(listener)
}

func Run(listener net.Listener) {

	config := internal.InitializeConfig(os.Args[1:])
	db := internal.InitializeDB()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			break
		}

		go handleClient(conn, db, config)
	}
}

func handleClient(conn net.Conn, db internal.DB, config internal.Config) {
	defer func() {
		fmt.Println("Closing connection to ", conn.RemoteAddr().String())
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection: ", err.Error())
		}
	}()
	fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr().String())

	buf := make([]byte, 128)

	for {
		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			} else {
				fmt.Println("Error reading request: ", err.Error())
				return
			}
		}
		fmt.Println("Client sent: ", string(buf))

		response := commands.HandleRequest(&buf, db, config)
		fmt.Printf("Writing response to client: %s", string(response))

		_, err = conn.Write(response)
		if err != nil {
			fmt.Println("Error writing to client: ", err.Error())
			return
		}
	}
}
