package time

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

type Encoder struct {
	idMap   map[byte][]byte
	timeMap map[byte]uint16
	conn    io.Writer
}

func NewEncoder(data []byte, conn io.Writer) {
	fmt.Println("sending data: " + string(data))
	e := &Encoder{}
	e.conn = conn
	e.idMap = make(map[byte][]byte)
	e.timeMap = make(map[byte]uint16)
	e.divideData(data)
	e.sendAllIDS()
	go e.sendNextBatch()
	for {
		if len(e.timeMap) == 0 {
			time.Sleep(time.Millisecond)
			return
		}
	}
}

func (e *Encoder) divideData(data []byte) {
	offset := int(len(data) / 256)
	for i := range 256 {
		e.idMap[byte(i)] = data[i*offset : (i+1)*offset]
	}
}

func (e *Encoder) sendAllIDS() {
	p := make([]byte, 256)
	currentTime := currentTime()
	for i := range 256 {
		p[i] = byte(i)
		e.timeMap[byte(i)] = currentTime
	}
	pack := Packet{IDs: p}
	pack.Send(e.conn)
	// fmt.Println("sent all IDs")
}

func (e *Encoder) sendNextBatch() {
	p := Packet{IDs: make([]byte, 0, len(e.timeMap))}
	currentTime := currentTime()
	for i, t := range e.timeMap {
		var num uint16
		val := e.idMap[byte(i)]
		if len(val) == 1 {
			num = uint16(val[0])
		} else if len(val) >= 2 {
			num = binary.BigEndian.Uint16(val[:2])
		} else {
			num = 0
		}
		// fmt.Println(i, t, num, currentTime-t)
		if currentTime-t >= num {
			p.IDs = append(p.IDs, byte(i))   // add new correct ID
			e.timeMap[byte(i)] = currentTime // update time map
			if len(val) <= 2 {
				delete(e.timeMap, byte(i)) // remove the ID from the map if it has no more data to send
			} else {
				e.idMap[byte(i)] = val[2:] // propagate the data required to send for the specific ID
			}

		}
	}

	if len(p.IDs) != 0 {
		// fmt.Println("sent batch of", len(p.IDs), "IDs")
		p.Send(e.conn)
	}

	if len(e.timeMap) == 0 {
		p.IDs = []byte{255}
		p.Send(e.conn) // signal end
		return         // exit the loop if there are no more IDs to send
	}

	time.Sleep(time.Millisecond)
	e.sendNextBatch()
}
