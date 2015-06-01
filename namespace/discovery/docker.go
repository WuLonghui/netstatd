package discovery

import (
	"fmt"
	"log"
	"os"

	"github.com/fsouza/go-dockerclient"

	. "netstatd/namespace"
)

type DockerDiscovery struct {
	client *docker.Client
}

func NewDockerDiscovery() *DockerDiscovery {
	dockerHost := os.Getenv("DOCKER_HOST")
	if dockerHost == "" {
		dockerHost = "unix:///run/docker.sock"
	}

	client, err := docker.NewClient(dockerHost)
	if err != nil {
		log.Fatal(err)
	}

	return &DockerDiscovery{
		client: client,
	}
}

func (d DockerDiscovery) ListAllNamespaces() []*Namespace {
	namespaces := make([]*Namespace, 0)
	containers, err := d.client.ListContainers(docker.ListContainersOptions{Filters: map[string][]string{"status": {"paused", "running"}}})
	if err != nil {
		log.Println("error docker list containers, %v", err)
		return namespaces
	}

	for _, c := range containers {
		container, err := d.client.InspectContainer(c.ID)
		if err != nil {
			log.Println("error docker inspect container, %v", err)
			continue
		}
		namespace := NewNamespace(container.State.Pid, fmt.Sprintf("docker<%s>", c.ID[0:12]))
		namespaces = append(namespaces, namespace)
	}

	return namespaces
}

func (d DockerDiscovery) GetNamespace(id string) (*Namespace, error) {
	container, err := d.client.InspectContainer(id)
	if err != nil {
		return nil, err
	}

	namespace := NewNamespace(container.State.Pid, fmt.Sprintf("docker<%s>", container.ID[0:12]))
	return namespace, nil
}
