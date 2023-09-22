package scanner

import (
	"fmt"
	"github.com/kaveh-ahangar/cfscanner/internal/config"
	"github.com/kaveh-ahangar/cfscanner/internal/logger"
	"net"
	"strconv"
	"strings"
)

func portScanTechnique(activeIPChan chan<- string) {
	if config.G.Verbose {
		logger.Log("Starting port scan", "Info")
	}

	if config.G.IpList != "" {
		ipList := strings.Split(config.G.IpList, "\n")
		for _, ip := range ipList {
			portScanIP(ip, activeIPChan)
		}
	} else if config.G.Ip != "" {
		portScanIP(config.G.Ip, activeIPChan)
	}
}

func portScanIP(ip string, activeIPChan chan<- string) {
	timeout := config.G.PortscanTimeout

	if config.G.Verbose {
		logger.Log(fmt.Sprintf("Scanning ports on %s with a timeout of %s...", ip, timeout), "Info")
	}

	for _, portStr := range strings.Split(config.G.Ports, ",") {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			logger.Log(fmt.Sprintf("Invalid port number: %s", portStr), "Error")
			continue
		}

		target := fmt.Sprintf("%s:%d", ip, port)

		conn, err := net.DialTimeout("tcp", target, timeout)
		if err != nil {
			if config.G.Verbose {
				logger.Log(fmt.Sprintf("Port %d on %s is closed", port, ip), "Info")
			}
			continue
		}
		conn.Close()

		if config.G.Verbose {
			logger.Log(fmt.Sprintf("Port %d on %s is open", port, ip), "Info")
		}

		activeIPChan <- ip
	}
}
