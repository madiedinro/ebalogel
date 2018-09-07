package adapter

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// StartSocket start socket consumer
func StartSocket(listen string) chan string {
	ch := make(chan string)

	go func() {
		serverAddr, err := net.ResolveUDPAddr("udp", listen)

		conn, err := net.ListenUDP("udp", serverAddr)
		if err != nil {
			fmt.Println("UDP Socket / error listening:", err.Error())
			os.Exit(1)
		}
		// Close the listener when the application closes.
		defer conn.Close()
		fmt.Printf("-> UDP Socket listen: %s\n", listen)

		for {
			data, err := bufio.NewReader(conn).ReadString('\n')

			if err != nil {
				fmt.Println("UDP Socket / error reading:", err.Error())
				continue
			}

			ch <- data
		}

	}()
	return ch
}
