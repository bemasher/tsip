package tsip

import (
	"fmt"
	"io"
)

const (
	DLE = 0x10
	ETX = 0x03
)

type PacketID byte

// Represents a TSIP packet.
type Packet struct {
	ID      byte
	SubCode byte // This field is for convenience only, it is ignored in the Write method.
	Data    []byte
}

// Return a string representation of the packet, trimming the Data field to 24 bytes.
func (p Packet) String() string {
	if len(p.Data) > 24 {
		return fmt.Sprintf("{ID:0x%02X SubCode:0x%02X Data:%02X...}", p.ID, p.SubCode, p.Data[:24])
	}
	return fmt.Sprintf("{ID:0x%02X SubCode:0x%02X Data:%02X}", p.ID, p.SubCode, p.Data)
}

// Read a packet from r, discarding data until a valid packet is read. No
// length checking is done, so be careful to only use this with a trusted data
// source.
func (p *Packet) Read(r io.Reader) (err error) {
	b := make([]byte, 1)

	// Scan until we find a valid packet start. DLE followed by neither DLE or ETX.
	var prevDLE bool
	for {
		_, err = r.Read(b)
		if err != nil {
			return
		}
		if b[0] != DLE && b[0] != ETX && prevDLE {
			break
		}
		prevDLE = b[0] == DLE
	}

	p.ID = b[0]

	var dleCount int

loop:
	for {
		// Get a byte.
		_, err = r.Read(b)
		if err != nil {
			return
		}

		switch b[0] {
		case DLE:
			// If DLE, increment counter.
			dleCount++
		case ETX:
			// Exit if ETX is prefixed by an odd number of DLE's
			if dleCount&1 == 1 {
				break loop
			}
		default:
			// Reset DLE count if we get a non-DLE byte.
			dleCount = 0
		}

		// Only store a DLE if it is an even numbered DLE.
		// Non-DLE bytes will be written anyway since dleCount will be 0.
		if dleCount&1 == 0 {
			p.Data = append(p.Data, b[0])
		}
	}

	// Set the SubCode if we have enough data to do so.
	if len(p.Data) > 0 {
		p.SubCode = p.Data[0]
	}

	return
}

// Write the packet to w, escaping any DLE's encountered in Data.
func (p *Packet) Write(w io.Writer) error {
	if p.ID == DLE || p.ID == ETX {
		return fmt.Errorf("invalid id: %02X", p.ID)
	}

	if _, err := w.Write([]byte{DLE, p.ID}); err != nil {
		return err
	}
	for _, b := range p.Data {
		if b != DLE {
			if _, err := w.Write([]byte{b}); err != nil {
				return err
			}
		} else {
			if _, err := w.Write([]byte{DLE, DLE}); err != nil {
				return err
			}
		}
	}
	if _, err := w.Write([]byte{DLE, ETX}); err != nil {
		return err
	}
	return nil
}
