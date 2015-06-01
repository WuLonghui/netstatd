package discovery

import (
	. "netstatd/namespace"
)

type Discovery interface {
	GetNamespace(containerId string) (*Namespace, error)
	ListAllNamespaces() []*Namespace
}
