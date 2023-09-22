package config

import (
	"flag"
	"fmt"
	"github.com/kaveh-ahangar/cfscanner/internal/logger"
	"github.com/spf13/pflag"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var G Config

const (
	defaultPorts = "443,8443,2096"
)

type Config struct {
	Ip              string
	IpList          string
	Cidr            string
	CidrList        string
	OutputFile      string
	Ports           string
	Threads         int
	Verbose         bool
	IsSilent        bool
	Help            bool
	PingTimeout     time.Duration
	PortscanTimeout time.Duration
	FullMode        bool
	PingMode        bool
	PortscanMode    bool
	TempDir         string
}

func init() {
	G = Config{
		Ports:           defaultPorts,
		Threads:         5,
		PingTimeout:     400 * time.Millisecond,
		PortscanTimeout: 400 * time.Millisecond,
	}
}

func InitFromFlags() {
	parseFlags()
	next()
}

func next() {
	validateInput()

	setupTempDir()
	convertCIDRtoIPList()
}

func parseFlags() {
	// Create a new Config struct with default values
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine) // Use pflag alongside the standard flag package
	pflag.StringVarP(&G.Ip, "ip", "i", "", "single IP")
	pflag.StringVarP(&G.IpList, "ip-list", "I", "", "list of IP")
	pflag.StringVarP(&G.Cidr, "cidr", "c", "", "list of CIDR")
	pflag.StringVarP(&G.CidrList, "cidr-list", "C", "", "list of CIDR")
	pflag.StringVarP(&G.OutputFile, "output", "o", "", "output file")
	pflag.StringVarP(&G.Ports, "ports", "p", defaultPorts, "Comma-separated list of ports. default: 80,443,22")
	pflag.DurationVar(&G.PingTimeout, "ping-timeout", 400*time.Millisecond, "Ping timeout in milliseconds")
	pflag.DurationVar(&G.PortscanTimeout, "portscan-timeout", 400*time.Millisecond, "Port-scan timeout in milliseconds")
	pflag.BoolVarP(&G.FullMode, "full", "", false, "Runs full mode")
	pflag.BoolVarP(&G.PingMode, "ping", "", true, "Runs only ping mode")
	pflag.BoolVarP(&G.PortscanMode, "portscan", "", false, "Runs only portscan mode")
	pflag.IntVarP(&G.Threads, "threads", "t", 5, "number of threads")
	pflag.BoolVarP(&G.Verbose, "verbose", "v", false, "verbose mode. If set, it shows according to which technique the IP is active.")
	pflag.BoolVarP(&G.IsSilent, "silent", "s", false, "silent mode")
	pflag.BoolVarP(&G.Help, "help", "h", false, "print this help menu")
	pflag.Parse()
	pflag.Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "ip", "ip-list", "cidr", "cidr-list", "output", "ports":
			// Validate and process flag values as needed
		case "verbose":
			G.Verbose = true
		case "silent":
			logger.Silent()
		case "help":
			G.Help = true
		case "full":
			G.FullMode = true
		case "ping":
			G.PingMode = true
		case "portscan":
			G.PortscanMode = true
		}
	})
}

func validateInput() {
	if G.Cidr == "" && G.CidrList == "" && G.Ip == "" && G.IpList == "" {
		logger.Log("You must specify at least one of the following flags: (-c | --cidr), (-i | --ip), (-I | --ip-list), (-C | --cidr-list).", "Error")
		logger.Log("Use -h or --help for more information", "Info")
		os.Exit(1)
	}

	if (G.Ip != "" && G.IpList != "") || (G.Ip != "" && G.Cidr != "") || (G.Ip != "" && G.CidrList != "") || (G.IpList != "" && G.Cidr != "") || (G.IpList != "" && G.CidrList != "") || (G.Cidr != "" && G.CidrList != "") {
		logger.Log("Incompatible flags detected. You can only use one of the following flags: (-i | --ip), (-I | --ip-list), (-c | --cidr), (-C | --cidr-list).", "Error")
		logger.Log("Use -h or --help for more information.", "Info")
		os.Exit(1)
	}

	if (G.CidrList != "" || G.IpList != "") && G.OutputFile == "" {
		logger.Log("You must specify an output file when using -ip-list or -cidr-list flags.", "Error")
		logger.Log("Use -h or --help for more information.", "Info")
		os.Exit(1)
	}

	if G.Threads <= 0 {
		logger.Log("Number of threads must be greater than 0.", "Error")
		os.Exit(1)
	}

	if G.PingTimeout <= 0 {
		logger.Log("Ping timeout must be greater than 0.", "Error")
		os.Exit(1)
	}

	if G.PortscanTimeout <= 0 {
		logger.Log("Portscan timeout must be greater than 0.", "Error")
		os.Exit(1)
	}

	if G.PortscanTimeout < G.PingTimeout {
		logger.Log("Portscan timeout cannot be less than ping timeout.", "Error")
		os.Exit(1)
	}

	if G.Help {
		printHelp()
		os.Exit(0)
	}
}

func setupTempDir() {
	G.TempDir = filepath.Join(os.TempDir(), "goscan-"+strconv.Itoa(rand.Int()))
	if err := os.MkdirAll(G.TempDir, os.ModePerm); err != nil {
		logger.Log(fmt.Sprintf("Failed to create temporary directory: %v", err), "Error")
		os.Exit(1)
	}
}
