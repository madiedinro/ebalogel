package adapter

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kshvakov/clickhouse"
)

var (
	dbconn   *sql.DB
	logsstmt *sql.Stmt
	logstx   *sql.Tx
)

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
