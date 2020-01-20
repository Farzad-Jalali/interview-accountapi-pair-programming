package dockercompose

import (
	"os"
	"strconv"

	"github.com/phayes/freeport"
)

type freePortMap = map[string]int

func newFreePortMap(portNames ...DynamicPort) freePortMap {

	ports := make(map[string]int, len(portNames))

	for _, portName := range portNames {
		freePort, err := freeport.GetFreePort()
		if err != nil {
			panic(err)
		}
		ports[portName.PortName] = freePort
		os.Setenv(portName.PortName, strconv.Itoa(freePort))
	}

	return ports

}
