package config

import "fmt"

func printHelp() {
	helpText := `
    Goscan - A simple network scanner tool in Go
    Usage:
      goscan [flags]
    Flags:
      -c, --cidr string           List of CIDR
      -C, --cidr-list string      List of CIDR
      -i, --ip string             Single IP
      -I, --ip-list string        List of IP
      -o, --output string         Output file
      -p, --ports string           Comma-separated list of ports. default: 80,443,22
          --ping-timeout duration Ping timeout in milliseconds (default 400ms)
          --portscan-timeout duration
                                   Port-scan timeout in milliseconds (default 400ms)
          --full                  Runs full mode
          --ping                  Runs only ping mode
          --portscan              Runs only portscan mode
          --ptr                   Runs PTR scan
      -s, --silent                Silent mode
      -t, --threads int           Number of threads (default 5)
      -v, --verbose               Verbose mode. If set, it shows according to which technique the IP is active.
      -h, --help                  Print this help menu
    `
	fmt.Println(helpText)
}
