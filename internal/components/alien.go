package components

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Alien - describes all alien operations
type Alien interface {
	// Wait - alien has detected clash, waits in its own city
	Wait()
	// Kill - kills alien
	Kill()
	// Move - tries to move alien
	Move()
	// GoBack - returns alien to origin city
	GoBack()
	// GetName - gets alien name
	GetName() int
	// GetIteration - gets current iteration
	GetIteration() int
	// IsDead - is alien dead
	IsDead() bool
	// IsStuck - is alien stuck
	IsStuck() bool
	// GetOrigin - gets origin city
	GetOrigin() City
	// SetOrigin - sets origin city
	SetOrigin(city City)
}

// MadAlien - one alien implementation
type MadAlien struct {
	name      int
	iteration int
	origin    City
	dead      bool
	stuck     bool
	strategy  DirectionStrategy
}

// SetOrigin - sets origin city
func (a *MadAlien) SetOrigin(city City) {
	a.origin = city
}

// GetOrigin - gets origin city
func (a *MadAlien) GetOrigin() City {
	return a.origin
}

// GetName - gets alien name
func (a *MadAlien) GetName() int {
	return a.name
}

// GetIteration - gets current iteration
func (a *MadAlien) GetIteration() int {
	return a.iteration
}

// IsDead - is alien dead
func (a *MadAlien) IsDead() bool {
	return a.dead
}

// IsStuck - is alien stuck
func (a *MadAlien) IsStuck() bool {
	return a.stuck
}

// Move - tries to move alien
func (a *MadAlien) Move() {
	a.iteration++
	direction := a.strategy.Direction(a)
	if strings.Compare(direction, "wait") == 0 {
		a.Wait()
		log.Debugf("alien %d is waiting\n", a.GetName())
		return
	}
	if strings.Compare(direction, "stuck") == 0 {
		a.stuck = true
		log.Debugf("alien %d is being stuck\n", a.GetName())
		return
	}
	a.GetOrigin().ClearOccupied()
	next := a.GetOrigin().GetNeighbors()[direction]

	if next == nil {
		// this should not happen with regular RandomDirectionStrategy
		panic(fmt.Sprintf("non existing neighbor for city %s and direction %s", a.GetOrigin().GetName(), direction))
	}

	log.Debugf("alien %d tries to move from %s to %s\n", a.name, a.origin.GetName(), next.GetName())
	next.Queue(a)
}

// Wait - alien has detected clash, waits in its own city
func (a *MadAlien) Wait() {
	a.origin.ClearOccupied()
	a.origin.Queue(a)
}

// Kill - kills alien
func (a *MadAlien) Kill() {
	a.dead = true
	a.origin.ClearOccupied()
}

// GoBack - returns alien to origin city
func (a *MadAlien) GoBack() {
	log.Debugf("alien %d reporting! mad aliens already fighting! going back to %s\n", a.name, a.origin.GetName())
	a.GetOrigin().Return(a)
}

func (a *MadAlien) String() string {
	return fmt.Sprint(a.name)
}

// NewAlien - constructor
func NewAlien(name int, origin City, strategy DirectionStrategy) *MadAlien {
	return &MadAlien{name: name, origin: origin, strategy: strategy, dead: false}
}
