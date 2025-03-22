package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udpServer, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatalf("error udpAddr: %s\n", err)
	}

	conn, err := net.DialUDP(udpServer.Network(), nil, udpServer)
	if err != nil {
		log.Fatalf("error creating conn: %s\n", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error reading: %s\n", err)
		}

		_, err = conn.Write([]byte(data))
		if err != nil {
			log.Fatalf("error writing: %s\n", err)
		}
	}
}
