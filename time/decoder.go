package time

import (
	"io"
)

type Decoder struct {
	timeMap        map[byte]uint16
	decodedMessage []byte
}

func NewDecoder(conn io.Reader) []byte {
	d := &Decoder{}

	d.decodedMessage = make([]byte, 256)

	for {
		p := ReadPacket(conn)
		if d.Decode(p) != nil {
			return d.decodedMessage
		}
	}
}

func (d *Decoder) Decode(p Packet) []byte {
	if p.IDs[0] == 255 {
		return d.decodedMessage
	}

	t := currentTime()

	if d.timeMap == nil {
		d.timeMap = make(map[byte]uint16)
		for id := range 256 {
			d.timeMap[byte(id)] = t
		}
		return nil
	}

	for _, id := range p.IDs {
		data := t - d.timeMap[id]
		d.decodedMessage[id] = byte(data) + 1
		// d.decodedMessage[id+1] = byte(data>>8) + 1
	}
	return nil
}
