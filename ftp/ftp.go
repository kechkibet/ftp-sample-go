package ftp

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Fprintf(conn, "220 Welcome to the TMS\r\n")
	serialNumber := ""
	isLoggedIn := false

	var dataConn net.Conn // Data connection specific to this FTP session
	var dataListener net.Listener
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		command := scanner.Text()
		parts := strings.Split(command, " ")
		cmd := strings.ToUpper(parts[0])

		fmt.Printf("Received command: %s\n", command)

		switch cmd {
		case "USER":
			if len(parts) < 2 {
				fmt.Fprintf(conn, "501 Syntax error in parameters or arguments.\r\n")
			} else {
				if isValidSerialNumber(parts[1]) {
					serialNumber = parts[1]
					fmt.Fprintf(conn, "331 User name okay, need password.\r\n")
				} else {
					fmt.Fprintf(conn, "530 Not logged in. Invalid user name.\r\n")
				}
			}
		case "PASS":
			if len(parts) < 2 {
				fmt.Fprintf(conn, "501 Syntax error in parameters or arguments.\r\n")
			} else {
				if parts[1] == "test" {
					isLoggedIn = true
					fmt.Fprintf(conn, "230 User logged in, proceed.\r\n")
				} else {
					fmt.Fprintf(conn, "530 Not logged in. Invalid password.\r\n")
				}
			}
		case "PWD":
			if !isLoggedIn {
				fmt.Fprintf(conn, "530 Not logged in.\r\n")
			} else {
				fmt.Fprintf(conn, "257 \"/\" is the current directory.\r\n")
			}
		case "TYPE":
			if len(parts) < 2 {
				fmt.Fprintf(conn, "501 Syntax error in parameters or arguments.\r\n")
			} else {
				typeCode := strings.ToUpper(parts[1])
				if typeCode == "I" { //Passive mode
					fmt.Fprintf(conn, "200 Type set to I.\r\n")
				} else {
					fmt.Fprintf(conn, "504 Command not implemented for that parameter.\r\n")
				}
			}
		case "NLST":
			if !isLoggedIn {
				fmt.Fprintf(conn, "530 Not logged in.\r\n")
			} else {
				if dataConn != nil { // Check if dataConn is set (not nil)
					files := listFilesBySerialNumber(serialNumber)

					// Send the "150 Opening ASCII mode data connection" response
					fmt.Fprintf(conn, "150 Opening ASCII mode data connection for file list.\r\n")

					// Send the list of files over dataConn
					for _, file := range files {
						fmt.Fprintf(dataConn, "%s\r\n", file)
					}

					// Close the data connection when the listing is complete
					dataConn.Close()
					dataConn = nil // Set dataConn to nil
					dataListener.Close()

					// Send the "226 Transfer complete" response to the control connection
					fmt.Fprintf(conn, "226 Transfer complete.\r\n")
				} else {
					fmt.Fprintf(conn, "425 Can't open data connection. Use PASV command first.\r\n")
				}
			}
		case "CWD":
			if len(parts) < 2 {
				fmt.Fprintf(conn, "501 Syntax error in parameters or arguments.\r\n")
			} else {
				serialNumber := strings.TrimLeft(parts[1], "/")
				if isValidSerialNumber(serialNumber) {
					fmt.Fprintf(conn, "250 Serial number changed to %s\r\n", serialNumber)
				} else {
					fmt.Fprintf(conn, "550 Requested action not taken. Invalid serial number.\r\n")
				}
			}
		case "PASV":
			if !isLoggedIn {
				fmt.Fprintf(conn, "530 Not logged in.\r\n")
			} else {
				// Define the port range string
				portRange := "50000-50001" //TODO- repace wit env variable

				// Split the string into min and max parts
				portRangeParts := strings.Split(portRange, "-")
				if len(portRangeParts) != 2 {
					fmt.Fprintf(conn, "500 Internal error: Invalid port range\r\n")
					continue
				}

				// Convert the parts to integers
				minPort, errMin := strconv.Atoi(portRangeParts[0])
				maxPort, errMax := strconv.Atoi(portRangeParts[1])
				if errMin != nil || errMax != nil || minPort > maxPort {
					fmt.Fprintf(conn, "500 Internal error: Invalid port range\r\n")
					continue
				}

				// Generate a random port within the range
				dataPort := minPort + rand.Intn(maxPort-minPort+1)

				// Calculate the IP address in the required format
				//host := "167.99.33.213" // Replace with your server's IP address
				host := "127.0.0.1" // Replace with your server's IP address
				publicIp := "167.99.33.213"

				var err error

				// Start listening on the passive data port
				dataListener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", host, dataPort))
				if err != nil {
					fmt.Println("Error listening on passive data port:", err)
					fmt.Fprintf(conn, "425 Can't open passive data connection\r\n")
					continue
				}
				//defer dataListener.Close() //TODO - this is concerning...

				// Inform the client about the address and port
				portHigh := dataPort / 256
				portLow := dataPort % 256
				pasvResponse := fmt.Sprintf("227 Entering Passive Mode (%s,%d,%d).\r\n", strings.Replace(publicIp, ".", ",", -1), portHigh, portLow)
				fmt.Fprintf(conn, pasvResponse)

				// Accept the client's connection on the data port and store it in the dataConn variable
				dataConn, err = dataListener.Accept()
				if err != nil {
					fmt.Println("Error accepting passive data connection:", err)
					fmt.Fprintf(conn, "425 Can't open passive data connection\r\n")
					continue
				}
				fmt.Println("Data connection accepted sucessfully")
			}
		case "RETR":
			if !isLoggedIn {
				fmt.Fprintf(conn, "530 Not logged in.\r\n")
			} else {
				if dataConn != nil { // Check if dataConn is set (not nil)
					if len(parts) < 2 {
						fmt.Fprintf(conn, "501 Syntax error in parameters or arguments.\r\n")
					} else {
						fileName := parts[1]
						content, err := getFileBytes(fileName)

						if err != nil {
							fmt.Fprintf(conn, "550 %s\r\n", err.Error()) // You might want to log the error but send a generic error message to the client
						} else {
							// Send the "150 Opening binary mode data connection" response
							fmt.Fprintf(conn, "150 Opening binary mode data connection for %s.\r\n", fileName)

							// Send the file content over dataConn
							_, err = dataConn.Write(content)

							if err != nil {
								fmt.Fprintf(conn, "426 Connection closed; transfer aborted due to %s.\r\n", err.Error())
							} else {
								// Send the "226 Transfer complete" response to the control connection
								fmt.Fprintf(conn, "226 Transfer complete.\r\n")
							}
						}
					}

					// Close the data connection after the transfer
					dataConn.Close()
					dataConn = nil // Set dataConn to nil
					dataListener.Close()
				} else {
					fmt.Fprintf(conn, "425 Can't open data connection. Use PASV command first.\r\n")
				}
			}
		case "QUIT":
			fmt.Fprintf(conn, "221 Goodbye.\r\n")
			conn.Close()
			break
		default:
			fmt.Fprintf(conn, "502 Command not implemented.\r\n")
		}
	}
	println("Client Disconnected")
}
