// Package testcontainermongo provides an easy way to start a mongo testcontainer using docker
package testcontainermongo

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// New setup a mongo testcontainer
func New(ctx context.Context, tag, init string) (testcontainers.Container, string, error) {
	const (
		name = "test_db"
		user = "root"
		pass = "password1234"
	)

	// Create mongoQL container request
	req := testcontainers.ContainerRequest{
		Image: "mongo:" + tag,
		Env: map[string]string{
			"MONGO_INITDB_DATABASE":      name,
			"MONGO_INITDB_ROOT_USERNAME": user,
			"MONGO_INITDB_ROOT_PASSWORD": pass,
		},
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp"),
	}
	if init != "" {
		req.Mounts = testcontainers.Mounts(testcontainers.ContainerMount{
			Source: testcontainers.GenericBindMountSource{
				HostPath: init,
			},
			Target: testcontainers.ContainerMountTarget("/docker-entrypoint-initdb.d/init.sh"),
		})
	}

	// Start mongoQL container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to start mongo container: %w", err)
	}

	// Get host and port of mongoQL container
	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get host: %w", err)
	}

	port, err := container.MappedPort(ctx, "27017")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get port: %w", err)
	}

	conn := fmt.Sprintf("mongo://%v:%v@%v:%v/%v", user, pass, host, port.Port(), name)

	// Create db connection string and connect
	return container, conn, nil
}

// Terminate terminates the container in a defer friendly way
func Terminate(ctx context.Context, c testcontainers.Container) {
	err := c.Terminate(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to terminate container: %v", err))
	}
}
