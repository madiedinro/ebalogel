package adapter

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/gliderlabs/logspout/router"
	"github.com/madiedinro/ebaloger/types"
)

// Version is the running version of logspout
var Version string

var ch = make(chan types.BaseMsg)

// StartLogspout start socket consumer
func StartLogspout() chan types.BaseMsg {

	os.Setenv("ROUTESPATH", "data/")

	router.AdapterFactories.Register(newLogspoutAdapter, "ebaloger")

	fmt.Printf("# logspout %s by gliderlabs\n", Version)
	fmt.Printf("# adapters: %s\n", strings.Join(router.AdapterFactories.Names(), " "))
	fmt.Printf("# options : ")
	fmt.Printf("backlog:%s ", getEnv("BACKLOG", ""))
	fmt.Printf("persist:%s\n", getEnv("ROUTESPATH", "/mnt/routes"))

	route := new(router.Route)
	route.Adapter = "ebaloger"
	router.Routes.Add(route)

	var jobs []string
	for _, job := range router.Jobs.All() {
		err := job.Setup()
		if err != nil {
			fmt.Println("!!", err)
			os.Exit(1)
		}
		if job.Name() != "" {
			jobs = append(jobs, job.Name())
		}
	}
	fmt.Printf("# jobs    : %s\n", strings.Join(jobs, " "))

	routes, _ := router.Routes.GetAll()
	if len(routes) > 0 {
		fmt.Println("# routes  :")
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintln(w, "#   ADAPTER\tADDRESS\tCONTAINERS\tSOURCES\tOPTIONS")
		for _, route := range routes {
			fmt.Fprintf(w, "#   %s\t%s\t%s\t%s\t%s\n",
				route.Adapter,
				route.Address,
				route.FilterID+route.FilterName+strings.Join(route.FilterLabels, ","),
				strings.Join(route.FilterSources, ","),
				route.Options)
		}
		w.Flush()
	} else {
		fmt.Println("# routes  : none")
	}

	for _, job := range router.Jobs.All() {
		job := job
		go func() {
			log.Fatalf("%s ended: %s", job.Name(), job.Run())
		}()
	}

	return ch
}

// NewLogspoutAdapter returns a configured raw.Adapter
func newLogspoutAdapter(route *router.Route) (router.LogAdapter, error) {
	return &Adapter{
		channel: ch,
	}, nil
}

// Adapter is a simple adapter that streams log output to a connection without any templating
type Adapter struct {
	channel chan types.BaseMsg
}

// Stream sends log data to a connection
func (a *Adapter) Stream(logstream chan *router.Message) {
	for message := range logstream {
		a.channel <- types.BaseMsg{
			ContainerID:   message.Container.ID,
			ContainerName: message.Container.Name,
			Data:          message.Data,
			Source:        message.Source,
			Time:          message.Time,
		}
	}
}
