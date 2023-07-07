// Package testcontainermongo provides an easy way to start a mongo testcontainer using docker
package testcontainermongo

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Option func(testcontainers.ContainerRequest) testcontainers.ContainerRequest

// WithInit adds a path to a folder containing sql files to be executed on startup
func WithInit(init string) Option {
	return func(req testcontainers.ContainerRequest) testcontainers.ContainerRequest {
		req.Mounts = testcontainers.Mounts(testcontainers.ContainerMount{
			Source: testcontainers.GenericBindMountSource{
				HostPath: init,
			},
			Target: testcontainers.ContainerMountTarget("/docker-entrypoint-initdb.d"),
		})
		return req
	}
}

// WithDb adds a database name to the container request
func WithDb(db string) Option {
	return func(req testcontainers.ContainerRequest) testcontainers.ContainerRequest {
		if req.Env == nil {
			req.Env = make(map[string]string)
		}
		req.Env["MONGO_INITDB_DATABASE"] = db
		return req
	}
}

// WithAuth adds a username and password to the container request
func WithAuth(user, pass string) Option {
	return func(req testcontainers.ContainerRequest) testcontainers.ContainerRequest {
		if req.Env == nil {
			req.Env = make(map[string]string)
		}
		req.Env["MONGO_INITDB_ROOT_USERNAME"] = user
		req.Env["MONGO_INITDB_ROOT_PASSWORD"] = pass
		return req
	}
}

// WithEnv replaces the environment variables of the container request
func WithEnv(env map[string]string) Option {
	return func(req testcontainers.ContainerRequest) testcontainers.ContainerRequest {
		req.Env = env
		return req
	}
}

func connectionString(host, port string, env map[string]string) string {
	if env == nil {
		return fmt.Sprintf("mongodb://%s:%s", host, port)
	}
	db, dbOk := env["MONGO_INITDB_DATABASE"]
	user, userOk := env["MONGO_INITDB_ROOT_USERNAME"]
	password, passwordOk := env["MONGO_INITDB_ROOT_PASSWORD"]

	credentials := user
	if passwordOk {
		credentials += fmt.Sprintf(":%s", password)
	}
	if userOk {
		credentials += "@"
	}

	if dbOk {
		return fmt.Sprintf("mongodb://%s%s:%s/%s", credentials, host, port, db)
	}
	return fmt.Sprintf("mongodb://%s%s:%s", credentials, host, port)
}

// New setup a mongo testcontainer
func New(ctx context.Context, tag string, opts ...Option) (testcontainers.Container, string, error) {
	const (
		name = "test_db"
		user = "root"
		pass = "example"
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

	// Apply config options
	for _, opt := range opts {
		req = opt(req)
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

	conn := connectionString(host, port.Port(), req.Env)

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
