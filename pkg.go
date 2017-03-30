package pkg

import (
	"bytes"
	"errors"
	"fmt"

	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/packer/packer"
)

type Pkg struct {
	spec   *Spec
	cancel bool
}

// Prepare is called with a set of configurations to setup the
// internal state of the provisioner. The multiple configurations
// should be merged in some sane way.
func (p *Pkg) Prepare(input ...interface{}) error {
	return mapstructure.Decode(input, p.spec)
}

// Provision is called to actually provision the machine. A UI is
// given to communicate with the user, and a communicator is given that
// is guaranteed to be connected to some machine so that provisioning
// can be done.
func (p *Pkg) Provision(ui packer.Ui, comm packer.Communicator) error {
	sys, err := ProbeSystem(comm)
	if err != nil {
		return err
	}

	var cmd []string

	switch sys.id {
	case "fedora":
		cmd = append(cmd, "dnf -y install")
	case "debian":
		fallthrough
	case "ubuntu":
		cmd = append(cmd, "apt -y install")
	default:
		return fmt.Errorf("os '%s' is not supported", sys.id)
	}

	if p.spec.file != "" {
		switch sys.id {
		case "fedora":
			_, err := runCommand(comm, strings.Join(append(cmd, p.spec.file), " "))
			if err == nil {
				return nil
			}
		}
	}

	if p.cancel {
		return nil
	}

	if p.spec.name != "" {
		_, err := runCommand(comm, strings.Join(append(cmd, p.spec.name), " "))
		if err == nil {
			return nil
		}
	}

	return errors.New("installation failed")
}

// Cancel is called to cancel the provisioning. This is usually called
// while Provision is still being called. The Provisioner should act
// to stop its execution as quickly as possible in a race-free way.
func (p *Pkg) Cancel() {
	p.cancel = true
}

func runCommand(comm packer.Communicator, command string) (string, error) {
	var cmd packer.RemoteCmd
	cmd.Command = command

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := comm.Start(&cmd); err != nil {
		return "", err
	}

	cmd.Wait()

	if cmd.ExitStatus != 0 {
		return "", fmt.Errorf("command '%s' failed", command)
	}

	return stdout.String(), nil
}
