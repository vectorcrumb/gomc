package client

import (
	"fmt"
	"testing"
)

func TestClientStateInitialCondition(t *testing.T) {
	var state ClientState
	if state.State != Uninitialized {
		t.Errorf("State should be Uninitialized, but is %v", state.State)
	}
}

func TestServerListPingAdvanceState(t *testing.T) {
	var tests = []struct {
		state ServerListPing
		next  ServerListPing
	}{
		{Uninitialized, Handshake},
		{Handshake, Request},
		{Request, Response},
		{Response, Ping},
		{Ping, Pong},
		{Pong, Request},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("State %v", tt.state)
		t.Run(testname, func(t *testing.T) {
			next := AdvanceServerListPingState(tt.state)
			if next != tt.next {
				t.Errorf("Got %v. Expected %v", next, tt.next)
			}
		})
	}
}

func TestClientStateAdvanceServerListPing(t *testing.T) {
	var state ClientState
	state.AdvanceServerListPingState()
	if state.State != Handshake {
		t.Errorf("State should be Handshake, but is %v", state.State)
	}
}
