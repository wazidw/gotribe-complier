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

func DockerRun(image string, code string, dest string, cmd string, langTimeout int64, memory int64, cpuset string) string {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.DefaultLogger.Error("NewClientWithOpts:", err)
	}
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:           image,
		Cmd:             []string{"sh", "-c", cmd},
		Tty:             false,
		AttachStderr:    true,
		AttachStdout:    true,
		NetworkDisabled: true,
	}, &container.HostConfig{
		Resources: container.Resources{
			Memory:     memory, // Minimum memory limit allowed is 6MB.
			CpusetCpus: cpuset,
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
	if err := cli.CopyToContainer(ctx, resp.ID, ".", &buf, types.CopyToContainerOptions{}); err != nil {
		log.DefaultLogger.Errorf("CopyToContainer err:%v:", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.DefaultLogger.Errorf("ContainerStart err:%v:", err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	timeout := time.NewTimer(time.Duration(langTimeout) * time.Second)
	select {
	case waitBody := <-statusCh:
		log.DefaultLogger.Info("waitBody err:", waitBody.StatusCode)
		break
	case errC := <-errCh:
		log.DefaultLogger.Errorf("ContainerWait statusCh err:%v:", errC)
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
	log.DefaultLogger.Info("output:", string(output))
	return string(output)
}
