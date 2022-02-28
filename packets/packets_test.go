package packets

import (
	"fmt"
	"testing"

	"lalvarez.me/mcgo/vartypes"
)

func TestStringPacketParse(t *testing.T) {
	var tests = []struct {
		buffer []byte
		sp     StringPacket
	}{
		{[]byte{9, 108, 111, 99, 97, 108, 104, 111, 115, 116}, StringPacket{length: vartypes.Varint{N: 9}, s: "localhost"}},
	}

	for _, tt := range tests {
		testname := tt.sp.s
		t.Run(testname, func(t *testing.T) {
			sp, _, err := ParseStringPacket(tt.buffer)
			if err != nil {
				t.Errorf("Got error parsing string: %v", err)
			}
			if sp.length.N != tt.sp.length.N || sp.s != tt.sp.s {
				t.Errorf("Got length %d and string %s. Expected length %d and string %s", sp.length.N, sp.s, tt.sp.length.N, tt.sp.s)
			}
		})
	}
}

func TestHeaderPacketParse(t *testing.T) {
	var tests = []struct {
		buffer []byte
		hp     HeaderPacket
	}{
		{[]byte{16, 0}, HeaderPacket{PacketLength: vartypes.Varint{N: 16}, PacketID: vartypes.Varint{N: 0}}},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("Header L%d-ID%d", tt.hp.PacketLength.N, tt.hp.PacketID.N)
		t.Run(testname, func(t *testing.T) {
			hp, _, err := ParseHeaderPacket(tt.buffer)
			if err != nil {
				t.Errorf("Got error parsing header packet: %v", err)
			}
			if hp.PacketLength.N != tt.hp.PacketLength.N || hp.PacketID.N != tt.hp.PacketID.N {
				t.Errorf("Got packet length %d and ID %d. Expected length %d and ID %d", hp.PacketLength.N, hp.PacketID.N, tt.hp.PacketLength.N, tt.hp.PacketID.N)
			}
		})
	}
}

func TestHandshakePacketParse(t *testing.T) {
	buff := []byte{245, 5, 9, 108, 111, 99, 97, 108, 104, 111, 115, 116, 99, 221, 1}
	hsp := HandshakePacket{
		ProtocolVersion: vartypes.Varint{N: 757},
		ServerAddress:   StringPacket{length: vartypes.Varint{N: 9}, s: "localhost"},
		ServerPort:      25565,
		NextState:       vartypes.Varint{N: 1},
	}
	phsp, nbytes, err := ParseHandshakePacket(buff)
	if err != nil {
		t.Errorf("Got error parsing handshake packet: %v", err)
	}
	if nbytes != len(buff) {
		t.Errorf("Got %d bytes parsed. Expected %d", nbytes, len(buff))
	}
	if phsp.ProtocolVersion.N != hsp.ProtocolVersion.N {
		t.Errorf("Got protocol version %d. Expected %d", phsp.ProtocolVersion.N, hsp.ProtocolVersion.N)
	}
	if phsp.ServerAddress.length.N != hsp.ServerAddress.length.N || phsp.ServerAddress.s != hsp.ServerAddress.s {
		t.Errorf("Got server address length %d and string %s. Expected length %d and string %s", phsp.ServerAddress.length.N, phsp.ServerAddress.s, hsp.ServerAddress.length.N, hsp.ServerAddress.s)
	}
	if phsp.ServerPort != hsp.ServerPort {
		t.Errorf("Got server port %d. Expected %d", phsp.ServerPort, hsp.ServerPort)
	}
	if phsp.NextState.N != hsp.NextState.N {
		t.Errorf("Got next state %d. Expected %d", phsp.NextState.N, hsp.NextState.N)
	}
}

func TestHeaderPacketStringFormat(t *testing.T) {
	hp := HeaderPacket{PacketLength: vartypes.Varint{N: 16}, PacketID: vartypes.Varint{N: 0}}
	expStr := "PACKET Header\n\tLen: 16\n\tID: 0"
	if hp.String() != expStr {
		t.Errorf("Got string:\n%s\nExpected string:\n%s", hp.String(), expStr)
	}
}

func TestStringPacketStringFormat(t *testing.T) {
	sp := StringPacket{length: vartypes.Varint{N: 9}, s: "localhost"}
	expStr := "localhost"
	if sp.String() != expStr {
		t.Errorf("Got string:\n%s\nExpected string:\n%s", sp.String(), expStr)
	}
}
