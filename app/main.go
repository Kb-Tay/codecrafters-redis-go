package main

import (
	"fmt"
	"net"
	"os"

	respprotocol "github.com/codecrafters-io/redis-starter-go/app/resp_protocol"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func handle(conn net.Conn) {
	defer conn.Close()

	var req string
	var b [1024]byte
	for {
		n, readErr := conn.Read(b[:])
		req += string(b[:n])

		for len(req) > 0 {
			args, bytesParsed, err := respprotocol.Parse(req)
			if err != nil {
				return
			}
			if bytesParsed == 0 {
				break
			}

			for _, arg := range(args) {
				if arg == "PING" {
					if _, err := conn.Write([]byte("+PONG\r\n")); err != nil {
						return
					}
				}
			}

			req = req[bytesParsed:]
		}
		if readErr != nil {
			return
		}
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment the code below to pass the first stage
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	handle(conn)
}
