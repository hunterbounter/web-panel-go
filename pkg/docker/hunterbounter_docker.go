package docker

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

// DockerManager is a struct for managing Docker operations.
type DockerManager struct{}

// NewDockerManager creates a new DockerManager instance.
func NewDockerManager() *DockerManager {
	return &DockerManager{}
}

// ListContainers lists all containers.
func (d *DockerManager) ListContainers() (string, error) {
	cmd := exec.Command("docker", "ps", "-a")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func (d *DockerManager) ListRunningContainersViaImageName(imageName string) ([]string, error) {
	cmd := exec.Command("docker", "ps", "--filter", "ancestor="+imageName, "--format", "{{.ID}}")
	log.Println("Full command: ", cmd.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	// Split the output by newlines to get individual container IDs
	output := strings.TrimSpace(out.String())
	if output == "" {
		return []string{}, nil
	}

	containerIDs := strings.Split(output, "\n")

	// Check if containerIDs actually contains IDs
	if len(containerIDs) == 0 || (len(containerIDs) == 1 && containerIDs[0] == "") {
		return []string{}, nil
	}

	return containerIDs, nil
}

// StartContainer starts a container by ID.
func (d *DockerManager) StartContainer(containerID string) error {
	cmd := exec.Command("docker", "start", containerID)
	return cmd.Run()
}

// StopContainer stops a container by ID with a specified timeout.
func (d *DockerManager) StopContainer(containerID string) error {
	cmd := exec.Command("docker", "stop", "-t", "0", containerID)
	log.Println("Full command: ", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Command output:", string(output))
		return err
	}
	return nil
}

// RemoveContainer removes a container by ID.
func (d *DockerManager) RemoveContainer(containerID string) error {
	cmd := exec.Command("docker", "rm", "-f", containerID)
	return cmd.Run()
}

// RunContainer runs a new container with specified options.
func (d *DockerManager) RunContainer(imageName, user, port, dns string) (string, error) {
	cmd := exec.Command("docker", "run", "-u", user, "--dns", dns, "-d", imageName)
	//cmd := exec.Command("docker", "run", "-u", user, "-p", port, "--dns", dns, "-d", imageName)
	log.Println("Full command: ", cmd.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
