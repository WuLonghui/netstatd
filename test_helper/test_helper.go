package testhelper

import (
	. "github.com/onsi/gomega"

	"github.com/fsouza/go-dockerclient"
)

var dockerClient *docker.Client

func InitDockerClient() {
	if dockerClient == nil {
		var err error
		dockerHost := "unix:///run/docker.sock"
		dockerClient, err = docker.NewClient(dockerHost)
		Expect(err).ShouldNot(HaveOccurred())
	}
}

func DockerRunContainer() *docker.Container {
	config := docker.Config{
		Image: "ubuntu",
		Cmd:   []string{"/bin/sh", "-c", "while true; do echo hello world; sleep 1; done"},
	}
	opts := docker.CreateContainerOptions{Config: &config}

	container, err := dockerClient.CreateContainer(opts)
	Expect(err).ShouldNot(HaveOccurred())

	err = dockerClient.StartContainer(container.ID, &docker.HostConfig{})
	Expect(err).ShouldNot(HaveOccurred())

	container, err = dockerClient.InspectContainer(container.ID)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(container.State.Running).Should(BeTrue())

	return container
}

func DockerCleanAllContainers() {
	containers, err := dockerClient.ListContainers(docker.ListContainersOptions{All: true})
	Expect(err).ShouldNot(HaveOccurred())
	for _, c := range containers {
		opts := docker.RemoveContainerOptions{ID: c.ID, Force: true}
		err := dockerClient.RemoveContainer(opts)
		Expect(err).ShouldNot(HaveOccurred())
	}
}
