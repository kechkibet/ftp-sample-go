package main

import (
	"fmt"
	"net"
	"tms/ftp"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:2121")
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("FTP server started on port 2121")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err)
			continue
		}

		// Extract the client's IP address and port
		clientAddr := conn.RemoteAddr()
		clientIP := clientAddr.(*net.TCPAddr).IP
		clientPort := clientAddr.(*net.TCPAddr).Port

		fmt.Printf("Client connected from IP: %s, Port: %d\n", clientIP, clientPort)

		// Now, you can pass conn to your FTP handling function
		go ftp.HandleConnection(conn)
	}
}
