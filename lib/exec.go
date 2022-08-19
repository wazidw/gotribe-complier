package lib

import (
	"archive/tar"
	"bytes"
	"context"
	"gotribe/compiler/lib/log"
	"io/ioutil"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func DockerRun(image string, code string, dest string, cmd string) string {
	// log.DefaultLogger.Info("DockerRun-------------:")
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.DefaultLogger.Error("NewClientWithOpts:", err)
	}
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        image,
		Cmd:          []string{"sh", "-c", cmd},
		Tty:          true,
		AttachStderr: true,
		AttachStdout: true,
	}, &container.HostConfig{
		Resources: container.Resources{
			Memory: 100 * 1024 * 1024, // Minimum memory limit allowed is 6MB.
		},
	}, nil, nil, "")
	if err != nil {
		log.DefaultLogger.Error("ContainerCreate:", err)
	}

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	err = tw.WriteHeader(&tar.Header{
		Name: dest,             // filename
		Mode: 0777,             // permissions
		Size: int64(len(code)), // filesize
	})
	if err != nil {
		log.DefaultLogger.Error("docker copy err:", err)
	}
	tw.Write([]byte(code))
	tw.Close()

	// use &buf as argument for content in CopyToContainer
	cli.CopyToContainer(ctx, resp.ID, ".", &buf, types.CopyToContainerOptions{})

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.DefaultLogger.Errorf("ContainerStart err:%v:", err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	timeout := time.NewTimer(5 * time.Second)
	select {
	case waitBody := <-statusCh:
		log.DefaultLogger.Errorf("waitBody err:%v:", waitBody.StatusCode)
		break
	case errC := <-errCh:
		log.DefaultLogger.Errorf("statusCh err:%v:", errC)
	case <-timeout.C:
		log.DefaultLogger.Error("execute timeout")
		cli.ContainerKill(ctx, resp.ID, "SIGKILL")
		log.DefaultLogger.Error("ContainerKill")
		return "execute timeout"
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		log.DefaultLogger.Errorf("ContainerLogs err:%v:", err)
	}

	defer out.Close()
	output, _ := ioutil.ReadAll(out)
	log.DefaultLogger.Info("ContainerRemove err:%v:", output)
	return string(output)
}
