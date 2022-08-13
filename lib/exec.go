package lib

import (
	"archive/tar"
	"bytes"
	"context"
	"gotribe/compiler/lib/log"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func DockerRun(image string, code string, dest string, cmd string) string {
	// log.DefaultLogger.Info("DockerRun-------------:")
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{"sh", "-c", cmd},
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	err = tw.WriteHeader(&tar.Header{
		Name: dest,             // filename
		Mode: 0777,             // permissions
		Size: int64(len(code)), // filesize
	})
	if err != nil {
		// panic(err)
		log.DefaultLogger.Error("docker copy err:", err)
	}
	tw.Write([]byte(code))
	tw.Close()

	// use &buf as argument for content in CopyToContainer
	cli.CopyToContainer(context.Background(), resp.ID, ".", &buf, types.CopyToContainerOptions{})

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
	defer out.Close()
	output, _ := ioutil.ReadAll(out)
	log.DefaultLogger.Errorf("output-------------%q+:", output)
	return string(output)
}
