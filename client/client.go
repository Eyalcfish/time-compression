package client

import (
	"io"
	"strings"
	"time/comp/time"
)

func Activate(conn io.Writer) {
	data := strings.Repeat("cd", 128)
	time.NewEncoder([]byte(data), conn)
}
