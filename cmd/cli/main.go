package main

import (
	"fmt"
	"github.com/kaveh-ahangar/cfscanner/internal/config"
	"github.com/kaveh-ahangar/cfscanner/internal/logger"
	"github.com/kaveh-ahangar/cfscanner/internal/scanner"
	"os"
)

func main() {
	config.InitFromFlags()

	fmt.Println("here")

	scannerEngine := scanner.NewScanner()
	scannerEngine.Run()

	logger.Log("End, good bye:)", "Print")
	removeTempDir(config.G.TempDir)
}

func removeTempDir(tempDir string) {
	if tempDir != "" {
		os.RemoveAll(tempDir)
	}
}
