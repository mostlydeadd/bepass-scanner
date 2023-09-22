package scanner

import (
	"fmt"
	"github.com/kaveh-ahangar/cfscanner/internal/config"
	"github.com/kaveh-ahangar/cfscanner/internal/logger"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"strings"
	"time"
)

func pingTechnique(activeIPChan chan<- string) {
	if config.G.Verbose {
		logger.Log("Starting ping scan", "Info")
	}

	if config.G.IpList != "" {
		ipList := strings.Split(config.G.IpList, "\n")
		for _, ip := range ipList {
			pingIP(ip, activeIPChan)
		}
	} else if config.G.Ip != "" {
		pingIP(config.G.Ip, activeIPChan)
	}
}

func pingIP(ip string, activeIPChan chan<- string) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		logger.Log(fmt.Sprintf("Invalid IP address: %s", ip), "Error")
		return
	}

	timeout := config.G.PingTimeout
	if config.G.Verbose {
		logger.Log(fmt.Sprintf("Pinging %s with a timeout of %s...", ip, timeout), "Info")
	}

	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		logger.Log(fmt.Sprintf("Failed to create ICMP packet listener for %s: %v", ip, err), "Error")
		return
	}
	defer conn.Close()

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte(""),
		},
	}

	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		logger.Log(fmt.Sprintf("Failed to marshal ICMP message for %s: %v", ip, err), "Error")
		return
	}

	startTime := time.Now()

	_, err = conn.WriteTo(msgBytes, &net.IPAddr{IP: parsedIP})
	if err != nil {
		logger.Log(fmt.Sprintf("Failed to send ICMP packet to %s: %v", ip, err), "Error")
		return
	}

	conn.SetReadDeadline(time.Now().Add(timeout))

	response := make([]byte, 1500)
	_, _, err = conn.ReadFrom(response)
	if err != nil {
		logger.Log(fmt.Sprintf("No response from %s: %v", ip, err), "Info")
		return
	}

	elapsedTime := time.Since(startTime)
	if config.G.Verbose {
		logger.Log(fmt.Sprintf("Received ICMP reply from %s in %s", ip, elapsedTime), "Info")
	}

	activeIPChan <- ip
}
