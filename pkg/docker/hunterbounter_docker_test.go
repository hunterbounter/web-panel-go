package docker

import (
	"fmt"
	"os"
	"testing"
)

// Usage : go test -run TestNewDockerContainer -v
func TestNewDockerContainer(t *testing.T) {
	manager := NewDockerManager()

	// Example usage: Run a new container with specified options
	imageName := "hunter_bounter_zapv1"
	user := "root"
	port := "5002:5002"
	dns := "1.1.1.1"

	containerID, err := manager.RunContainer(imageName, user, port, dns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running container: %v\n", err)
		return
	}
	fmt.Printf("New container ID: %s\n", containerID)
}

// go test -run TestDockerManager_StopContainer -v
func TestDockerManager_StopContainer(t *testing.T) {
	manager := NewDockerManager()

	// Example usage: Stop a container by ID
	containerID := "1b3fb3400225"
	err := manager.StopContainer(containerID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error stopping container: %v\n", err)
		return
	}
	fmt.Printf("Container stopped: %s\n", containerID)
}
