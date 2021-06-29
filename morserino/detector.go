package morserino

import (
	"strings"

	"go.bug.st/serial/enumerator"
)

// Detector retrieves all connected morserino devices port names.
type Detector interface {
	Detect() ([]string, error)
}

func NewPortDetector(input string) Detector {
	switch strings.ToLower(input) {
	case "auto":
		return &autoDetector{listPortsFn: enumerator.GetDetailedPortsList}
	default:
		return &namedPortDetector{name: input}
	}
}

type namedPortDetector struct {
	name string
}

func (npd *namedPortDetector) Detect() ([]string, error) {
	// TODO(jly): naive implementation, but we could add some sanity checks here.
	// For example, actually check that the pointed device is a morserino device.
	return []string{npd.name}, nil
}

// Expected morserino devices identifiers.
const (
	deviceVID = "10C4"
	devicePID = "EA60"
)

type autoDetector struct {
	listPortsFn func() ([]*enumerator.PortDetails, error)
}

func (ad *autoDetector) Detect() ([]string, error) {
	ports, err := ad.listPortsFn()

	if err != nil {
		return nil, err
	}

	var morserinoPorts []string

	for _, port := range ports {
		if port.IsUSB && port.VID == deviceVID && port.PID == devicePID {
			morserinoPorts = append(morserinoPorts, port.Name)
		}
	}

	return morserinoPorts, nil
}
