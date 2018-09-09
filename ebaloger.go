package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/madiedinro/ebaloger/adapter"
	melody "gopkg.in/olahol/melody.v1"
)

//ElState conf
type ElState struct {
	LogsListen    string
	HTTPListen    string
	SocketListen  string
	ClickHouseDSN string
	HTTPServer    *gin.Engine
	WSServer      *melody.Melody
	SnowNode      *snowflake.Node
}

func main() {

	s := new(ElState)
	s.HTTPListen = getEnv("HTTP_LISTEN", ":8080")
	s.SocketListen = getEnv("SOCK_LISTEN", ":8090")
	s.LogsListen = "80"
	s.ClickHouseDSN = getEnv("CH_DSN", "8091")

	s.initShowFlake()
	s.startHTTP()
	s.startWS()

	chLogs := adapter.StartLogspout()
	chSocket := adapter.StartSocket(s.SocketListen)

	for {
		select {
		case msg := <-chLogs:
			msg.ID = s.id64()
			data, err := json.Marshal(msg)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
				continue
			}
			s.WSServer.Broadcast(data)
		case msg := <-chSocket:
			s.WSServer.Broadcast([]byte(msg))
		}
	}
}

func (s *ElState) id64() uint64 {
	id := s.SnowNode.Generate()
	return uint64(id.Int64())
}

func (s *ElState) initShowFlake() {
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.SnowNode = node
}

func (s *ElState) startHTTP() {
	s.HTTPServer = gin.Default()
	go s.HTTPServer.Run(s.HTTPListen)
}

func (s *ElState) startWS() {
	s.WSServer = melody.New()
	s.HTTPServer.GET("/ws", func(gc *gin.Context) {
		s.WSServer.HandleRequest(gc.Writer, gc.Request)
	})
}

func toJSON(value interface{}) string {
	bytes, err := json.Marshal(value)
	if err != nil {
		log.Println("error marshalling to JSON: ", err)
		return "null"
	}
	return string(bytes)
}

func getEnv(env string, def string) string {
	envVal := os.Getenv(env)
	if envVal != "" {
		return envVal
	}
	return def
}
