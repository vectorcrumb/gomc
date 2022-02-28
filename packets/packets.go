package packets

import (
	"errors"
	"fmt"

	"lalvarez.me/mcgo/vartypes"
)

const (
	PACKET_ID_HANDSHAKE     = 0x00
	HANDSHAKE_PACKET_MAXLEN = 267
)

type StringPacket struct {
	length vartypes.Varint
	s      string
}

type HeaderPacket struct {
	PacketLength vartypes.Varint
	PacketID     vartypes.Varint
}

type HandshakePacket struct {
	ProtocolVersion vartypes.Varint
	ServerAddress   StringPacket
	ServerPort      uint16
	NextState       vartypes.Varint
}

func (p StringPacket) String() string {
	return p.s
}

func (p HeaderPacket) String() string {
	return fmt.Sprintf("PACKET Header\n\tLen: %v\n\tID: %v", p.PacketLength, p.PacketID)
}

func (p HandshakePacket) String() string {
	return fmt.Sprintf(`PACKET Handshake
	Version: %v
	Address: %s:%v
	NextState: %v`, p.ProtocolVersion, p.ServerAddress.s, p.ServerPort, p.NextState)
}

func ParseStringPacket(buffer []byte) (StringPacket, int, error) {
	var pktString StringPacket
	var offset int = 0
	var off int = 0
	var err error
	pktString.length, off, err = vartypes.ReadVarint(buffer[offset:])
	if err != nil {
		return pktString, 0, errors.New("error parsing string length")
	}
	offset += off
	pktString.s = string(buffer[offset : offset+int(pktString.length.N)])
	offset += int(pktString.length.N)

	return pktString, offset, nil
}

func ParseHeaderPacket(buffer []byte) (HeaderPacket, int, error) {
	var pktHeader HeaderPacket
	var offset int = 0
	// Parse PacketLength
	pktLength, off, err := vartypes.ReadVarint(buffer[offset:])
	if err != nil {
		return pktHeader, 0, errors.New("error parsing packet length")
	}
	offset = offset + off
	// Parse PacketID
	pktID, off, err := vartypes.ReadVarint(buffer[offset:])
	if err != nil {
		return pktHeader, 0, errors.New("error parsing packet ID")
	}
	offset += off
	// Write fields to header
	pktHeader.PacketLength = pktLength
	pktHeader.PacketID = pktID

	return pktHeader, offset, nil
}

func ParseHandshakePacket(buffer []byte) (HandshakePacket, int, error) {
	var pktHandshake HandshakePacket
	var offset int = 0
	var off int = 0
	var err error
	// Parse protocol version
	pktHandshake.ProtocolVersion, off, err = vartypes.ReadVarint(buffer[offset:])
	if err != nil {
		fmt.Println("Error parsing Protocol Version: ", err)
		return pktHandshake, offset, errors.New("error parsing protocol version")
	}
	offset += off
	// Parse server address
	pktHandshake.ServerAddress, off, err = ParseStringPacket(buffer[offset:])
	if err != nil {
		fmt.Println("Error parsing Protocol Version: ", err)
		return pktHandshake, offset, errors.New("error parsing server address")
	}
	offset += off
	// Parse server port
	pktHandshake.ServerPort = uint16(buffer[offset])<<8 | uint16(buffer[offset+1])
	offset += 2
	// Parse next state
	pktHandshake.NextState, off, err = vartypes.ReadVarint(buffer[offset:])
	if err != nil {
		fmt.Println("Error parsing next state: ", err)
		return pktHandshake, offset, errors.New("error parsing next state")
	}
	offset += off

	return pktHandshake, offset, nil
}
