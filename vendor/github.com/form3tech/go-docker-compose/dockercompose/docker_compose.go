package dockercompose

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/lookup"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
	retry "github.com/giantswarm/retry-go"
	log "github.com/sirupsen/logrus"
)

const (
	defaultTimeoutDuration = 2 * time.Minute
)

type dockerCompose struct {
	p                   project.APIProject
	c                   project.Context
	containerPorts      freePortMap
	ctxContext          *ctx.Context
	transientContainers []string
}

type DynamicPort struct {
	PortName      string
	ContainerName string
	InternalPort  int
}

type DockerCompose interface {
	Start(w *waitFor) error
	Stop(services ...string)
	GetDynamicContainerPort(portVariable string) int
}

func NewDynamicPorts(ports ...string) ([]DynamicPort, error) {

	result := []DynamicPort{}

	for _, port := range ports {
		s := strings.Split(port, ":")
		if len(s) != 3 {
			return nil, fmt.Errorf("you need to supply ports in format PORT_NAME:CONTAINER_NAME:PORT eg SQS_PORT:sqs-1:1212")
		}
		internalPort, err := strconv.Atoi(s[2])
		if err != nil {
			return nil, fmt.Errorf("port must be an int")
		}
		result = append(result, DynamicPort{PortName: s[0], ContainerName: s[1], InternalPort: internalPort})
	}

	return result, nil

}

func (dc *dockerCompose) getPermanentServices(transientContainers []string) []string {
	services := make([]string, 0)
	for _, service := range dc.p.(*project.Project).ServiceConfigs.Keys() {
		isTransiant := false
		for _, transiant := range transientContainers {
			if service == transiant {
				isTransiant = true
				break
			}
		}

		if !isTransiant {
			services = append(services, service)
		}
	}
	return services
}

func (dc *dockerCompose) AreContainersAlreadyRunning(transientContainers ...string) bool {
	services := dc.getPermanentServices(transientContainers)
	allContainers, _ := dc.p.Containers(context.Background(), project.Filter{State: project.AnyState}, services...)
	runningContainers, _ := dc.p.Containers(context.Background(), project.Filter{State: project.Running}, services...)

	return len(runningContainers) > 0 && len(allContainers) == len(runningContainers)
}

func NewDockerCompose(repositoryAuth RepositoryAuth, composeFilePath, projectName string, ports []DynamicPort, transientContainers ...string) (DockerCompose, error) {

	dc := &dockerCompose{
		transientContainers: transientContainers,
	}

	c := project.Context{
		ComposeFiles:      []string{composeFilePath},
		ProjectName:       projectName,
		EnvironmentLookup: &lookup.OsEnvLookup{},
	}

	dc.ctxContext = &ctx.Context{
		Context: c,
	}

	if repositoryAuth != nil {
		configDir, err := repositoryAuth.GenerateAuthConfig()

		if err != nil {
			return nil, err
		}
		dc.ctxContext.ConfigDir = configDir

		defer os.RemoveAll(configDir)
	}

	dc.containerPorts = newFreePortMap(ports...)
	p, err := docker.NewProject(dc.ctxContext, nil)

	dc.c = c
	dc.p = p

	if err != nil {
		return nil, fmt.Errorf("could not create docker project: %v", err)
	}

	if dc.AreContainersAlreadyRunning(dc.transientContainers...) {

		client := dc.ctxContext.ClientFactory.Create(nil)

		portMap := make(freePortMap, len(ports))

		allContainers, err := dc.p.Containers(context.Background(), project.Filter{State: project.AnyState})

		if err != nil {
			return nil, err
		}

		for _, containerId := range allContainers {
			containerInfo, _ := client.ContainerInspect(context.Background(), containerId)

			for port, portBinding := range containerInfo.NetworkSettings.Ports {
				if portBinding == nil || len(portBinding) == 0 {
					continue
				}
				for _, p := range ports {
					expectedContainerName := "/" + projectName + "_" + p.ContainerName + "_1"
					if expectedContainerName == containerInfo.Name && p.InternalPort == port.Int() {
						internalPort, err := nat.ParsePort(portBinding[0].HostPort)
						if err != nil {
							return nil, err
						}
						portMap[p.PortName] = internalPort
						os.Setenv(p.PortName, fmt.Sprintf("%d", internalPort))
					}
				}

			}
		}

		dc.containerPorts = portMap

	}

	return dc, nil

}

func (dc *dockerCompose) Start(w *waitFor) error {

	if dc.AreContainersAlreadyRunning(dc.transientContainers...) {
		return nil
	}

	dc.p.Down(context.Background(), options.Down{})
	err := dc.p.Up(context.Background(), options.Up{})

	if err != nil {
		return fmt.Errorf("could not start docker compose, error: %v", err)
	}

	services := dc.getPermanentServices(dc.transientContainers)
	containers, err := dc.p.Containers(context.Background(), project.Filter{State: project.AnyState}, services...)

	if err != nil {
		return err
	}

	numberOfContainers := len(containers)

	err = retry.Do(func() error {
		cs, err := dc.p.Containers(context.Background(), project.Filter{State: project.Running}, services...)

		if err != nil {
			return err
		}

		if len(cs) != numberOfContainers {
			return errors.New("container not started")
		}

		log.Infof("%d containers started with ids %s", len(cs), strings.Join(cs, ","))

		return nil

	}, retry.MaxTries(120), retry.Sleep(time.Duration(250*time.Millisecond)))

	if err != nil {
		return fmt.Errorf("could not query containers, error: %v", err)
	}

	containerIds, err := dc.p.Containers(context.Background(), project.Filter{State: project.AnyState})

	if err != nil {
		return fmt.Errorf("could not get container ids, error: %v", err)
	}

	err = w.waitForContainers(dc.c.ProjectName, dc.ctxContext.ClientFactory.Create(nil), containerIds)

	return err
}

func (dc *dockerCompose) Stop(services ...string) {
	dc.p.Down(context.Background(), options.Down{}, services...)
}

func (dc *dockerCompose) GetDynamicContainerPort(portVariable string) int {
	return dc.containerPorts[portVariable]
}
