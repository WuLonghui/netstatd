package netstatd

import (
	"fmt"

	"github.com/coreos/go-namespaces/namespace"
)

const (
	CURRENT_NS_PID = 0
)

type Namespace struct {
	Pid int
}

func NewNamespace(pid int) *Namespace {
	return &Namespace{
		Pid: pid,
	}
}

func (n Namespace) String() string {
	return fmt.Sprintf("pid: %d", n.Pid)
}

func (n Namespace) Set() error {
	if n.Pid == CURRENT_NS_PID {
		return nil
	}

	fd, err := namespace.OpenProcess(n.Pid, namespace.CLONE_NEWNET)
	defer namespace.Close(fd)

	if err != nil {
		return err
	}

	// Join the container namespace
	errno := namespace.Setns(fd, namespace.CLONE_NEWNET)
	if errno != 0 {
		return fmt.Errorf("error setting namespace")
	}

	return nil
}

func (n Namespace) Exist() bool {
	fd, err := namespace.OpenProcess(n.Pid, namespace.CLONE_NEWNET)
	defer namespace.Close(fd)

	if err != nil {
		return false
	}

	return true
}
