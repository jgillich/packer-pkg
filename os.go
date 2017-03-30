package pkg

import (
	"errors"

	"strings"

	"github.com/mitchellh/packer/packer"
)

type System struct {
	id      string
	version string
}

func ProbeSystem(comm packer.Communicator) (*System, error) {
	rel, err := runCommand(comm, "cat /etc/os-release")
	if err == nil {
		return osRelease(rel)
	}

	return nil, errors.New("could not determine operating system")
}

// os-release specification https://www.freedesktop.org/software/systemd/man/os-release.html
func osRelease(rel string) (*System, error) {
	m := map[string]string{}
	for _, line := range strings.Split(rel, "\n") {
		split := strings.SplitN(line, "=", 2)

		if len(split) != 2 {
			return nil, errors.New("malformed os-release")
		}

		key := split[0]
		key = strings.Trim(key, " ")

		value := split[1]
		value = strings.Trim(value, " ")
		value = strings.Trim(value, "\"")

		m[key] = value
	}

	return &System{
		id:      m["ID"],
		version: m["VERSION_ID"],
	}, nil
}
