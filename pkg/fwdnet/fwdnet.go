package fwdnet

import (
	"errors"
	"fmt"
	"github.com/txn2/kubefwd/pkg/fwdIp"
	"net"
	"os"
	"os/exec"
)

// ReadyInterface prepares a local IP address on
// the loopback interface.
func ReadyInterface(svcName string, podName string, clusterN int, namespaceN int, port string, interfaceName string) (net.IP, error) {
	if len(interfaceName) == 0 {
		interfaceName = "lo"
	}
	ip, _ := fwdIp.GetIp(svcName, podName, clusterN, namespaceN)

	networkInterface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return net.IP{}, err
	}

	addrs, err := networkInterface.Addrs()
	if err != nil {
		return net.IP{}, err
	}

	// check the addresses already assigned to the interface
	for _, addr := range addrs {

		// found a match
		if addr.String() == ip.String()+"/8" {
			// found ip, now check for unused port
			conn, err := net.Dial("tcp", ip.String()+":"+port)
			if err != nil {
				return ip, nil
			}
			_ = conn.Close()
		}
	}

	// ip is not in the list of addrs for networkInterface
	cmd := "ifconfig"
	args := []string{interfaceName, "alias", ip.String(), "up"}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Println("Cannot ifconfig " + interfaceName + " alias " + ip.String() + " up")
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", ip.String()+":"+port)
	if err != nil {
		return ip, nil
	}
	_ = conn.Close()

	return net.IP{}, errors.New("unable to find an available IP/Port")
}

// RemoveInterfaceAlias can remove the Interface alias after port forwarding.
// if -alias command get err, just print the error and continue.
func RemoveInterfaceAlias(ip net.IP, interfaceName string) {
	cmd := "ifconfig"
	args := []string{interfaceName, "-alias", ip.String()}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		// suppress for now
		// @todo research alternative to ifconfig
		// @todo suggest ifconfig or alternative
		// @todo research libs for interface management
		//fmt.Println("Cannot ifconfig lo0 -alias " + ip.String() + "\r\n" + err.Error())
	}
}
