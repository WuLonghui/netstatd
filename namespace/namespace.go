package namespace

import (
	"fmt"

	"github.com/coreos/go-namespaces/namespace"
)

const (
	CURRENT_NAMESPACE_PID = 0
)

type Namespace struct {
	Pid int    `json:"pid"`
	Tag string `json:"tag"`
}

func NewNamespace(pid int, tag string) *Namespace {
	return &Namespace{
		Pid: pid,
		Tag: tag,
	}
}

func (n Namespace) String() string {
	return fmt.Sprintf("pid: %d, tag: %s", n.Pid, n.Tag)
}

func (n Namespace) Set() error {
	if n.Pid == CURRENT_NAMESPACE_PID {
		return nil
	}

	fd, err := namespace.OpenProcess(n.Pid, namespace.CLONE_NEWNET)
	defer namespace.Close(fd)

	if err != nil {
		return err
	}

	// Join the namespace
	errno := namespace.Setns(fd, namespace.CLONE_NEWNET)
	if errno != 0 {
		return fmt.Errorf("error setting net namespace")
	}

	return nil
}

func (n Namespace) Exist() bool {
	if n.Pid == CURRENT_NAMESPACE_PID {
		return true
	}

	fd, err := namespace.OpenProcess(n.Pid, namespace.CLONE_NEWNET)
	defer namespace.Close(fd)

	if err != nil {
		return false
	}

	return true
}
