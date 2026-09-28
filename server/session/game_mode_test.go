package session

import (
	"testing"

	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TestGameTypeFromMode(t *testing.T) {
	for mode, want := range map[world.GameMode]int32{
		world.GameModeSurvival:  packet.GameTypeSurvival,
		world.GameModeCreative:  packet.GameTypeCreative,
		world.GameModeAdventure: packet.GameTypeSurvival,
		world.GameModeSpectator: packet.GameTypeSpectator,
	} {
		if got := gameTypeFromMode(mode); got != want {
			t.Errorf("gameTypeFromMode(%T) = %v, want %v", mode, got, want)
		}
	}
}
