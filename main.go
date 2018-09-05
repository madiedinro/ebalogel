package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func getEnv(env string, def string) string {
	envVal := os.Getenv(env)
	if envVal != "" {
		return envVal
	}
	return def
}

func main() {

	httpListen := getEnv("HTTP_LISTEN", ":8080")
	sockListen := getEnv("SOCK_LISTEN", ":5000")

	r := gin.Default()
	m := melody.New()

	r.GET("/logs", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	go listenSocket(m, sockListen)

	r.Run(httpListen)
	os.Exit(1)
}

func listenSocket(m *melody.Melody, listen string) {
	serverAddr, err := net.ResolveUDPAddr("udp", listen)

	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer conn.Close()
	fmt.Printf("Listening on %s", ":3333")

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		fmt.Print(message)
		m.Broadcast([]byte(message))
	}
}
