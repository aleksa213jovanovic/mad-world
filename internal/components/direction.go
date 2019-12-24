package components

import (
	"math/rand"
	"time"
)

// DirectionStrategy - describes alien behavior
type DirectionStrategy interface {
	// Direction - gets next Direction
	Direction(alien Alien) string
}

// BufferedDirectionStrategy - both cities and aliens are loaded from files
type BufferedDirectionStrategy struct {
	directions []string
}

// Direction - gets next Direction
func (bds *BufferedDirectionStrategy) Direction(alien Alien) string {
	if len(bds.directions) == 0 {
		return "stuck"
	}
	direction := bds.directions[0]
	bds.directions = bds.directions[1:]
	return direction
}

// NewBufferedDirectionStrategy - constructor
func NewBufferedDirectionStrategy(directions []string) *BufferedDirectionStrategy {
	return &BufferedDirectionStrategy{directions: directions}
}

// RandomDirectionStrategy - cities are loaded from file, aliens pick random directions
type RandomDirectionStrategy struct {
	random     *rand.Rand
}

// Direction - gets next Direction
func (rds *RandomDirectionStrategy) Direction(alien Alien) string {
	var directions []string
	for direction, neighbor := range alien.GetOrigin().GetNeighbors() {
		if !neighbor.IsDestroyed() {
			directions = append(directions, direction)
		}
	}
	if len(directions) == 0 {
		return "stuck"
	}
	return directions[rds.random.Intn(len(directions))]
}

// NewRandomDirectionStrategy - constructor
func NewRandomDirectionStrategy() *RandomDirectionStrategy {
	random := rand.New(rand.NewSource(time.Now().Unix()))
	return &RandomDirectionStrategy{random: random}
}
