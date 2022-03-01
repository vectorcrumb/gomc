package client

type ServerListPing int

const (
	Uninitialized ServerListPing = iota
	Handshake
	Request
	Response
	Ping
	Pong
)

type ClientState struct {
	ClientAddress string
	State         ServerListPing
}

func (slp ServerListPing) String() string {
	switch slp {
	case Uninitialized:
		return "Uninitialized"
	case Handshake:
		return "Handshake"
	case Request:
		return "Request"
	case Response:
		return "Response"
	case Ping:
		return "Ping"
	case Pong:
		return "Pong"
	default:
		return "Unknown"
	}
}

func (state *ClientState) AdvanceServerListPingState() {
	state.State = AdvanceServerListPingState(state.State)
}

func AdvanceServerListPingState(state ServerListPing) ServerListPing {
	switch state {
	case Uninitialized:
		return Handshake
	case Handshake:
		return Request
	case Request:
		return Response
	case Response:
		return Ping
	case Ping:
		return Pong
	case Pong:
		return Request
	default:
		return Uninitialized
	}
}
