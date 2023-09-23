package scanner

import (
	"bufio"
	"fmt"
	"github.com/kaveh-ahangar/cfscanner/internal/config"
	"github.com/kaveh-ahangar/cfscanner/internal/logger"
	"net"
	"net/netip"
	"os"
	"strings"
)

func convertCIDRtoIPList() {
	if config.G.Cidr != "" {
		logger.Log(fmt.Sprintf("Converting %s CIDR to IP list", config.G.Cidr), "Info")
		var err error
		p, err := netip.ParsePrefix(config.G.Cidr)
		if err != nil {
			err = fmt.Errorf("invalid cidr: %s, error %v", config.G.Cidr, err)
		}

		p = p.Masked()

		addr := p.Addr()
		for {
			if !p.Contains(addr) {
				break
			}
			config.G.IpList += addr.String() + "\n"
			addr = addr.Next()
		}
		config.G.IpList = strings.TrimSpace(config.G.IpList)
	}

	if config.G.CidrList != "" {
		file, err := os.Open(config.G.CidrList)
		if err != nil {
			logger.Log(fmt.Sprintf("Failed to open CIDR list file: %v", err), "Error")
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			cidr := scanner.Text()
			ip, ipnet, err := net.ParseCIDR(cidr)
			if err != nil {
				logger.Log(fmt.Sprintf("Failed to parse CIDR from list: %v", err), "Error")
				continue
			}

			for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
				config.G.IpList += ip.String() + "\n"
			}
		}
		config.G.IpList = strings.TrimSpace(config.G.IpList)
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
