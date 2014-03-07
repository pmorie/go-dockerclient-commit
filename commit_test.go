package main

import (
	"log"
	"testing"

	"github.com/fsouza/go-dockerclient"
)

func removeContainer(client *docker.Client, id string) {
	client.RemoveContainer(docker.RemoveContainerOptions{ID: id})
}

func removeImage(client *docker.Client, name string) {
	client.RemoveImage(name)
}

func exerciseCommit(setCmd bool, t *testing.T) {
	testImage := "go-dc/commit-test"
	dockerClient, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		t.Error(err.Error())
	}
	removeImage(dockerClient, testImage)

	config := docker.Config{Image: "pmorie/sti-fake", Cmd: []string{"/bin/true"}}
	container, err := dockerClient.CreateContainer(docker.CreateContainerOptions{Name: "", Config: &config})
	if err != nil {
		t.Error(err.Error())
	}
	defer removeContainer(dockerClient, container.ID)

	err = dockerClient.StartContainer(container.ID, &docker.HostConfig{})
	if err != nil {
		t.Error(err.Error())
	}

	exitCode, err := dockerClient.WaitContainer(container.ID)
	if err != nil {
		t.Error(err.Error())
	}

	if exitCode != 0 {
		t.Error("Unexpected exit code")
	}

	config = docker.Config{}
	if setCmd {
		config.Cmd = []string{"/usr/bin/run"}
	}

	log.Printf("Commiting container with config: %+v\n", config)
	createdImage, err := dockerClient.CommitContainer(docker.CommitContainerOptions{Container: container.ID, Repository: testImage, Run: &config})
	if err != nil {
		t.Error(err.Error())
	}

	log.Printf("Committed image: %+v\n", createdImage)

	config = docker.Config{Image: testImage, AttachStdout: false, AttachStdin: false}
	container, err = dockerClient.CreateContainer(docker.CreateContainerOptions{Name: "", Config: &config})
	if err != nil {
		t.Error(err.Error())
	}
	defer removeContainer(dockerClient, container.ID)

	err = dockerClient.StartContainer(container.ID, &docker.HostConfig{})
	if err != nil {
		t.Error(err.Error())
	}

	exitCode, err = dockerClient.WaitContainer(container.ID)
	if err != nil {
		t.Error(err.Error())
	}

	if exitCode != 0 {
		t.Error("Unexpected exit code from created image")
	}

}

func TestCommitWithoutCmd(t *testing.T) {
	exerciseCommit(false, t)
}

func TestCommitWithCmd(t *testing.T) {
	exerciseCommit(true, t)
}
