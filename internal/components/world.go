package components

import (
	log "github.com/sirupsen/logrus"
	"sort"
)

// World - struct that holds all cities and aliens
type World struct {
	iterations int
	aliens     []Alien
	cities     []City
	alienMap   map[int]Alien
	cityMap    map[string]City
}

// RunSimulation - runs complete simulation for the world we created
func (world *World) RunSimulation() {
	for {
		if len(world.aliens) == 0 {
			log.Info("no aliens provided")
			return
		}
		if world.allDeadOrStuck() {
			log.Info("all aliens are dead or stuck")
			return
		}
		if world.allReachedIterations() {
			log.Info("all aliens have reached max iterations")
			return
		}
		for _, alien := range world.aliens {
			if alien.IsDead() || alien.GetIteration() == world.iterations {
				continue
			}
			alien.Move()
		}

		// loop while there are collisions
		for world.atLeastOneCollision() {
			// sort cities by their incoming aliens number
			sort.Slice(world.cities, func(i, j int) bool {
				return len(world.cities[i].GetIncoming()) > len(world.cities[j].GetIncoming())
			})

			for _, city := range world.cities {
				if city.IsDestroyed() {
					continue
				}
				if !city.AliensCollide() {
					continue
				}

				firstAlien := city.GetIncoming()[0]
				secondAlien := city.GetIncoming()[1]
				firstAlien.Kill()
				secondAlien.Kill()
				city.Destroy()

				log.Infof("%s has been destroyed by alien %d and alien %d!\n", city.GetName(), firstAlien.GetName(), secondAlien.GetName())

				// return the of the aliens to their origin city: they avoided the war for now
				if len(city.GetIncoming()) > 2 {
					for _, alien := range city.GetIncoming()[2:] {
						alien.GoBack()
					}
				}
			}
		}

		// the rest of the aliens have successfully occupied the city
		for _, city := range world.cities {
			if city.IsDestroyed() {
				continue
			}
			if len(city.GetIncoming()) > 0 {
				city.Occupy(city.GetIncoming()[0])
				city.GetOccupied().SetOrigin(city)
				city.ClearIncoming()
			}
		}
	}
}

// allDead - checks whether all aliens are dead
func (world *World) allDeadOrStuck() bool {
	for _, alien := range world.aliens {
		if !alien.IsDead() && !alien.IsStuck() {
			return false
		}
	}
	return true
}

// allReachedIterations - checks whether all aliens have reached max iterations
func (world *World) allReachedIterations() bool {
	for _, alien := range world.aliens {
		if !alien.IsDead() && alien.GetIteration() < world.iterations {
			return false
		}
	}
	return true
}

// atLeastOneCollision - checks whether at least one city has collision
func (world *World) atLeastOneCollision() bool {
	for _, city := range world.cities {
		if city.IsDestroyed() {
			continue
		}
		if city.AliensCollide() {
			return true
		}
	}
	return false
}

// NewWorld - constructor
func NewWorld(iterations int, alienMap map[int]Alien, cityMap map[string]City) *World {
	aliens := make([]Alien, 0, len(alienMap))
	for _, a := range alienMap {
		aliens = append(aliens, a)
	}
	sort.Slice(aliens, func(i, j int) bool {
		return aliens[i].GetName() < aliens[j].GetName()
	})

	cities := make([]City, 0, len(cityMap))
	for _, c := range cityMap {
		cities = append(cities, c)
	}
	sort.Slice(cities, func(i, j int) bool {
		return cities[i].GetName() < cities[j].GetName()
	})
	return &World{iterations: iterations, alienMap: alienMap, cityMap: cityMap, aliens: aliens, cities: cities}
}
