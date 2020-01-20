package dockercompose

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/ahmetalpbalkan/dlog"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

type containerToWaitFor struct {
	containerName string
	waitFunc      func(chan error)
}

type waitFor struct {
	client       client.APIClient
	projectName  string
	containerIds []string
	containers   []*containerToWaitFor
	timeout      time.Duration
}

func WaitForContainersToStart() *waitFor {
	return WaitForContainersToStartWithTimeout(defaultTimeoutDuration)
}

func WaitForContainersToStartWithTimeout(timeout time.Duration) *waitFor {
	return &waitFor{
		containers: []*containerToWaitFor{},
		timeout:    timeout,
	}
}

func (w *waitFor) Container(containerName string, waitFunc func(finish chan error)) *waitFor {
	c := &containerToWaitFor{
		containerName: containerName,
		waitFunc:      waitFunc,
	}
	w.containers = append(w.containers, c)
	return w
}

func (w *waitFor) ContainerLogLine(containerName, logLine string, logSubscribers ...LogSubscriber) *waitFor {
	logSubscriber := newLogSubscriber(logSubscribers...)
	waitFunc := func(finish chan error) {
		containerId, e := w.getContainerIdFromName(containerName)
		if e != nil {
			finish <- e
		}
		e = w.waitForLogLine(logLine, containerId, logSubscriber)
		finish <- e
	}
	return w.Container(containerName, waitFunc)
}

func (w *waitFor) getContainerIdFromName(name string) (string, error) {

	for _, c := range w.containerIds {
		j, _ := w.client.ContainerInspect(context.Background(), c)
		if j.Name == ("/" + w.projectName + "_" + name) {
			return c, nil
		}
	}

	return "", nil

}

func (w *waitFor) waitForContainers(projectName string, client client.APIClient, containerIds []string) error {

	w.client = client
	w.containerIds = containerIds
	w.projectName = projectName

	var results []error
	finished := 0
	finish := make(chan error)

	for _, c := range w.containers {
		go c.waitFunc(finish)
	}

	for {
		res := <-finish
		finished++
		if res != nil {
			results = append(results, res)
		}

		if finished == len(w.containers) {
			break
		}
	}

	if len(results) > 0 {
		errorString := ""
		for i, error := range results {
			errorString += error.Error()
			if i != len(results)-1 {
				errorString += ", "
			}
		}
		return fmt.Errorf("continers didnt start: %s", errorString)
	}

	log.Info("All containers started")
	return nil
}

func (w *waitFor) waitForLogLine(line, containerId string, logSubscriber LogSubscriber) error {

	timeout := time.After(w.timeout)
	found := make(chan int)

	go func() {
		client := &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					conn, err := net.Dial("unix", "/var/run/docker.sock")
					if err != nil {
						return nil, fmt.Errorf("cannot connect docker socket: %v", err)
					}
					return conn, nil
				}}}
		url := fmt.Sprintf("http://-/containers/%s/logs?stdout=1&stderr=1&follow=1", containerId)
		resp, err := client.Get(url)
		if err != nil {
			log.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("unexpected status code: %s", resp.Status)
		}

		r := dlog.NewReader(resp.Body)
		s := bufio.NewScanner(r)
		for s.Scan() {
			logSubscriber.OnNext(s.Text())
			if strings.Contains(s.Text(), line) {
				found <- 1
				break
			}
		}
		if err := s.Err(); err != nil {
			log.Fatalf("read error: %v", err)
		}
	}()

	select {
	case <-timeout:
		return errors.New("timed out waiting for container to start")
	case <-found:
		log.Info("container started found log line")
		return nil
	}

}
