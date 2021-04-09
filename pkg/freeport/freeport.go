package freeport

import (
	"net"
	"sync"
)

var (
	freePorts = map[string]int{}
	lock      = &sync.RWMutex{}
)

func Get(ip, podKey string) (int, error) {
	port := getPortByKey(podKey)
	if port != nil {
		return *port, nil
	}
	addr, err := net.ResolveTCPAddr("tcp", ip+":0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	p := l.Addr().(*net.TCPAddr).Port
	savePort(podKey, p)
	return p, nil
}

func getPortByKey(podKey string) *int {
	lock.RLock()
	defer lock.RUnlock()
	key, found := freePorts[podKey]
	if found {
		return &key
	}
	return nil
}
func savePort(podKey string, port int) {
	lock.Lock()
	defer lock.Unlock()
	freePorts[podKey] = port
}
