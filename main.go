package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/kshvakov/clickhouse"
	melody "gopkg.in/olahol/melody.v1"
)

type baseMsg struct {
	ContainerID   string    `json:"id"`
	ContainerName string    `json:"name"`
	Time          time.Time `json:"time"`
	Source        string    `json:"source"`
	Data          string    `json:"data"`
}

var (
	dbconn   *sql.DB
	logsstmt *sql.Stmt
	logstx   *sql.Tx
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
	sockListen := getEnv("SOCK_LISTEN", ":8090")
	ChDSN := getEnv("CH_DSN", "")

	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	id := node.Generate()
	fmt.Printf("Generated id:%s\n", id)

	r := gin.Default()
	m := melody.New()

	r.GET("/logs", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	go listenSocket(m, sockListen)

	if ChDSN != "" {
		go initCh(ChDSN)
	}

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
	fmt.Printf("Listening on %s\n", listen)

	for {
		data, err := bufio.NewReader(conn).ReadBytes('\n')
		var rec baseMsg
		if err = json.Unmarshal(data, &rec); err != nil {
			fmt.Println("Error reading:", err.Error())
		}

		fmt.Println(data)
		fmt.Print("Record:", rec)

		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		m.Broadcast(data)
	}
}

func initCh(dsn string) {
	fmt.Println("Initializing CH")
	dbconn, err := sql.Open("clickhouse", dsn)
	if err != nil {
		fmt.Println("Error connecting DB:", err.Error())
		os.Exit(1)
	}

	if err := dbconn.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		return
	}

	fmt.Println("Starting ticker")

	ticker := time.NewTicker(time.Millisecond * 5000)

	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()
	// time.Sleep(time.Millisecond * 1500)
	// ticker.Stop()
	// fmt.Println("Ticker stopped")

	logstx, _ = dbconn.Begin()
	logsstmt, err = logstx.Prepare("INSERT INTO logs (id, date, dateTime, service, source, level, levelName, logger, message, data) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	// id, date, dateTime, service, source, level, levelName, logger, message, data

	if err := logstx.Commit(); err != nil {
		log.Fatal(err)
	}

	defer logsstmt.Close()
}
