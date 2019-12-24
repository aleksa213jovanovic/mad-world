package internal

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

// PlainFormatter - Basic formatter for logrus
type PlainFormatter struct {
	TimestampFormat string
	LevelDesc       []string
}

// Format - formats one log message
func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprint(entry.Message)), nil
}

// captureOutput - captures standard output to string
func captureOutput(f func()) string {
	formatter := new(PlainFormatter)
	log.SetFormatter(formatter)
	log.SetLevel(log.DebugLevel)

	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stdout)
	return buf.String()
}

// TestWorlds - run many tests provided in ../tests directory
func TestWorlds(t *testing.T) {
	files, err := ioutil.ReadDir("../tests")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Printf("Testing: %s\n", f.Name())
		citiesPath := fmt.Sprintf("../tests/%s/cities.txt", f.Name())
		aliensPath := fmt.Sprintf("../tests/%s/aliens.txt", f.Name())
		content, err := ioutil.ReadFile(fmt.Sprintf("../tests/%s/expected.txt", f.Name()))
		expected := string(content)

		builder := NewFileWorldBuilder(citiesPath, aliensPath)

		builder, err = builder.BuildCities()
		if err != nil {
			assert.Equal(t, expected, err.Error())
			continue
		}
		builder, err = builder.BuildAliens()
		if err != nil {
			assert.Equal(t, expected, err.Error())
			continue
		}

		world, err := builder.Build(10)
		if err != nil {
			assert.Equal(t, expected, err.Error())
			continue
		}

		assert.Equal(t, expected, captureOutput(func() {
			world.RunSimulation()
		}))
	}
}
