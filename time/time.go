package time

import (
	"fmt"
	"io"
	"net"
	"time"
)

type Packet struct {
	IDs []byte
}

func StartBackgroundServer() {
	listener, err := net.Listen("tcp", "127.0.0.1:5234")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}

	fmt.Println("Server: Accepted connection from client!")

	defer conn.Close()
	select {}
}

func GetConnectionReaderAndWriter() (io.Reader, io.Writer) {
	readerConn, writerConn := net.Pipe()

	return readerConn, writerConn
}

func ReadPacket(r io.Reader) Packet {
	buf := make([]byte, 256)

	n, err := r.Read(buf)
	if err != nil {
		fmt.Println("Read error:", err)
		return Packet{}
	}

	return Packet{IDs: buf[:n]}
}

func (p *Packet) Send(w io.Writer) error {
	payloadLen := len(p.IDs)

	buf := make([]byte, payloadLen)

	copy(buf, p.IDs)

	_, err := w.Write(buf)
	return err
}

func currentTime() uint16 {
	now := time.Now()

	intervalsFromSeconds := now.Second() * 100

	intervalsFromNano := now.Nanosecond() / 10000000

	totalIntervals := intervalsFromSeconds + intervalsFromNano

	return uint16(totalIntervals)
}
