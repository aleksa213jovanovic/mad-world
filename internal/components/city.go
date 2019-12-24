package components

import "fmt"

// City - describes all city operations
type City interface {
	// Return - returns alien to the city after failed occupation
	Return(alien Alien)
	// Occupy - occupies the city
	Occupy(alien Alien)
	// Queue - alien tries to occupy the city
	// If city is being occupied by more aliens, it waits in the queue
	Queue(alien Alien)
	// Destroy - destroys the city
	Destroy()
	// ClearOccupied - clears occupation
	ClearOccupied()
	// ClearIncoming - clears incoming aliens
	ClearIncoming()
	// AliensCollide - are there more than 2 aliens in the queue
	AliensCollide() bool
	// AddNeighbor - adds one neighbor
	AddNeighbor(direction string, city City)
	// GetName - gets city name
	GetName() string
	// GetNeighbors - get all the neighbors
	GetNeighbors() map[string]City
	// GetOccupied - gets alien that occupied the city
	GetOccupied() Alien
	// GetIncoming - gets incoming alliens queue
	GetIncoming() []Alien
	// GetOccupied - gets alien that occupied the city
	IsDestroyed() bool
}

// MadCity - one City implementation
type MadCity struct {
	name     string
	occupied Alien
	incoming []Alien

	neighbors map[string]City
	destroyed bool
}

// Queue - alien tries to occupy the city
// If city is being occupied by more aliens, it waits in the queue
func (c *MadCity) Queue(alien Alien) {
	c.incoming = append(c.incoming, alien)
}

// Return - returns alien to the city after failed occupation
func (c *MadCity) Return(alien Alien) {
	c.incoming = append([]Alien{alien}, c.incoming...)
}

// ClearIncoming - clears incoming aliens
func (c *MadCity) ClearIncoming() {
	c.incoming = nil
}

// ClearOccupied - clears occupation
func (c *MadCity) ClearOccupied() {
	c.occupied = nil
}

// AddNeighbor - adds one neighbor
func (c *MadCity) AddNeighbor(direction string, city City) {
	c.neighbors[direction] = city
}

// Destroy - destroys the city
func (c *MadCity) Destroy() {
	c.destroyed = true
}

// AliensCollide - are there more than 2 aliens in the queue
func (c *MadCity) AliensCollide() bool {
	return len(c.incoming) > 1
}

// GetName - gets city name
func (c *MadCity) GetName() string {
	return c.name
}

// GetOccupied - gets alien that occupied the city
func (c *MadCity) GetOccupied() Alien {
	return c.occupied
}

// IsDestroyed - is city destroyed?
func (c *MadCity) IsDestroyed() bool {
	return c.destroyed
}

// GetIncoming - gets incoming aliens queue
func (c *MadCity) GetIncoming() []Alien {
	return c.incoming
}

// GetNeighbors - get all the neighbors
func (c *MadCity) GetNeighbors() map[string]City {
	return c.neighbors
}

// Occupy - occupies the city
func (c *MadCity) Occupy(alien Alien) {
	c.occupied = alien
}

func (c *MadCity) String() string {
	return fmt.Sprint(c.name)
}

// NewCity constructor
func NewCity(name string) *MadCity {
	return &MadCity{name: name, neighbors: make(map[string]City), destroyed: false, incoming: []Alien{}}
}
