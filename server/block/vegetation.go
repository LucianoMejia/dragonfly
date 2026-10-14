package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// This file has decorative plants and cave blocks of newer versions. They
// keep the states of vanilla maps (so they load and look right) and are
// placed like vanilla, but don't grow, tilt or spread. The ones that can be
// waterlogged in vanilla hold a water source, which flows out of them; the
// land plants break when water flows into them.

// placePlant places b on the clicked face. It's the placement shared by the
// plants of this file.
func placePlant(b world.Block, pos cube.Pos, face cube.Face, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, b)
	if !used {
		return false
	}
	place(tx, pos, b, user, ctx)
	return placed(ctx)
}

// placeTall places the lower and upper half of a two-block plant.
func placeTall(lower, upper world.Block, pos cube.Pos, face cube.Face, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, lower)
	if !used || !replaceableWith(tx, pos.Side(cube.FaceUp), upper) {
		return false
	}
	place(tx, pos, lower, user, ctx)
	place(tx, pos.Side(cube.FaceUp), upper, user, ctx)
	return placed(ctx)
}

// plantBreakInfo breaks at once and drops the item.
func plantBreakInfo(drop world.Item) BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, oneOf(drop))
}

// multiFaceBit is the bit of the face in multi_face_direction_bits (down, up,
// south, west, north, east).
func multiFaceBit(f cube.Face) uint8 {
	switch f {
	case cube.FaceDown:
		return 1 << 0
	case cube.FaceUp:
		return 1 << 1
	case cube.FaceSouth:
		return 1 << 2
	case cube.FaceWest:
		return 1 << 3
	case cube.FaceNorth:
		return 1 << 4
	}
	return 1 << 5
}

// --- Dripleaf --------------------------------------------------------------

// SmallDripleaf is a small two-block high plant of the lush caves.
type SmallDripleaf struct {
	empty
	transparent
	sourceWaterDisplacer

	// Facing is the direction the leaves face.
	Facing cube.Direction
	// UpperPart is set for the upper half.
	UpperPart bool
}

// SideClosed ...
func (SmallDripleaf) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool { return false }

// UseOnBlock ...
func (s SmallDripleaf) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	s.Facing = user.Rotation().Direction().Opposite()
	s.UpperPart = false
	return placeTall(s, SmallDripleaf{Facing: s.Facing, UpperPart: true}, pos, face, tx, user, ctx)
}

// NeighbourUpdateTick breaks the plant when its other half is gone.
func (s SmallDripleaf) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	other := pos.Side(cube.FaceUp)
	if s.UpperPart {
		other = pos.Side(cube.FaceDown)
	}
	if o, ok := tx.Block(other).(SmallDripleaf); !ok || o.UpperPart == s.UpperPart {
		breakBlockNoDrops(s, pos, tx)
	}
}

// BreakInfo ...
func (s SmallDripleaf) BreakInfo() BreakInfo { return plantBreakInfo(SmallDripleaf{}) }

// EncodeItem ...
func (SmallDripleaf) EncodeItem() (string, int16) { return "minecraft:small_dripleaf_block", 0 }

// EncodeBlock ...
func (s SmallDripleaf) EncodeBlock() (string, map[string]any) {
	return "minecraft:small_dripleaf_block", map[string]any{"minecraft:cardinal_direction": s.Facing.String(), "upper_block_bit": s.UpperPart}
}

func allSmallDripleaf() (b []world.Block) {
	for _, d := range cube.Directions() {
		b = append(b, SmallDripleaf{Facing: d}, SmallDripleaf{Facing: d, UpperPart: true})
	}
	return
}

// DripleafTilt is how much the leaf of a big dripleaf is tilted.
type DripleafTilt uint8

func (t DripleafTilt) String() string {
	switch t {
	case 1:
		return "unstable"
	case 2:
		return "partial_tilt"
	case 3:
		return "full_tilt"
	}
	return "none"
}

// BigDripleaf is a tall plant of the lush caves with a leaf on top (the
// head) over a stem.
type BigDripleaf struct {
	empty
	transparent
	sourceWaterDisplacer

	// Facing is the direction the leaf faces.
	Facing cube.Direction
	// Head is set for the block with the leaf; the ones under it are stem.
	Head bool
	// Tilt is how tilted the leaf is (kept as it is: it doesn't tilt).
	Tilt DripleafTilt
}

// SideClosed ...
func (BigDripleaf) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool { return false }

// UseOnBlock ...
func (b BigDripleaf) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	return placePlant(BigDripleaf{Facing: user.Rotation().Direction().Opposite(), Head: true}, pos, face, tx, user, ctx)
}

// BreakInfo ...
func (b BigDripleaf) BreakInfo() BreakInfo { return plantBreakInfo(BigDripleaf{Head: true}) }

// EncodeItem ...
func (BigDripleaf) EncodeItem() (string, int16) { return "minecraft:big_dripleaf", 0 }

// EncodeBlock ...
func (b BigDripleaf) EncodeBlock() (string, map[string]any) {
	return "minecraft:big_dripleaf", map[string]any{
		"minecraft:cardinal_direction": b.Facing.String(),
		"big_dripleaf_head":            b.Head,
		"big_dripleaf_tilt":            b.Tilt.String(),
	}
}

func allBigDripleaf() (b []world.Block) {
	for _, d := range cube.Directions() {
		for t := DripleafTilt(0); t <= 3; t++ {
			b = append(b, BigDripleaf{Facing: d, Tilt: t}, BigDripleaf{Facing: d, Tilt: t, Head: true})
		}
	}
	return
}

// --- Cave blocks -------------------------------------------------------------

// GlowLichen is a glowing lichen that grows on the faces of blocks.
type GlowLichen struct {
	empty
	transparent
	sourceWaterDisplacer

	// Faces are the faces it's on (multi_face_direction_bits).
	Faces uint8
}

// SideClosed ...
func (GlowLichen) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool { return false }

// LightEmissionLevel ...
func (GlowLichen) LightEmissionLevel() uint8 { return 7 }

// UseOnBlock ...
func (g GlowLichen) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	return placePlant(GlowLichen{Faces: multiFaceBit(face.Opposite())}, pos, face, tx, user, ctx)
}

// BreakInfo ...
func (g GlowLichen) BreakInfo() BreakInfo { return plantBreakInfo(GlowLichen{}) }

// EncodeItem ...
func (GlowLichen) EncodeItem() (string, int16) { return "minecraft:glow_lichen", 0 }

// EncodeBlock ...
func (g GlowLichen) EncodeBlock() (string, map[string]any) {
	return "minecraft:glow_lichen", map[string]any{"multi_face_direction_bits": int32(g.Faces)}
}

// SculkVein is a thin sculk that grows on the faces of blocks.
type SculkVein struct {
	empty
	transparent
	sourceWaterDisplacer

	// Faces are the faces it's on (multi_face_direction_bits).
	Faces uint8
}

// SideClosed ...
func (SculkVein) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool { return false }

// UseOnBlock ...
func (s SculkVein) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	return placePlant(SculkVein{Faces: multiFaceBit(face.Opposite())}, pos, face, tx, user, ctx)
}

// BreakInfo ...
func (s SculkVein) BreakInfo() BreakInfo { return plantBreakInfo(SculkVein{}) }

// EncodeItem ...
func (SculkVein) EncodeItem() (string, int16) { return "minecraft:sculk_vein", 0 }

// EncodeBlock ...
func (s SculkVein) EncodeBlock() (string, map[string]any) {
	return "minecraft:sculk_vein", map[string]any{"multi_face_direction_bits": int32(s.Faces)}
}

func allMultiFace() (b []world.Block) {
	for f := uint8(0); f < 64; f++ {
		b = append(b, GlowLichen{Faces: f}, SculkVein{Faces: f})
	}
	return
}

// HangingRoots hang from the ceiling of lush caves.
type HangingRoots struct {
	empty
	transparent
	sourceWaterDisplacer
}

// SideClosed ...
func (HangingRoots) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool { return false }

// UseOnBlock ...
func (h HangingRoots) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	return placePlant(h, pos, face, tx, user, ctx)
}

// BreakInfo ...
func (h HangingRoots) BreakInfo() BreakInfo { return plantBreakInfo(h) }

// EncodeItem ...
func (HangingRoots) EncodeItem() (string, int16) { return "minecraft:hanging_roots", 0 }

// EncodeBlock ...
func (HangingRoots) EncodeBlock() (string, map[string]any) { return "minecraft:hanging_roots", nil }

// DripstoneThickness is the part of a pointed dripstone column.
type DripstoneThickness uint8

func (t DripstoneThickness) String() string {
	switch t {
	case 1:
		return "merge"
	case 2:
		return "frustum"
	case 3:
		return "middle"
	case 4:
		return "base"
	}
	return "tip"
}

// PointedDripstone is a stalactite (hanging) or stalagmite of dripstone.
type PointedDripstone struct {
	empty
	transparent
	sourceWaterDisplacer

	// Thickness is the part of the column.
	Thickness DripstoneThickness
	// Hanging is set for stalactites.
	Hanging bool
}

// SideClosed ...
func (PointedDripstone) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool { return false }

// UseOnBlock places a stalactite under a block, or a stalagmite on top.
func (p PointedDripstone) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	return placePlant(PointedDripstone{Hanging: face == cube.FaceDown}, pos, face, tx, user, ctx)
}

// BreakInfo ...
func (p PointedDripstone) BreakInfo() BreakInfo {
	return newBreakInfo(1.5, alwaysHarvestable, pickaxeEffective, oneOf(PointedDripstone{}))
}

// EncodeItem ...
func (PointedDripstone) EncodeItem() (string, int16) { return "minecraft:pointed_dripstone", 0 }

// EncodeBlock ...
func (p PointedDripstone) EncodeBlock() (string, map[string]any) {
	return "minecraft:pointed_dripstone", map[string]any{"dripstone_thickness": p.Thickness.String(), "hanging": p.Hanging}
}

func allPointedDripstone() (b []world.Block) {
	for t := DripstoneThickness(0); t <= 4; t++ {
		b = append(b, PointedDripstone{Thickness: t}, PointedDripstone{Thickness: t, Hanging: true})
	}
	return
}

// MangrovePropagule is the sapling of the mangrove tree.
type MangrovePropagule struct {
	empty
	transparent
	sourceWaterDisplacer

	// Hanging is set when it hangs from mangrove leaves.
	Hanging bool
	// Stage is its growth stage, 0-4.
	Stage int
}

// SideClosed ...
func (MangrovePropagule) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool { return false }

// UseOnBlock ...
func (m MangrovePropagule) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	return placePlant(MangrovePropagule{Hanging: face == cube.FaceDown}, pos, face, tx, user, ctx)
}

// BreakInfo ...
func (m MangrovePropagule) BreakInfo() BreakInfo { return plantBreakInfo(MangrovePropagule{}) }

// EncodeItem ...
func (MangrovePropagule) EncodeItem() (string, int16) { return "minecraft:mangrove_propagule", 0 }

// EncodeBlock ...
func (m MangrovePropagule) EncodeBlock() (string, map[string]any) {
	return "minecraft:mangrove_propagule", map[string]any{"hanging": m.Hanging, "propagule_stage": int32(m.Stage)}
}

func allMangrovePropagule() (b []world.Block) {
	for s := 0; s <= 4; s++ {
		b = append(b, MangrovePropagule{Stage: s}, MangrovePropagule{Stage: s, Hanging: true})
	}
	return
}

// --- Vines -------------------------------------------------------------------

// CaveVinesBerries is where a cave vine has glow berries.
type CaveVinesBerries uint8

// CaveVines hang from the ceiling of lush caves, with or without glow berries.
type CaveVines struct {
	empty
	transparent

	// Age is the growth age, 0-25.
	Age int
	// Berries: 0 none, 1 on a body block, 2 on the head (the lowest block).
	Berries CaveVinesBerries
}

// LightEmissionLevel ...
func (c CaveVines) LightEmissionLevel() uint8 {
	if c.Berries != 0 {
		return 14
	}
	return 0
}

// UseOnBlock ...
func (c CaveVines) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	return placePlant(CaveVines{}, pos, face, tx, user, ctx)
}

// HasLiquidDrops ...
func (CaveVines) HasLiquidDrops() bool { return true }

// BreakInfo ...
func (c CaveVines) BreakInfo() BreakInfo { return plantBreakInfo(CaveVines{}) }

// EncodeItem ...
func (CaveVines) EncodeItem() (string, int16) { return "minecraft:cave_vines", 0 }

// EncodeBlock ...
func (c CaveVines) EncodeBlock() (string, map[string]any) {
	name := "minecraft:cave_vines"
	switch c.Berries {
	case 1:
		name = "minecraft:cave_vines_body_with_berries"
	case 2:
		name = "minecraft:cave_vines_head_with_berries"
	}
	return name, map[string]any{"growing_plant_age": int32(c.Age)}
}

// NetherVines are the twisting vines (growing up) or weeping vines (hanging).
type NetherVines struct {
	empty
	transparent

	// Weeping is set for weeping vines; twisting vines otherwise.
	Weeping bool
	// Age is the growth age, 0-25.
	Age int
}

// UseOnBlock ...
func (n NetherVines) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	return placePlant(NetherVines{Weeping: n.Weeping}, pos, face, tx, user, ctx)
}

// HasLiquidDrops ...
func (NetherVines) HasLiquidDrops() bool { return true }

// BreakInfo ...
func (n NetherVines) BreakInfo() BreakInfo { return plantBreakInfo(NetherVines{Weeping: n.Weeping}) }

// EncodeItem ...
func (n NetherVines) EncodeItem() (string, int16) {
	if n.Weeping {
		return "minecraft:weeping_vines", 0
	}
	return "minecraft:twisting_vines", 0
}

// EncodeBlock ...
func (n NetherVines) EncodeBlock() (string, map[string]any) {
	if n.Weeping {
		return "minecraft:weeping_vines", map[string]any{"weeping_vines_age": int32(n.Age)}
	}
	return "minecraft:twisting_vines", map[string]any{"twisting_vines_age": int32(n.Age)}
}

func allCaveAndNetherVines() (b []world.Block) {
	for age := 0; age <= 25; age++ {
		for berries := CaveVinesBerries(0); berries <= 2; berries++ {
			b = append(b, CaveVines{Age: age, Berries: berries})
		}
		b = append(b, NetherVines{Age: age}, NetherVines{Age: age, Weeping: true})
	}
	return
}

// --- Moss --------------------------------------------------------------------

// MossBlock is a block of moss (pale moss in the pale garden).
type MossBlock struct {
	solid

	// Pale is set for pale moss.
	Pale bool
}

// BreakInfo ...
func (m MossBlock) BreakInfo() BreakInfo {
	return newBreakInfo(0.1, alwaysHarvestable, nothingEffective, oneOf(m))
}

// SoilFor lets plants grow on moss, like on dirt.
func (MossBlock) SoilFor(block world.Block) bool {
	switch block.(type) {
	case ShortGrass, Fern, DoubleTallGrass, Flower, DoubleFlower, PinkPetals, Bush, DryGrass, Wildflowers, LeafLitter, Azalea, Eyeblossom, Torchflower:
		return true
	}
	return false
}

// EncodeItem ...
func (m MossBlock) EncodeItem() (string, int16) {
	if m.Pale {
		return "minecraft:pale_moss_block", 0
	}
	return "minecraft:moss_block", 0
}

// EncodeBlock ...
func (m MossBlock) EncodeBlock() (string, map[string]any) {
	name, _ := m.EncodeItem()
	return name, nil
}

// PaleMossSide is how far pale moss carpet climbs a side: 0 none, 1 short,
// 2 tall.
type PaleMossSide uint8

func (s PaleMossSide) String() string {
	switch s {
	case 1:
		return "short"
	case 2:
		return "tall"
	}
	return "none"
}

// PaleMossCarpet is a carpet of pale moss that can climb the sides.
type PaleMossCarpet struct {
	carpet
	transparent

	// UpperPart is set for the block of carpet that only covers sides above.
	UpperPart bool
	// North, East, South and West are how far it climbs each side.
	North, East, South, West PaleMossSide
}

// UseOnBlock ...
func (p PaleMossCarpet) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	return placePlant(PaleMossCarpet{}, pos, face, tx, user, ctx)
}

// HasLiquidDrops ...
func (PaleMossCarpet) HasLiquidDrops() bool { return true }

// BreakInfo ...
func (p PaleMossCarpet) BreakInfo() BreakInfo { return plantBreakInfo(PaleMossCarpet{}) }

// EncodeItem ...
func (PaleMossCarpet) EncodeItem() (string, int16) { return "minecraft:pale_moss_carpet", 0 }

// EncodeBlock ...
func (p PaleMossCarpet) EncodeBlock() (string, map[string]any) {
	return "minecraft:pale_moss_carpet", map[string]any{
		"upper_block_bit":             p.UpperPart,
		"pale_moss_carpet_side_north": p.North.String(),
		"pale_moss_carpet_side_east":  p.East.String(),
		"pale_moss_carpet_side_south": p.South.String(),
		"pale_moss_carpet_side_west":  p.West.String(),
	}
}

func allPaleMossCarpet() (b []world.Block) {
	for _, upper := range []bool{false, true} {
		for n := PaleMossSide(0); n <= 2; n++ {
			for e := PaleMossSide(0); e <= 2; e++ {
				for s := PaleMossSide(0); s <= 2; s++ {
					for w := PaleMossSide(0); w <= 2; w++ {
						b = append(b, PaleMossCarpet{UpperPart: upper, North: n, East: e, South: s, West: w})
					}
				}
			}
		}
	}
	return
}

// PaleHangingMoss hangs from the pale oak leaves.
type PaleHangingMoss struct {
	empty
	transparent

	// Tip is set for the lowest block.
	Tip bool
}

// UseOnBlock ...
func (p PaleHangingMoss) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	return placePlant(PaleHangingMoss{Tip: true}, pos, face, tx, user, ctx)
}

// HasLiquidDrops ...
func (PaleHangingMoss) HasLiquidDrops() bool { return true }

// BreakInfo ...
func (p PaleHangingMoss) BreakInfo() BreakInfo { return plantBreakInfo(PaleHangingMoss{Tip: true}) }

// EncodeItem ...
func (PaleHangingMoss) EncodeItem() (string, int16) { return "minecraft:pale_hanging_moss", 0 }

// EncodeBlock ...
func (p PaleHangingMoss) EncodeBlock() (string, map[string]any) {
	return "minecraft:pale_hanging_moss", map[string]any{"tip": p.Tip}
}

// --- Land plants -------------------------------------------------------------

// Bush is a small green bush, or a firefly bush (which glows a little).
type Bush struct {
	empty
	transparent
	replaceable

	// Firefly is set for the firefly bush.
	Firefly bool
}

// LightEmissionLevel ...
func (b Bush) LightEmissionLevel() uint8 {
	if b.Firefly {
		return 2
	}
	return 0
}

// UseOnBlock ...
func (b Bush) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	return placePlant(b, pos, face, tx, user, ctx)
}

// HasLiquidDrops ...
func (Bush) HasLiquidDrops() bool { return true }

// BreakInfo ...
func (b Bush) BreakInfo() BreakInfo { return plantBreakInfo(b) }

// EncodeItem ...
func (b Bush) EncodeItem() (string, int16) {
	if b.Firefly {
		return "minecraft:firefly_bush", 0
	}
	return "minecraft:bush", 0
}

// EncodeBlock ...
func (b Bush) EncodeBlock() (string, map[string]any) {
	name, _ := b.EncodeItem()
	return name, nil
}

// DryGrass is the short or tall dry grass of the badlands and deserts.
type DryGrass struct {
	empty
	transparent
	replaceable

	// Tall is set for tall dry grass.
	Tall bool
}

// UseOnBlock ...
func (d DryGrass) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	return placePlant(d, pos, face, tx, user, ctx)
}

// HasLiquidDrops ...
func (DryGrass) HasLiquidDrops() bool { return true }

// BreakInfo ...
func (d DryGrass) BreakInfo() BreakInfo { return plantBreakInfo(d) }

// EncodeItem ...
func (d DryGrass) EncodeItem() (string, int16) {
	if d.Tall {
		return "minecraft:tall_dry_grass", 0
	}
	return "minecraft:short_dry_grass", 0
}

// EncodeBlock ...
func (d DryGrass) EncodeBlock() (string, map[string]any) {
	name, _ := d.EncodeItem()
	return name, nil
}

// groundCoverStates are the states shared by wildflowers and leaf litter:
// how many there are and the way they face.
func groundCoverStates(f func(growth int, d cube.Direction) world.Block) (b []world.Block) {
	for g := 0; g <= 7; g++ {
		for _, d := range cube.Directions() {
			b = append(b, f(g, d))
		}
	}
	return
}

// Wildflowers are small flowers that cover the ground, like pink petals.
type Wildflowers struct {
	empty
	transparent

	// Growth is how many flowers there are, 0-7 (0-3 in game).
	Growth int
	// Facing is the direction they face.
	Facing cube.Direction
}

// UseOnBlock adds a flower to the ones already there, or places them.
func (w Wildflowers) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	if existing, ok := tx.Block(pos).(Wildflowers); ok && existing.Growth < 3 {
		existing.Growth++
		place(tx, pos, existing, user, ctx)
		return placed(ctx)
	}
	return placePlant(Wildflowers{Facing: user.Rotation().Direction().Opposite()}, pos, face, tx, user, ctx)
}

// HasLiquidDrops ...
func (Wildflowers) HasLiquidDrops() bool { return true }

// BreakInfo ...
func (w Wildflowers) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, simpleDrops(item.NewStack(Wildflowers{}, w.Growth+1)))
}

// EncodeItem ...
func (Wildflowers) EncodeItem() (string, int16) { return "minecraft:wildflowers", 0 }

// EncodeBlock ...
func (w Wildflowers) EncodeBlock() (string, map[string]any) {
	return "minecraft:wildflowers", map[string]any{"growth": int32(w.Growth), "minecraft:cardinal_direction": w.Facing.String()}
}

// LeafLitter is a layer of fallen leaves on the ground.
type LeafLitter struct {
	empty
	transparent

	// Growth is how many leaves there are, 0-7 (0-3 in game).
	Growth int
	// Facing is the direction they face.
	Facing cube.Direction
}

// UseOnBlock adds leaves to the ones already there, or places them.
func (l LeafLitter) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	if existing, ok := tx.Block(pos).(LeafLitter); ok && existing.Growth < 3 {
		existing.Growth++
		place(tx, pos, existing, user, ctx)
		return placed(ctx)
	}
	return placePlant(LeafLitter{Facing: user.Rotation().Direction().Opposite()}, pos, face, tx, user, ctx)
}

// HasLiquidDrops ...
func (LeafLitter) HasLiquidDrops() bool { return true }

// BreakInfo ...
func (l LeafLitter) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, simpleDrops(item.NewStack(LeafLitter{}, l.Growth+1)))
}

// EncodeItem ...
func (LeafLitter) EncodeItem() (string, int16) { return "minecraft:leaf_litter", 0 }

// EncodeBlock ...
func (l LeafLitter) EncodeBlock() (string, map[string]any) {
	return "minecraft:leaf_litter", map[string]any{"growth": int32(l.Growth), "minecraft:cardinal_direction": l.Facing.String()}
}

// Azalea is a small bush that grows into an azalea tree, flowering or not.
type Azalea struct {
	transparent
	sourceWaterDisplacer

	// Flowering is set for the flowering azalea.
	Flowering bool
}

// Model ...
func (Azalea) Model() world.BlockModel { return solid{}.Model() }

// SideClosed ...
func (Azalea) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool { return false }

// UseOnBlock ...
func (a Azalea) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	return placePlant(a, pos, face, tx, user, ctx)
}

// BreakInfo ...
func (a Azalea) BreakInfo() BreakInfo { return plantBreakInfo(a) }

// EncodeItem ...
func (a Azalea) EncodeItem() (string, int16) {
	if a.Flowering {
		return "minecraft:flowering_azalea", 0
	}
	return "minecraft:azalea", 0
}

// EncodeBlock ...
func (a Azalea) EncodeBlock() (string, map[string]any) {
	name, _ := a.EncodeItem()
	return name, nil
}

// Eyeblossom is the flower of the pale garden, open (at night) or closed.
type Eyeblossom struct {
	empty
	transparent

	// Open is set for the open eyeblossom.
	Open bool
}

// UseOnBlock ...
func (e Eyeblossom) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	return placePlant(e, pos, face, tx, user, ctx)
}

// HasLiquidDrops ...
func (Eyeblossom) HasLiquidDrops() bool { return true }

// BreakInfo ...
func (e Eyeblossom) BreakInfo() BreakInfo { return plantBreakInfo(e) }

// EncodeItem ...
func (e Eyeblossom) EncodeItem() (string, int16) {
	if e.Open {
		return "minecraft:open_eyeblossom", 0
	}
	return "minecraft:closed_eyeblossom", 0
}

// EncodeBlock ...
func (e Eyeblossom) EncodeBlock() (string, map[string]any) {
	name, _ := e.EncodeItem()
	return name, nil
}

// Torchflower is the flower grown from the seeds the sniffer digs up.
type Torchflower struct {
	empty
	transparent
}

// UseOnBlock ...
func (t Torchflower) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	return placePlant(t, pos, face, tx, user, ctx)
}

// HasLiquidDrops ...
func (Torchflower) HasLiquidDrops() bool { return true }

// BreakInfo ...
func (t Torchflower) BreakInfo() BreakInfo { return plantBreakInfo(t) }

// EncodeItem ...
func (Torchflower) EncodeItem() (string, int16) { return "minecraft:torchflower", 0 }

// EncodeBlock ...
func (Torchflower) EncodeBlock() (string, map[string]any) { return "minecraft:torchflower", nil }

// CactusFlower is the flower that grows on top of cactus.
type CactusFlower struct {
	empty
	transparent
}

// UseOnBlock ...
func (c CactusFlower) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	return placePlant(c, pos, face, tx, user, ctx)
}

// HasLiquidDrops ...
func (CactusFlower) HasLiquidDrops() bool { return true }

// BreakInfo ...
func (c CactusFlower) BreakInfo() BreakInfo { return plantBreakInfo(c) }

// EncodeItem ...
func (CactusFlower) EncodeItem() (string, int16) { return "minecraft:cactus_flower", 0 }

// EncodeBlock ...
func (CactusFlower) EncodeBlock() (string, map[string]any) { return "minecraft:cactus_flower", nil }

// PitcherPlant is a two-block high flower grown from the pitcher pod.
type PitcherPlant struct {
	empty
	transparent

	// UpperPart is set for the upper half.
	UpperPart bool
}

// UseOnBlock ...
func (p PitcherPlant) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	return placeTall(PitcherPlant{}, PitcherPlant{UpperPart: true}, pos, face, tx, user, ctx)
}

// NeighbourUpdateTick breaks the plant when its other half is gone.
func (p PitcherPlant) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	other := pos.Side(cube.FaceUp)
	if p.UpperPart {
		other = pos.Side(cube.FaceDown)
	}
	if o, ok := tx.Block(other).(PitcherPlant); !ok || o.UpperPart == p.UpperPart {
		breakBlockNoDrops(p, pos, tx)
	}
}

// HasLiquidDrops ...
func (PitcherPlant) HasLiquidDrops() bool { return true }

// BreakInfo ...
func (p PitcherPlant) BreakInfo() BreakInfo { return plantBreakInfo(PitcherPlant{}) }

// EncodeItem ...
func (PitcherPlant) EncodeItem() (string, int16) { return "minecraft:pitcher_plant", 0 }

// EncodeBlock ...
func (p PitcherPlant) EncodeBlock() (string, map[string]any) {
	return "minecraft:pitcher_plant", map[string]any{"upper_block_bit": p.UpperPart}
}

// SweetBerryBush is the bush of the sweet berries (it doesn't grow or hurt).
type SweetBerryBush struct {
	empty
	transparent

	// Growth is its growth stage, 0-7.
	Growth int
}

// HasLiquidDrops ...
func (SweetBerryBush) HasLiquidDrops() bool { return true }

// BreakInfo ...
func (s SweetBerryBush) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, simpleDrops())
}

// EncodeBlock ...
func (s SweetBerryBush) EncodeBlock() (string, map[string]any) {
	return "minecraft:sweet_berry_bush", map[string]any{"growth": int32(s.Growth)}
}

// FrogSpawn is the eggs frogs lay on the water.
type FrogSpawn struct {
	empty
	transparent
}

// UseOnBlock ...
func (f FrogSpawn) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	return placePlant(f, pos, face, tx, user, ctx)
}

// BreakInfo ...
func (f FrogSpawn) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, simpleDrops())
}

// EncodeItem ...
func (FrogSpawn) EncodeItem() (string, int16) { return "minecraft:frog_spawn", 0 }

// EncodeBlock ...
func (FrogSpawn) EncodeBlock() (string, map[string]any) { return "minecraft:frog_spawn", nil }

// RootedDirt is dirt with roots, found over azalea trees.
type RootedDirt struct {
	solid
}

// SoilFor ...
func (RootedDirt) SoilFor(block world.Block) bool {
	switch block.(type) {
	case ShortGrass, Fern, DoubleTallGrass, Flower, DoubleFlower, Bush, DryGrass, Azalea, Sapling:
		return true
	}
	return false
}

// BreakInfo ...
func (r RootedDirt) BreakInfo() BreakInfo {
	return newBreakInfo(0.5, alwaysHarvestable, shovelEffective, oneOf(r))
}

// EncodeItem ...
func (RootedDirt) EncodeItem() (string, int16) { return "minecraft:dirt_with_roots", 0 }

// EncodeBlock ...
func (RootedDirt) EncodeBlock() (string, map[string]any) { return "minecraft:dirt_with_roots", nil }

// --- Coral fans --------------------------------------------------------------

// CoralFan is a fan of coral on the floor, or on a wall (WallFacing set).
// Like the coral, it lives in water.
type CoralFan struct {
	empty
	transparent
	sourceWaterDisplacer

	// Type is the type of coral.
	Type CoralType
	// Dead is set for dead coral.
	Dead bool
	// Wall is set for a fan on the side of a block, facing WallFacing.
	Wall bool
	// WallFacing is the direction of a wall fan (coral_direction, 0-3).
	WallFacing int
	// Rotated turns a floor fan 90 degrees (coral_fan_direction).
	Rotated bool
}

// SideClosed ...
func (CoralFan) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool { return false }

// UseOnBlock places a floor fan, or a wall fan on the side of a block.
func (c CoralFan) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	fan := CoralFan{Type: c.Type, Dead: c.Dead}
	switch face {
	case cube.FaceUp, cube.FaceDown:
		fan.Rotated = user.Rotation().Direction().Face().Axis() == cube.X
	default:
		fan.Wall = true
		fan.WallFacing = coralWallDirection(face)
	}
	return placePlant(fan, pos, face, tx, user, ctx)
}

// coralWallDirection is the coral_direction of a fan on the face of a block.
func coralWallDirection(f cube.Face) int {
	switch f {
	case cube.FaceEast:
		return 1
	case cube.FaceNorth:
		return 2
	case cube.FaceSouth:
		return 3
	}
	return 0 // west
}

// BreakInfo ...
func (c CoralFan) BreakInfo() BreakInfo {
	return plantBreakInfo(CoralFan{Type: c.Type, Dead: c.Dead})
}

// EncodeItem ...
func (c CoralFan) EncodeItem() (string, int16) {
	if c.Dead {
		return "minecraft:dead_" + c.Type.String() + "_coral_fan", 0
	}
	return "minecraft:" + c.Type.String() + "_coral_fan", 0
}

// EncodeBlock ...
func (c CoralFan) EncodeBlock() (string, map[string]any) {
	prefix := "minecraft:"
	if c.Dead {
		prefix += "dead_"
	}
	if c.Wall {
		return prefix + c.Type.String() + "_coral_wall_fan", map[string]any{"coral_direction": int32(c.WallFacing)}
	}
	dir := int32(0)
	if c.Rotated {
		dir = 1
	}
	return prefix + c.Type.String() + "_coral_fan", map[string]any{"coral_fan_direction": dir}
}

func allCoralFans() (b []world.Block) {
	for _, t := range CoralTypes() {
		for _, dead := range []bool{false, true} {
			b = append(b, CoralFan{Type: t, Dead: dead}, CoralFan{Type: t, Dead: dead, Rotated: true})
			for d := 0; d <= 3; d++ {
				b = append(b, CoralFan{Type: t, Dead: dead, Wall: true, WallFacing: d})
			}
		}
	}
	return
}

// allVegetation returns every state of the blocks of this file.
func allVegetation() (b []world.Block) {
	b = append(b, allSmallDripleaf()...)
	b = append(b, allBigDripleaf()...)
	b = append(b, allMultiFace()...)
	b = append(b, allPointedDripstone()...)
	b = append(b, allMangrovePropagule()...)
	b = append(b, allCaveAndNetherVines()...)
	b = append(b, allPaleMossCarpet()...)
	b = append(b, allCoralFans()...)
	b = append(b, groundCoverStates(func(g int, d cube.Direction) world.Block { return Wildflowers{Growth: g, Facing: d} })...)
	b = append(b, groundCoverStates(func(g int, d cube.Direction) world.Block { return LeafLitter{Growth: g, Facing: d} })...)
	for g := 0; g <= 7; g++ {
		b = append(b, SweetBerryBush{Growth: g})
	}
	b = append(b,
		HangingRoots{}, MossBlock{}, MossBlock{Pale: true}, PaleHangingMoss{}, PaleHangingMoss{Tip: true},
		Bush{}, Bush{Firefly: true}, DryGrass{}, DryGrass{Tall: true}, Azalea{}, Azalea{Flowering: true},
		Eyeblossom{}, Eyeblossom{Open: true}, Torchflower{}, CactusFlower{}, PitcherPlant{}, PitcherPlant{UpperPart: true},
		FrogSpawn{}, RootedDirt{},
	)
	return
}

// vegetationItems are the items of the blocks of this file.
func vegetationItems() []world.Item {
	items := []world.Item{
		SmallDripleaf{}, BigDripleaf{Head: true}, GlowLichen{}, SculkVein{}, HangingRoots{}, PointedDripstone{},
		MangrovePropagule{}, CaveVines{}, NetherVines{}, NetherVines{Weeping: true}, MossBlock{}, MossBlock{Pale: true},
		PaleMossCarpet{}, PaleHangingMoss{Tip: true}, Bush{}, Bush{Firefly: true}, DryGrass{}, DryGrass{Tall: true},
		Wildflowers{}, LeafLitter{}, Azalea{}, Azalea{Flowering: true}, Eyeblossom{}, Eyeblossom{Open: true},
		Torchflower{}, CactusFlower{}, PitcherPlant{}, FrogSpawn{}, RootedDirt{},
	}
	for _, t := range CoralTypes() {
		items = append(items, CoralFan{Type: t}, CoralFan{Type: t, Dead: true})
	}
	return items
}
