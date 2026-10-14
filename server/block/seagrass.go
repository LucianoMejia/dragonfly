package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Seagrass is a non-solid plant that grows underwater. It's always in a
// water source: the water flows out of it like from any source, and stays
// when the seagrass is broken. Tall seagrass is two blocks high.
type Seagrass struct {
	empty
	transparent
	sourceWaterDisplacer

	// Tall is set for the two-block high seagrass.
	Tall bool
	// UpperPart is set for the upper half of tall seagrass.
	UpperPart bool
}

// SideClosed always returns false: water flows through seagrass.
func (Seagrass) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

// keepWater makes sure there's a water source in the seagrass, so that it
// flows out (maps from vanilla may have seagrass without the water layer).
func (s Seagrass) keepWater(pos cube.Pos, tx *world.Tx) {
	if liquid, ok := tx.Liquid(pos); ok {
		if _, water := liquid.(Water); water && liquid.LiquidDepth() == 8 {
			return
		}
	}
	tx.SetLiquid(pos, Water{Still: true, Depth: 8})
}

// NeighbourUpdateTick ...
func (s Seagrass) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	s.keepWater(pos, tx)
	if !s.Tall {
		return
	}
	if s.UpperPart {
		if below, ok := tx.Block(pos.Side(cube.FaceDown)).(Seagrass); !ok || !below.Tall || below.UpperPart {
			breakBlockNoDrops(s, pos, tx)
		}
	} else if above, ok := tx.Block(pos.Side(cube.FaceUp)).(Seagrass); !ok || !above.Tall || !above.UpperPart {
		breakBlockNoDrops(s, pos, tx)
	}
}

// UseOnBlock places seagrass in a water source, over a solid block.
func (s Seagrass) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, s)
	if !used {
		return false
	}
	below := pos.Side(cube.FaceDown)
	if !tx.Block(below).Model().FaceSolid(below, cube.FaceUp, tx) {
		return false
	}
	liquid, ok := tx.Liquid(pos)
	if !ok {
		return false
	}
	if _, water := liquid.(Water); !water || liquid.LiquidDepth() < 8 {
		return false
	}
	place(tx, pos, Seagrass{}, user, ctx)
	return placed(ctx)
}

// BreakInfo ...
func (s Seagrass) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, shearsEffective, func(t item.Tool, _ []item.Enchantment) []item.Stack {
		if t.ToolType() == item.TypeShears {
			return []item.Stack{item.NewStack(Seagrass{}, 1)}
		}
		return nil
	})
}

// CompostChance ...
func (Seagrass) CompostChance() float64 {
	return 0.3
}

// EncodeItem ...
func (Seagrass) EncodeItem() (name string, meta int16) {
	return "minecraft:seagrass", 0
}

// EncodeBlock ...
func (s Seagrass) EncodeBlock() (string, map[string]any) {
	typ := "default"
	if s.Tall {
		typ = "double_bot"
		if s.UpperPart {
			typ = "double_top"
		}
	}
	return "minecraft:seagrass", map[string]any{"sea_grass_type": typ}
}

// allSeagrass returns every state of seagrass.
func allSeagrass() []world.Block {
	return []world.Block{Seagrass{}, Seagrass{Tall: true}, Seagrass{Tall: true, UpperPart: true}}
}
