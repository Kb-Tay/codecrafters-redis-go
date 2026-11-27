package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/app/internal/store"
)

func handle(conn net.Conn, store store.Store) {
	defer conn.Close()

	addr := conn.RemoteAddr()
	log.Printf("Client connected: %s", addr)
	defer log.Printf("Client disconnected: %s", addr)

	var req string
	var b [1024]byte

	for {
		n, readErr := conn.Read(b[:])
		req += string(b[:n])

		for len(req) > 0 {
			args, bytesParsed, err := resp.Parse(req)
			log.Println(args)
			if err != nil {
				return
			}
			if bytesParsed == 0 {
				break
			}

			i := 0 
			for i < len(args) {
				if args[i] == "PING" {
					if _, err := conn.Write([]byte("+PONG\r\n")); err != nil {
						return
					}
				}

				if args[i] == "ECHO" {
					// might need consider case where the ECHO arg is not fully parsed
					if _, err := conn.Write([]byte(resp.BulkStr(args[i+1]))); err != nil {
						return
					}

					i += 1
				}

				if args[i] == "SET" {
					store.Set(args[i+1], args[i+2])
					if _, err := conn.Write([]byte("+OK\r\n")); err != nil {
						return
					}
				}

				if args[i] == "GET" {
					v := store.Get(args[i+1])
					if _, err := conn.Write([]byte(resp.BulkStr(v))); err != nil {
						return
					}
				}

				i += 1
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
	
	var store = store.New()

	for {
		conn, err := l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				os.Exit(1)	
			}

			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handle(conn, store)
	}
}
