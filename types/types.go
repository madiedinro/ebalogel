package types

import "time"

// BaseMsg bla bla
type BaseMsg struct {
	ID            uint64    `json:"id"`
	ContainerID   string    `json:"cid"`
	ContainerName string    `json:"cname"`
	Time          time.Time `json:"time"`
	Source        string    `json:"source"`
	Data          string    `json:"data"`
}
