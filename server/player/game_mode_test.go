package player

import (
	"context"
	"testing"
	"time"

	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// withPlayer runs f with a new player without a session, in the game mode passed.
func withPlayer(t *testing.T, mode world.GameMode, f func(p *Player)) {
	t.Helper()
	w := world.Config{Entities: entity.DefaultRegistry}.New()
	t.Cleanup(func() { _ = w.Close() })
	err := w.Do(func(tx *world.Tx) {
		p := tx.AddEntity(world.EntitySpawnOpts{Position: mgl64.Vec3{0, 64, 0}}.New(Type, Config{Name: "test", GameMode: mode})).(*Player)
		f(p)
	}).Wait(context.Background())
	if err != nil {
		t.Fatalf("world task failed: %v", err)
	}
}

func TestSetGameModeSpectatorFlying(t *testing.T) {
	withPlayer(t, world.GameModeSurvival, func(p *Player) {
		p.SetGameMode(world.GameModeSpectator)
		if !p.Flying() {
			t.Error("player in spectator mode is not flying")
		}
		if !p.Invisible() {
			t.Error("player in spectator mode is not invisible")
		}
		p.StopFlying()
		p.SetGameMode(world.GameModeSpectator)
		if !p.Flying() {
			t.Error("player switching to spectator mode again is not flying")
		}

		p.SetGameMode(world.GameModeCreative)
		if !p.Flying() {
			t.Error("player switching from spectator to creative mode stopped flying")
		}
		if p.Invisible() {
			t.Error("player switching from spectator to creative mode is still invisible")
		}

		p.SetGameMode(world.GameModeSpectator)
		p.SetGameMode(world.GameModeSurvival)
		if p.Flying() {
			t.Error("player in survival mode is still flying")
		}
		if p.Invisible() {
			t.Error("player in survival mode is still invisible")
		}
	})
}

func TestSetGameModeKeepsInvisibilityEffect(t *testing.T) {
	withPlayer(t, world.GameModeSurvival, func(p *Player) {
		p.AddEffect(effect.New(effect.Invisibility, 1, time.Minute))
		p.SetGameMode(world.GameModeSpectator)
		p.SetGameMode(world.GameModeSurvival)
		if !p.Invisible() {
			t.Error("player with the invisibility effect became visible after leaving spectator mode")
		}
	})
}

func TestSpectatorCannotUseItems(t *testing.T) {
	withPlayer(t, world.GameModeSpectator, func(p *Player) {
		p.SetHeldItems(item.NewStack(item.Snowball{}, 16), item.Stack{})
		p.UseItem()
		if held, _ := p.HeldItems(); held.Count() != 16 {
			t.Errorf("player in spectator mode used a snowball: %v left", held.Count())
		}
		for e := range p.Tx().Entities() {
			if e.H().Type() == entity.SnowballType {
				t.Error("player in spectator mode threw a snowball")
			}
		}
	})
}
