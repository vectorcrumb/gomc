package main

import (
	"fmt"
	"net"
	"os"

	"lalvarez.me/mcgo/packets"
)

const (
	CONN_HOST        = "localhost"
	CONN_PORT        = "25565"
	CONN_TYPE        = "tcp"
	CONN_BUFFER_SIZE = 8192
)

func main() {

	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening on port", CONN_PORT, "with error:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	// Iterate over connections
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(2)
		}
		// Handle connections
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	// Show client IP and port
	fmt.Printf("\nServing %s\n", conn.RemoteAddr().String())
	// Create a buffer for receiving client data
	buf := make([]byte, CONN_BUFFER_SIZE)
	nbytes, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
	} else {
		fmt.Printf("Read %d bytes\n", nbytes)
	}
	handleServerListPing(buf, nbytes)
	// Write back to client
	conn.Write([]byte("Message received.\n"))
	// Kill connection
	conn.Close()
}

/*
1) Handshake (Client -> Server) Max length is 267 bytes

*/

func handleServerListPing(buffer []byte, nbytes int) {
	handshakeLen := packets.HANDSHAKE_PACKET_MAXLEN
	if handshakeLen > nbytes {
		handshakeLen = nbytes
	}
	reducedBuffer := buffer[:handshakeLen]
	fmt.Printf("% 02X\n", reducedBuffer)
	fmt.Printf("%v\n", reducedBuffer)
	fmt.Printf("Got %d bytes for handshake\n\n", len(reducedBuffer))
	// Parse Header Packet
	pktHeader, off, err := packets.ParseHeaderPacket(reducedBuffer)
	if err != nil {
		fmt.Println("Error parsing packet header in server list ping: ", err.Error())
	}
	fmt.Println(pktHeader)

	if pktHeader.PacketID.N == packets.PACKET_ID_HANDSHAKE {
		// Parse Handshake Packet
		pktHandshake, _, err := packets.ParseHandshakePacket(reducedBuffer[off:])
		if err != nil {
			fmt.Println("Error parsing packet handshake in server list ping: ", err.Error())
		}
		fmt.Println(pktHandshake)
	}

}

// func structIterator() {
// 	fields := reflect.VisibleFields(reflect.TypeOf(struct{ vartypes.Varint }{}))
// 	for _, field := range fields {
// 		fmt.Printf("Key: %s\tType: %s\n", field.Name, field.Type)
// 		if field.Type.Name() == "int32" {
// 			fmt.Printf("hi!")
// 		}
// 	}
// }

// func parseHandshake(buffer []byte) (HandshakePacket, error) {
// 	fmt.Printf("%v\n", buffer)
// 	fmt.Printf("Got %d bytes for handshake\n", len(buffer))
// 	var offset int16 = 0
// 	var handshakePacket HandshakePacket

// 	hsFields := reflect.VisibleFields(reflect.TypeOf(struct{ HandshakePacket }{}))
// 	for _, field := range hsFields {
// 		fmt.Printf("\tKey: %s\t Type: %s\n", field.Name, field.Type)
// 		if field.Type.Name() == "PacketHeader" {
// 			parsePacketHeader(buffer)
// 		}
// 		if field.Type.Name() == "Varint" {
// 			fmt.Printf("\t\tField is of type Varint!\n")
// 		}
// 	}

// 	return handshakePacket, nil
// }
