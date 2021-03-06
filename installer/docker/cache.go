package docker

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// RestoreComposerCache - Create a volume on the local client containing the backup
func RestoreComposerCache() {
	fmt.Println()
	fmt.Println("Building a Tokaido composer cache volume")
	ctx := context.Background()

	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to find Docker installed on your system. Have you run the Docker installer included in this package?")
		log.Fatal(err)
	}

	options := volumetypes.VolumeCreateBody{
		Name: "tok_composer_cache",
	}
	_, err = cli.VolumeCreate(ctx, options)
	if err != nil {
		log.Fatal(err)
	}

	ex, err := os.Executable()
	if err != nil {
		fmt.Println("Unable to determine the running binary directory: ", err.Error())
		os.Exit(1)
	}
	exPath := filepath.Dir(ex)
	cp := filepath.Join(exPath, "composer")

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "loomchild/volume-backup",
		Cmd:   []string{"restore", "tok_composer_cache.tar.bz2"},
		Tty:   true,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			mount.Mount{
				Source:   "tok_composer_cache",
				Target:   "/volume",
				Type:     "volume",
				ReadOnly: false,
			},
			mount.Mount{
				Source:   cp,
				Target:   "/backup",
				Type:     "bind",
				ReadOnly: false,
			},
		},
	}, &network.NetworkingConfig{}, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
	if err != nil {
		panic(err)
	}
}
