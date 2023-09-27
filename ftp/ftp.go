package ftp

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Fprintf(conn, "220 Welcome to the TMS\r\n")
	serialNumber := ""

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		command := scanner.Text()
		parts := strings.Split(command, " ")
		cmd := strings.ToUpper(parts[0])

		fmt.Printf("Received command: %s\n", command)

		switch cmd {
		case "USER":
			fmt.Fprintf(conn, "331 User name okay, need password.\r\n")
		case "PASS":
			fmt.Fprintf(conn, "230 User logged in, proceed.\r\n")
		case "CWD":
			if len(parts) < 2 {
				fmt.Fprintf(conn, "501 Syntax error in parameters or arguments.\r\n")
			} else {
				if isValidSerialNumber(parts[1]) {
					serialNumber = parts[1]
					fmt.Fprintf(conn, "250 Serial number changed to %s\r\n", serialNumber)
				} else {
					fmt.Fprintf(conn, "550 Requested action not taken. Invalid serial number.\r\n")
				}
			}
		case "LIST":
			if serialNumber == "" {
				fmt.Fprintf(conn, "550 Requested action not taken. No serial number selected.\r\n")
			} else {
				files := listFilesBySerialNumber(serialNumber)
				for _, file := range files {
					fmt.Fprintf(conn, "%s\r\n", file)
				}
				fmt.Fprintf(conn, "226 Listing completed.\r\n")
			}
		case "RETR":
			if len(parts) < 2 {
				fmt.Fprintf(conn, "501 Syntax error in parameters or arguments.\r\n")
				break
			}

			fileName := parts[1]
			content, err := getFileBytes(fileName)

			if err != nil {
				fmt.Fprintf(conn, "550 %s\r\n", err.Error()) // You might want to log the error but send a generic error message to the client
				break
			}

			fmt.Fprintf(conn, "150 Opening binary mode data connection for %s.\r\n", fileName)
			_, err = conn.Write(content)

			if err != nil {
				fmt.Fprintf(conn, "426 Connection closed; transfer aborted due to %s.\r\n", err.Error())
				break
			}

			fmt.Fprintf(conn, "226 Transfer complete.\r\n")
		case "QUIT":
			fmt.Fprintf(conn, "221 Goodbye.\r\n")
			conn.Close()
			break
		default:
			fmt.Fprintf(conn, "502 Command not implemented.\r\n")
		}

	}
}
