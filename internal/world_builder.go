package internal

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/cosmos/mad_alien_invasion/internal/components"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// WorldBuilder - World instance builder
type WorldBuilder interface {
	// BuildAliens - build aliens
	BuildAliens() (WorldBuilder, error)
	// BuildCities - build cities
	BuildCities() (WorldBuilder, error)
	// Build - build a whole new world
	Build(maxIterations int) (*components.World, error)
}

// FileWorldBuilder - reads both cities and aliens from files
type FileWorldBuilder struct {
	aliensPath        string
	citiesPath        string
	citiesInitialized bool
	alienMap          map[int]components.Alien
	cityMap           map[string]components.City
}

// BuildAliens - reads aliens and directions from file path
func (builder *FileWorldBuilder) BuildAliens() (WorldBuilder, error) {
	if !builder.citiesInitialized {
		return nil, errors.New("cities should be initialized first")
	}
	file, err := os.Open(builder.aliensPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var parsed []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}
		parsed = strings.Fields(scanner.Text())
		pair := strings.Split(parsed[0], "=")
		city := builder.cityMap[pair[1]]
		if city == nil {
			return nil, errors.New(fmt.Sprintf("non existing city %s", pair[1]))
		}
		name, err := strconv.Atoi(pair[0])
		if err != nil {
			return nil, err
		}
		directionStrategy := components.NewBufferedDirectionStrategy(parsed[1:])
		alien := components.NewAlien(name, city, directionStrategy)
		builder.alienMap[name] = alien
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return builder, nil
}

// BuildCities - reads cities and directions from file path
func (builder *FileWorldBuilder) BuildCities() (WorldBuilder, error) {
	cities, err := loadCitiesFromFile(builder.citiesPath)
	if err != nil {
		return nil, err
	}
	builder.cityMap = cities
	builder.citiesInitialized = true
	return builder, nil
}

// Build - builds a whole new world
func (builder *FileWorldBuilder) Build(maxIterations int) (*components.World, error) {
	if len(builder.alienMap) > len(builder.cityMap) {
		logrus.Error("there are more aliens than cities")
		return nil, errors.New("there are more aliens than cities")
	}
	return components.NewWorld(maxIterations, builder.alienMap, builder.cityMap), nil
}

// NewFileWorldBuilder - constructor
func NewFileWorldBuilder(citiesPath string, aliensPath string) WorldBuilder {
	return &FileWorldBuilder{citiesPath: citiesPath, aliensPath: aliensPath, alienMap: make(map[int]components.Alien), cityMap: make(map[string]components.City)}
}

// DefaultWorldBuilder - reads cities from file and generates random aliens
type DefaultWorldBuilder struct {
	aliens            int
	citiesPath        string
	citiesInitialized bool

	alienMap map[int]components.Alien
	cityMap  map[string]components.City
}

// BuildAliens - reads aliens and directions from file path
func (builder *DefaultWorldBuilder) BuildAliens() (WorldBuilder, error) {
	if !builder.citiesInitialized {
		return nil, errors.New("cities not initialized")
	}

	if builder.aliens <= 0 {
		return nil, errors.New("number of aliens is lesser than 0")
	}
	if builder.aliens > len(builder.cityMap) {
		return nil, errors.New("number of aliens is greater than number of cities")
	}

	cities := make([]components.City, 0)
	for _, city := range builder.cityMap {
		cities = append(cities, city)
	}

	for i := 1; i <= builder.aliens; i++ {
		var index int
		if len(cities) == 1 {
			index = 0
		} else {
			s := rand.NewSource(time.Now().Unix())
			r := rand.New(s)
			index = r.Intn(len(cities) - 1)
		}

		builder.alienMap[i] = components.NewAlien(i, cities[index], components.NewRandomDirectionStrategy())
		cities = append(cities[:index], cities[index+1:]...)
	}
	return builder, nil
}

// BuildCities - reads cities and directions from file path
func (builder *DefaultWorldBuilder) BuildCities() (WorldBuilder, error) {
	cities, err := loadCitiesFromFile(builder.citiesPath)
	if err != nil {
		return nil, err
	}
	builder.cityMap = cities
	builder.citiesInitialized = true
	return builder, nil
}

// Build - builds a whole new world
func (builder *DefaultWorldBuilder) Build(maxIterations int) (*components.World, error) {
	if len(builder.alienMap) > len(builder.cityMap) {
		return nil, errors.New("there are more aliens than cities")
	}
	return components.NewWorld(maxIterations, builder.alienMap, builder.cityMap), nil
}

// NewDefaultWorldBuilder - constructor
func NewDefaultWorldBuilder(citiesPath string, aliens int) WorldBuilder {
	return &DefaultWorldBuilder{citiesPath: citiesPath, aliens: aliens, alienMap: make(map[int]components.Alien), cityMap: make(map[string]components.City)}
}

// loadCitiesFromFile - loads city structure from file
func loadCitiesFromFile(path string) (map[string]components.City, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var cities = make(map[string]components.City)
	var parsed []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}
		parsed = strings.Fields(scanner.Text())
		city := cities[parsed[0]]
		if city == nil {
			city = components.NewCity(parsed[0])
		}
		for _, directionText := range parsed[1:] {
			pair := strings.Split(directionText, "=")
			if len(pair) != 2 {
				return nil, errors.New("wrong cities format")
			}
			direction := strings.ToLower(pair[0])
			if !(contains( []string{"north", "south", "west", "east", "wait"}, direction)) {
				return nil, errors.New("wrong city direction format")
			}

			neighbor := cities[pair[1]]
			if neighbor == nil {
				neighbor = components.NewCity(pair[1])
				cities[pair[1]] = neighbor
			}
			city.AddNeighbor(strings.ToLower(direction), neighbor)
		}
		cities[parsed[0]] = city
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cities, nil
}

// contains - checks whether slice contains an element
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}