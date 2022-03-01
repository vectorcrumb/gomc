package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"lalvarez.me/mcgo/client"
	"lalvarez.me/mcgo/packets"
)

const (
	CONN_HOST        = "localhost"
	CONN_PORT        = "25565"
	CONN_TYPE        = "tcp"
	CONN_BUFFER_SIZE = 8192
	MAX_CLIENTS      = 10
)

func main() {

	clientList := make(map[string]*client.ClientState, 0)

	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening on port", CONN_PORT, "with error:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(2)
		}
		// Get connection IP and check if the client is already connected
		connAddress := strings.Split(conn.RemoteAddr().String(), ":")[0]
		var clientState *client.ClientState
		var ok bool
		if clientState, ok = clientList[connAddress]; ok {
			fmt.Printf("Client %s already connected\n", connAddress)
		} else {
			fmt.Printf("Client %s not connected. Creating a new client state object\n", connAddress)
			clientState = new(client.ClientState)
			clientList[connAddress] = clientState
		}
		go handleRequest(conn, clientState)
	}
}

func handleRequest(conn net.Conn, clientState *client.ClientState) {
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
	// Handle the server list ping
	if clientState.State == client.Uninitialized || clientState.State == client.Handshake {
		handleServerListPing(buf, nbytes)
		clientState.AdvanceServerListPingState()
	}

	// Write back to client
	conn.Write([]byte("Message received.\n"))
	// Kill connection
	fmt.Println()
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
	// If the packet is a handshake packet, parse it
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
