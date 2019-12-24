package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/cosmos/mad_alien_invasion/internal"
	"os"
)

// PlainFormatter - Basic formatter for logrus
type PlainFormatter struct {
	TimestampFormat string
	LevelDesc []string
}

// Format = formats one log message
func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprint(entry.Message)), nil
}

// main - application entry-point
func main() {
	iterations := flag.Int("iterations", 10000, "maximum iterations per alien")
	pathToFile := flag.String("path", "", "path to file containing cities")
	numberOfAliens := flag.Int( "aliens", 0, "number of aliens")
	debug := flag.Bool( "debug", false, "number of aliens")
	flag.Parse()

	formatter := new(PlainFormatter)
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	if *pathToFile == "" {
		log.Fatalln("please input path to file containing cities")
	}

	if *numberOfAliens <= 0 {
		log.Fatalln("please input correct number of aliens")
	}

	builder := internal.NewDefaultWorldBuilder(*pathToFile, *numberOfAliens)
	builder, err := builder.BuildCities()
	if err != nil {
		log.Fatalln(err)
	}
	builder, err = builder.BuildAliens()
	if err != nil {
		log.Fatalln(err)
	}

	world, err := builder.Build(*iterations)
	if err != nil {
		log.Fatalln(err)
	}
	world.RunSimulation()
}