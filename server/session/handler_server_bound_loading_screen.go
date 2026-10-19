package session

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"sync/atomic"
)

// ServerBoundLoadingScreenHandler handles loading screen updates from the clients. It is used to ensure that
// the server knows when the client is loading a screen, and when it is done loading it.
type ServerBoundLoadingScreenHandler struct {
	currentID  atomic.Uint32
	expectedID atomic.Uint32
}

// Handle ...
func (h *ServerBoundLoadingScreenHandler) Handle(p packet.Packet, s *Session, _ *world.Tx, c Controllable) error {
	pk := p.(*packet.ServerBoundLoadingScreen)
	v, ok := pk.LoadingScreenID.Value()
	expected := h.expectedID.Load()

	switch {
	case !ok || expected == 0:
		return nil
	case v != expected:
		return fmt.Errorf("expected loading screen ID %d, got %d", expected, v)
	case pk.Type == packet.LoadingScreenTypeEnd:
		s.changingDimension.Store(false)
		h.expectedID.Store(0)
		// The client can reset the local player's abilities to the defaults of its game type when the dimension
		// change ends, dropping what was sent during it (e.g. being able to fly in a survival game type with
		// flying allowed, set right after moving to another world). Send them again now that it's done.
		s.SendGameMode(c)
	}

	return nil
}
