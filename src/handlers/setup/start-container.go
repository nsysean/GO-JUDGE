package setup_services

import (
	"context"
	"log"

	models "app/src/models"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func StartContainer(params models.Submission, cli *client.Client) string {
	ctx := context.Background()
	log.Println("Starting container")

	container, err := cli.ContainerCreate(ctx, &container.Config{
		Image: params.Language,
		Cmd: []string{"tail", "-f", "/dev/null"},
	}, &container.HostConfig{
		NetworkMode: "none",
		Resources:   container.Resources{Memory: 128e7},
	}, nil, nil, "",
	)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, container.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return container.ID
}
