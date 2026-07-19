package main

import (
	"time"
	"time/comp/client"
	"time/comp/server"
	times "time/comp/time"
)

func main() {
	go times.StartBackgroundServer()

	time.Sleep(100 * time.Millisecond)
	read, write := times.GetConnectionReaderAndWriter()
	go server.Activate(read)
	client.Activate(write)
	time.Sleep(1000 * time.Millisecond)
}
