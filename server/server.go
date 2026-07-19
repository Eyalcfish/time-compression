package server

import (
	"fmt"
	"io"
	"time/comp/time"
)

func Activate(conn io.Reader) {
	data := time.NewDecoder(conn)
	fmt.Println("recieved data: " + string(data))
}
