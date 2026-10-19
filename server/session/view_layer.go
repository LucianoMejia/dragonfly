package session

import "github.com/df-mc/dragonfly/server/world"

// ViewLayer returns the session's ViewLayer. The layer may be used to override how entities are viewed
// by this session, such as with a different name tag or visibility state.
func (s *Session) ViewLayer() *world.ViewLayer {
	return s.viewLayer
}

// ViewNameTag overwrites the public name tag of the entity and immediately refreshes it for this session.
func (s *Session) ViewNameTag(entity world.Entity, nameTag string) {
	if s.viewLayer == nil {
		return
	}
	s.viewLayer.ViewNameTag(entity, nameTag)
}

// ViewPublicNameTag removes the name tag override from the entity and immediately refreshes it for this session.
func (s *Session) ViewPublicNameTag(entity world.Entity) {
	if s.viewLayer == nil {
		return
	}
	s.viewLayer.ViewPublicNameTag(entity)
}

// ViewAlwaysShowNameTag overrides whether the entity's name tag is shown at all distances for this session.
func (s *Session) ViewAlwaysShowNameTag(entity world.Entity, alwaysShow bool) {
	if s.viewLayer == nil {
		return
	}
	s.viewLayer.ViewAlwaysShowNameTag(entity, alwaysShow)
}

// ViewPublicAlwaysShowNameTag removes the always-show name tag override from the entity for this session.
func (s *Session) ViewPublicAlwaysShowNameTag(entity world.Entity) {
	if s.viewLayer == nil {
		return
	}
	s.viewLayer.ViewPublicAlwaysShowNameTag(entity)
}

// ViewScoreTag overwrites the public score tag of the entity and immediately refreshes it for this session.
func (s *Session) ViewScoreTag(entity world.Entity, scoreTag string) {
	if s.viewLayer == nil {
		return
	}
	s.viewLayer.ViewScoreTag(entity, scoreTag)
}

// ViewPublicScoreTag removes the score tag override from the entity and immediately refreshes it for this session.
func (s *Session) ViewPublicScoreTag(entity world.Entity) {
	if s.viewLayer == nil {
		return
	}
	s.viewLayer.ViewPublicScoreTag(entity)
}

// ViewVisibility overwrites the public visibility of the entity and immediately refreshes it for this session.
func (s *Session) ViewVisibility(entity world.Entity, level world.VisibilityLevel) {
	if s.viewLayer == nil {
		return
	}
	s.viewLayer.ViewVisibility(entity, level)
}

// RemoveViewLayer removes all overrides for the entity and immediately refreshes it for this session.
func (s *Session) RemoveViewLayer(entity world.Entity) {
	if s.viewLayer == nil {
		return
	}
	s.viewLayer.Remove(entity)
}

// ViewLayerEntityChanged refreshes the entity metadata for this session if the entity is currently visible.
func (s *Session) ViewLayerEntityChanged(e world.Entity) {
	if s.entityHidden(e) || !s.viewingEntity(e.H()) {
		return
	}
	s.ViewEntityState(e)
}

// viewingEntity checks if this session currently has a runtime ID assigned to the entity handle.
func (s *Session) viewingEntity(handle *world.EntityHandle) bool {
	s.entityMutex.RLock()
	_, ok := s.entityRuntimeIDs[handle]
	s.entityMutex.RUnlock()
	return ok
}

// ViewPowered overrides the powered flag of the entity for this session (see world.ViewLayer.ViewPowered).
func (s *Session) ViewPowered(entity world.Entity, powered bool) {
	if s.viewLayer == nil {
		return
	}
	s.viewLayer.ViewPowered(entity, powered)
}

// ViewPublicPowered removes the powered flag override from the entity for this session.
func (s *Session) ViewPublicPowered(entity world.Entity) {
	if s.viewLayer == nil {
		return
	}
	s.viewLayer.ViewPublicPowered(entity)
}

// ViewMarkVariant overrides the mark variant of the entity for this session (see world.ViewLayer.ViewMarkVariant).
func (s *Session) ViewMarkVariant(entity world.Entity, markVariant int32) {
	if s.viewLayer == nil {
		return
	}
	s.viewLayer.ViewMarkVariant(entity, markVariant)
}

// ViewPublicMarkVariant removes the mark variant override from the entity for this session.
func (s *Session) ViewPublicMarkVariant(entity world.Entity) {
	if s.viewLayer == nil {
		return
	}
	s.viewLayer.ViewPublicMarkVariant(entity)
}
